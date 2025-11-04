package color

import "math"

// StaticColor is a static color implementation
type StaticColor struct {
	hue        float32
	saturation float32
	brightness float32
	alpha      float32
	argb       int32
}

// NewStaticColor creates a new static color
func NewStaticColor(hue, saturation, brightness, alpha float32) *StaticColor {
	c := &StaticColor{
		hue:        float32(math.Mod(float64(hue), 360.0)),
		saturation: clampf32(saturation, 0, 1),
		brightness: clampf32(brightness, 0, 1),
		alpha:      clampf32(alpha, 0, 1),
	}
	c.argb = HSBAToRGB(c.hue, c.saturation, c.brightness, c.alpha)
	return c
}

func (c *StaticColor) Hue() float32        { return c.hue }
func (c *StaticColor) Saturation() float32 { return c.saturation }
func (c *StaticColor) Brightness() float32 { return c.brightness }
func (c *StaticColor) Alpha() float32      { return c.alpha }
func (c *StaticColor) ARGB() int32         { return c.argb }
func (c *StaticColor) RGBA() int32         { return (c.argb << 8) | (c.argb>>24)&0xFF }
func (c *StaticColor) Transparent() bool   { return c.alpha == 0 }
func (c *StaticColor) Red() int32          { return (c.argb >> 16) & 0xFF }
func (c *StaticColor) Green() int32        { return (c.argb >> 8) & 0xFF }
func (c *StaticColor) Blue() int32         { return c.argb & 0xFF }
func (c *StaticColor) AlphaInt() int32     { return (c.argb >> 24) & 0xFF }

// Helper functions
func clampf32(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
