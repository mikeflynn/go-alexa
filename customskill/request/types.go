package request

// Request Types\
type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    Body    `json:"request"`
}

type Session struct {
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

type Body struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Intent    Intent `json:"intent,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type Intent struct {
	Name  string          `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

type Slot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
