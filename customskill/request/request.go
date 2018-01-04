package request

import (
	"encoding/json"
	"errors"
)

var jsonUnmarshal = json.Unmarshal // Used to enable unit testing

// BootstrapFromJSON bootstraps the Metadata and the specific Request from a JSON payload.
// It returns a Metadata pointer, a Request pointer, and an error.
// The error will be non-nil if there are issues unmarshaling the JSON or if the request type is not supported.
// If the error is non-nil the Metadata and Request pointers will be nil.
func BootstrapFromJSON(data []byte) (*Metadata, interface{}, error) {
	var efu envelope

	if err := jsonUnmarshal(data, &efu); err != nil {
		return nil, nil, errors.New("failed to unmarshal elements common to all request envelopes: " + err.Error())
	}

	switch efu.Request.Type {
	case "LaunchRequest":
		var env launchRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.New("failed to unmarshal launch request envelope: " + err.Error())
		}
		return &efu.Metadata, &env.Request, nil
	case "IntentRequest":
		var env intentRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.New("failed to unmarshal intent request envelope: " + err.Error())
		}
		return &efu.Metadata, &env.Request, nil
	case "SessionEndedRequest":
		var env sessionEndedRequestEnvelope
		if err := jsonUnmarshal(data, &env); err != nil {
			return nil, nil, errors.New("failed to unmarshal session ended request envelope: " + err.Error())
		}
		return &efu.Metadata, &env.Request, nil
	default:
		return nil, nil, errors.New("request type " + efu.Request.Type + " not supported")
	}
}
