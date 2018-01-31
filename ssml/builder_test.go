package ssml

import (
	"sync"
	"testing"
	"time"

	"github.com/mikeflynn/go-alexa/ssml/amazoneffect"
	"github.com/mikeflynn/go-alexa/ssml/emphasis"
	"github.com/mikeflynn/go-alexa/ssml/pause"
	"github.com/mikeflynn/go-alexa/ssml/prosody"
	"github.com/pkg/errors"
)

func TestNewBuilder_ReturnsEmptySSML(t *testing.T) {
	b, err := NewBuilder()

	if err != nil {
		t.Fatalf("new error mismatch: got: %+v, wanted: nil", err)
	}

	got, errs := b.Build()
	if !errsEqual(nil, errs) {
		t.Errorf("build error mismatch: got: %+v, wanted: nil", errs)
	}

	wanted := "<speak></speak>"
	if got != wanted {
		t.Errorf("output mismatch: got: %s, wanted: %s", got, wanted)
	}
}

func TestBuilder_AppendPlainSpeech(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendPlainSpeech("hello ").AppendPlainSpeech("world")

	got, errs := b.Build()
	if !errsEqual(nil, errs) {
		t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
	}

	wanted := "<speak>hello world</speak>"
	if got != wanted {
		t.Errorf("output mismatch: got: %s, wanted: %s", got, wanted)
	}
}

func TestBuilder_AppendAmazonEffect(t *testing.T) {
	tests := []struct {
		name   string
		effect amazoneffect.Effect
		wanted string
	}{
		{
			name:   "whispered",
			effect: amazoneffect.Whispered,
			wanted: `<speak><amazon:effect name="whispered">text1</amazon:effect><amazon:effect name="whispered">text2</amazon:effect></speak>`,
		},
		{
			name:   "custom",
			effect: amazoneffect.Effect("custom"),
			wanted: `<speak><amazon:effect name="custom">text1</amazon:effect><amazon:effect name="custom">text2</amazon:effect></speak>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBuilder()

			b.AppendAmazonEffect(test.effect, "text1").AppendAmazonEffect(test.effect, "text2")

			got, errs := b.Build()
			if !errsEqual(nil, errs) {
				t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
			}

			if got != test.wanted {
				t.Errorf("output mismatch: got: %s, wanted: %s", got, test.wanted)
			}
		})
	}
}

func TestBuilder_AppendAudio(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		errs   []error
		wanted string
	}{
		{
			name:   "happyPath",
			src:    "https://domain.tld",
			errs:   nil,
			wanted: `<speak><audio src="https://domain.tld"/><audio src="https://domain.tld"/></speak>`,
		},
		{
			name: "nonHTTPSUrl",
			src:  "http://domain.tld",
			errs: []error{
				errors.New("unsupported URL scheme type: must be https"),
				errors.New("unsupported URL scheme type: must be https"),
			},
			wanted: `<speak></speak>`,
		},
		{
			name: "badUrl",
			src:  "%notarealurl",
			errs: []error{
				errors.New(`failed to parse src into a valid URL: parse %notarealurl: invalid URL escape "%no"`),
				errors.New(`failed to parse src into a valid URL: parse %notarealurl: invalid URL escape "%no"`),
			},
			wanted: `<speak></speak>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBuilder()

			b.AppendAudio(test.src).AppendAudio(test.src)

			got, errs := b.Build()
			if !errsEqual(errs, test.errs) {
				t.Errorf("error mismatch: got: %+v, wanted: %+v", errs, test.errs)
			}

			if got != test.wanted {
				t.Errorf("output mismatch: got: %s, wanted: %s", got, test.wanted)
			}
		})
	}
}

