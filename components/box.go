package components

import (
	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
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
	// 按 Style 绘制（优先），否则使用构造色
	radius := b.Radius
	if st, ok := any(b).(interface{ GetStyle() goui.Style }); ok {
		s := st.GetStyle()
		if s.BorderRadius != nil && s.BorderRadius.Unit == goui.UnitPoint {
			radius = s.BorderRadius.Value
		}
		// 解析背景色
		bgR, bgG, bgB, bgA := b.R, b.G, b.B, b.A
		if s.BackgroundColor != nil {
			if r, gc, bc, ac, okc := goui.ParseColor(*s.BackgroundColor); okc {
				bgR, bgG, bgB, bgA = r, gc, bc, ac
			}
		}

		// 计算统一边框宽度（取最大 px 值）
		var bw float32
		if s.Border != nil {
			dims := []*goui.Dimension{s.Border.Left, s.Border.Right, s.Border.Top, s.Border.Bottom, s.Border.Horizontal, s.Border.Vertical, s.Border.All}
			for _, d := range dims {
				if d == nil {
					continue
				}
				if d.Unit == goui.UnitPoint && d.Value > bw {
					bw = d.Value
				}
			}
		}

		// 如果有边框颜色且宽度>0，则先画外框为边框色，再画内框为背景色
		if s.BorderColor != nil && bw > 0 {
			if r, gc, bc, ac, okc := goui.ParseColor(*s.BorderColor); okc {
				outer := thorvg.NewShape()
				if outer != nil {
					outer.AppendRect(l.Left, l.Top, l.Width, l.Height, radius, radius, false)
					outer.SetFillColor(r, gc, bc, ac)
					canvas.Push(&outer.Paint)
					outer.Release()
				}
				// 内矩形
				innerLeft := l.Left + bw
				innerTop := l.Top + bw
				innerW := l.Width - 2*bw
				innerH := l.Height - 2*bw
				if innerW > 0 && innerH > 0 {
					inner := thorvg.NewShape()
					if inner != nil {
						innerRadius := radius - bw
						if innerRadius < 0 {
							innerRadius = 0
						}
						inner.AppendRect(innerLeft, innerTop, innerW, innerH, innerRadius, innerRadius, false)
						inner.SetFillColor(bgR, bgG, bgB, bgA)
						canvas.Push(&inner.Paint)
						inner.Release()
					}
				}
				return
			}
		}
		// 无边框：单个背景色矩形
		shape := thorvg.NewShape()
		if shape == nil {
			return
		}
		shape.AppendRect(l.Left, l.Top, l.Width, l.Height, radius, radius, false)
		shape.SetFillColor(bgR, bgG, bgB, bgA)
		canvas.Push(&shape.Paint)
		shape.Release()
		return
	}
	shape := thorvg.NewShape()
	if shape == nil {
		return
	}
	shape.AppendRect(l.Left, l.Top, l.Width, l.Height, radius, radius, false)
	shape.SetFillColor(b.R, b.G, b.B, b.A)
	canvas.Push(&shape.Paint)
	shape.Release()
}
