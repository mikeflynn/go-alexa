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

// WordRole is used as the role argument in the AppendPartOfSpeech method. This should
// be one of the constants defined in the Amazon SSML Reference docs.
type WordRole string

const (
	// PresentSimple is used to pronounce the word as a verb
	PresentSimple WordRole = "amazon:VB"
	// PastParticle is used to pronounce the word as a past particle
	PastParticle WordRole = "amazon:VBD"
	// Noun is used to pronounce the word as a noun
	Noun WordRole = "amazon:NN"
	// AlternateSense is used to select the alternate sense for a specific word. According
	// to the Amazon SSML Reference:
	// 		" Use the non-default sense of the word. For example, the noun "bass" is pronounced
	//		differently depending on meaning. The "default" meaning is the lowest part of the
	// 		musical range. The alternate sense (which is still a noun) is a freshwater fish.
	// 		Specifying <speak><w role="amazon:SENSE_1">bass</w>"</speak> renders the non-default
	// 		pronunciation (freshwater fish)."
	AlternateSense WordRole = "amazon:SENSE_1"
)

// PhoneticAlphabet represents the alphabet to be used when appending phonemes
type PhoneticAlphabet string

const (
	// Ipa is the International Phonetic Alphabet
	Ipa PhoneticAlphabet = "ipa"
	// XSampa is the Extended Speech Assesment Methods Phonetic Alphabet
	XSampa PhoneticAlphabet = "x-sampa"
)

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

// AppendPartOfSpeech is used to explictily define the part of speech for a word that is being
// appended to the text to speech output sent in a skill server response.
func (builder *SSMLTextBuilder) AppendPartOfSpeech(role WordRole, text string) *SSMLTextBuilder {

	if role != "" {
		builder.buffer.WriteString(fmt.Sprintf("<w role=\"%s\">%s</w>", role, text))
	}

	return builder
}

// AppendSubstitution provides a way to indicate an alternate pronunciation for a piece of text.
func (builder *SSMLTextBuilder) AppendSubstitution(text, alias string) *SSMLTextBuilder {

	builder.buffer.WriteString(fmt.Sprintf("<sub alias=\"%s\">%s</sub>", alias, text))

	return builder
}

// AppendSayAs is used to provide additional information about how the text string being appended
// should be interpreted. For example this can be used to interpret the string as a list
// of individual characters or to read out digits one at a time. The format string is
// ignored unless the interpret-as argument is `date`. Refer to the SSML referene for valid
// values for the interpretAs parameter.
func (builder *SSMLTextBuilder) AppendSayAs(interpretAs, format, text string) *SSMLTextBuilder {

	if interpretAs == "date" {
		builder.buffer.WriteString(fmt.Sprintf("<say-as interpret-as=\"%s\" format=\"%s\">%s</say-as>",
			interpretAs, format, text))
	} else if interpretAs != "" {
		builder.buffer.WriteString(fmt.Sprintf("<say-as interpret-as=\"%s\">%s</say-as>", interpretAs, text))
	}

	return builder
}

// AppendPhoneme is used to specify a phonetic pronunciation for a piece of text to be appended
// to the response.
func (builder *SSMLTextBuilder) AppendPhoneme(alphabet PhoneticAlphabet, phoneme, text string) *SSMLTextBuilder {

	if phoneme != "" && text != "" && alphabet != PhoneticAlphabet("") {
		builder.buffer.WriteString(fmt.Sprintf("<phoneme alphabet=\"%s\" ph=\"%s\">%s</phoneme>", alphabet, phoneme, text))
	}

	return builder
}

// Build will construct the appropriate speech string including any SSML
// tags that were added to the Builder.
func (builder *SSMLTextBuilder) Build() string {
	return fmt.Sprintf("<speak>%s</speak>", builder.buffer.String())
}