func TestBuilder_AppendBreak(t *testing.T) {
	tests := []struct {
		name   string
		param  interface{}
		errs   []error
		wanted string
	}{
		{
			name:   "default",
			param:  pause.Default,
			errs:   nil,
			wanted: `<speak><break strength="medium"/><break strength="medium"/></speak>`,
		},
		{
			name:   "none",
			param:  pause.None,
			errs:   nil,
			wanted: `<speak><break strength="none"/><break strength="none"/></speak>`,
		},
		{
			name:   "x-weak",
			param:  pause.XWeak,
			errs:   nil,
			wanted: `<speak><break strength="x-weak"/><break strength="x-weak"/></speak>`,
		},
		{
			name:   "weak",
			param:  pause.Weak,
			errs:   nil,
			wanted: `<speak><break strength="weak"/><break strength="weak"/></speak>`,
		},
		{
			name:   "medium",
			param:  pause.Medium,
			errs:   nil,
			wanted: `<speak><break strength="medium"/><break strength="medium"/></speak>`,
		},
		{
			name:   "strong",
			param:  pause.Strong,
			errs:   nil,
			wanted: `<speak><break strength="strong"/><break strength="strong"/></speak>`,
		},
		{
			name:   "x-strong",
			param:  pause.XStrong,
			errs:   nil,
			wanted: `<speak><break strength="x-strong"/><break strength="x-strong"/></speak>`,
		},
		{
			name:   "custom",
			param:  pause.Strength("custom"),
			errs:   nil,
			wanted: `<speak><break strength="custom"/><break strength="custom"/></speak>`,
		},
		{
			name:   "time",
			param:  time.Second,
			errs:   nil,
			wanted: `<speak><break time="1000ms"/><break time="1000ms"/></speak>`,
		},
		{
			name:  "invalidType",
			param: 4,
			errs: []error{
				errors.New("unsupported parameter type: must be either pause.Strength or time.Duration"),
				errors.New("unsupported parameter type: must be either pause.Strength or time.Duration"),
			},
			wanted: `<speak></speak>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBuilder()

			b.AppendBreak(test.param).AppendBreak(test.param)

			got, errs := b.Build()
			if !errsEqual(errs, test.errs) {
				t.Errorf("error mismatch: got: %+v, wanted: %+v", errs, test.errs)
			}

			if got != test.wanted {
				t.Errorf("output mismatch: got: %s, wanted: %s", got, test.wanted)
			}
		})
	}
}

func TestBuilder_AppendEmphasis(t *testing.T) {
	tests := []struct {
		name   string
		level  emphasis.Level
		text   string
		wanted string
	}{
		{
			name:   "default",
			level:  emphasis.Default,
			wanted: `<speak><emphasis level="moderate">text1</emphasis><emphasis level="moderate">text2</emphasis></speak>`,
		},
		{
			name:   "strong",
			level:  emphasis.Strong,
			wanted: `<speak><emphasis level="strong">text1</emphasis><emphasis level="strong">text2</emphasis></speak>`,
		},
		{
			name:   "moderate",
			level:  emphasis.Moderate,
			wanted: `<speak><emphasis level="moderate">text1</emphasis><emphasis level="moderate">text2</emphasis></speak>`,
		},
		{
			name:   "reduced",
			level:  emphasis.Reduced,
			wanted: `<speak><emphasis level="reduced">text1</emphasis><emphasis level="reduced">text2</emphasis></speak>`,
		},
		{
			name:   "reduced",
			level:  emphasis.Level("custom"),
			wanted: `<speak><emphasis level="custom">text1</emphasis><emphasis level="custom">text2</emphasis></speak>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBuilder()

			b.AppendEmphasis(test.level, "text1").AppendEmphasis(test.level, "text2")

			got, errs := b.Build()
			if !errsEqual(nil, errs) {
				t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
			}

			if got != test.wanted {
				t.Errorf("output mismatch: got: %s, wanted: %s", got, test.wanted)
			}
		})
	}
}

func TestBuilder_AppendParagraph(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendParagraph("text1").AppendParagraph("text2")

	got, errs := b.Build()
	if !errsEqual(nil, errs) {
		t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
	}

	wanted := `<speak><p>text1</p><p>text2</p></speak>`
	if got != wanted {
		t.Errorf("got: %s, wanted: %s", got, wanted)
	}
}

