package ui

type FontWeight int32

// Possible values.
const (
	FontWeightInvisible FontWeight = iota * 100
	FontWeightThin
	FontWeightExtraLight
	FontWeightLight
	FontWeightRegular
	FontWeightMedium
	FontWeightSemiBold
	FontWeightBold
	FontWeightExtraBold
	FontWeightBlack
	FontWeightExtraBlack
)

type FontSpacing byte

const (
	FontSpacingUltraCondensed FontSpacing = iota + 1
	FontSpacingExtraCondensed
	FontSpacingCondensed
	FontSpacingSemiCondensed
	FontSpacingStandard
	FontSpacingSemiExpanded
	FontSpacingExpanded
	FontSpacingExtraExpanded
	FontSpacingUltraExpanded
)

type FontSlant byte
const (
	FontSlantNormal FontSlant = iota
	FontSlantItalic
	FontSlantOblique
)

// FontFaceDescriptor holds information necessary to construct a FontFace.
type FontFaceDescriptor struct {
	Family  string      `json:"family"`
	Weight  FontWeight  `json:"weight"`
	Spacing FontSpacing `json:"spacing"`
	Slant   FontSlant   `json:"slant"`
}
