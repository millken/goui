package component

import (
	"github.com/millken/go/thorvg"
	gouicolor "github.com/millken/goui/color"
	"github.com/millken/goui/event"
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
	// text & layout
	fontFamily string
	fontSize   float64
	paddingX   float32
	paddingY   float32
	// hover styles (optional)
	hoverBG gouicolor.Color
	// persistent default background to restore on hover leave
	defaultBG gouicolor.Color
	// active styles and state
	activeBG gouicolor.Color
	pressed  bool
	// disabled state
	disabled   bool
	disabledBG gouicolor.Color
	disabledFG gouicolor.Color
	savedFG    gouicolor.Color
}

func NewButton(id ...string) *Button {
	btn := &Button{
		BaseComponent: NewBaseComponent(id...),
		bg:            gouicolor.NewColor(230, 230, 230, 255),
		fg:            gouicolor.NewColor(30, 30, 30, 255),
		radius:        4,
		fontFamily:    "OPPO Sans 4.0",
		fontSize:      14.0,
		paddingX:      12,
		paddingY:      8,
		hoverBG:       gouicolor.NewColor(210, 210, 210, 255),
		defaultBG:     gouicolor.NewColor(230, 230, 230, 255),
		activeBG:      gouicolor.NewColor(200, 200, 200, 255),
		disabledBG:    gouicolor.NewColor(235, 235, 235, 255),
		disabledFG:    gouicolor.NewColor(160, 160, 160, 255),
		savedFG:       gouicolor.NewColor(30, 30, 30, 255),
	}
	// 默认最小高度，避免高度为 0
	btn.GetYogaNode().SetMinHeight(36)
	// 默认自适应宽度并居中（在列容器交叉轴上不被拉伸）
	btn.GetYogaNode().SetAlignSelf(yoga.AlignCenter)

	// 默认 hover 行为：鼠标移入变色，移出恢复
	btn.OnHover(func(hover bool) {
		if btn.disabled {
			return
		}
		if hover {
			btn.SetBackgroundColor(btn.hoverBG)
		} else {
			// 恢复为进入 hover 前记录的默认背景色
			btn.SetBackgroundColor(btn.defaultBG)
		}
	})

	// active（按下）行为：
	// - 当鼠标在按钮内按下（Down && Count==1）进入 pressed 状态，使用 activeBG
	// - 移动过程中若已 pressed，保持 inside 时 activeBG，移出时恢复 defaultBG
	// - 释放时（!Down && Count==1）结束 pressed，根据 inside 决定 hoverBG 或 defaultBG
	btn.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		me, ok := e.(*event.MouseEvent)
		if !ok {
			return false
		}
		if btn.disabled {
			return false
		}
		// 命中测试
		r := btn.Bounds()
		inside := me.X >= float32(r.Min.X) && me.X < float32(r.Max.X) && me.Y >= float32(r.Min.Y) && me.Y < float32(r.Max.Y)

		// 鼠标按下（开始 active）
		if me.Down && me.Count == 1 {
			if inside {
				btn.pressed = true
				btn.SetBackgroundColor(btn.activeBG)
				return false
			}
			return false
		}
		// 鼠标移动时（Count==0）：若处于 pressed，根据是否在内部切换样式
		if !me.Down && me.Count == 0 {
			if btn.pressed {
				if inside {
					btn.SetBackgroundColor(btn.activeBG)
				} else {
					btn.SetBackgroundColor(btn.defaultBG)
				}
				return false
			}
			return false
		}
		// 鼠标释放（结束 active）
		if !me.Down && me.Count == 1 {
			if btn.pressed {
				btn.pressed = false
				if inside {
					// 保持 hover 态颜色（视觉上更自然）
					btn.SetBackgroundColor(btn.hoverBG)
				} else {
					btn.SetBackgroundColor(btn.defaultBG)
				}
				return false
			}
			return false
		}
		return false
	}))
	return btn
}

func (b *Button) SetText(text string) {
	b.text = text
	b.updateIntrinsicWidth()
	b.updateIntrinsicHeight()
	b.MarkDirty(false) // 文本改变时标记为脏
}

func (b *Button) SetBackgroundColor(c gouicolor.Color) {
	if b.bg != c {
		b.bg = c
		// 非悬停状态下更新默认背景色，保证 hover 离开时恢复到“原始”颜色
		if !b.hovered {
			b.defaultBG = c
		}
		b.MarkDirty(false)
	}
}

func (b *Button) SetTextColor(c gouicolor.Color) {
	if b.fg != c {
		b.fg = c
		b.MarkDirty(false)
	}
}

// OnHover 为按钮注册悬停回调，便于外部自定义悬停样式
func (b *Button) OnHover(handler func(hover bool)) {
	b.BaseComponent.OnHover(handler)
}

