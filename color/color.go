package color

type Color interface {
	RGBA() (r uint8, g uint8, b uint8, a uint8)
}

type color struct {
	r, g, b, a uint8
}

func NewColor(r, g, b, a uint8) Color {
	return &color{r: r, g: g, b: b, a: a}
}

func (c *color) RGBA() (r uint8, g uint8, b uint8, a uint8) {
	return c.r, c.g, c.b, c.a
}

var (
	BLACK       = NewColor(0, 0, 0, 255)
	WHITE       = NewColor(255, 255, 255, 255)
	RED         = NewColor(255, 0, 0, 255)
	GREEN       = NewColor(0, 255, 0, 255)
	BLUE        = NewColor(0, 0, 255, 255)
	TRANSPARENT = NewColor(0, 0, 0, 0)
)
