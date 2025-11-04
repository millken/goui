package component

import (
	"github.com/millken/go/thorvg"
	gouicolor "github.com/millken/goui/color"
	gouifont "github.com/millken/goui/font"
	"github.com/millken/yoga"
)

type Button struct {
	*BaseComponent
	text string
	// basic styles
	bg     gouicolor.Color
	fg     gouicolor.Color
	radius float32
}

func NewButton(id ...string) *Button {
	btn := &Button{
		BaseComponent: NewBaseComponent(id...),
		bg:            gouicolor.NewColor(230, 230, 230, 255),
		fg:            gouicolor.NewColor(30, 30, 30, 255),
		radius:        4,
	}
	// 默认最小高度，避免高度为 0
	btn.GetYogaNode().SetMinHeight(36)
	// 默认自适应宽度并居中（在列容器交叉轴上不被拉伸）
	btn.GetYogaNode().SetAlignSelf(yoga.AlignCenter)
	return btn
}

func (b *Button) SetText(text string) {
	b.text = text
	b.updateIntrinsicWidth()
	b.MarkDirty(false) // 文本改变时标记为脏
}

func (b *Button) SetBackgroundColor(c gouicolor.Color) {
	if b.bg != c {
		b.bg = c
		b.MarkDirty(false)
	}
}

func (b *Button) SetTextColor(c gouicolor.Color) {
	if b.fg != c {
		b.fg = c
		b.MarkDirty(false)
	}
}

// 根据文本测量结果计算按钮的建议宽度（文本宽 + 水平内边距*2）
func (b *Button) updateIntrinsicWidth() {
	if b.text == "" {
		return
	}
	tm := gouifont.MeasureText(b.text)
	paddingX := float32(12)
	width := float32(tm.Width) + paddingX*2
	if width < 48 {
		width = 48
	}
	b.GetYogaNode().SetWidth(width)
}

// Paint 绘制按钮（背景 + 文本），仅在脏时绘制
func (b *Button) Paint(vg *thorvg.Canvas) bool {
	if !b.IsDirty() {
		// 递归子组件绘制结果
		return b.BaseComponent.Paint(vg)
	}

	// 几何
	l := b.ComputedLayout()
	r := b.Bounds()
	x, y := float32(r.Min.X), float32(r.Min.Y)

	// 背景
	if b.bg != gouicolor.TRANSPARENT {
		bg := thorvg.NewShape()
		if bg != nil {
			bg.AppendRect(x, y, l.Width, l.Height, b.radius, b.radius, false)
			rr, gg, bb, aa := b.bg.RGBA()
			bg.SetFillColor(rr, gg, bb, aa)
			vg.Push(&bg.Paint)
			bg.Release()
		}
	}

	// 文本
	if b.text != "" {
		tm := gouifont.MeasureText(b.text)
		tx := x + (l.Width-float32(tm.Width))/2
		// 以文本的高度为准，简单将文本块在按钮内垂直居中
		ty := y + (l.Height-float32(tm.Height))/2

		text := thorvg.NewText()
		if text != nil {
			text.SetText(b.text)
			text.SetSize(14.0)
			cr, cg, cb, _ := b.fg.RGBA()
			text.SetColor(cr, cg, cb)
			text.SetFont("OPPO Sans 4.0")
			text.Translate(tx, ty)
			vg.Push(&text.Paint)
			text.Release()
		}
	}

	b.ClearDirty()

	// 再绘制子组件（如果有）
	childRendered := b.BaseComponent.Paint(vg)
	return true || childRendered
}
