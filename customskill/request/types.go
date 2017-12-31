package request

type Metadata struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Context Context `json:"context"`
}

/* Request types */
type LaunchRequest struct {
	Common
}

type IntentRequest struct {
	Common
	DialogState string `json:"dialogState"`
	Intent      Intent `json:"intent,omitempty"`
}

type SessionEndedRequest struct {
	Common
	Error  Error  `json:"error"`
	Reason string `json:"reason,omitempty"`
}

/* Types shared across request types */

// Common contains properties common to all requests
type Common struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Locale    string `json:"locale"`
}

type Session struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Application Application            `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        User                   `json:"user"`
}

type Context struct {
	System      System      `json:"System"`
	AudioPlayer AudioPlayer `json:"AudioPlayer"`
}

type System struct {
	Application Application `json:"application"`
	User        User        `json:"user"`
	Device      Device      `json:"device"`
	ApiEndpoint string      `json:"apiEndpoint"`
}

type Application struct {
	ApplicationID string `json:"applicationId"`
}

type Device struct {
	DeviceId            string                 `json:"deviceId"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

type AudioPlayer struct {
	Token                string `json:"token"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	PlayerActivity       string `json:"playerActivity"`
}

type User struct {
	UserID      string      `json:"userId"`
	AccessToken string      `json:"accessToken,omitempty"`
	Permissions Permissions `json:"permissions"`
}

type Permissions struct {
	ConsentToken string `json:"consentToken"`
}

type Intent struct {
	Name               string          `json:"name"`
	ConfirmationStatus string          `json:"confirmationStatus"`
	Slots              map[string]Slot `json:"slots"`
}

type Slot struct {
	Name               string      `json:"name"`
	Value              string      `json:"value"`
	ConfirmationStatus string      `json:"confirmationStatus"`
	Resolutions        interface{} `json:"resolutions"` // TODO: Improve this
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

/* Unexported types used for unmarshalling requests */
type envelope struct {
	Metadata
	Request Common `json:"request"`
}

type launchRequestEnvelope struct {
	Request LaunchRequest `json:"request"`
}

type intentRequestEnvelope struct {
	Request IntentRequest `json:"request"`
}

type sessionEndedRequestEnvelope struct {
	Request SessionEndedRequest `json:"request"`
}
