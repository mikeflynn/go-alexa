package response

type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              Body                   `json:"response"`
}

type Body struct {
	OutputSpeech     *Payload  `json:"outputSpeech,omitempty"`
	Card             *Payload  `json:"card,omitempty"`
	Reprompt         *Reprompt `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool      `json:"shouldEndSession"`
}

type Reprompt struct {
	OutputSpeech Payload `json:"outputSpeech,omitempty"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}
