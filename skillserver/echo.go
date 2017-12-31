package skillserver

import (
	"encoding/json"
	"errors"
	"time"
)

// Request Functions
// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", this.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) VerifyAppID(myAppID string) bool {
	if this.Session.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) GetSessionID() string {
	return this.Session.SessionID
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) GetUserID() string {
	return this.Session.User.UserID
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) GetRequestType() string {
	return this.Request.Type
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) GetIntentName() string {
	if this.GetRequestType() == "IntentRequest" {
		return this.Request.Intent.Name
	}

	return this.GetRequestType()
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) GetSlotValue(slotName string) (string, error) {
	if _, ok := this.Request.Intent.Slots[slotName]; ok {
		return this.Request.Intent.Slots[slotName].Value, nil
	}

	return "", errors.New("Slot name not found.")
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoRequest) AllSlots() map[string]EchoSlot {
	return this.Request.Intent.Slots
}

// Response Functions
// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func NewEchoResponse() *EchoResponse {
	er := &EchoResponse{
		Version: "1.0",
		Response: EchoRespBody{
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	return er
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) OutputSpeech(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "PlainText",
		Text: text,
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) Card(title string, content string) *EchoResponse {
	return this.SimpleCard(title, content)
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) OutputSpeechSSML(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "SSML",
		SSML: text,
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) StandardCard(title string, content string, smallImg string, largeImg string) *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		this.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		this.Response.Card.Image.LargeImageURL = largeImg
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) LinkAccountCard() *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type: "LinkAccount",
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) Reprompt(text string) *EchoResponse {
	this.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) RepromptSSML(text string) *EchoResponse {
	this.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "SSML",
			Text: text,
		},
	}

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) EndSession(flag bool) *EchoResponse {
	this.Response.ShouldEndSession = flag

	return this
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
func (this *EchoResponse) String() ([]byte, error) {
	jsonStr, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// Request Types

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoReqBody struct {
	Type      string     `json:"type"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Intent    EchoIntent `json:"intent,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoIntent struct {
	Name  string              `json:"name"`
	Slots map[string]EchoSlot `json:"slots"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response Types

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload `json:"card,omitempty"`
	Reprompt         *EchoReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool             `json:"shouldEndSession"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoReprompt struct {
	OutputSpeech EchoRespPayload `json:"outputSpeech,omitempty"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Deprecated: Please use the github.com/mikeflynn/go-alexa/customskill package
type EchoRespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}