func TestBuilder_AppendProsody(t *testing.T) {
	tests := []struct {
		name   string
		rate   interface{}
		pitch  interface{}
		volume interface{}
		errs   []error
		wanted string
	}{
		{
			name:   "x-slow, x-low, & silent",
			rate:   prosody.RateXSlow,
			pitch:  prosody.PitchXLow,
			volume: prosody.VolumeSilent,
			errs:   nil,
			wanted: `<speak><prosody rate="x-slow" pitch="x-low" volume="silent">text1</prosody><prosody rate="x-slow" pitch="x-low" volume="silent">text2</prosody></speak>`,
		},
		{
			name:   "slow, low, & x-soft",
			rate:   prosody.RateSlow,
			pitch:  prosody.PitchLow,
			volume: prosody.VolumeXSoft,
			errs:   nil,
			wanted: `<speak><prosody rate="slow" pitch="low" volume="x-soft">text1</prosody><prosody rate="slow" pitch="low" volume="x-soft">text2</prosody></speak>`,
		},
		{
			name:   "medium, medium, & soft",
			rate:   prosody.RateMedium,
			pitch:  prosody.PitchMedium,
			volume: prosody.VolumeSoft,
			errs:   nil,
			wanted: `<speak><prosody rate="medium" pitch="medium" volume="soft">text1</prosody><prosody rate="medium" pitch="medium" volume="soft">text2</prosody></speak>`,
		},
		{
			name:   "fast, high, & medium",
			rate:   prosody.RateFast,
			pitch:  prosody.PitchHigh,
			volume: prosody.VolumeMedium,
			errs:   nil,
			wanted: `<speak><prosody rate="fast" pitch="high" volume="medium">text1</prosody><prosody rate="fast" pitch="high" volume="medium">text2</prosody></speak>`,
		},
		{
			name:   "x-fast, x-high, & loud",
			rate:   prosody.RateXFast,
			pitch:  prosody.PitchXHigh,
			volume: prosody.VolumeLoud,
			errs:   nil,
			wanted: `<speak><prosody rate="x-fast" pitch="x-high" volume="loud">text1</prosody><prosody rate="x-fast" pitch="x-high" volume="loud">text2</prosody></speak>`,
		},
		{
			name:   "x-fast, x-high, & x-loud",
			rate:   prosody.RateXFast,
			pitch:  prosody.PitchXHigh,
			volume: prosody.VolumeXLoud,
			errs:   nil,
			wanted: `<speak><prosody rate="x-fast" pitch="x-high" volume="x-loud">text1</prosody><prosody rate="x-fast" pitch="x-high" volume="x-loud">text2</prosody></speak>`,
		},
		{
			name:   "custom",
			rate:   prosody.Rate("custom rate"),
			pitch:  prosody.Pitch("custom pitch"),
			volume: prosody.Volume("custom volume"),
			errs:   nil,
			wanted: `<speak><prosody rate="custom rate" pitch="custom pitch" volume="custom volume">text1</prosody><prosody rate="custom rate" pitch="custom pitch" volume="custom volume">text2</prosody></speak>`,
		},
		{
			name:   "onlyRate",
			rate:   prosody.RateXSlow,
			pitch:  nil,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody rate="x-slow">text1</prosody><prosody rate="x-slow">text2</prosody></speak>`,
		},
		{
			name:   "onlyPitch",
			rate:   nil,
			pitch:  prosody.PitchXLow,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody pitch="x-low">text1</prosody><prosody pitch="x-low">text2</prosody></speak>`,
		},
		{
			name:   "onlyVolume",
			rate:   nil,
			pitch:  nil,
			volume: prosody.VolumeSilent,
			errs:   nil,
			wanted: `<speak><prosody volume="silent">text1</prosody><prosody volume="silent">text2</prosody></speak>`,
		},
		{
			name:   "percentageRate",
			rate:   110,
			pitch:  nil,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody rate="110%">text1</prosody><prosody rate="110%">text2</prosody></speak>`,
		},
		{
			name:   "positivePercentagePitch",
			rate:   nil,
			pitch:  10,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody pitch="+10%">text1</prosody><prosody pitch="+10%">text2</prosody></speak>`,
		},
		{
			name:   "negativePercentagePitch",
			rate:   nil,
			pitch:  -10,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody pitch="-10%">text1</prosody><prosody pitch="-10%">text2</prosody></speak>`,
		},
		{
			name:   "positiveDbVolume",
			rate:   nil,
			pitch:  nil,
			volume: 4,
			errs:   nil,
			wanted: `<speak><prosody volume="+4dB">text1</prosody><prosody volume="+4dB">text2</prosody></speak>`,
		},
		{
			name:   "negativeDbVolume",
			rate:   nil,
			pitch:  nil,
			volume: -4,
			errs:   nil,
			wanted: `<speak><prosody volume="-4dB">text1</prosody><prosody volume="-4dB">text2</prosody></speak>`,
		},
		{
			name:   "allNil",
			rate:   nil,
			pitch:  nil,
			volume: nil,
			errs:   nil,
			wanted: `<speak><prosody>text1</prosody><prosody>text2</prosody></speak>`,
		},
		{
			name:   "invalidRateType",
			rate:   true,
			pitch:  nil,
			volume: nil,
			errs: []error{
				errors.New("unsupported rate type: must be either prosody.Rate or int"),
				errors.New("unsupported rate type: must be either prosody.Rate or int"),
			},
			wanted: `<speak></speak>`,
		},
		{
			name:   "invalidPitchType",
			rate:   nil,
			pitch:  true,
			volume: nil,
			errs: []error{
				errors.New("unsupported pitch type: must be either prosody.Pitch or int"),
				errors.New("unsupported pitch type: must be either prosody.Pitch or int"),
			},
			wanted: `<speak></speak>`,
		},
		{
			name:   "invalidVolumeType",
			rate:   nil,
			pitch:  nil,
			volume: true,
			errs: []error{
				errors.New("unsupported volume type: must be either prosody.Volume or int"),
				errors.New("unsupported volume type: must be either prosody.Volume or int"),
			},
			wanted: `<speak></speak>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBuilder()

			b.AppendProsody(test.rate, test.pitch, test.volume, "text1").AppendProsody(test.rate, test.pitch, test.volume, "text2")

			got, errs := b.Build()
			if !errsEqual(errs, test.errs) {
				t.Errorf("error mismatch: got: %+v, wanted: %+v", errs, test.errs)
			}

			if got != test.wanted {
				t.Errorf("output mismatch: got: %s, wanted: %s", got, test.wanted)
			}
		})
	}
}

