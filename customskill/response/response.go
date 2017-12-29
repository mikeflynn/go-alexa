package response

import (
	"encoding/json"
	"fmt"
)

var jsonMarshal = json.Marshal // Used to enable unit testing

func New() *Envelope {
	e := &Envelope{
		Version: "1.0",
		Response: Response{
			ShouldEndSession: Bool(true),
		},
		SessionAttributes: make(map[string]interface{}),
	}

	return e
}

func (e *Envelope) SetOutputSpeech(text string) *Envelope {
	e.Response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return e
}

func (e *Envelope) SetOutputSpeechSSML(text string) *Envelope {
	e.Response.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: text,
	}

	return e
}

func (e *Envelope) SetReprompt(text string) *Envelope {
	e.Response.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}

	return e
}

func (e *Envelope) SetRepromptSSML(text string) *Envelope {
	e.Response.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "SSML",
			SSML: text,
		},
	}

	return e
}

func (e *Envelope) SetCard(title, content string) *Envelope {
	return e.SetSimpleCard(title, content)
}

func (e *Envelope) SetSimpleCard(title, content string) *Envelope {
	e.Response.Card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return e
}

func (e *Envelope) SetStandardCard(title, content, smallImg, largeImg string) *Envelope {
	e.Response.Card = &Card{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if (smallImg != "") || (largeImg != "") {
		e.Response.Card.Image = &Image{
			SmallImageURL: smallImg,
			LargeImageURL: largeImg,
		}
	}

	return e
}

func (e *Envelope) SetLinkAccountCard() *Envelope {
	e.Response.Card = &Card{
		Type: "LinkAccount",
	}

	return e
}

func (e *Envelope) SetEndSession(flag *bool) *Envelope {
	e.Response.ShouldEndSession = flag

	return e
}

func (e *Envelope) String() string {
	b, err := jsonMarshal(e)
	if err != nil {
		return fmt.Sprintf("failed to marshal JSON: %v", err)
	}

	return string(b)
}

func Bool(bool bool) *bool {
	return &bool
}
