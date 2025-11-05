package component

import (
	"github.com/millken/go/thorvg"
	gouicolor "github.com/millken/goui/color"
	"github.com/millken/goui/event"
	gouifont "github.com/millken/goui/font"
	gouiwindow "github.com/millken/goui/window"
)

type Button struct {
	*BaseComponent
	text string
	// basic styles are now sourced from BaseComponent.Style
	// kept only for transitional logic fields below
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
	// initialize default visual Style
	btn.UpdateStyle(func(s *Style) {
		bgHex := colorToHex(btn.defaultBG)
		fgHex := colorToHex(btn.savedFG)
		s.BackgroundColor = &bgHex
		s.Color = &fgHex
		s.BorderRadius = Px(4)
	})
	// 默认最小高度，避免高度为 0
	btn.SetMinHeight(36)
	// 默认自适应宽度并居中（在列容器交叉轴上不被拉伸）
	btn.SetAlignSelf(AlignCenter)

	// 默认 hover 行为：鼠标移入变色，移出恢复
	btn.OnHover(func(hover bool) {
		if btn.disabled {
			return
		}
		if hover {
			btn.SetBackgroundColor(btn.hoverBG)
			// 切换鼠标为手型
			if rc, ok := btn.Root().(RootComponent); ok && rc != nil {
				if w := rc.GetWindow(); w != nil {
					w.SetCursor(gouiwindow.CursorPointer)
				}
			}
		} else {
			// 恢复为进入 hover 前记录的默认背景色
			btn.SetBackgroundColor(btn.defaultBG)
			// 恢复默认光标
			if rc, ok := btn.Root().(RootComponent); ok && rc != nil {
				if w := rc.GetWindow(); w != nil {
					w.SetCursor(gouiwindow.CursorDefault)
				}
			}
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
	// update Style.BackgroundColor and dirty
	r, g, bl, a := c.RGBA()
	hex := colorToHex(c)
	b.UpdateStyle(func(s *Style) { s.BackgroundColor = &hex })
	// update stored defaults for hover lifecycle
	if !b.hovered {
		b.defaultBG = gouicolor.NewColor(r, g, bl, a)
	}
}

func (b *Button) SetTextColor(c gouicolor.Color) {
	hex := colorToHex(c)
	b.UpdateStyle(func(s *Style) { s.Color = &hex })
	// keep savedFG in sync for disable/enable restore
	if !b.disabled {
		b.savedFG = c
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
	b.UpdateStyle(func(s *Style) { s.BorderRadius = Px(r) })
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
		// apply disabled colors via Style; keep savedFG as previously set normal color
		bgHex := colorToHex(b.disabledBG)
		fgHex := colorToHex(b.disabledFG)
		b.UpdateStyle(func(st *Style) {
			st.BackgroundColor = &bgHex
			st.Color = &fgHex
		})
	} else {
		// restore defaults via Style
		bgHex := colorToHex(b.defaultBG)
		fgHex := colorToHex(b.savedFG)
		b.UpdateStyle(func(st *Style) {
			st.BackgroundColor = &bgHex
			st.Color = &fgHex
		})
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
		b.savedFG = gouicolor.NewColor(255, 255, 255, 255)
		b.defaultBG = gouicolor.NewColor(33, 150, 243, 255)
		b.hoverBG = gouicolor.NewColor(30, 136, 229, 255)
		b.activeBG = gouicolor.NewColor(25, 118, 210, 255)
	case VariantDanger:
		b.savedFG = gouicolor.NewColor(255, 255, 255, 255)
		b.defaultBG = gouicolor.NewColor(244, 67, 54, 255)
		b.hoverBG = gouicolor.NewColor(229, 57, 53, 255)
		b.activeBG = gouicolor.NewColor(211, 47, 47, 255)
	default: // Neutral
		b.savedFG = gouicolor.NewColor(30, 30, 30, 255)
		b.defaultBG = gouicolor.NewColor(230, 230, 230, 255)
		b.hoverBG = gouicolor.NewColor(210, 210, 210, 255)
		b.activeBG = gouicolor.NewColor(200, 200, 200, 255)
	}
	if !b.disabled {
		// push style immediately if not disabled
		bgHex := colorToHex(b.defaultBG)
		fgHex := colorToHex(b.savedFG)
		b.UpdateStyle(func(s *Style) {
			s.BackgroundColor = &bgHex
			s.Color = &fgHex
		})
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
	// 使用便捷方法以自动标记布局脏
	b.SetWidth(width)
}

// 根据字体和内边距计算推荐高度（多行暂不自动扩展）
func (b *Button) updateIntrinsicHeight() {
	tm := gouifont.MeasureText("Hg", gouifont.WithFontSize(b.fontSize)) // 用大写高度近似行高
	h := float32(tm.Height) + b.paddingY*2
	if h < 36 {
		h = 36
	}
	// 使用便捷方法以自动标记布局脏
	b.SetMinHeight(h)
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

	// 背景（从 Style 读取背景色和圆角）
	s := b.Style()
	if s.BackgroundColor != nil {
		if rr, gg, bb, aa, ok := gouicolor.ParseColor(*s.BackgroundColor); ok {
			bg := thorvg.NewShape()
			if bg != nil {
				var radius float32
				if s.BorderRadius != nil && s.BorderRadius.Unit == UnitPoint {
					radius = s.BorderRadius.Value
				}
				bg.AppendRect(x, y, l.Width, l.Height, radius, radius, false)
				bg.SetFillColor(rr, gg, bb, aa)
				vg.Push(&bg.Paint)
				bg.Release()
			}
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
			// 使用 ComputedStyle 的字体与颜色
			cs := b.ComputedStyle()
			size := float32(b.fontSize)
			if cs.FontSize != nil && *cs.FontSize > 0 {
				size = float32(*cs.FontSize)
			}
			text.SetSize(size)
			if cs.Color != nil {
				if cr, cg, cb, _, ok := gouicolor.ParseColor(*cs.Color); ok {
					text.SetColor(cr, cg, cb)
				}
			}
			if cs.Font != nil && *cs.Font != "" {
				text.SetFont(*cs.Font)
			} else {
				text.SetFont(b.fontFamily)
			}
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

// colorToHex converts a color.Color to #RRGGBBAA string
func colorToHex(c gouicolor.Color) string {
	r, g, b, a := c.RGBA()
	hexdigits := "0123456789ABCDEF"
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
