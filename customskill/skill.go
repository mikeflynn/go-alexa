package customskill

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
	"github.com/pkg/errors"
)

var (
	requestBootstrapFromJSON = request.BootstrapFromJSON
	jsonMarshal              = json.Marshal
	ioCopy                   = io.Copy
	httpError                = http.Error
	sHandle                  = func(s *Skill, w io.Writer, b []byte) error { return s.Handle(w, b) }
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
		return errors.Errorf("failed to bootstrap request from JSON payload: %v", err)
	}

	if !s.applicationIDIsValid(m.Session.Application.ApplicationID) {
		return errors.Errorf("application ID %s is not in the list of valid application IDs: %v", m.Session.Application.ApplicationID, s.ValidApplicationIDs)
	}

	switch e.(type) {
	case *request.LaunchRequest:
		if s.OnLaunch == nil {
			return errors.Errorf("no OnLaunch handler defined")
		}
		lr := e.(*request.LaunchRequest)
		resp, sess, err = s.OnLaunch(lr)
		if err != nil {
			return errors.Errorf("OnLaunch handler failed: %v", err)
		}
	case *request.IntentRequest:
		if s.OnIntent == nil {
			return errors.Errorf("no OnIntent handler defined")
		}
		ir := e.(*request.IntentRequest)
		resp, sess, err = s.OnIntent(ir, &m.Session)
		if err != nil {
			return errors.Errorf("OnIntent handler failed: %v", err)
		}
	case *request.SessionEndedRequest:
		if s.OnSessionEnded == nil {
			return errors.Errorf("no OnSessionEnded handler defined")
		}
		ser := e.(*request.SessionEndedRequest)
		if err = s.OnSessionEnded(ser); err != nil {
			return errors.Errorf("OnSessionEnded handler failed: %v", err)
		}
		// A skill cannot return a response to SessionEndedRequest.
		return nil
	default:
		return errors.Errorf("unsupported request type: %T", e)
	}

	jsonB, err := jsonMarshal(response.NewEnvelope(resp, sess))
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

func (s *Skill) DefaultHTTPHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := ioCopy(buf, r.Body); err != nil {
		httpError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := sHandle(s, w, buf.Bytes()); err != nil {
		httpError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Skill) applicationIDIsValid(appID string) bool {
	for _, str := range s.ValidApplicationIDs {
		if str == appID {
			return true
		}
	}
	return false
}
