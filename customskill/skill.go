package customskill

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
	"github.com/pkg/errors"
)

var jsonMarshal = json.Marshal

// Handle parses a JSON payload, calls the appropriate request handler, serializes the response, and writes it to the provided writer.
func (s *Skill) Handle(w io.Writer, b []byte) error {
	var (
		resp *response.Envelope
		err  error
	)

	m, e, err := request.BootstrapFromJSON(b)
	if err != nil {
		return errors.Errorf("failed to bootstrap request from JSON payload: %v", err)
	}

	if !s.applicationIDIsValid(m.Session.Application.ApplicationID) {
		return errors.Errorf("application ID %s is not in the list of valid application IDs: %v", m.Session.Application.ApplicationID, s.ValidApplicationIDs)
	}

	switch reflect.TypeOf(e) {
	case reflect.TypeOf(&request.LaunchRequest{}):
		if s.OnLaunch == nil {
			return errors.Errorf("no OnLaunch handler defined")
		}
		lr := e.(*request.LaunchRequest)
		resp, err = s.OnLaunch(lr)
		if err != nil {
			return errors.Errorf("OnLaunch handler failed: %v", err)
		}
	case reflect.TypeOf(&request.IntentRequest{}):
		if s.OnIntent == nil {
			return errors.Errorf("no OnIntent handler defined")
		}
		ir := e.(*request.IntentRequest)
		resp, err = s.OnIntent(ir, &m.Session)
		if err != nil {
			return errors.Errorf("OnIntent handler failed: %v", err)
		}
	case reflect.TypeOf(&request.SessionEndedRequest{}):
		if s.OnSessionEnded == nil {
			return errors.Errorf("no OnSessionEnded handler defined")
		}
		ser := e.(*request.SessionEndedRequest)
		if err = s.OnSessionEnded(ser); err != nil {
			return errors.Errorf("OnSessionEnded handler failed: %v", err)
		}
		// A skill cannot return a response to SessionEndedRequest.
		return nil
	}

	jsonB, err := jsonMarshal(resp)
	if err != nil {
		return errors.Errorf("failed to marshal response: %v", err)
	}
	n, err := w.Write(jsonB)
	if err != nil {
		return errors.Errorf("failed to write response: %v", err)
	}
	if n != len(jsonB) {
		return errors.Errorf("failed to completely write response: %d of %d bytes written", n, len(jsonB))
	}
	return nil
}

func (s *Skill) applicationIDIsValid(appID string) bool {
	for _, str := range s.ValidApplicationIDs {
		if str == appID {
			return true
		}
	}
	return false
}
