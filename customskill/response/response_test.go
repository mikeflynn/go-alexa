package response

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetOutputSpeech(t *testing.T) {
	got := New().SetOutputSpeech("TestOutputSpeech")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: "TestOutputSpeech",
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetOutputSpeechSSML(t *testing.T) {
	got := New().SetOutputSpeechSSML("TestOutputSpeechSSML")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			OutputSpeech: &OutputSpeech{
				Type: "SSML",
				SSML: "TestOutputSpeechSSML",
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetReprompt(t *testing.T) {
	got := New().SetReprompt("TestReprompt")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Reprompt: &Reprompt{
				OutputSpeech: &OutputSpeech{
					Type: "PlainText",
					Text: "TestReprompt",
				},
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetRepromptSSML(t *testing.T) {
	got := New().SetRepromptSSML("TestRepromptSSML")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Reprompt: &Reprompt{
				OutputSpeech: &OutputSpeech{
					Type: "SSML",
					SSML: "TestRepromptSSML",
				},
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetCard(t *testing.T) {
	got := New().SetCard("TestTitle", "TestContent")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Card: &Card{
				Type:    "Simple",
				Title:   "TestTitle",
				Content: "TestContent",
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetSimpleCard(t *testing.T) {
	got := New().SetSimpleCard("TestTitle", "TestContent")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Card: &Card{
				Type:    "Simple",
				Title:   "TestTitle",
				Content: "TestContent",
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetStandardCard(t *testing.T) {
	got := New().SetStandardCard("TestTitle", "TestContent", "TestSmallImageURL", "TestLargeImageURL")
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Card: &Card{
				Type:    "Standard",
				Title:   "TestTitle",
				Content: "TestContent",
				Image: &Image{
					SmallImageURL: "TestSmallImageURL",
					LargeImageURL: "TestLargeImageURL",
				},
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetLinkAccountCard(t *testing.T) {
	got := New().SetLinkAccountCard()
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			Card: &Card{
				Type: "LinkAccount",
			},
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetEndSession(t *testing.T) {
	got := New().SetEndSession(true)
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			ShouldEndSession: true,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}

	got = New().SetEndSession(false)
	want = &Envelope{
		Version: "1.0",
		Response: Response{
			ShouldEndSession: false,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_String_Success(t *testing.T) {
	got := New().String()

	b, err := json.Marshal(New())
	if err != nil {
		t.Error(err)
	}
	want := string(b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_String_Failure(t *testing.T) {
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("dummy error")
	}
	got := New().String()
	want := "failed to marshal JSON: dummy error"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}
