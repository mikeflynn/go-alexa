package skillserver

import (
	"bytes"
	"fmt"
)

/**
 * Details about the Speech Synthesis Markup Language (SSML) can be found on this page:
 * https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/speech-synthesis-markup-language-ssml-reference
 */

// Helper Types

type SSMLTextBuilder struct {
	buffer *bytes.Buffer
}

func NewSSMLTextBuilder() *SSMLTextBuilder {
	return &SSMLTextBuilder{bytes.NewBufferString("")}
}

func (builder *SSMLTextBuilder) AppendPlainSpeech(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(text)

	return builder
}

func (builder *SSMLTextBuilder) AppendAmazonEffect(text, name string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", name, text))

	return builder
}

func (builder *SSMLTextBuilder) AppendAudio(src string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", src))

	return builder
}

func (builder *SSMLTextBuilder) AppendBreak(strength, time string) *SSMLTextBuilder {

	if strength == "" {
		// The default strength is medium
		strength = "medium"
	}

	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\" time=\"%s\"/>", strength, time))

	return builder
}

func (builder *SSMLTextBuilder) AppendEmphasis(text, level string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))

	return builder
}

func (builder *SSMLTextBuilder) AppendParagraph(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))

	return builder
}

func (builder *SSMLTextBuilder) AppendProsody(text, rate, pitch, volume string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<prosody rate=\"%s\" pitch=\"%s\" volume=\"%s\">%s</prosody>", rate, pitch, volume, text))

	return builder
}

func (builder *SSMLTextBuilder) AppendSentence(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))

	return builder
}

func (builder *SSMLTextBuilder) AppendSubstitution(text, alias string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))

	return builder
}

func (builder *SSMLTextBuilder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
