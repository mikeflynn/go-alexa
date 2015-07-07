package skillserver

import (
	"bytes"
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
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type EchoApplication struct {
	AppID   string
	Handler func(http.ResponseWriter, *http.Request)
}

var Applications = map[string]EchoApplication{}

func Run(apps map[string]EchoApplication, port string) {
	Applications = apps

	router := mux.NewRouter()

	// /echo/* Endpoints
	echoRouter := mux.NewRouter()
	for uri, meta := range Applications {
		echoRouter.HandleFunc(uri, meta.Handler).Methods("POST")
	}

	router.PathPrefix("/echo/").Handler(negroni.New(
		negroni.HandlerFunc(validateRequest),
		negroni.HandlerFunc(verifyJSON),
		negroni.Wrap(echoRouter),
	))

	/*
		// /* Endpoints
		pageRouter := mux.NewRouter()
		pageRouter.HandleFunc("/", HomePage)
		pageRouter.HandleFunc("/about", AboutPage)

		router.PathPrefix("/").Handler(negroni.New(
			negroni.Wrap(pageRouter),
		))
	*/

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + port)
}

func GetEchoRequest(r *http.Request) *EchoRequest {
	return context.Get(r, "echoRequest").(*EchoRequest)
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
	if !echoReq.VerifyAppID(Applications[r.URL.Path].AppID) {
		HTTPError(w, "Echo AppID mismatch!", "Bad Request", 400)
		return
	}

	context.Set(r, "echoRequest", echoReq)

	next(w, r)
}

// Run all mandatory Amazon security checks on the request.
func validateRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Check for debug bypass flag
	devFlag := r.URL.Query().Get("_dev")

	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) && devFlag == "" {
		HTTPError(w, "", "Not Authorized", 401)
		return
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil && devFlag == "" {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil && devFlag == "" {
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
		return
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil && devFlag == "" {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		if devFlag == "" {
			HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
			return
		}
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
	if err != nil && devFlag == "" {
		HTTPError(w, "Signature match failed.", "Not Authorized", 401)
		return
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
	if !strings.HasSuffix(path, "/echo.api/echo-api-cert.pem") {
		return false
	}

	if !strings.HasPrefix(path, "https://s3.amazonaws.com/echo.api/") && !strings.HasPrefix(path, "https://s3.amazonaws.com:443/echo.api/") {
		return false
	}

	return true
}
