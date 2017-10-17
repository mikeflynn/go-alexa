package ssml

import (
	"bytes"
	"fmt"
	"time"

	"net/url"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
)

// NewBuilder returns an empty new SSML builder.
func NewBuilder() (*Builder, error) {
	return &Builder{bytes.NewBufferString("")}, nil
}

// AppendPlainSpeech appends raw text to the builder's internal SSML string.
func (builder *Builder) AppendPlainSpeech(text string) (*Builder, error) {
	builder.buffer.WriteString(text)
	return builder, nil
}

// AppendAmazonEffect appends an AmazonEffect to the builder's internal SSML string.
// Valid Effects can be found in the amazoneffect sub-package
func (builder *Builder) AppendAmazonEffect(effect amazoneffect.Effect, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", effect, text))
	return builder, nil
}

// AppendAmazonEffect appends an audio element to the builder's internal SSML string.
// src must be a valid HTTPS url.
// It returns the builder pointer and an error if the src is an invalid URL.
func (builder *Builder) AppendAudio(src string) (*Builder, error) {
	u, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("src failed to parse into a valid URL: %v", err)
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("src must be a HTTPS URL. Scheme %s not valid", u.Scheme)
	}
	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", u.String()))
	return builder, nil
}

// AppendAmazonEffect appends a break/pause element to the builder's internal SSML string.
// strengthOrDuration must either be of type Strength (from the pause sub-package) or time.Duration
// It returns the builder pointer and an error if strengthOrDuration is of an invalid type
func (builder *Builder) AppendBreak(strengthOrDuration interface{}) (*Builder, error) {
	strength, ok := strengthOrDuration.(pause.Strength)
	if !ok {
		duration, ok := strengthOrDuration.(time.Duration)
		if !ok {
			return builder, fmt.Errorf("unsupported parameter type. must be either pause.Strength or time.Duration")
		}
		builder.buffer.WriteString(fmt.Sprintf("<break time=\"%dms\"/>", duration.Nanoseconds()/1e6))
		return builder, nil
	}
	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\"/>", strength))
	return builder, nil
}

// AppendEmphasis appends an emphasis element to the builder's internal SSML string.
// It returns the builder pointer and a nil error.
func (builder *Builder) AppendEmphasis(level emphasis.Level, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))
	return builder, nil
}

// AppendParagraph appends a paragraph element to the builder's internal SSML string.
// It returns the builder pointer and a nil error.
func (builder *Builder) AppendParagraph(text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))
	return builder, nil
}

// AppendProsody appends a prosody element to the builder's internal SSML string.
// It returns the builder pointer and an error if any parameters fall outside their accepted ranges.
// TODO: Validate parameters
func (builder *Builder) AppendProsody(rate prosody.Rate, pitch prosody.Pitch, volume prosody.Volume, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<prosody rate=\"%s\" pitch=\"%s\" volume=\"%s\">%s</prosody>", rate, pitch, volume, text))
	return builder, nil
}

// AppendSentence appends a sentence element to the builder's internal SSML string.
// It returns the builder pointer and a nil error.
func (builder *Builder) AppendSentence(text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))
	return builder, nil
}

// AppendSubstitution appends a substitution element to the builder's internal SSML string.
// It returns the builder pointer and a nil error.
func (builder *Builder) AppendSubstitution(alias, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))
	return builder, nil
}

// Build builds the SSML string.
// It returns the SSML string.
func (builder *Builder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
