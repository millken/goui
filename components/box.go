package components

import (
	"github.com/millken/go/thorvg"
)

// Box 是一个简单的矩形组件，用于占位和布局验证
type Box struct {
	Base
	R, G, B, A uint8
	Radius     float32
}

func NewBox(r, g, b, a uint8) *Box {
	return &Box{Base: NewBase(), R: r, G: g, B: b, A: a}
}

func (b *Box) Paint(canvas *thorvg.Canvas) {
	if !b.IsVisible() {
		return
	}
	l := b.GetComputedLayout()
	shape := thorvg.NewShape()
	if shape == nil {
		return
	}
	shape.AppendRect(l.Left, l.Top, l.Width, l.Height, b.Radius, b.Radius, false)
	shape.SetFillColor(b.R, b.G, b.B, b.A)
	canvas.Push(&shape.Paint)
	shape.Release()
}