func TestBuilder_AppendSentence(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendSentence("text1")
	b.AppendSentence("text2")

	got, errs := b.Build()
	if !errsEqual(nil, errs) {
		t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
	}

	wanted := `<speak><s>text1</s><s>text2</s></speak>`
	if got != wanted {
		t.Errorf("output mismatch: got: %s, wanted: %s", got, wanted)
	}
}

func TestBuilder_AppendSubstitution(t *testing.T) {
	b, _ := NewBuilder()

	b.AppendSubstitution("alias1", "text1")
	b.AppendSubstitution("alias2", "text2")

	got, errs := b.Build()
	if !errsEqual(nil, errs) {
		t.Errorf("error mismatch: got: %+v, wanted: nil", errs)
	}

	wanted := `<speak><sub alias="alias1">text1</sub><sub alias="alias2">text2</sub></speak>`
	if got != wanted {
		t.Errorf("output mismatch: got: %s, wanted: %s", got, wanted)
	}
}

func TestAppendErrorConcurrency(t *testing.T) {
	var (
		numOfWorkers         = 100
		numOfErrorsPerWorker = 100
	)

	b, _ := NewBuilder()
	wg := sync.WaitGroup{}

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < numOfErrorsPerWorker; i++ {
				b.appendError(errors.New("error"))
			}
			wg.Done()
		}()
	}

	wg.Wait()

	_, errs := b.Build()
	got := len(errs)
	wanted := numOfWorkers * numOfErrorsPerWorker
	if got != wanted {
		t.Errorf("error count mismatch: got: %d, wanted: %d", got, wanted)
	}
}

func errsEqual(errs1, errs2 []error) bool {
	if len(errs1) != len(errs2) {
		return false
	}
	for i, _ := range errs1 {
		if errs1[i].Error() != errs2[i].Error() {
			return false
		}
	}
	return true
}
