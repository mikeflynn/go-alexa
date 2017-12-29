package response

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestNewEnvelope(t *testing.T) {
	got := NewEnvelope(nil, nil)
	want := &envelope{
		Version:           "1.0",
		Response:          nil,
		SessionAttributes: nil,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_String_Success(t *testing.T) {
	got := NewEnvelope(nil, nil).String()
	want := `{"version":"1.0","response":null}`

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
}

func TestEnvelope_String_Failure(t *testing.T) {
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("dummy error")
	}
	got := NewEnvelope(nil, nil).String()
	want := "failed to marshal JSON: dummy error"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("request mismatch:\n\tgot:    %#v\n\twanted: %#v", got, want)
	}
	// Restore the mocked jsonMarshal
	jsonMarshal = json.Marshal
}
