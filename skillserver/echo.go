package skillserver

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mikeflynn/go-alexa/skillserver/dialog"
)

// ConfirmationStatus represents the status of either a dialog or slot confirmation.
type ConfirmationStatus string

const (
	// ConfConfirmed indicates the intent or slot has been confirmed by the end user.
	ConfConfirmed ConfirmationStatus = "CONFIRMED"

	// ConfDenied means the end user indicated the intent or slot should NOT proceed.
	ConfDenied ConfirmationStatus = "DENIED"

	// ConfNone means there has been not acceptance or denial of the intent or slot.
	ConfNone ConfirmationStatus = "NONE"
)

// Request Functions

// VerifyTimestamp will parse the timestamp in the EchoRequest and verify that it is in the correct
// format and is not too old. True will be returned if the timestamp is valid; false otherwise.
func (r *EchoRequest) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", r.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

// VerifyAppID check that the incoming application ID matches the application ID provided
// when running the server. This is a step required for skill certification.
func (r *EchoRequest) VerifyAppID(myAppID string) bool {
	if r.Session.Application.ApplicationID == myAppID ||
		r.Context.System.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

// GetSessionID is a convenience method for getting the session ID out of an EchoRequest.
func (r *EchoRequest) GetSessionID() string {
	return r.Session.SessionID
}

// GetUserID is a convenience method for getting the user identifier out of an EchoRequest.
func (r *EchoRequest) GetUserID() string {
	return r.Session.User.UserID
}

// GetRequestType is a convenience method for getting the request type out of an EchoRequest.
func (r *EchoRequest) GetRequestType() string {
	return r.Request.Type
}

// GetIntentName is a convenience method for getting the intent name out of an EchoRequest.
func (r *EchoRequest) GetIntentName() string {
	if r.GetRequestType() == "IntentRequest" {
		return r.Request.Intent.Name
	}

	return r.GetRequestType()
}

// GetSlotValue is a convenience method for getting the value of the specified slot out of an EchoRequest
// as a string. An error is returned if a slot with that value is not found in the request.
func (r *EchoRequest) GetSlotValue(slotName string) (string, error) {
	slot, err := r.GetSlot(slotName)

	if err != nil {
		return "", err
	}

	return slot.Value, nil
}

// GetSlot will return an EchoSlot from the EchoRequest with the given name.
func (r *EchoRequest) GetSlot(slotName string) (EchoSlot, error) {
	if _, ok := r.Request.Intent.Slots[slotName]; ok {
		return r.Request.Intent.Slots[slotName], nil
	}

	return EchoSlot{}, errors.New("slot name not found")
}

// AllSlots will return a map of all the slots in the EchoRequest mapped by their name.
func (r *EchoRequest) AllSlots() map[string]EchoSlot {
	return r.Request.Intent.Slots
}

// Locale returns the locale specified in the request.
func (r *EchoRequest) Locale() string {
	return r.Request.Locale
}

// Response Functions

// NewEchoResponse will construct a new response instance with the required metadata and an empty speech string.
// By default the response will indicate that the session should be ended. Use the `EndSession(bool)` method if the
// session should be left open.
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

// OutputSpeech will replace any existing text that should be spoken with this new value. If the output
// needs to be constructed in steps or special speech tags need to be used, see the `SSMLTextBuilder`.
func (r *EchoResponse) OutputSpeech(text string) *EchoResponse {
	r.Response.OutputSpeech = &EchoRespPayload{
		Type: "PlainText",
		Text: text,
	}

	return r
}

// Card will add a card to the Alexa app's response with the provided title and content strings.
func (r *EchoResponse) Card(title string, content string) *EchoResponse {
	return r.SimpleCard(title, content)
}

// OutputSpeechSSML will add the text string provided and indicate the speech type is SSML in the response.
// This should only be used if the text to speech string includes special SSML tags.
func (r *EchoResponse) OutputSpeechSSML(text string) *EchoResponse {
	r.Response.OutputSpeech = &EchoRespPayload{
		Type: "SSML",
		SSML: text,
	}

	return r
}

// SimpleCard will indicate that a card should be included in the Alexa companion app as part of the response.
// The card will be shown with the provided title and content.
func (r *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	r.Response.Card = &EchoRespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

// StandardCard will indicate that a card should be shown in the Alexa companion app as part of the response.
// The card shown will include the provided title and content as well as images loaded from the locations provided
// as remote locations.
func (r *EchoResponse) StandardCard(title string, content string, smallImg string, largeImg string) *EchoResponse {
	r.Response.Card = &EchoRespPayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		r.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		r.Response.Card.Image.LargeImageURL = largeImg
	}

	return r
}

// LinkAccountCard is used to indicate that account linking still needs to be completed to continue
// using the Alexa skill. This will force an account linking card to be shown in the user's companion app.
func (r *EchoResponse) LinkAccountCard() *EchoResponse {
	r.Response.Card = &EchoRespPayload{
		Type: "LinkAccount",
	}

	return r
}

// Reprompt will send a prompt back to the user, this could be used to request additional information from the user.
func (r *EchoResponse) Reprompt(text string) *EchoResponse {
	r.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

// RepromptSSML is similar to the `Reprompt` method but should be used when the prompt
// to the user should include special speech tags.
func (r *EchoResponse) RepromptSSML(text string) *EchoResponse {
	r.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "SSML",
			Text: text,
		},
	}

	return r
}

