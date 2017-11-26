package customskill

import (
	"testing"

	"io"

	"strings"

	"encoding/json"

	"bytes"

	"io/ioutil"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
	"github.com/pkg/errors"
)

func TestSkill_Handle(t *testing.T) {

	var tests = []struct {
		name                string
		skill               *Skill
		w                   io.ReadWriter
		b                   string
		jsonMarshal         func(v interface{}) ([]byte, error)
		partialErrorMessage *string
		written             string
	}{
		{
			name: "happy-path-launch-request",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnLaunch: func(request *request.LaunchRequest) (*response.Envelope, error) {
					return response.New(), nil
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
			written: `{"version":"1.0","response":{"shouldEndSession":true}}`,
		},
		{
			name: "happy-path-intent-request",
			skill: &Skill{
				ValidApplicationIDs: []string{"testApplicationId"},
				OnIntent: func(intentRequest *request.IntentRequest, session *request.Session) (*response.Envelope, error) {
					return response.New(), nil
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
			written: `{"version":"1.0","response":{"shouldEndSession":true}}`,
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
				OnLaunch: func(request *request.LaunchRequest) (*response.Envelope, error) {
					return nil, errors.New("dummy error")
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
				OnIntent: func(intentRequest *request.IntentRequest, session *request.Session) (*response.Envelope, error) {
					return nil, errors.New("dummy error")
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
				OnLaunch: func(request *request.LaunchRequest) (*response.Envelope, error) {
					return nil, nil
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
				OnLaunch: func(request *request.LaunchRequest) (*response.Envelope, error) {
					return nil, nil
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
				OnLaunch: func(request *request.LaunchRequest) (*response.Envelope, error) {
					return nil, nil
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
			partialErrorMessage: strPointer("failed to completely write response: 0 of 4 bytes written"),
		},
	}

	for _, test := range tests {
		jsonMarshal = json.Marshal
		if test.jsonMarshal != nil {
			jsonMarshal = test.jsonMarshal
		}
		err := test.skill.Handle(test.w, []byte(test.b))
		if !errorContains(err, test.partialErrorMessage) {
			t.Errorf("%s: error mismatch:\n\tgot:    %v\n\texpected: it to contain %s", test.name, err, pointerStr(test.partialErrorMessage))
			continue
		}

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
