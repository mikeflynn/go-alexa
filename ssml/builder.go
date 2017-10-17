package ssml

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

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
// It returns a nil error.
func (builder *Builder) AppendPlainSpeech(text string) error {
	builder.buffer.WriteString(text)
	return nil
}

// AppendAmazonEffect appends an AmazonEffect to the builder's internal SSML string.
// Valid Effects can be found in the amazoneffect sub-package
// It returns a nil error.
func (builder *Builder) AppendAmazonEffect(effect amazoneffect.Effect, text string) error {
	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", effect, text))
	return nil
}

// AppendAmazonEffect appends an audio element to the builder's internal SSML string.
// src must be a valid HTTPS url.
// It returns an error if the src is an invalid URL.
func (builder *Builder) AppendAudio(src string) error {
	u, err := url.Parse(src)
	if err != nil {
		return fmt.Errorf("src failed to parse into a valid URL: %v", err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("src must be a HTTPS URL. Scheme %s not valid", u.Scheme)
	}
	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", u.String()))
	return nil
}

// AppendAmazonEffect appends a break/pause element to the builder's internal SSML string.
// strengthOrDuration must either be of type Strength (from the pause sub-package) or time.Duration
// It returns an error if strengthOrDuration is of an invalid type
func (builder *Builder) AppendBreak(strengthOrDuration interface{}) error {
	strength, ok := strengthOrDuration.(pause.Strength)
	if !ok {
		duration, ok := strengthOrDuration.(time.Duration)
		if !ok {
			return fmt.Errorf("unsupported parameter type. must be either pause.Strength or time.Duration")
		}
		builder.buffer.WriteString(fmt.Sprintf("<break time=\"%dms\"/>", duration.Nanoseconds()/1e6))
		return nil
	}
	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\"/>", strength))
	return nil
}

// AppendEmphasis appends an emphasis element to the builder's internal SSML string.
// It returns a nil error.
func (builder *Builder) AppendEmphasis(level emphasis.Level, text string) error {
	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))
	return nil
}

// AppendParagraph appends a paragraph element to the builder's internal SSML string.
// It returns a nil error.
func (builder *Builder) AppendParagraph(text string) error {
	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))
	return nil
}

// AppendProsody appends a prosody element to the builder's internal SSML string.
// rate must either be nil or a Rate (from the prosody.Rate sub-package) or an int. If nil no rate value is included
// in the prosody element.
// pitch must either be nil or a Pitch (from the prosody.Pitch sub-package) or an int. If nil no pitch value is
// included in the prosody element.
// volume must either be nil or a Volume (from the prosody.Volume sub-package) or an int. If nil no volume value is
// included in the prosody element.
// It returns an error if a parameter is of an invalid type.
func (builder *Builder) AppendProsody(rate, pitch, volume interface{}, text string) error {
	src := ""
	if rate != nil {
		rateStr, ok := rate.(prosody.Rate)
		if !ok {
			ratePercent, ok := rate.(int)
			if !ok {
				return fmt.Errorf("unsupported rate type. must be either prosody.Rate or int")
			}
			src += fmt.Sprintf("rate=\"%d%%\" ", ratePercent)
		} else {
			src += fmt.Sprintf("rate=\"%s\" ", rateStr)
		}
	}

	if pitch != nil {
		pitchStr, ok := pitch.(prosody.Pitch)
		if !ok {
			pitchPercent, ok := pitch.(int)
			if !ok {
				return fmt.Errorf("unsupported pitch type. must be either prosody.Pitch or int")
			}
			sign := ""
			if pitchPercent > 0 {
				sign = "+"
			}
			src += fmt.Sprintf("pitch=\"%s%d%%\" ", sign, pitchPercent)
		} else {
			src += fmt.Sprintf("pitch=\"%s\" ", pitchStr)
		}
	}

	if volume != nil {
		volumeStr, ok := volume.(prosody.Volume)
		if !ok {
			volumeDb, ok := volume.(int)
			if !ok {
				return fmt.Errorf("unsupported volume type. must be either prosody.Volume or int")
			}
			sign := ""
			if volumeDb > 0 {
				sign = "+"
			}
			src += fmt.Sprintf("volume=\"%s%ddB\" ", sign, volumeDb)
		} else {
			src += fmt.Sprintf("volume=\"%s\" ", volumeStr)
		}
	}

	builder.buffer.WriteString(fmt.Sprintf("<prosody %s>%s</prosody>", strings.TrimSpace(src), text))
	return nil
}

// AppendSentence appends a sentence element to the builder's internal SSML string.
// It returns a nil error.
func (builder *Builder) AppendSentence(text string) error {
	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))
	return nil
}

// AppendSubstitution appends a substitution element to the builder's internal SSML string.
// It returns a nil error.
func (builder *Builder) AppendSubstitution(alias, text string) error {
	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))
	return nil
}

// Build builds the SSML string.
// It returns the SSML string.
func (builder *Builder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