// EndSession is a convenience method for setting the flag in the response that will
// indicate if the session between the end user's device and the skillserver should be closed.
func (r *EchoResponse) EndSession(flag bool) *EchoResponse {
	r.Response.ShouldEndSession = flag

	return r
}

// RespondToIntent is used to Delegate/Elicit/Confirm a dialog or an entire intent with
// user of alexa. The func takes in name of the dialog, updated intent/intent to confirm
// if any and optional slot value. It prepares a Echo Response to be returned.
// Multiple directives can be returned by calling the method in chain
// (eg. RespondToIntent(...).RespondToIntent(...), each RespondToIntent call appends the
// data to Directives array and will return the same at the end.
func (r *EchoResponse) RespondToIntent(name dialog.Type, intent *EchoIntent, slot *EchoSlot) *EchoResponse {
	directive := EchoDirective{Type: name}
	if intent != nil && name == dialog.ConfirmIntent {
		directive.IntentToConfirm = intent.Name
	} else {
		directive.UpdatedIntent = intent
	}

	if slot != nil {
		if name == dialog.ElicitSlot {
			directive.SlotToElicit = slot.Name
		} else if name == dialog.ConfirmSlot {
			directive.SlotToConfirm = slot.Name
		}
	}
	r.Response.Directives = append(r.Response.Directives, &directive)
	return r
}

func (r *EchoResponse) String() ([]byte, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// Request Types

// EchoRequest represents all fields sent from the Alexa service to the skillserver.
// Convenience methods are provided to pull commonly used properties out of the request.
type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
	Context EchoContext `json:"context"`
}

// EchoSession contains information about the ongoing session between the Alexa server and
// the skillserver. This session is stored as part of each request.
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

// EchoContext contains information about the context in which the request was sent.
// This could be information about the device from which the request was sent or about the invoked Alexa application.
type EchoContext struct {
	System struct {
		Device struct {
			DeviceID string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
	} `json:"System,omitempty"`
}

// EchoReqBody contains all data related to the type of request sent.
type EchoReqBody struct {
	Type        string     `json:"type"`
	RequestID   string     `json:"requestId"`
	Timestamp   string     `json:"timestamp"`
	Intent      EchoIntent `json:"intent,omitempty"`
	Reason      string     `json:"reason,omitempty"`
	Locale      string     `json:"locale,omitempty"`
	DialogState string     `json:"dialogState,omitempty"`
}

// EchoIntent represents the intent that is sent as part of an EchoRequest. This includes
// the name of the intent configured in the Alexa developers dashboard as well as any slots
// and the optional confirmation status if one is needed to complete an intent.
type EchoIntent struct {
	Name               string              `json:"name"`
	Slots              map[string]EchoSlot `json:"slots"`
	ConfirmationStatus ConfirmationStatus  `json:"confirmationStatus"`
}

// EchoSlot represents variable values that can be sent that were specified by the end user
// when invoking the Alexa application.
type EchoSlot struct {
	Name               string             `json:"name"`
	Value              string             `json:"value"`
	Resolutions        EchoResolution     `json:"resolutions"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus"`
}

// EchoResolution contains the results of entity resolutions when it relates to slots and how
// the values are resolved. The resolutions will be organized by authority, for custom slots
// the authority will be the custom slot type that was defined.
// Find more information here: https://developer.amazon.com/docs/custom-skills/define-synonyms-and-ids-for-slot-type-values-entity-resolution.html#intentrequest-changes
type EchoResolution struct {
	ResolutionsPerAuthority []EchoResolutionPerAuthority `json:"resolutionsPerAuthority"`
}

// EchoResolutionPerAuthority contains information about a single slot resolution from a single
// authority. The values silce will contain all possible matches for different slots.
// These resolutions are most interesting when working with synonyms.
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

// EchoResponse represents the information that should be sent back to the Alexa service
// from the skillserver.
type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

// EchoRespBody contains the body of the response to be sent back to the Alexa service.
// This includes things like the text that should be spoken or any cards that should
// be shown in the Alexa companion app.
type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload `json:"card,omitempty"`
	Reprompt         *EchoReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool             `json:"shouldEndSession"`
	Directives       []*EchoDirective `json:"directives,omitempty"`
}

// EchoReprompt contains speech that should be spoken back to the end user to retrieve
// additional information or to confirm an action.
type EchoReprompt struct {
	OutputSpeech EchoRespPayload `json:"outputSpeech,omitempty"`
}

// EchoRespImage represents a single image with two variants that should be returned as part
// of a response. Small and Large image sizes can be provided.
type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// EchoRespPayload contains the interesting parts of the Echo response including text to be spoken,
// card attributes, and images.
type EchoRespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}

// EchoDirective includes information about intents and slots that should be confirmed or elicted from the user.
// The type value can be used to delegate the action to the Alexa service. In this case, a pre-configured prompt
// will be used from the developer console.
type EchoDirective struct {
	Type            dialog.Type `json:"type"`
	UpdatedIntent   *EchoIntent `json:"updatedIntent,omitempty"`
	SlotToConfirm   string      `json:"slotToConfirm,omitempty"`
	SlotToElicit    string      `json:"slotToElicit,omitempty"`
	IntentToConfirm string      `json:"intentToConfirm,omitempty"`
}
