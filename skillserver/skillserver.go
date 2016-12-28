package skillserver

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type EchoApplication struct {
	AppID          string
	Handler        func(http.ResponseWriter, *http.Request)
	OnLaunch       func(*EchoRequest, *EchoResponse)
	OnIntent       func(*EchoRequest, *EchoResponse)
	OnSessionEnded func(*EchoRequest, *EchoResponse)
}

type StdApplication struct {
	Methods string
	Handler func(http.ResponseWriter, *http.Request)
}

var Applications = map[string]interface{}{}

func Run(apps map[string]interface{}, port string) {
	router := mux.NewRouter()
	Init(apps, router)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + port)
}

func RunSSL(apps map[string]interface{}, port, cert, key string) error {
	router := mux.NewRouter()
	Init(apps, router)

	err := http.ListenAndServeTLS(port, cert, key, router)
	return err
}

func Init(apps map[string]interface{}, router *mux.Router) {
	Applications = apps

	// /echo/* Endpoints
	echoRouter := mux.NewRouter()
	// /* Endpoints
	pageRouter := mux.NewRouter()

	for uri, meta := range Applications {
		switch app := meta.(type) {
		case EchoApplication:
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				echoReq := r.Context().Value("echoRequest").(*EchoRequest)
				echoResp := NewEchoResponse()

				if echoReq.GetRequestType() == "LaunchRequest" {
					if app.OnLaunch != nil {
						app.OnLaunch(echoReq, echoResp)
					}
				} else if echoReq.GetRequestType() == "IntentRequest" {
					if app.OnIntent != nil {
						app.OnIntent(echoReq, echoResp)
					}
				} else if echoReq.GetRequestType() == "SessionEndedRequest" {
					if app.OnSessionEnded != nil {
						app.OnSessionEnded(echoReq, echoResp)
					}
				} else {
					http.Error(w, "Invalid request.", http.StatusBadRequest)
				}

				json, _ := echoResp.String()
				w.Header().Set("Content-Type", "application/json;charset=UTF-8")
				w.Write(json)
			}

			if app.Handler != nil {
				handlerFunc = app.Handler
			}

			echoRouter.HandleFunc(uri, handlerFunc).Methods("POST")
		case StdApplication:
			pageRouter.HandleFunc(uri, app.Handler).Methods(app.Methods)
		}
	}

	router.PathPrefix("/echo/").Handler(negroni.New(
		negroni.HandlerFunc(validateRequest),
		negroni.HandlerFunc(verifyJSON),
		negroni.Wrap(echoRouter),
	))

	router.PathPrefix("/").Handler(negroni.New(
		negroni.Wrap(pageRouter),
	))
}

func GetEchoRequest(r *http.Request) *EchoRequest {
	return r.Context().Value("echoRequest").(*EchoRequest)
}

func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Println(logMsg)
	}

	http.Error(w, err, errCode)
}

// Decode the JSON request and verify it.
func verifyJSON(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var echoReq *EchoRequest
	err := json.NewDecoder(r.Body).Decode(&echoReq)
	if err != nil {
		HTTPError(w, err.Error(), "Bad Request", 400)
		return
	}

	// Check the timestamp
	if !echoReq.VerifyTimestamp() && r.URL.Query().Get("_dev") == "" {
		HTTPError(w, "Request too old to continue (>150s).", "Bad Request", 400)
		return
	}

	// Check the app id
	if !echoReq.VerifyAppID(Applications[r.URL.Path].(EchoApplication).AppID) {
		HTTPError(w, "Echo AppID mismatch!", "Bad Request", 400)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "echoRequest", echoReq))

	next(w, r)
}

// Run all mandatory Amazon security checks on the request.
func validateRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Check for debug bypass flag
	devFlag := r.URL.Query().Get("_dev")

	isDev := devFlag != ""

	if !isDev {
		certURL := r.Header.Get("SignatureCertChainUrl")

		// Verify certificate URL
		if !verifyCertURL(certURL) && devFlag == "" {
			HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", 401)
			return
		}

		// Fetch certificate data
		certContents, err := readCert(certURL)
		if err != nil {
			HTTPError(w, err.Error(), "Not Authorized", 401)
			return
		}

		// Decode certificate data
		block, _ := pem.Decode(certContents)
		if block == nil {
			HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
			return
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			HTTPError(w, err.Error(), "Not Authorized", 401)
			return
		}

		// Check the certificate date
		if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
			HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
			return
		}

		// Check the certificate alternate names
		foundName := false
		for _, altName := range cert.Subject.Names {
			if altName.Value == "echo-api.amazon.com" {
				foundName = true
			}
		}

		if !foundName && devFlag == "" {
			HTTPError(w, "Amazon certificate invalid.", "Not Authorized", 401)
			return
		}

		// Verify the key
		publicKey := cert.PublicKey
		encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

		// Make the request body SHA1 and verify the request with the public key
		var bodyBuf bytes.Buffer
		hash := sha1.New()
		_, err = io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
		if err != nil {
			HTTPError(w, err.Error(), "Internal Error", 500)
			return
		}
		//log.Println(bodyBuf.String())
		r.Body = ioutil.NopCloser(&bodyBuf)

		err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
		if err != nil {
			HTTPError(w, "Signature match failed.", "Not Authorized", 401)
			return
		}
	}

	next(w, r)
}

func readCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("Could not download Amazon cert file.")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("Could not read Amazon cert file.")
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}
