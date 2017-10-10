package skillserver

import (
	"encoding/json"
	"errors"
	"time"
)

type DialogType string

const (
	DialogDelegate      DialogType = "Dialog.Delegate"
	DialogElicitSlot    DialogType = "Dialog.ElicitSlot"
	DialogConfirmSlot   DialogType = "Dialog.ConfirmSlot"
	DialogConfirmIntent DialogType = "Dialog.ConfirmIntent"
)

type DialogState string

const (
	DialogStarted    DialogState = "STARTED"
	DialogInProgress DialogState = "IN_PROGRESS"
	DialogCompleted  DialogState = "COMPLETED"
)

type ConfirmationStatus string

const (
	ConfirmationConfirmed ConfirmationStatus = "CONFIRMED"
	ConfirmationDenied    ConfirmationStatus = "DENIED"
)

// Request Functions
func (this *EchoRequest) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", this.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

func (this *EchoRequest) VerifyAppID(myAppID string) bool {
	if this.Session.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

func (this *EchoRequest) GetSessionID() string {
	return this.Session.SessionID
}

func (this *EchoRequest) GetUserID() string {
	return this.Session.User.UserID
}

func (this *EchoRequest) GetRequestType() string {
	return this.Request.Type
}

func (this *EchoRequest) GetIntentName() string {
	if this.GetRequestType() == "IntentRequest" {
		return this.Request.Intent.Name
	}

	return this.GetRequestType()
}

func (this *EchoRequest) GetSlotValue(slotName string) (string, error) {
	if _, ok := this.Request.Intent.Slots[slotName]; ok {
		return this.Request.Intent.Slots[slotName].Value, nil
	}

	return "", errors.New("Slot name not found.")
}

func (this *EchoRequest) AllSlots() map[string]EchoSlot {
	return this.Request.Intent.Slots
}

func (this *EchoRequest) GetDialogState() DialogState {
	return this.Request.DialogState
}

func (this *EchoRequest) GetIntentConfirmationStatus() ConfirmationStatus {
	return this.Request.Intent.ConfirmationStatus
}

// Response Functions
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

func (this *EchoResponse) OutputSpeech(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "PlainText",
		Text: text,
	}

	return this
}

func (this *EchoResponse) Card(title string, content string) *EchoResponse {
	return this.SimpleCard(title, content)
}

func (this *EchoResponse) OutputSpeechSSML(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "SSML",
		SSML: text,
	}

	return this
}

func (this *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return this
}

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

func (this *EchoResponse) LinkAccountCard() *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type: "LinkAccount",
	}

	return this
}

func (this *EchoResponse) Reprompt(text string) *EchoResponse {
	this.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return this
}

func (this *EchoResponse) RepromptSSML(text string) *EchoResponse {
	this.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "SSML",
			Text: text,
		},
	}

	return this
}

func (this *EchoResponse) EndSession(flag bool) *EchoResponse {
	this.Response.ShouldEndSession = flag

	return this
}

func (this *EchoResponse) DialogDelegate(updatedIntent *EchoIntent) *EchoResponse {
	this.appendDirective(&DialogDirective{
		Type:          DialogDelegate,
		UpdatedIntent: updatedIntent,
	})
	this.EndSession(false)

	return this
}

func (this *EchoResponse) ElicitSlot(slotToElicit string, updatedIntent *EchoIntent) *EchoResponse {
	this.appendDirective(&DialogDirective{
		Type:          DialogElicitSlot,
		SlotToElicit:  slotToElicit,
		UpdatedIntent: updatedIntent,
	})
	this.EndSession(false)

	return this
}

func (this *EchoResponse) ConfirmSlot(slotToConfirm string, updatedIntent *EchoIntent) *EchoResponse {
	this.appendDirective(&DialogDirective{
		Type:          DialogConfirmSlot,
		SlotToConfirm: slotToConfirm,
		UpdatedIntent: updatedIntent,
	})
	this.EndSession(false)

	return this
}

func (this *EchoResponse) ConfirmIntent(intentToConfirm string, updatedIntent *EchoIntent) *EchoResponse {
	this.appendDirective(&DialogDirective{
		Type:            DialogConfirmIntent,
		IntentToConfirm: intentToConfirm,
		UpdatedIntent:   updatedIntent,
	})
	this.EndSession(false)

	return this
}

func (this *EchoResponse) appendDirective(directive *DialogDirective) {
	if this.Response.Directives == nil {
		this.Response.Directives = make([]*DialogDirective, 0, 5)
	}

	this.Response.Directives = append(this.Response.Directives, directive)
}

func (this *EchoResponse) String() ([]byte, error) {
	jsonStr, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// Request Types

type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
}

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

type EchoReqBody struct {
	Type        string      `json:"type"`
	RequestID   string      `json:"requestId"`
	Timestamp   string      `json:"timestamp"`
	Intent      EchoIntent  `json:"intent,omitempty"`
	Reason      string      `json:"reason,omitempty"`
	DialogState DialogState `json:"dialogState,omitempty"`
}

type EchoIntent struct {
	Name               string              `json:"name"`
	ConfirmationStatus ConfirmationStatus  `json:"confirmationStatus,omitempty"`
	Slots              map[string]EchoSlot `json:"slots"`
}

type EchoSlot struct {
	Name               string             `json:"name"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus,omitempty"`
	Value              string             `json:"value"`
}

// Response Types

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

type DialogDirective struct {
	Type            DialogType  `json:"type,omitempty"`
	SlotToConfirm   string      `json:"slotToConfirm,omitempty"`
	SlotToElicit    string      `json:"slotToElicit,omitempty"`
	IntentToConfirm string      `json:"intentToConfirm,omitempty"`
	UpdatedIntent   *EchoIntent `json:"updatedIntent,omitempty"`
}

type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload   `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload   `json:"card,omitempty"`
	Reprompt         *EchoReprompt      `json:"reprompt,omitempty"` // Pointer so it's dropped ifempty in JSON response.
	Directives       []*DialogDirective `json:"directives,omitempty"`
	ShouldEndSession bool               `json:"shouldEndSession"`
}

type EchoReprompt struct {
	OutputSpeech EchoRespPayload `json:"outputSpeech,omitempty"`
}

type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type EchoRespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}
