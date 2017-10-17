package ssml

import (
	"testing"

	"time"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
)

func TestNewBuilder_ReturnsEmptySSML(t *testing.T) {
	b, err := NewBuilder()

	if err != nil {
		t.Fatalf("failed to get new builder: expected no error, got :%v", err)
	}

	actual := b.Build()
	expected := "<speak></speak>"
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
	}
}

func TestBuilder_AppendPlainSpeech(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendPlainSpeech("hello ")
	b.AppendPlainSpeech("world")

	actual := b.Build()
	expected := "<speak>hello world</speak>"
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
	}
}

func TestBuilder_AppendAmazonEffect(t *testing.T) {
	tests := []struct {
		name     string
		effect   amazoneffect.Effect
		expected string
	}{
		{
			name:     "whispered",
			effect:   amazoneffect.Whispered,
			expected: `<speak><amazon:effect name="whispered">text1</amazon:effect><amazon:effect name="whispered">text2</amazon:effect></speak>`,
		},
		{
			name:     "custom",
			effect:   amazoneffect.Effect("custom"),
			expected: `<speak><amazon:effect name="custom">text1</amazon:effect><amazon:effect name="custom">text2</amazon:effect></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendAmazonEffect(test.effect, "text1")
		b.AppendAmazonEffect(test.effect, "text2")

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendAudio(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendAudio("source1")
	b.AppendAudio("source2")

	actual := b.Build()
	expected := `<speak><audio src="source1"/><audio src="source2"/></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
	}
}

func TestBuilder_AppendBreak(t *testing.T) {
	tests := []struct {
		name     string
		param    interface{}
		err      bool
		expected string
	}{
		{
			name:     "default",
			param:    pause.Default,
			err:      false,
			expected: `<speak><break strength="medium"/><break strength="medium"/></speak>`,
		},
		{
			name:     "none",
			param:    pause.None,
			err:      false,
			expected: `<speak><break strength="none"/><break strength="none"/></speak>`,
		},
		{
			name:     "x-weak",
			param:    pause.XWeak,
			err:      false,
			expected: `<speak><break strength="x-weak"/><break strength="x-weak"/></speak>`,
		},
		{
			name:     "weak",
			param:    pause.Weak,
			err:      false,
			expected: `<speak><break strength="weak"/><break strength="weak"/></speak>`,
		},
		{
			name:     "medium",
			param:    pause.Medium,
			err:      false,
			expected: `<speak><break strength="medium"/><break strength="medium"/></speak>`,
		},
		{
			name:     "strong",
			param:    pause.Strong,
			err:      false,
			expected: `<speak><break strength="strong"/><break strength="strong"/></speak>`,
		},
		{
			name:     "x-strong",
			param:    pause.XStrong,
			err:      false,
			expected: `<speak><break strength="x-strong"/><break strength="x-strong"/></speak>`,
		},
		{
			name:     "custom",
			param:    pause.Strength("custom"),
			err:      false,
			expected: `<speak><break strength="custom"/><break strength="custom"/></speak>`,
		},
		{
			name:     "time",
			param:    time.Second,
			err:      false,
			expected: `<speak><break time="1000ms"/><break time="1000ms"/></speak>`,
		},
		{
			name:     "invalidType",
			param:    4,
			err:      true,
			expected: `<speak></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		_, err := b.AppendBreak(test.param)
		if (err != nil) != test.err {
			t.Errorf("%s: error mismatch: expected %t, got %v", test.name, test.err, err)
		}

		_, err = b.AppendBreak(test.param)
		if (err != nil) != test.err {
			t.Errorf("%s: error mismatch: expected %t, got %v", test.name, test.err, err)
		}

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendEmphasis(t *testing.T) {
	tests := []struct {
		name     string
		level    emphasis.Level
		text     string
		expected string
	}{
		{
			name:     "default",
			level:    emphasis.Default,
			expected: `<speak><emphasis level="moderate">text1</emphasis><emphasis level="moderate">text2</emphasis></speak>`,
		},
		{
			name:     "strong",
			level:    emphasis.Strong,
			expected: `<speak><emphasis level="strong">text1</emphasis><emphasis level="strong">text2</emphasis></speak>`,
		},
		{
			name:     "moderate",
			level:    emphasis.Moderate,
			expected: `<speak><emphasis level="moderate">text1</emphasis><emphasis level="moderate">text2</emphasis></speak>`,
		},
		{
			name:     "reduced",
			level:    emphasis.Reduced,
			expected: `<speak><emphasis level="reduced">text1</emphasis><emphasis level="reduced">text2</emphasis></speak>`,
		},
		{
			name:     "reduced",
			level:    emphasis.Level("custom"),
			expected: `<speak><emphasis level="custom">text1</emphasis><emphasis level="custom">text2</emphasis></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendEmphasis(test.level, "text1")
		b.AppendEmphasis(test.level, "text2")

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendParagraph(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendParagraph("text1")
	b.AppendParagraph("text2")

	actual := b.Build()
	expected := `<speak><p>text1</p><p>text2</p></speak>`
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestBuilder_AppendProsody(t *testing.T) {
	tests := []struct {
		name     string
		rate     prosody.Rate
		pitch    prosody.Pitch
		volume   prosody.Volume
		expected string
	}{
		{
			name:     "x-slow, x-low, & silent",
			rate:     prosody.RateXSlow,
			pitch:    prosody.PitchXLow,
			volume:   prosody.VolumeSilent,
			expected: `<speak><prosody rate="x-slow" pitch="x-low" volume="silent">text1</prosody><prosody rate="x-slow" pitch="x-low" volume="silent">text2</prosody></speak>`,
		},
		{
			name:     "slow, low, & x-soft",
			rate:     prosody.RateSlow,
			pitch:    prosody.PitchLow,
			volume:   prosody.VolumeXSoft,
			expected: `<speak><prosody rate="slow" pitch="low" volume="x-soft">text1</prosody><prosody rate="slow" pitch="low" volume="x-soft">text2</prosody></speak>`,
		},
		{
			name:     "medium, medium, & soft",
			rate:     prosody.RateMedium,
			pitch:    prosody.PitchMedium,
			volume:   prosody.VolumeSoft,
			expected: `<speak><prosody rate="medium" pitch="medium" volume="soft">text1</prosody><prosody rate="medium" pitch="medium" volume="soft">text2</prosody></speak>`,
		},
		{
			name:     "fast, high, & medium",
			rate:     prosody.RateFast,
			pitch:    prosody.PitchHigh,
			volume:   prosody.VolumeMedium,
			expected: `<speak><prosody rate="fast" pitch="high" volume="medium">text1</prosody><prosody rate="fast" pitch="high" volume="medium">text2</prosody></speak>`,
		},
		{
			name:     "x-fast, x-high, & loud",
			rate:     prosody.RateXFast,
			pitch:    prosody.PitchXHigh,
			volume:   prosody.VolumeLoud,
			expected: `<speak><prosody rate="x-fast" pitch="x-high" volume="loud">text1</prosody><prosody rate="x-fast" pitch="x-high" volume="loud">text2</prosody></speak>`,
		},
		{
			name:     "x-fast, x-high, & x-loud",
			rate:     prosody.RateXFast,
			pitch:    prosody.PitchXHigh,
			volume:   prosody.VolumeXLoud,
			expected: `<speak><prosody rate="x-fast" pitch="x-high" volume="x-loud">text1</prosody><prosody rate="x-fast" pitch="x-high" volume="x-loud">text2</prosody></speak>`,
		},
		{
			name:     "custom",
			rate:     prosody.Rate("custom rate"),
			pitch:    prosody.Pitch("custom pitch"),
			volume:   prosody.Volume("custom volume"),
			expected: `<speak><prosody rate="custom rate" pitch="custom pitch" volume="custom volume">text1</prosody><prosody rate="custom rate" pitch="custom pitch" volume="custom volume">text2</prosody></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendProsody(test.rate, test.pitch, test.volume, "text1")
		b.AppendProsody(test.rate, test.pitch, test.volume, "text2")

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendSentence(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendSentence("text1")
	b.AppendSentence("text2")

	actual := b.Build()
	expected := `<speak><s>text1</s><s>text2</s></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
	}
}

func TestBuilder_AppendSubstitution(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendSubstitution("alias1", "text1")
	b.AppendSubstitution("alias2", "text2")

	actual := b.Build()
	expected := `<speak><sub alias="alias1">text1</sub><sub alias="alias2">text2</sub></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
	}
}
