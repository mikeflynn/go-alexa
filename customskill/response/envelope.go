package response

import "fmt"

func NewEnvelope(resp *Response, session map[string]interface{}) *envelope {
	return &envelope{
		Version:           "1.0",
		Response:          resp,
		SessionAttributes: session,
	}
}

func (e *envelope) String() string {
	b, err := jsonMarshal(e)
	if err != nil {
		return fmt.Sprintf("failed to marshal JSON: %v", err)
	}

	return string(b)
}
