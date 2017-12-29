package customskill

import (
	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
)

// A Skill represents an Alexa custom skill.
type Skill struct {
	ValidApplicationIDs []string
	OnLaunch            func(*request.LaunchRequest) (*response.Envelope, error)
	OnIntent            func(*request.IntentRequest, *request.Session) (*response.Envelope, error)
	OnSessionEnded      func(endedRequest *request.SessionEndedRequest) error
}