// SetActiveBackgroundColor 自定义 active 态背景色
func (b *Button) SetActiveBackgroundColor(c gouicolor.Color) {
	b.activeBG = c
	if b.pressed {
		b.SetBackgroundColor(c)
	}
}

// SetPadding 设置按钮内容与边框的水平/垂直内边距
func (b *Button) SetPadding(padX, padY float32) {
	if b.paddingX == padX && b.paddingY == padY {
		return
	}
	b.paddingX, b.paddingY = padX, padY
	b.updateIntrinsicWidth()
	b.updateIntrinsicHeight()
	b.MarkDirty(false)
}

// SetCornerRadius 设置圆角半径
func (b *Button) SetCornerRadius(r float32) {
	if b.radius != r {
		b.radius = r
		b.MarkDirty(false)
	}
}

// SetTextSize 设置文本字号（影响测量与绘制）
func (b *Button) SetTextSize(size float64) {
	if b.fontSize != size {
		b.fontSize = size
		b.updateIntrinsicWidth()
		b.updateIntrinsicHeight()
		b.MarkDirty(false)
	}
}

// SetFontFamily 设置字体名称（用于绘制）
func (b *Button) SetFontFamily(name string) {
	if b.fontFamily != name {
		b.fontFamily = name
		b.updateIntrinsicWidth()
		b.updateIntrinsicHeight()
		b.MarkDirty(false)
	}
}

// SetDisabled 设置禁用状态
func (b *Button) SetDisabled(disabled bool) {
	if b.disabled == disabled {
		return
	}
	b.disabled = disabled
	if disabled {
		b.pressed = false
		b.savedFG = b.fg
		// 直接赋值以避免 defaultBG 被覆盖
		b.bg = b.disabledBG
		b.fg = b.disabledFG
		b.MarkDirty(false)
	} else {
		b.bg = b.defaultBG
		b.fg = b.savedFG
		b.MarkDirty(false)
	}
}

// ButtonVariant 预设风格
type ButtonVariant int

const (
	VariantNeutral ButtonVariant = iota
	VariantPrimary
	VariantDanger
)

// SetVariant 应用预设配色方案
func (b *Button) SetVariant(v ButtonVariant) {
	switch v {
	case VariantPrimary:
		b.fg = gouicolor.NewColor(255, 255, 255, 255)
		b.defaultBG = gouicolor.NewColor(33, 150, 243, 255)
		b.hoverBG = gouicolor.NewColor(30, 136, 229, 255)
		b.activeBG = gouicolor.NewColor(25, 118, 210, 255)
	case VariantDanger:
		b.fg = gouicolor.NewColor(255, 255, 255, 255)
		b.defaultBG = gouicolor.NewColor(244, 67, 54, 255)
		b.hoverBG = gouicolor.NewColor(229, 57, 53, 255)
		b.activeBG = gouicolor.NewColor(211, 47, 47, 255)
	default: // Neutral
		b.fg = gouicolor.NewColor(30, 30, 30, 255)
		b.defaultBG = gouicolor.NewColor(230, 230, 230, 255)
		b.hoverBG = gouicolor.NewColor(210, 210, 210, 255)
		b.activeBG = gouicolor.NewColor(200, 200, 200, 255)
	}
	if !b.disabled {
		b.bg = b.defaultBG
	}
	b.updateIntrinsicWidth()
	b.updateIntrinsicHeight()
	b.MarkDirty(false)
}

// 根据文本测量结果计算按钮的建议宽度（文本宽 + 水平内边距*2）
func (b *Button) updateIntrinsicWidth() {
	if b.text == "" {
		return
	}
	tm := gouifont.MeasureText(b.text, gouifont.WithFontSize(b.fontSize))
	width := float32(tm.Width) + b.paddingX*2
	if width < 48 {
		width = 48
	}
	b.GetYogaNode().SetWidth(width)
}

// 根据字体和内边距计算推荐高度（多行暂不自动扩展）
func (b *Button) updateIntrinsicHeight() {
	tm := gouifont.MeasureText("Hg", gouifont.WithFontSize(b.fontSize)) // 用大写高度近似行高
	h := float32(tm.Height) + b.paddingY*2
	if h < 36 {
		h = 36
	}
	b.GetYogaNode().SetMinHeight(h)
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
		tm := gouifont.MeasureText(b.text, gouifont.WithFontSize(b.fontSize))
		tx := x + (l.Width-float32(tm.Width))/2
		// 以文本的高度为准，简单将文本块在按钮内垂直居中
		ty := y + (l.Height-float32(tm.Height))/2

		text := thorvg.NewText()
		if text != nil {
			text.SetText(b.text)
			text.SetSize(float32(b.fontSize))
			cr, cg, cb, _ := b.fg.RGBA()
			text.SetColor(cr, cg, cb)
			text.SetFont(b.fontFamily)
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
