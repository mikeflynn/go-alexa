package customskill

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
)

func TestSkill_Handle(t *testing.T) {

	var tests = []struct {
		name                     string
		skill                    *Skill
		w                        io.ReadWriter
		b                        string
		jsonMarshal              func(v interface{}) ([]byte, error)
		requestBootstrapFromJSON func(data []byte) (*request.Metadata, interface{}, error)
		partialErrorMessage      *string
		written                  string
	}{
		{
			name: "happy-path-launch-request",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
					sessAttrs := make(map[string]interface{})
					sessAttrs["name"] = "happy-path-launch-request"
					return response.New(), sessAttrs, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			w:       bytes.NewBuffer(nil),
			written: `{"version":"1.0","sessionAttributes":{"name":"happy-path-launch-request"},"response":{"shouldEndSession":true}}`,
		},
		{
			name: "happy-path-intent-request",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnIntent: func(intentRequest *request.IntentRequest, session *request.Session) (*response.Response, map[string]interface{}, error) {
					sessAttrs := make(map[string]interface{})
					sessAttrs["name"] = "happy-path-intent-request"
					return response.New(), sessAttrs, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "IntentRequest"
				}
			}`,
			w:       bytes.NewBuffer(nil),
			written: `{"version":"1.0","sessionAttributes":{"name":"happy-path-intent-request"},"response":{"shouldEndSession":true}}`,
		},
		{
			name: "happy-path-intent-request-with-session-attributes",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnIntent: func(intentRequest *request.IntentRequest, session *request.Session) (*response.Response, map[string]interface{}, error) {
					session.Attributes["name"] = "happy-path-intent-request-with-session-attributes"
					return response.New(), session.Attributes, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					},
					"attributes": {}
				},
				"request": {
					"type": "IntentRequest"
				}
			}`,
			w:       bytes.NewBuffer(nil),
			written: `{"version":"1.0","sessionAttributes":{"name":"happy-path-intent-request-with-session-attributes"},"response":{"shouldEndSession":true}}`,
		},
		{
			name: "happy-path-session-ended-request",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnSessionEnded: func(endedRequest *request.SessionEndedRequest) error {
					return nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "SessionEndedRequest"
				}
			}`,
			w:       bytes.NewBuffer(nil),
			written: ``,
		},
		{
			name:                "invalid-request-returns-error",
			b:                   "",
			partialErrorMessage: strPointer("failed to bootstrap request from JSON payload"),
		},
		{
			name: "unsupported-request-type-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
			},
			requestBootstrapFromJSON: func(data []byte) (*request.Metadata, interface{}, error) {
				return &request.Metadata{
					Session: request.Session{
						Application: request.Application{
							ApplicationID: "testApplicationId",
						},
					},
				}, nil, nil
			},
			partialErrorMessage: strPointer("unsupported request type: <nil>"),
		},
		{
			name:  "invalid-application-id-returns-error",
			skill: &Skill{},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			partialErrorMessage: strPointer("application ID testApplicationId is not in the list of valid application IDs"),
		},
		{
			name: "undefined-launch-request-handler-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			partialErrorMessage: strPointer("no OnLaunch handler defined"),
		},
		{
			name: "on-launch-request-handler-error-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
					return nil, nil, errors.New("dummy error")
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			partialErrorMessage: strPointer("OnLaunch handler failed"),
		},
		{
			name: "undefined-intent-request-handler-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "IntentRequest"
				}
			}`,
			partialErrorMessage: strPointer("no OnIntent handler defined"),
		},
		{
			name: "on-intent-request-handler-error-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnIntent: func(intentRequest *request.IntentRequest, session *request.Session) (*response.Response, map[string]interface{}, error) {
					return nil, nil, errors.New("dummy error")
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "IntentRequest"
				}
			}`,
			partialErrorMessage: strPointer("OnIntent handler failed"),
		},
		{
			name: "undefined-session-ended-request-handler-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "SessionEndedRequest"
				}
			}`,
			partialErrorMessage: strPointer("no OnSessionEnded handler defined"),
		},
		{
			name: "on-session-ended-request-handler-error-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnSessionEnded: func(endedRequest *request.SessionEndedRequest) error {
					return errors.New("dummy error")
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "SessionEndedRequest"
				}
			}`,
			partialErrorMessage: strPointer("OnSessionEnded handler failed"),
		},
		{
			name: "responses-which-cannot-be-marshalled-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
					return nil, nil, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			jsonMarshal: func(v interface{}) ([]byte, error) {
				return nil, errors.New("dummy error")
			},
			partialErrorMessage: strPointer("failed to marshal response"),
		},
		{
			name: "writer-which-cannot-be-written-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
					return nil, nil, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			w: &badReadWriter{
				err: errors.New("dummy error"),
			},
			partialErrorMessage: strPointer("failed to write response"),
		},
		{
			name: "writer-which-partially-writes-returns-error",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
					return nil, nil, nil
				},
			},
			b: `
			{
				"session": {
					"application": {
						"applicationId": "testApplicationId"
					}
				},
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			w: &badReadWriter{
				n: 0,
			},
			partialErrorMessage: strPointer("failed to completely write response: 0 of 33 bytes written"),
		},
	}

	for _, test := range tests {
		// Override mocked functions
		if test.requestBootstrapFromJSON != nil {
			requestBootstrapFromJSON = test.requestBootstrapFromJSON
		}
		if test.jsonMarshal != nil {
			jsonMarshal = test.jsonMarshal
		}
		err := test.skill.Handle(test.w, []byte(test.b))
		if !errorContains(err, test.partialErrorMessage) {
			t.Errorf("%s: error mismatch:\n\tgot:    %v\n\twanted: it to contain '%s'", test.name, err, pointerStr(test.partialErrorMessage))
			continue
		}

		// Restore mocked functions
		requestBootstrapFromJSON = request.BootstrapFromJSON
		jsonMarshal = json.Marshal

		if test.partialErrorMessage != nil {
			continue
		}

		b, err := ioutil.ReadAll(test.w)
		if err != nil {
			t.Errorf("%s: failed to read test writer: %v", test.name, err)
		}

		if string(b) != test.written {
			t.Errorf("%s: write mismatch:\n\tgot:    %v\n\texpected:%s", test.name, string(b), test.written)

		}

	}
	/*
	   	s := Skill{
	   		ValidApplicationIDs: []string{"Test"},
	   		OnIntent: func(request *request.IntentRequest, session *request.Session) (*response.Envelope, error) {
	   			fmt.Printf("Request: %#v\n", *request)
	   			fmt.Printf("Session: %#v\n", *session)

	   			resp := response.NewResponse()
	   			resp.OutputSpeech("Hello World")
	   			resp.EndSession(true)
	   			resp.SessionAttributes = session.Attributes
	   			resp.SimpleCard("title", "content")

	   			return resp, nil
	   		},
	   	}

	   	buf := bytes.NewBuffer(nil)
	   	if err := s.Handle(buf, []byte(`{
	   		"version": "testVersion",
	           "session": {
	               "application": {
	                   "applicationId": "Test"
	               }
	           },
	   		"request": {
	   	"type": "IntentRequest",
	   	"dialogState": "testDialogState"
	   	}
	   	}`)); err != nil {
	   		fmt.Printf("Got an error! %v", err)
	   	}

	   	fmt.Printf("Written response: %s", buf.String())*/
}

