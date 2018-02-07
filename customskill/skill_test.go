package customskill

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
	"github.com/pkg/errors"
)

func TestSkill_Handle(t *testing.T) {
	handleTests(t)
}

func BenchmarkSkill_Handle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handleTests(b)
	}
}

func handleTests(t testingiface) {
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
				OnLaunch: func(request *request.LaunchRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnIntent: func(intentRequest *request.IntentRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnIntent: func(intentRequest *request.IntentRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
					metadata.Session.Attributes["name"] = "happy-path-intent-request-with-session-attributes"
					return response.New(), metadata.Session.Attributes, nil
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
				OnSessionEnded: func(endedRequest *request.SessionEndedRequest, metadata *request.Metadata) error {
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
			name: "nil-request-type-returns-error",
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
		}, {
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
				}, 0, nil
			},
			partialErrorMessage: strPointer("unsupported request type: int"),
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
				OnLaunch: func(request *request.LaunchRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnIntent: func(intentRequest *request.IntentRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnSessionEnded: func(endedRequest *request.SessionEndedRequest, metadata *request.Metadata) error {
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
				OnLaunch: func(request *request.LaunchRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnLaunch: func(request *request.LaunchRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
				OnLaunch: func(request *request.LaunchRequest, metadata *request.Metadata) (*response.Response, map[string]interface{}, error) {
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
		t.Logf("Testing: %s", test.name)
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
}

/* Test helper functions */

type testingiface interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

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
