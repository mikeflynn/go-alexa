package ssml

import (
	"testing"
	"time"
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
	b, _ := NewBuilder()

	b.AppendAmazonEffect("effect1", "effect1text")
	b.AppendAmazonEffect("effect2", "effect2text")

	actual := b.Build()
	expected := `<speak><amazon:effect name="effect1">effect1text</amazon:effect><amazon:effect name="effect2">effect2text</amazon:effect></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
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
		strength BreakStrength
		duration time.Duration
		expected string
	}{
		{
			name:     "default",
			strength: StrengthDefault,
			duration: time.Second,
			expected: `<speak><break strength="medium" time="1000ms"/><break strength="medium" time="1000ms"/></speak>`,
		},
		{
			name:     "none",
			strength: StrengthNone,
			duration: time.Second / 2,
			expected: `<speak><break strength="none" time="500ms"/><break strength="none" time="500ms"/></speak>`,
		},
		{
			name:     "x-weak",
			strength: StrengthXWeak,
			duration: time.Second * 2,
			expected: `<speak><break strength="x-weak" time="2000ms"/><break strength="x-weak" time="2000ms"/></speak>`,
		},
		{
			name:     "weak",
			strength: StrengthWeak,
			duration: time.Second * 3,
			expected: `<speak><break strength="weak" time="3000ms"/><break strength="weak" time="3000ms"/></speak>`,
		},
		{
			name:     "medium",
			strength: StrengthMedium,
			duration: time.Second * 4,
			expected: `<speak><break strength="medium" time="4000ms"/><break strength="medium" time="4000ms"/></speak>`,
		},
		{
			name:     "strong",
			strength: StrengthStrong,
			duration: time.Second * 5,
			expected: `<speak><break strength="strong" time="5000ms"/><break strength="strong" time="5000ms"/></speak>`,
		},
		{
			name:     "x-strong",
			strength: StrengthXStrong,
			duration: time.Second * 6,
			expected: `<speak><break strength="x-strong" time="6000ms"/><break strength="x-strong" time="6000ms"/></speak>`,
		},
		{
			name:     "custom",
			strength: BreakStrength("custom"),
			duration: time.Second * 7,
			expected: `<speak><break strength="custom" time="7000ms"/><break strength="custom" time="7000ms"/></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendBreak(test.strength, test.duration)
		b.AppendBreak(test.strength, test.duration)

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendEmphasis(t *testing.T) {
	tests := []struct {
		name     string
		level    EmphasisLevel
		text     string
		expected string
	}{
		{
			name:     "default",
			level:    EmphasisDefault,
			text:     "test1",
			expected: `<speak><emphasis level="moderate">test1</emphasis><emphasis level="moderate">test1</emphasis></speak>`,
		},
		{
			name:     "strong",
			level:    EmphasisStrong,
			text:     "test2",
			expected: `<speak><emphasis level="strong">test2</emphasis><emphasis level="strong">test2</emphasis></speak>`,
		},
		{
			name:     "moderate",
			level:    EmphasisModerate,
			text:     "test3",
			expected: `<speak><emphasis level="moderate">test3</emphasis><emphasis level="moderate">test3</emphasis></speak>`,
		},
		{
			name:     "reduced",
			level:    EmphasisReduced,
			text:     "test4",
			expected: `<speak><emphasis level="reduced">test4</emphasis><emphasis level="reduced">test4</emphasis></speak>`,
		},
		{
			name:     "reduced",
			level:    EmphasisLevel("custom"),
			text:     "test5",
			expected: `<speak><emphasis level="custom">test5</emphasis><emphasis level="custom">test5</emphasis></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendEmphasis(test.level, test.text)
		b.AppendEmphasis(test.level, test.text)

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
	b, _ := NewBuilder()

	b.AppendProsody("rate1", "pitch1", "volume1", "text1")
	b.AppendProsody("rate2", "pitch2", "volume2", "text2")

	actual := b.Build()
	expected := `<speak><prosody rate="rate1" pitch="pitch1" volume="volume1">text1</prosody><prosody rate="rate2" pitch="pitch2" volume="volume2">text2</prosody></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
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
