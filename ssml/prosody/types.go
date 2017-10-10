package prosody

type Rate string

const (
	RateXSlow  Rate = "x-slow"
	RateSlow   Rate = "slow"
	RateMedium Rate = "medium"
	RateFast   Rate = "fast"
	RateXFast  Rate = "x-fast"
)

type Pitch string

const (
	PitchXLow   Pitch = "x-low"
	PitchLow    Pitch = "low"
	PitchMedium Pitch = "medium"
	PitchHigh   Pitch = "high"
	PitchXHigh  Pitch = "x-high"
)

type Volume string

const (
	VolumeSilent Volume = "silent"
	VolumeXSoft  Volume = "x-soft"
	VolumeSoft   Volume = "soft"
	VolumeMedium Volume = "medium"
	VolumeLoud   Volume = "loud"
	VolumeXLoud  Volume = "x-loud"
)
