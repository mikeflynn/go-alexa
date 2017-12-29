package request

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestRequest_BootstrapFromJSON(t *testing.T) {
	var tests = []struct {
		name                string
		payload             string
		partialErrorMessage *string
		metadata            *Metadata
		request             interface{}
		jsonUnmarshal       func([]byte, interface{}) error
	}{
		{
			name: "happy-path-launch-request",
			payload: `{
				"version": "testVersion",
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			partialErrorMessage: nil,
			metadata: &Metadata{
				Version: "testVersion",
				Session: Session{}},
			request: &LaunchRequest{
				Common: Common{
					Type: "LaunchRequest",
				},
			},
			jsonUnmarshal: func(data []byte, v interface{}) error {
				return json.Unmarshal(data, v)
			},
		},
		{
			name: "happy-path-intent-request",
			payload: `{
				"version": "testVersion",
				"request": {
					"type": "IntentRequest",
					"dialogState": "testDialogState",
					"intent": {
						"name": "testIntent",
						"confirmationStatus": "testConfirmationStatus",
						"slots": {
							"testSlot": {
								"name":  "testSlotName",
								"value": "testSlotValue",
								"confirmationStatus": "testConfirmationStatus",
								"resolutions": "testResolutions"
							}
						}
					}
				}
			}`,
			partialErrorMessage: nil,
			metadata: &Metadata{
				Version: "testVersion",
				Session: Session{}},
			request: &IntentRequest{
				Common: Common{
					Type: "IntentRequest",
				},
				DialogState: "testDialogState",
				Intent: Intent{
					Name:               "testIntent",
					ConfirmationStatus: "testConfirmationStatus",
					Slots: map[string]Slot{
						"testSlot": {
							Name:               "testSlotName",
							Value:              "testSlotValue",
							ConfirmationStatus: "testConfirmationStatus",
							Resolutions:        "testResolutions",
						},
					},
				},
			},
			jsonUnmarshal: func(data []byte, v interface{}) error {
				return json.Unmarshal(data, v)
			},
		},
		{
			name: "happy-path-session-ended-request",
			payload: `{
				"version": "testVersion",
				"request": {
					"type": "SessionEndedRequest",
                    "error": {
						"type": "testType",
						"message": "testMessage"
					},
					"reason": "testReason"
				}
			}`,
			partialErrorMessage: nil,
			metadata: &Metadata{
				Version: "testVersion",
				Session: Session{}},
			request: &SessionEndedRequest{
				Common: Common{
					Type: "SessionEndedRequest",
				},
				Error: Error{
					Type:    "testType",
					Message: "testMessage",
				},
				Reason: "testReason",
			},
			jsonUnmarshal: func(data []byte, v interface{}) error {
				return json.Unmarshal(data, v)
			},
		},
		{
			name: "unsupported-type-returns-error",
			payload: `{
				"request": {
					"type": "UnknownRequestType"
				}
			}`,
			partialErrorMessage: strPointer("request type UnknownRequestType not supported"),
			metadata:            nil,
			request:             nil,
			jsonUnmarshal: func(data []byte, v interface{}) error {
				return json.Unmarshal(data, v)
			},
		},
		{
			name:                "invalid-envelope-common-returns-error",
			payload:             `{}`,
			partialErrorMessage: strPointer("failed to unmarshal elements common to all request envelopes:"),
			metadata:            nil,
			request:             nil,
			jsonUnmarshal: func(data []byte, v interface{}) error {
				return errors.New("dummy error")
			},
		},
		{
			name: "invalid-launch-request-returns-error",
			payload: `{
				"request": {
					"type": "LaunchRequest"
				}
			}`,
			partialErrorMessage: strPointer("failed to unmarshal launch request envelope:"),
			metadata:            nil,
			request:             nil,
			jsonUnmarshal: func(data []byte, v interface{}) error {
				if _, ok := v.(*launchRequestEnvelope); ok {
					return errors.New("dummy error")
				}
				return json.Unmarshal(data, v)
			},
		},
		{
			name: "invalid-intent-request-returns-error",
			payload: `{
				"request": {
					"type": "IntentRequest"
				}
			}`,
			partialErrorMessage: strPointer("failed to unmarshal intent request envelope:"),
			metadata:            nil,
			request:             nil,
			jsonUnmarshal: func(data []byte, v interface{}) error {
				if _, ok := v.(*intentRequestEnvelope); ok {
					return errors.New("dummy error")
				}
				return json.Unmarshal(data, v)
			},
		},
		{
			name: "invalid-session-ended-request-returns-error",
			payload: `{
				"request": {
					"type": "SessionEndedRequest"
				}
			}`,
			partialErrorMessage: strPointer("failed to unmarshal session ended request envelope:"),
			metadata:            nil,
			request:             nil,
			jsonUnmarshal: func(data []byte, v interface{}) error {
				if _, ok := v.(*sessionEndedRequestEnvelope); ok {
					return errors.New("dummy error")
				}
				return json.Unmarshal(data, v)
			},
		},
	}

	for _, test := range tests {
		jsonUnmarshal = test.jsonUnmarshal

		m, r, err := BootstrapFromJSON([]byte(test.payload))
		if !errorContains(err, test.partialErrorMessage) {
			t.Errorf("%s: error mismatch:\n\tgot:    %v\n\texpected: it to contain %s", test.name, err, pointerStr(test.partialErrorMessage))
			continue
		}

		if test.partialErrorMessage != nil {
			continue
		}

		if !reflect.DeepEqual(*m, *test.metadata) {
			t.Errorf("%s: metadata mismatch:\n\tgot:    %#v\n\twanted: %#v", test.name, *m, *test.metadata)
		}

		if !reflect.DeepEqual(r, test.request) {
			t.Errorf("%s: request mismatch:\n\tgot:    %#v\n\twanted: %#v", test.name, r, test.request)
		}
	}
}

/* Test helper functions */
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
