package component

import (
	"image"

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

// Bounds 返回 Box 的实际边界（基于 x,y,w,h），用于命中测试
func (b *Box) Bounds() image.Rectangle {
	return image.Rect(b.x, b.y, b.x+b.w, b.y+b.h)
}

func (b *Box) SetBackgroundColor(c color.Color) {
	if b.bgColor != c {
		b.bgColor = c
		b.MarkDirty(false) // 颜色改变时标记为脏
	}
}

// Paint 渲染Box组件
func (b *Box) Paint(vg *thorvg.Canvas) bool {
	// 首先渲染自己（如果需要）
	selfRendered := false
	if b.IsDirty() {
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
		selfRendered = true
	}

	// 然后渲染子组件
	childRendered := b.BaseComponent.Paint(vg)

	return selfRendered || childRendered
}
