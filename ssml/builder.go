package ssml

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
)

func NewBuilder() (*Builder, error) {
	return &Builder{bytes.NewBufferString("")}, nil
}

func (builder *Builder) AppendPlainSpeech(text string) (*Builder, error) {
	builder.buffer.WriteString(text)
	return builder, nil
}

func (builder *Builder) AppendAmazonEffect(effect amazoneffect.Effect, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", effect, text))
	return builder, nil
}

func (builder *Builder) AppendAudio(src string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", src))
	return builder, nil
}

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

func (builder *Builder) AppendEmphasis(level emphasis.Level, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))
	return builder, nil
}

func (builder *Builder) AppendParagraph(text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))
	return builder, nil
}

func (builder *Builder) AppendProsody(rate prosody.Rate, pitch prosody.Pitch, volume prosody.Volume, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<prosody rate=\"%s\" pitch=\"%s\" volume=\"%s\">%s</prosody>", rate, pitch, volume, text))
	return builder, nil
}

func (builder *Builder) AppendSentence(text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))
	return builder, nil
}

func (builder *Builder) AppendSubstitution(alias, text string) (*Builder, error) {
	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))
	return builder, nil
}

func (builder *Builder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
