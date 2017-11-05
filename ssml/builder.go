package ssml

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
	"github.com/pkg/errors"
)

// NewBuilder returns an empty new SSML builder.
func NewBuilder() (*Builder, error) {
	return &Builder{
		buffer: bytes.NewBufferString(""),
	}, nil
}

// AppendPlainSpeech appends raw text to the builder's internal SSML string.
// It will not append an error to the builder's internal error slice.
// It returns a pointer to the builder
func (builder *Builder) AppendPlainSpeech(text string) *Builder {
	builder.buffer.WriteString(text)
	return builder
}

// AppendAmazonEffect appends an AmazonEffect to the builder's internal SSML string.
// Valid Effects can be found in the amazoneffect sub-package.
// It will not append an error to the builder's internal error slice.
// It returns a pointer to the builder
func (builder *Builder) AppendAmazonEffect(effect amazoneffect.Effect, text string) *Builder {
	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", effect, text))
	return builder
}

// AppendAmazonEffect appends an audio element to the builder's internal SSML string.
// It will append an error to the builder's internal error slice if the src is an invalid URL or not a HTTPs URL.
// It returns a pointer to the builder.
func (builder *Builder) AppendAudio(src string) *Builder {
	u, err := url.Parse(src)
	if err != nil {
		return builder.appendError(fmt.Errorf("failed to parse src into a valid URL: %v", err))
	}
	if u.Scheme != "https" {
		return builder.appendError(errors.New("unsupported URL scheme type: must be https"))
	}
	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", u.String()))
	return builder
}

// AppendAmazonEffect appends a break/pause element to the builder's internal SSML string.
// strengthOrDuration must either be of type Strength (from the pause sub-package) or time.Duration.
// It will append an error to the builder's internal error slice if strengthOrDuration is of an invalid type.
// It returns a pointer to the builder.
func (builder *Builder) AppendBreak(strengthOrDuration interface{}) *Builder {
	strength, ok := strengthOrDuration.(pause.Strength)
	if !ok {
		duration, ok := strengthOrDuration.(time.Duration)
		if !ok {
			return builder.appendError(errors.New("unsupported parameter type: must be either pause.Strength or time.Duration"))
		}
		builder.buffer.WriteString(fmt.Sprintf("<break time=\"%dms\"/>", duration.Nanoseconds()/1e6))
		return builder
	}
	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\"/>", strength))
	return builder
}

// AppendEmphasis appends an emphasis element to the builder's internal SSML string.
// It will not append an error to the builder's internal error slice.
// It returns a pointer to the builder
func (builder *Builder) AppendEmphasis(level emphasis.Level, text string) *Builder {
	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))
	return builder
}

// AppendParagraph appends a paragraph element to the builder's internal SSML string.
// It will not append an error to the builder's internal error slice.
// It returns a pointer to the builder.
func (builder *Builder) AppendParagraph(text string) *Builder {
	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))
	return builder
}

// AppendProsody appends a prosody element to the builder's internal SSML string.
// rate must either be nil or a Rate (from the prosody.Rate sub-package) or an int. If nil no rate value is included
// in the prosody element.
// pitch must either be nil or a Pitch (from the prosody.Pitch sub-package) or an int. If nil no pitch value is
// included in the prosody element.
// volume must either be nil or a Volume (from the prosody.Volume sub-package) or an int. If nil no volume value is
// included in the prosody element.
// It returns an error if a parameter is of an invalid type.
func (builder *Builder) AppendProsody(rate, pitch, volume interface{}, text string) *Builder {
	src := ""
	if rate != nil {
		rateStr, ok := rate.(prosody.Rate)
		if !ok {
			ratePercent, ok := rate.(int)
			if !ok {
				return builder.appendError(errors.New("unsupported rate type: must be either prosody.Rate or int"))
			}
			src += fmt.Sprintf(" rate=\"%d%%\"", ratePercent)
		} else {
			src += fmt.Sprintf(" rate=\"%s\"", rateStr)
		}
	}

	if pitch != nil {
		pitchStr, ok := pitch.(prosody.Pitch)
		if !ok {
			pitchPercent, ok := pitch.(int)
			if !ok {
				return builder.appendError(errors.New("unsupported pitch type: must be either prosody.Pitch or int"))
			}

			if pitchPercent > 0 {
				src += fmt.Sprintf(" pitch=\"+%d%%\"", pitchPercent)
			} else {
				src += fmt.Sprintf(" pitch=\"%d%%\"", pitchPercent)
			}

		} else {
			src += fmt.Sprintf(" pitch=\"%s\"", pitchStr)
		}
	}

	if volume != nil {
		volumeStr, ok := volume.(prosody.Volume)
		if !ok {
			volumeDb, ok := volume.(int)
			if !ok {
				return builder.appendError(errors.New("unsupported volume type: must be either prosody.Volume or int"))
			}
			if volumeDb > 0 {
				src += fmt.Sprintf(" volume=\"+%ddB\"", volumeDb)
			} else {
				src += fmt.Sprintf(" volume=\"%ddB\"", volumeDb)
			}
		} else {
			src += fmt.Sprintf(" volume=\"%s\"", volumeStr)
		}
	}

	builder.buffer.WriteString(fmt.Sprintf("<prosody%s>%s</prosody>", src, text))
	return builder
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
func (builder *Builder) Build() (string, []error) {
	builder.lock.RLock()
	defer builder.lock.RUnlock()
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String()), builder.errs
}

func (builder *Builder) appendError(err error) *Builder {
	builder.lock.Lock()
	defer builder.lock.Unlock()
	builder.errs = append(builder.errs, err)
	return builder
}
