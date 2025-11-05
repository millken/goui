package component

import (
	"image"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui/color"
)

type Box struct {
	*BaseComponent
	x, y, w, h int
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
	}
	return box
}

func rgbToHex(r, g, b, a uint8) string {
	const hexdigits = "0123456789ABCDEF"
	out := make([]byte, 9)
	out[0] = '#'
	out[1] = hexdigits[r>>4]
	out[2] = hexdigits[r&0x0F]
	out[3] = hexdigits[g>>4]
	out[4] = hexdigits[g&0x0F]
	out[5] = hexdigits[b>>4]
	out[6] = hexdigits[b&0x0F]
	out[7] = hexdigits[a>>4]
	out[8] = hexdigits[a&0x0F]
	return string(out)
}

// Bounds 返回 Box 的实际边界（基于 x,y,w,h），用于命中测试
func (b *Box) Bounds() image.Rectangle {
	return image.Rect(b.x, b.y, b.x+b.w, b.y+b.h)
}

func (b *Box) SetBackgroundColor(c color.Color) {
	// 兼容旧 API：将颜色写入 Style.BackgroundColor
	r, g, bl, a := c.RGBA()
	hex := rgbToHex(r, g, bl, a)
	b.UpdateStyle(func(s *Style) { s.BackgroundColor = &hex })
}

// Paint 渲染Box组件
func (b *Box) Paint(vg *thorvg.Canvas) bool {
	// 首先渲染自己（如果需要）
	selfRendered := false
	if b.IsDirty() {
		// 从 Style 读取背景与圆角
		s := b.Style()
		if s.BackgroundColor != nil {
			if r, g, bl, a, ok := color.ParseColor(*s.BackgroundColor); ok {
				background := thorvg.NewShape()
				if background != nil {
					// 圆角（统一半径）
					var rr float32
					if s.BorderRadius != nil && s.BorderRadius.Unit == UnitPoint {
						rr = s.BorderRadius.Value
					}
					background.AppendRect(float32(b.x), float32(b.y), float32(b.w), float32(b.h), rr, rr, false)
					background.SetFillColor(r, g, bl, a)
					vg.Push(&background.Paint)
					background.Release()
				}
			}
		}
		selfRendered = true
	}

	// 然后渲染子组件
	childRendered := b.BaseComponent.Paint(vg)

	return selfRendered || childRendered
}
