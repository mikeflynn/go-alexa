package skillserver

import (
	"encoding/json"
	"errors"
	"github.com/mikeflynn/go-alexa/skillserver/dialog"
	"time"
)

type ConfirmationStatus string

const (
	CONF_CONFIRMED ConfirmationStatus = "CONFIRMED"
	CONF_DENIED    ConfirmationStatus = "DENIED"
	CONF_NONE      ConfirmationStatus = "NONE"
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
	if this.Session.Application.ApplicationID == myAppID ||
		this.Context.System.Application.ApplicationID == myAppID {
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
	slot, err := this.GetSlot(slotName)

	if err != nil {
		return "", err
	}

	return slot.Value, nil
}

func (this *EchoRequest) GetSlot(slotName string) (EchoSlot, error) {
	if _, ok := this.Request.Intent.Slots[slotName]; ok {
		return this.Request.Intent.Slots[slotName], nil
	}

	return EchoSlot{}, errors.New("Slot name not found.")
}

func (this *EchoRequest) AllSlots() map[string]EchoSlot {
	return this.Request.Intent.Slots
}

// Locale returns the locale specified in the request.
func (this *EchoRequest) Locale() string {
	return this.Request.Locale
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

// Delegate/Elicit/Confirm a dialog or an entire intent with user of alexa.
// The func takes in name of the dialog, updated intent/intent to confirm if any and optional slot value.
// It prepares a Echo Response to be returned.
// Multiple directives can be returned by calling the method in chain (eg. RespondToIntent(...).RespondToIntent(...)
//, each RespondToIntent call appends the data to Directives array and will return the same at the end.
func (this *EchoResponse) RespondToIntent(name dialog.Type, intent *EchoIntent, slot *EchoSlot) *EchoResponse {
	directive := EchoDirective{Type: name}
	if intent != nil && name == dialog.CONFIRM_INTENT {
		directive.IntentToConfirm = intent.Name
	} else {
		directive.UpdatedIntent = intent
	}

	if slot != nil {
		if name == dialog.ELICIT_SLOT {
			directive.SlotToElicit = slot.Name
		} else if name == dialog.CONFIRM_SLOT {
			directive.SlotToConfirm = slot.Name
		}
	}
	this.Response.Directives = append(this.Response.Directives, &directive)
	return this
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
	Context EchoContext `json:"context"`
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

type EchoContext struct {
	System struct {
		Device struct {
			DeviceId string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
	} `json:"System,omitempty"`
}

type EchoReqBody struct {
	Type        string     `json:"type"`
	RequestID   string     `json:"requestId"`
	Timestamp   string     `json:"timestamp"`
	Intent      EchoIntent `json:"intent,omitempty"`
	Reason      string     `json:"reason,omitempty"`
	Locale      string     `json:"locale,omitempty"`
	DialogState string     `json:"dialogState,omitempty"`
}

type EchoIntent struct {
	Name               string              `json:"name"`
	Slots              map[string]EchoSlot `json:"slots"`
	ConfirmationStatus ConfirmationStatus  `json:"confirmationStatus"`
}

type EchoSlot struct {
	Name               string             `json:"name"`
	Value              string             `json:"value"`
	Resolutions        EchoResolution     `json:"resolutions"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus"`
}

type EchoResolution struct {
	ResolutionsPerAuthority []EchoResolutionPerAuthority `json:"resolutionsPerAuthority"`
}

type EchoResolutionPerAuthority struct {
	Authority string `json:"authority"`
	Status    struct {
		Code string `json:"code"`
	} `json:"status"`
	Values []map[string]struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"values"`
}

// Response Types

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload `json:"card,omitempty"`
	Reprompt         *EchoReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool             `json:"shouldEndSession"`
	Directives       []*EchoDirective `json:"directives,omitempty"`
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

type EchoDirective struct {
	Type            dialog.Type `json:"type"`
	UpdatedIntent   *EchoIntent `json:"updatedIntent,omitempty"`
	SlotToConfirm   string      `json:"slotToConfirm,omitempty"`
	SlotToElicit    string      `json:"slotToElicit,omitempty"`
	IntentToConfirm string      `json:"intentToConfirm,omitempty"`
}
