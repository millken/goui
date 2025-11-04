package component

import (
	"github.com/millken/go/thorvg"
	"github.com/millken/goui/color"
)

type Box struct {
	*BaseComponent
	x, y, w, h int

	bgColor color.Color
}

// 确保Box实现了Component接口
var _ Component = (*Box)(nil)

func NewBox(x, y, w, h int) *Box {
	box := &Box{
		BaseComponent: NewBaseComponent(),
		x:             x,
		y:             y,
		w:             w,
		h:             h,
		bgColor:       color.TRANSPARENT,
	}
	return box
}

func (b *Box) SetBackgroundColor(c color.Color) {
	b.bgColor = c
}

// Render 渲染Box组件
func (b *Box) Render(vg *thorvg.Canvas) {
	// 绘制背景色
	if b.bgColor != color.TRANSPARENT {
		background := thorvg.NewShape()
		if background != nil {
			background.AppendRect(float32(b.x), float32(b.y), float32(b.w), float32(b.h), 0, 0, false)
			background.SetFillColor(b.bgColor.RGBA())
			vg.Push(&background.Paint)
			background.Release()
		}
	}
	b.BaseComponent.Render(vg)

}
