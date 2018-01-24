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
	want := &Response{
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetOutputSpeech(t *testing.T) {
	got := New().SetOutputSpeech("TestOutputSpeech")
	want := &Response{
		OutputSpeech: &OutputSpeech{
			Type: "PlainText",
			Text: "TestOutputSpeech",
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetOutputSpeechSSML(t *testing.T) {
	got := New().SetOutputSpeechSSML("TestOutputSpeechSSML")
	want := &Response{
		OutputSpeech: &OutputSpeech{
			Type: "SSML",
			SSML: "TestOutputSpeechSSML",
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetReprompt(t *testing.T) {
	got := New().SetReprompt("TestReprompt")
	want := &Response{
		Reprompt: &Reprompt{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: "TestReprompt",
			},
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetRepromptSSML(t *testing.T) {
	got := New().SetRepromptSSML("TestRepromptSSML")
	want := &Response{
		Reprompt: &Reprompt{
			OutputSpeech: &OutputSpeech{
				Type: "SSML",
				SSML: "TestRepromptSSML",
			},
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetCard(t *testing.T) {
	got := New().SetCard("TestTitle", "TestContent")
	want := &Response{
		Card: &Card{
			Type:    "Simple",
			Title:   "TestTitle",
			Content: "TestContent",
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetSimpleCard(t *testing.T) {
	got := New().SetSimpleCard("TestTitle", "TestContent")
	want := &Response{
		Card: &Card{
			Type:    "Simple",
			Title:   "TestTitle",
			Content: "TestContent",
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetStandardCard(t *testing.T) {
	got := New().SetStandardCard("TestTitle", "TestContent", "TestSmallImageURL", "TestLargeImageURL")
	want := &Response{
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetLinkAccountCard(t *testing.T) {
	got := New().SetLinkAccountCard()
	want := &Response{
		Card: &Card{
			Type: "LinkAccount",
		},
		ShouldEndSession: Bool(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_SetEndSession(t *testing.T) {
	var tests = []struct {
		name     string
		input    *bool
		expected *Response
		validate func(r *Response, t *testing.T)
	}{
		{
			name:  "nil-input",
			input: nil,
			expected: &Response{
				ShouldEndSession: nil,
			},
			validate: func(r *Response, t *testing.T) {
				b, err := json.Marshal(r)
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
			expected: &Response{
				ShouldEndSession: Bool(true),
			},
			validate: func(r *Response, t *testing.T) {
				b, err := json.Marshal(r)
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
			expected: &Response{
				ShouldEndSession: Bool(false),
			},
			validate: func(r *Response, t *testing.T) {
				b, err := json.Marshal(r)
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

func TestResponse_String_Success(t *testing.T) {
	got := New().String()
	want := `{"shouldEndSession":true}`

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestResponse_String_Failure(t *testing.T) {
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("dummy error")
	}
	got := New().String()
	want := "failed to marshal response to JSON: dummy error"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
	// Restore the mocked jsonMarshal
	jsonMarshal = json.Marshal
}
