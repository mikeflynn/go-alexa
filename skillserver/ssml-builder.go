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

// SSMLTextBuilder implements the builder pattern for constructing a speech string
// which may or may not contain SSML tags.
type SSMLTextBuilder struct {
	buffer *bytes.Buffer
}

// NewSSMLTextBuilder is a convenienve method for constructing a new SSMLTextBuilder
// instance that starts with no speech text added.
func NewSSMLTextBuilder() *SSMLTextBuilder {
	return &SSMLTextBuilder{bytes.NewBufferString("")}
}

// AppendPlainSpeech will append the supplied text as regular speech to be spoken by the Alexa device.
func (builder *SSMLTextBuilder) AppendPlainSpeech(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(text)

	return builder
}

// AppendAmazonEffect will add a new speech string with the provided effect name.
// Check the SSML reference page for a list of available effects.
func (builder *SSMLTextBuilder) AppendAmazonEffect(text, name string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<amazon:effect name=\"%s\">%s</amazon:effect>", name, text))

	return builder
}

// AppendAudio will append the playback of an MP3 file to the response. The audio playback
// will take place at the specific point in the text to speech response.
func (builder *SSMLTextBuilder) AppendAudio(src string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<audio src=\"%s\"/>", src))

	return builder
}

// AppendBreak will add a pause to the text to speech output. The default is a medium pause.
// Refer to the SSML reference for the available strength values.
func (builder *SSMLTextBuilder) AppendBreak(strength, time string) *SSMLTextBuilder {

	if strength == "" {
		// The default strength is medium
		strength = "medium"
	}

	builder.buffer.WriteString(fmt.Sprintf("<break strength=\"%s\" time=\"%s\"/>", strength, time))

	return builder
}

// AppendEmphasis will include a set of text to be spoken with the specific level of emphasis.
// Refer to the SSML reference for available emphasis level values.
func (builder *SSMLTextBuilder) AppendEmphasis(text, level string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<emphasis level=\"%s\">%s</emphasis>", level, text))

	return builder
}

// AppendParagraph will append the specific text as a new paragraph. Extra strong breaks will
// be used before and after this text.
func (builder *SSMLTextBuilder) AppendParagraph(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<p>%s</p>", text))

	return builder
}

// AppendProsody provides a way to modify the rate, pitch, and volume of a piece of spoken text.
func (builder *SSMLTextBuilder) AppendProsody(text, rate, pitch, volume string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<prosody rate=\"%s\" pitch=\"%s\" volume=\"%s\">%s</prosody>", rate, pitch, volume, text))

	return builder
}

// AppendSentence will indicate the provided text should be spoken as a new sentence. This text will
// include strong breaks before and after.
func (builder *SSMLTextBuilder) AppendSentence(text string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<s>%s</s>", text))

	return builder
}

// AppendSubstitution provides a way to indicate an alternate pronunciation for a piece of text.
func (builder *SSMLTextBuilder) AppendSubstitution(text, alias string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))

	return builder
}

// Build will construct the appropriate speech string including any SSML
// tags that were added to the Builder.
func (builder *SSMLTextBuilder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
