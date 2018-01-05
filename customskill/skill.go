package customskill

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
)

var (
	requestBootstrapFromJSON = request.BootstrapFromJSON
	jsonMarshal              = json.Marshal
)

// Handle parses a JSON payload, calls the appropriate request handler, serializes the response, and writes it to the provided writer.
func (s *Skill) Handle(w io.Writer, b []byte) error {
	var (
		resp *response.Response
		sess map[string]interface{}
		err  error
	)

	m, e, err := requestBootstrapFromJSON(b)
	if err != nil {
		return errors.New("failed to bootstrap request from JSON payload: " + err.Error())
	}

	if !s.applicationIDIsValid(m.Session.Application.ApplicationID) {
		return errors.New("application ID " + m.Session.Application.ApplicationID + " is not in the list of valid application IDs: " + strings.Join(s.ValidApplicationIDs, ","))
	}

	switch e.(type) {
	case *request.LaunchRequest:
		if s.OnLaunch == nil {
			return errors.New("no OnLaunch handler defined")
		}
		lr := e.(*request.LaunchRequest)
		resp, sess, err = s.OnLaunch(lr, m)
		if err != nil {
			return errors.New("OnLaunch handler failed: " + err.Error())
		}
	case *request.IntentRequest:
		if s.OnIntent == nil {
			return errors.New("no OnIntent handler defined")
		}
		ir := e.(*request.IntentRequest)
		resp, sess, err = s.OnIntent(ir, m)
		if err != nil {
			return errors.New("OnIntent handler failed: " + err.Error())
		}
	case *request.SessionEndedRequest:
		if s.OnSessionEnded == nil {
			return errors.New("no OnSessionEnded handler defined")
		}
		ser := e.(*request.SessionEndedRequest)
		if err = s.OnSessionEnded(ser, m); err != nil {
			return errors.New("OnSessionEnded handler failed:" + err.Error())
		}
		// A skill cannot return a response to SessionEndedRequest.
		return nil
	default:
		return errors.New("unsupported request type: " + getType(e))
	}

	jsonB, err := jsonMarshal(response.NewEnvelope(resp, sess))
	if err != nil {
		return errors.New("failed to marshal response: " + err.Error())
	}
	n, err := w.Write(jsonB)
	if err != nil {
		return errors.New("failed to write response: " + err.Error())
	}
	if n != len(jsonB) {
		return errors.New("failed to completely write response: " + strconv.Itoa(n) + " of " + strconv.Itoa(len(jsonB)) + " bytes written")
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

func getType(i interface{}) string {
	t := reflect.TypeOf(i)
	if t == nil {
		return "<nil>"
	}
	return t.String()
}
