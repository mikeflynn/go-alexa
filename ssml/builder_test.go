package ssml

import (
	"testing"
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
		strength string
		time     string
		expected string
	}{
		{
			name:     "blank_input",
			strength: "",
			time:     "",
			expected: `<speak><break strength="medium" time=""/><break strength="medium" time=""/></speak>`,
		},
	}

	for _, test := range tests {
		b, _ := NewBuilder()

		b.AppendBreak(test.strength, test.time)
		b.AppendBreak(test.strength, test.time)

		actual := b.Build()
		if actual != test.expected {
			t.Errorf("%s: output mismatch: expected %s, got %s", test.name, test.expected, actual)
		}
	}
}

func TestBuilder_AppendEmphasis(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendEmphasis("level1", "text1")
	b.AppendEmphasis("level2", "text2")

	actual := b.Build()
	expected := `<speak><emphasis level="level1">text1</emphasis><emphasis level="level2">text2</emphasis></speak>`
	if actual != expected {
		t.Errorf("output mismatch: expected %s, got %s", expected, actual)
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
