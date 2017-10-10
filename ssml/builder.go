package ssml

import (
	"bytes"
	"fmt"
)

/**
 * Details about the Speech Synthesis Markup Language (SSML) can be found on this page:
 * https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/speech-synthesis-markup-language-ssml-reference
 */

func NewBuilder() *builder {
	return &builder{bytes.NewBufferString("")}
}

func (builder *builder) AppendPlainSpeech(text string) *builder {
	builder.buffer.WriteString(text)
	return builder
}

func (builder *builder) AppendAmazonEffect(name, text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", name, text))
	return builder
}

func (builder *builder) AppendAudio(src string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", src))
	return builder
}

func (builder *builder) AppendBreak(strength, time string) *builder {
	if strength == "" {
		// The default strength is medium
		strength = "medium"
	}
	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\" time=\"%s\"/>", strength, time))
	return builder
}

func (builder *builder) AppendEmphasis(level, text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))
	return builder
}

func (builder *builder) AppendParagraph(text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))
	return builder
}

func (builder *builder) AppendProsody(rate, pitch, volume, text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<prosody rate=\"%s\" pitch=\"%s\" volume=\"%s\">%s</prosody>", rate, pitch, volume, text))
	return builder
}

func (builder *builder) AppendSentence(text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))
	return builder
}

func (builder *builder) AppendSubstitution(alias, text string) *builder {
	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))
	return builder
}

func (builder *builder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
