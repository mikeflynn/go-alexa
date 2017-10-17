package ssml

import (
	"fmt"

	"time"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
)

func ExampleNewBuilder_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak></speak>
}

func ExampleBuilder_AppendPlainSpeech_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append the plain speech.
	b.AppendPlainSpeech("Hello World!")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak>Hello World!</speak>
}

func ExampleBuilder_AppendAmazonEffect_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append some Whispered text.
	b.AppendAmazonEffect(amazoneffect.Whispered, "This is whispered")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><amazon:effect name="whispered">This is whispered</amazon:effect></speak>
}

func ExampleBuilder_AppendAudio_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append an audio url.
	b.AppendAudio("https://domain.tld/dummy.mp3")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><audio src="https://domain.tld/dummy.mp3"/></speak>
}

func ExampleBuilder_AppendBreak_Strength_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a strong break.
	b.AppendBreak(pause.Strong)

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><break strength="strong"/></speak>
}

func ExampleBuilder_AppendBreak_Time_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a one second break.
	b.AppendBreak(time.Second)

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><break time="1000ms"/></speak>
}

func ExampleBuilder_AppendEmphasis_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append some reduced emphasis text.
	b.AppendEmphasis(emphasis.Reduced, "reduced emphasis")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><emphasis level="reduced">reduced emphasis</emphasis></speak>
}

func ExampleBuilder_AppendParagraph_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a paragraph.
	b.AppendParagraph("sample paragraph")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><p>sample paragraph</p></speak>
}

func ExampleBuilder_AppendProsody_Rate_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a prosody element which just modified the rate.
	b.AppendProsody(prosody.RateMedium, nil, nil, "this is said at a medium rate")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><prosody rate="medium">this is said at a medium rate</prosody></speak>
}

func ExampleBuilder_AppendProsody_Constants_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a prosody element using pre-defined constants.
	b.AppendProsody(prosody.RateSlow, prosody.PitchLow, prosody.VolumeXLoud, "this is said slowly, in a low pitch, but loudly")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><prosody rate="slow" pitch="low" volume="x-loud">this is said slowly, in a low pitch, but loudly</prosody></speak>
}

func ExampleBuilder_AppendProsody_Int_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a prosody element using ints.
	b.AppendProsody(110, -10, 4, "this is said slightly quicker than normal (110%), in a lower pitch (-10%), loudly (+4dB)")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><prosody rate="110%" pitch="-10%" volume="+4dB">this is said slightly quicker than normal (110%), in a lower pitch (-10%), loudly (+4dB)</prosody></speak>
}

func ExampleBuilder_AppendSentence_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a sentence.
	b.AppendSentence("this is a sentence")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><s>this is a sentence</s></speak>
}

func ExampleBuilder_AppendSubstitution_output() {
	// Create a new builder. Ignore any error.
	b, _ := NewBuilder()

	// Append a substation element.
	b.AppendSubstitution("alias", "replacement")

	// Print the built string.
	fmt.Print(b.Build())
	// Output: <speak><sub alias="alias">replacement</sub></speak>
}
