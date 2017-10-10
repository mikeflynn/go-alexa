package ssml

import "bytes"

type builder struct {
	buffer *bytes.Buffer
}

type BreakStrength string

const (
	StrengthDefault BreakStrength = "medium"
	StrengthNone    BreakStrength = "none"
	StrengthXWeak   BreakStrength = "x-weak"
	StrengthWeak    BreakStrength = "weak"
	StrengthMedium  BreakStrength = "medium"
	StrengthStrong  BreakStrength = "strong"
	StrengthXStrong BreakStrength = "x-strong"
)

type EmphasisLevel string

const (
	EmphasisDefault  EmphasisLevel = "moderate"
	EmphasisStrong   EmphasisLevel = "strong"
	EmphasisModerate EmphasisLevel = "moderate"
	EmphasisReduced  EmphasisLevel = "reduced"
)

type ProsodyRate string

const (
	RateXSlow  ProsodyRate = "x-slow"
	RateSlow   ProsodyRate = "slow"
	RateMedium ProsodyRate = "medium"
	RateFast   ProsodyRate = "fast"
	RateXFast  ProsodyRate = "x-fast"
)

type ProsodyPitch string

const (
	PitchXLow   ProsodyPitch = "x-low"
	PitchLow    ProsodyPitch = "low"
	PitchMedium ProsodyPitch = "medium"
	PitchHigh   ProsodyPitch = "high"
	PitchXHigh  ProsodyPitch = "x-high"
)

type ProsodyVolume string

const (
	VolumeSilent ProsodyVolume = "silent"
	VolumeXSoft  ProsodyVolume = "x-soft"
	VolumeSoft   ProsodyVolume = "soft"
	VolumeMedium ProsodyVolume = "medium"
	VolumeLoud   ProsodyVolume = "loud"
	VolumeXLoud  ProsodyVolume = "x-loud"
)
