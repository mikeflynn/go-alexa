package ssml

import "bytes"

type builder struct {
	buffer *bytes.Buffer
}

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
