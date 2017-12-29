package response

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Envelope{
		Version: "1.0",
		Response: Response{
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
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
			ShouldEndSession: Bool(true),
		},
		SessionAttributes: make(map[string]interface{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_SetEndSession(t *testing.T) {
	var tests = []struct {
		name     string
		input    *bool
		expected *Envelope
		validate func(e *Envelope, T *testing.T)
	}{
		{
			name:  "nil-input",
			input: nil,
			expected: &Envelope{
				Version: "1.0",
				Response: Response{
					ShouldEndSession: nil,
				},
				SessionAttributes: make(map[string]interface{}),
			},
			validate: func(e *Envelope, t *testing.T) {
				b, err := json.Marshal(e)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if strings.Contains(string(b), "shouldEndSession") {
					t.Errorf("unexpected shouldEndSession in marshaled JSON: %s", string(b))
				}
			},
		},
		{
			name:  "true-input",
			input: Bool(true),
			expected: &Envelope{
				Version: "1.0",
				Response: Response{
					ShouldEndSession: Bool(true),
				},
				SessionAttributes: make(map[string]interface{}),
			},
			validate: func(e *Envelope, t *testing.T) {
				b, err := json.Marshal(e)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				want := `"shouldEndSession":true`
				if !strings.Contains(string(b), want) {
					t.Errorf("expected string %s in marshaled JSON: %s", want, string(b))
				}
			},
		},
		{
			name:  "false-input",
			input: Bool(false),
			expected: &Envelope{
				Version: "1.0",
				Response: Response{
					ShouldEndSession: Bool(false),
				},
				SessionAttributes: make(map[string]interface{}),
			},
			validate: func(e *Envelope, t *testing.T) {
				b, err := json.Marshal(e)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				want := `"shouldEndSession":false`
				if !strings.Contains(string(b), want) {
					t.Errorf("expected string %s in marshaled JSON: %s", want, string(b))
				}
			},
		},
	}

	for _, test := range tests {
		got := New().SetEndSession(test.input)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("%s: request mismatch:\n\tgot:    %#v\n\twanted: %#v", test.name, got, test.expected)
		}
		test.validate(got, t)
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
