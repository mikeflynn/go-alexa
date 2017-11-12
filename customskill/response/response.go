package response

import "encoding/json"

func NewResponse() *Response {
	er := &Response{
		Version: "1.0",
		Body: Body{
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	return er
}

func (r *Response) OutputSpeech(text string) *Response {
	r.Body.OutputSpeech = &Payload{
		Type: "PlainText",
		Text: text,
	}

	return r
}

func (r *Response) Card(title string, content string) *Response {
	return r.SimpleCard(title, content)
}

func (r *Response) OutputSpeechSSML(text string) *Response {
	r.Body.OutputSpeech = &Payload{
		Type: "SSML",
		SSML: text,
	}

	return r
}

func (r *Response) SimpleCard(title string, content string) *Response {
	r.Body.Card = &Payload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

func (r *Response) StandardCard(title string, content string, smallImg string, largeImg string) *Response {
	r.Body.Card = &Payload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		r.Body.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		r.Body.Card.Image.LargeImageURL = largeImg
	}

	return r
}

func (r *Response) LinkAccountCard() *Response {
	r.Body.Card = &Payload{
		Type: "LinkAccount",
	}

	return r
}

func (r *Response) Reprompt(text string) *Response {
	r.Body.Reprompt = &Reprompt{
		OutputSpeech: Payload{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

func (r *Response) RepromptSSML(text string) *Response {
	r.Body.Reprompt = &Reprompt{
		OutputSpeech: Payload{
			Type: "SSML",
			Text: text,
		},
	}

	return r
}

func (r *Response) EndSession(flag bool) *Response {
	r.Body.ShouldEndSession = flag

	return r
}

func (r *Response) String() ([]byte, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
