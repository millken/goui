package ui

import "gitee.com/millken/goui/internal/skia"

// FontMetrics flags
const (
	UnderlineThicknessIsValidFontMetricsFlag = 1 << iota
	UnderlinePositionIsValidFontMetricsFlag
	StrikeoutThicknessIsValidFontMetricsFlag
	StrikeoutPositionIsValidFontMetricsFlag
	BoundsInvalidFontMetricsFlag
)

// FontMetrics holds various metrics about a font.
type FontMetrics = skia.FontMetrics

type Font struct {
	// Size is the font size in points.
	Size float32
	// Family is the font family.
	Family string
	// Bold is whether the font is bold.
	Bold bool
	// Italic is whether the font is italic.
	Italic bool
	// Underline is whether the font is underlined.
	Underline bool
	// Strikeout is whether the font is strikeout.
	Strikeout bool
	// Hinting is the hinting level.
	//Hinting Hinting
}
