package skillserver

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
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

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// EchoApplication represents a single Alexa application server. This application type needs to include
// the application ID from the Alexa developer portal that will be making requests to the server. This AppId needs
// to be verified to ensure the requests are coming from the correct app. Handlers can also be provied for
// different types of requests sent by the Alexa Skills Kit such as OnLaunch or OnIntent.
type EchoApplication struct {
	AppID              string
	Handler            func(http.ResponseWriter, *http.Request)
	OnLaunch           func(*EchoRequest, *EchoResponse)
	OnIntent           func(*EchoRequest, *EchoResponse)
	OnSessionEnded     func(*EchoRequest, *EchoResponse)
	OnAudioPlayerState func(*EchoRequest, *EchoResponse)
}

// StdApplication is a type of application that allows the user to accept and manually process
// requests from an Alexa application on an existing HTTP server. Request validation and parsing
// will need to be done manually to ensure compliance with the requirements of the Alexa Skills Kit.
type StdApplication struct {
	Methods string
	Handler func(http.ResponseWriter, *http.Request)
}

type requestContextKey string

var (
	applications       = map[string]interface{}{}
	rootPrefix         = "/"
	echoPrefix         = "/echo/"
	insecureSkipVerify = false
)

// SetEchoPrefix provides a way to specify a single path prefix that all EchoApplications will share.SetEchoPrefix
// All incoming requests to an initialized EchoApplication will need to have a path that starts with this prefix.
func SetEchoPrefix(prefix string) {
	echoPrefix = prefix
}

// SetRootPrefix allows a single path prefix to be applied to the request path of all
// StdApplications. All requests to the StdApplications provided will need to begin with
// this prefix.
func SetRootPrefix(prefix string) {
	rootPrefix = prefix
}

// Run will initialize the apps provided and start an HTTP server listening on the specified port.
// Use an optional fourth boolean parameter to indicate whether to skip AWS cert's validation
func Run(apps map[string]interface{}, port string, insecureSkipVerify ...bool) {
	router := mux.NewRouter()
	isv := false
	if len(insecureSkipVerify) > 0 {
		isv = insecureSkipVerify[0]
	}
	initialize(apps, router, isv)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + port)
}

// RunSSL takes in a map of application, server port, certificate and key files, and
// tries to start a TLS server which alexa can directly pass commands to.
// It panics out with the error if the server couldn't be started. Or else the method blocks
// at ListenAndServeTLS line.
// If the server starts succcessfully and there are connection errors afterwards, they are
// logged to the stdout and no error is returned.
// For generating a testing cert and key, read the following:
// https://developer.amazon.com/docs/custom-skills/configure-web-service-self-signed-certificate.html
// Use an optional fourth boolean parameter to indicate whether to skip AWS cert's validation
func RunSSL(apps map[string]interface{}, port, cert, key string, insecureSkipVerify ...bool) {
	router := mux.NewRouter()

	isv := false
	if len(insecureSkipVerify) > 0 {
		isv = insecureSkipVerify[0]
	}
	initialize(apps, router, isv)

	// This is very limited TLS configuration which is required to connect alexa to our webservice.
	// The curve preferences are used by ECDSA/ECDHE algorithms for figuring out the matching algorithm
	// from alexa side starting from the strongest to the weakest.
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// If the connection throws errors related to crypt algorithm mismatch between server and client,
			// this line must be replaced by constants present in crypt/tls package for the value that works.
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS(cert, key))
}

func initialize(apps map[string]interface{}, router *mux.Router, isv bool) {
	insecureSkipVerify = isv
	if isv {
		log.Println("insecure skip verify, certs will not be checked")
	}
	applications = apps

	// /echo/* Endpoints
	echoRouter := mux.NewRouter()
	// /* Endpoints
	pageRouter := mux.NewRouter()

	hasPageRouter := false

	for uri, meta := range applications {
		switch app := meta.(type) {
		case EchoApplication:
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				echoReq := GetEchoRequest(r)
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
				} else if strings.HasPrefix(echoReq.GetRequestType(), "AudioPlayer.") {
					if app.OnAudioPlayerState != nil {
						app.OnAudioPlayerState(echoReq, echoResp)
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
			hasPageRouter = true
			pageRouter.HandleFunc(uri, app.Handler).Methods(app.Methods)
		}
	}

	router.PathPrefix(echoPrefix).Handler(negroni.New(
		negroni.HandlerFunc(validateRequest),
		negroni.HandlerFunc(verifyJSON),
		negroni.Wrap(echoRouter),
	))

	if hasPageRouter {
		router.PathPrefix(rootPrefix).Handler(negroni.New(
			negroni.Wrap(pageRouter),
		))
	}
}

// GetEchoRequest is a convenience method for retrieving and casting an `EchoRequest` out of a
// standard `http.Request`.
func GetEchoRequest(r *http.Request) *EchoRequest {
	return r.Context().Value(requestContextKey("echoRequest")).(*EchoRequest)
}

// HTTPError is a convenience method for logging a message and writing the provided error message
// and error code to the HTTP response.
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
	if !echoReq.VerifyAppID(applications[r.URL.Path].(EchoApplication).AppID) {
		HTTPError(w, "Echo AppID mismatch!", "Bad Request", 400)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), requestContextKey("echoRequest"), echoReq))

	next(w, r)
}

// Run all mandatory Amazon security checks on the request.
func validateRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	devFlag := r.URL.Query().Get("_dev")
	isDev := devFlag != ""
	if !isDev && !IsValidAlexaRequest(w, r) {
		log.Println("Request invalid")
		return
	}
	next(w, r)
}

// IsValidAlexaRequest handles all the necessary steps to validate that an incoming http.Request has actually come from
// the Alexa service. If an error occurs during the validation process, an http.Error will be written to the provided http.ResponseWriter.
// The required steps for request validation can be found on this page:
// --insecure-skip-verify flag will disable all validations
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#hosting-a-custom-skill-as-a-web-service
func IsValidAlexaRequest(w http.ResponseWriter, r *http.Request) bool {
	if insecureSkipVerify {
		return true
	}
	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", 401)
		return false
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
		return false
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		HTTPError(w, "Amazon certificate invalid.", "Not Authorized", 401)
		return false
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
		return false
	}
	//log.Println(bodyBuf.String())
	r.Body = ioutil.NopCloser(&bodyBuf)

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil {
		HTTPError(w, "Signature match failed.", "Not Authorized", 401)
		return false
	}

	return true
}

func readCert(certURL string) ([]byte, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil || certPool == nil {
		log.Println("Can't open system cert pools")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: certPool, InsecureSkipVerify: insecureSkipVerify},
	}
	hc := &http.Client{Timeout: 2 * time.Second, Transport: tr}

	cert, err := hc.Get(certURL)
	if err != nil {
		return nil, errors.New("could not download Amazon cert file: " + err.Error())
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("could not read Amazon cert file: " + err.Error())
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