func TestSkill_DefaultHTTPHandler(t *testing.T) {
	checkCounter := 0
	var tests = []struct {
		name      string
		checkNum  int
		w         http.ResponseWriter
		r         *http.Request
		ioCopy    func(dst io.Writer, src io.Reader) (written int64, err error)
		httpError func(w http.ResponseWriter, error string, code int)
		sHandle   func(s *Skill, w io.Writer, b []byte) error
	}{
		{
			name:     "happy-path",
			checkNum: 1,
			sHandle: func(s *Skill, w io.Writer, b []byte) error {
				checkCounter++
				got := string(b)
				want := "happy-path"
				if got != want {
					t.Errorf("[happy-path] bytes written mismatch: got: %s, want: %s", got, want)
				}
				return nil
			},
			r: httptest.NewRequest("", "http://domain.tld", bytes.NewReader([]byte("happy-path"))),
			httpError: func(w http.ResponseWriter, error string, code int) {
				t.Error("[happy-path] unexpected http.Error call")
			},
		},
		{
			name:     "io-copy-error-sends-http-error",
			checkNum: 1,
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, errors.New("dummy error")
			},
			r: httptest.NewRequest("", "http://domain.tld", bytes.NewReader(nil)),
			httpError: func(w http.ResponseWriter, error string, code int) {
				checkCounter++
				if w != nil {
					t.Errorf("[io-copy-error-sends-http-error] http.Error w argument mismatch: got: %v, want: <nil>", w)
				}
				if error != http.StatusText(http.StatusInternalServerError) {
					t.Errorf("[io-copy-error-sends-http-error] http.Error error argument mismatch: got: %s, want: %s", error, http.StatusText(http.StatusInternalServerError))
				}
				if code != http.StatusInternalServerError {
					t.Errorf("[io-copy-error-sends-http-error] http.Error code argument mismatch: got: %d, want: %d", code, http.StatusInternalServerError)
				}
			},
		},
		{
			name:     "s-Handle-error-sends-http-error",
			checkNum: 1,
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
			sHandle: func(s *Skill, w io.Writer, b []byte) error {
				return errors.New("dummy error")
			},
			r: httptest.NewRequest("", "http://domain.tld", bytes.NewReader(nil)),
			httpError: func(w http.ResponseWriter, error string, code int) {
				checkCounter++
				if w != nil {
					t.Errorf("[s-Handle-error-sends-http-error] http.Error w argument mismatch: got: %v, want: <nil>", w)
				}
				if error != http.StatusText(http.StatusInternalServerError) {
					t.Errorf("[s-Handle-error-sends-http-error] http.Error error argument mismatch: got: %s, want: %s", error, http.StatusText(http.StatusInternalServerError))
				}
				if code != http.StatusInternalServerError {
					t.Errorf("[s-Handle-error-sends-http-error] http.Error code argument mismatch: got: %d, want: %d", code, http.StatusInternalServerError)
				}
			},
		},
	}

	for _, test := range tests {
		// Override the helper functions for each test
		if test.ioCopy != nil {
			ioCopy = test.ioCopy
		}
		if test.httpError != nil {
			httpError = test.httpError
		}
		if test.sHandle != nil {
			sHandle = test.sHandle
		}

		// Call the tested function
		s := Skill{}
		s.DefaultHTTPHandler(test.w, test.r)

		// Check to ensure the desired number of checks were performed
		if checkCounter != test.checkNum {
			t.Errorf("[%s] check counter mismatch: got: %d, want: %d", test.name, test.checkNum, checkCounter)
		}

		// Reset the check counter
		checkCounter = 0

		// Restore the helper functions
		ioCopy = io.Copy
		httpError = http.Error
		sHandle = func(s *Skill, w io.Writer, b []byte) error { return s.Handle(w, b) }
	}
}

/* Test helper functions */

type badReadWriter struct {
	n   int
	err error
}

func (rw *badReadWriter) Write(p []byte) (n int, err error) {
	return rw.n, rw.err
}

func (rw *badReadWriter) Read(p []byte) (n int, err error) {
	return 0, nil
}

func errorContains(err error, message *string) bool {
	if err == nil {
		return message == nil
	}
	if message != nil {
		return strings.Contains(err.Error(), *message)
	}
	return false
}

func strPointer(s string) *string {
	return &s
}

func pointerStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
