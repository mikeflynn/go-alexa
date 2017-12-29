package response

import (
	"fmt"
)

func New() *Response {
	return &Response{
		ShouldEndSession: Bool(true),
	}
}

func (r *Response) SetOutputSpeech(text string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return r
}

func (r *Response) SetOutputSpeechSSML(text string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: text,
	}

	return r
}

func (r *Response) SetReprompt(text string) *Response {
	r.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

func (r *Response) SetRepromptSSML(text string) *Response {
	r.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "SSML",
			SSML: text,
		},
	}

	return r
}

func (r *Response) SetCard(title, content string) *Response {
	return r.SetSimpleCard(title, content)
}

func (r *Response) SetSimpleCard(title, content string) *Response {
	r.Card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

func (r *Response) SetStandardCard(title, content, smallImg, largeImg string) *Response {
	r.Card = &Card{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if (smallImg != "") || (largeImg != "") {
		r.Card.Image = &Image{
			SmallImageURL: smallImg,
			LargeImageURL: largeImg,
		}
	}

	return r
}

func (r *Response) SetLinkAccountCard() *Response {
	r.Card = &Card{
		Type: "LinkAccount",
	}

	return r
}

func (r *Response) SetEndSession(flag *bool) *Response {
	r.ShouldEndSession = flag

	return r
}

func (r *Response) String() string {
	b, err := jsonMarshal(r)
	if err != nil {
		return fmt.Sprintf("failed to marshal JSON: %v", err)
	}

	return string(b)
}
