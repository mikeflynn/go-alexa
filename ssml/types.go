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
