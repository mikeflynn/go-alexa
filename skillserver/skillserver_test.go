package skillserver

import (
	"github.com/stretchr/testify/assert"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var validJson = `{
  "session": {
    "sessionId": "SessionId.cabcd",
    "application": {
      "applicationId": "amzn1.ask.skill.XXXXXXX"
    },
    "attributes": {},
    "user": {
      "userId": "amzn1.ask.account.XXXXXXXXXXXXXXXXXX"
    },
    "new": true
  },
  "request": {
    "type": "IntentRequest",
    "requestId": "EdwRequestId.c8d5a2fe-408f-4caa-b237-35dc3fe6378a",
    "locale": "de-DE",
    "timestamp": "2016-12-15T14:20:02Z",
    "intent": {
      "name": "TestIntent",
      "slots": {}
    }
  },
  "version": "1.0"
}
`

func TestVerifyJsonInDevMode(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false

	verifyJsonHandler := VerifyJSON("amzn1.ask.skill.XXXXXXX", func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})
	testUrl, _ := url.Parse("/echo?_dev=true")
	testRequest := &http.Request{}
	bodyBuffer := bytes.NewBufferString(validJson)
	body := ioutil.NopCloser(bodyBuffer)
	testRequest.Body = body
	testRequest.URL = testUrl

	responseRecorder := httptest.NewRecorder()
	verifyJsonHandler(responseRecorder, testRequest)

	assert.True(handlerCalled)
	assert.NotEqual(responseRecorder.Code, 400)
}

func TestInvalidAppId(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false

	verifyJsonHandler := VerifyJSON("amzn1.ask.skill.YYYYYYYYYYYYY", func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})
	testUrl, _ := url.Parse("/echo?_dev=true")
	testRequest := &http.Request{}
	bodyBuffer := bytes.NewBufferString(validJson)
	body := ioutil.NopCloser(bodyBuffer)
	testRequest.Body = body
	testRequest.URL = testUrl

	responseRecorder := httptest.NewRecorder()
	verifyJsonHandler(responseRecorder, testRequest)

	assert.False(handlerCalled)
	assert.Equal(responseRecorder.Code, 400)
}
