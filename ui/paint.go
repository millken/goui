package ui

import (
	"runtime"

	"gitee.com/millken/goui/internal/skia"
)

// PaintStyle holds the type of painting to do.
type PaintStyle byte

// Possible values for PaintStyle.
const (
	Fill PaintStyle = iota
	Stroke
	StrokeAndFill
)

// StrokeCap holds the style for rendering the endpoint of a stroked line.
type StrokeCap byte

// Possible values for StrokeCap.
const (
	ButtCap StrokeCap = iota
	RoundCap
	SquareCap
)

// StrokeJoin holds the method for drawing the junction between connected line segments.
type StrokeJoin byte

// Possible values for StrokeJoin.
const (
	MiterJoin StrokeJoin = iota
	RoundJoin
	BevelJoin
)

// BlendMode holds the mode used for blending pixels.
type BlendMode byte

// Possible values for BlendMode.
const (
	ClearBlendMode BlendMode = iota
	SrcBlendMode
	DstBlendMode
	SrcOverBlendMode
	DstOverBlendMode
	SrcInBlendMode
	DstInBlendMode
	SrcOutBlendMode
	DstOutBlendMode
	SrcAtopBlendMode
	DstAtopBlendMode
	XorBlendMode
	PlusBlendMode
	ModulateBlendMode
	ScreenBlendMode
	OverlayBlendMode
	DarkenBlendMode
	LightenBlendMode
	ColorDodgeBlendMode
	ColorBurnBlendMode
	HardLightBlendMode
	SoftLightBlendMode
	DifferenceBlendMode
	ExclusionBlendMode
	MultiplyBlendMode
	HueBlendMode
	SaturationBlendMode
	ColorBlendMode
	LuminosityBlendMode
)

// FilterQuality holds the image filtering level. Lower settings draw faster, while higher settings look better when the
// image is scaled.
type FilterQuality byte

// Possible values for FilterQuality.
const (
	NoQuality FilterQuality = iota
	LowQuality
	MediumQuality
	HighQuality
)

type Paint struct {
	paint skia.Paint
}

func newPaint(paint skia.Paint) *Paint {
	p := &Paint{paint: paint}
	runtime.SetFinalizer(p, func(obj *Paint) {
		skia.PaintDelete(obj.paint)
	})
	return p
}

// NewPaint creates a new Paint.
func NewPaint() *Paint {
	p := newPaint(skia.PaintNew())
	p.SetAntialias(true)
	return p
}

// Clone the Paint.
func (p *Paint) Clone() *Paint {
	return newPaint(skia.PaintClone(p.paint))
}

// Reset the Paint back to its default state.
func (p *Paint) Reset() {
	skia.PaintReset(p.paint)
}

// Antialias returns true if pixels on the active edges of a path may be drawn with partial transparency.
func (p *Paint) Antialias() bool {
	return skia.PaintIsAntialias(p.paint)
}

// SetAntialias sets whether pixels on the active edges of a path may be drawn with partial transparency.
func (p *Paint) SetAntialias(enabled bool) {
	skia.PaintSetAntialias(p.paint, enabled)
}

// Color returns the current color.
func (p *Paint) Color() Color {
	return Color(skia.PaintGetColor(p.paint))
}

// SetColor sets the color.
func (p *Paint) SetColor(color Color) {
	skia.PaintSetColor(p.paint, skia.Color(color))
}

// Style returns the current PaintStyle.
func (p *Paint) Style() PaintStyle {
	return PaintStyle(skia.PaintGetStyle(p.paint))
}

// SetStyle sets the PaintStyle.
func (p *Paint) SetStyle(style PaintStyle) {
	skia.PaintSetStyle(p.paint, skia.PaintStyle(style))
}

func (p *Paint) paintOrNil() skia.Paint {
	if p == nil {
		return nil
	}
	return p.paint
}
