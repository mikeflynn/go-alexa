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
