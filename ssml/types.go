// package ssml provides a Speech Synthesis Markup Language (SSML) string builder for use in Alexa skills.
// Details about SSML can be found on this page: https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/speech-synthesis-markup-language-ssml-reference
package ssml

import "bytes"

type Builder struct {
	buffer *bytes.Buffer
}
