package request

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var jsonUnmarshal = json.Unmarshal // Used to enable unit testing

// BootstrapFromJSON boostraps the Metadata and the specific Request from a JSON payload.
// It returns a Metadata pointer, a Request pointer, and an error.
// The error will be non-nil if there are issues unmarshalling the JSON or if the request type is not supported.
// If the error is non-nil the Metadata and Request pointers will be nil.
func BootstrapFromJSON(data []byte) (*Metadata, interface{}, error) {
	var efu envelope

	if err := jsonUnmarshal(data, &efu); err != nil {
		return nil, nil, errors.Errorf("failed to unmarshal elements common to all request envelopes: %v", err)
	}

	switch efu.Request.Type {
	case "LaunchRequest":
		var env launchRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.Errorf("failed to unmarshal launch request envelope: %v", err)
		}
		return &efu.Metadata, &env.Request, nil
	case "IntentRequest":
		var env intentRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.Errorf("failed to unmarshal intent request envelope: %v", err)
		}
		return &efu.Metadata, &env.Request, nil
	case "SessionEndedRequest":
		var env sessionEndedRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.Errorf("failed to unmarshal session ended request envelope: %v", err)
		}
		return &efu.Metadata, &env.Request, nil
	default:
		return nil, nil, errors.Errorf("request type %s not supported", efu.Request.Type)
	}
}
