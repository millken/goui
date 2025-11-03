package components

import (
	"image"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/yoga"
)

// ButtonState 按钮状态
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHovered
	ButtonStatePressed
	ButtonStateDisabled
)

// Button 按钮组件（V2版本）
type Button struct {
	*goui.BaseComponent
	text      string
	state     ButtonState
	bgColor   *Color
	textColor *Color
	border    *BorderStyle
	padding   *Insets
}

// 确保Button实现了Component接口
var _ goui.Component = (*Button)(nil)

// Color 按钮颜色定义（避免与主包的Color冲突）
type Color struct {
	R, G, B, A uint8
}

// BorderStyle 边框样式
type BorderStyle struct {
	Width float32
	Color Color
}

// Insets 内边距
type Insets struct {
	Top, Right, Bottom, Left float32
}

// NewButton 创建按钮
func NewButton(id, text string) *Button {
	base := goui.NewBaseComponent(id)
	button := &Button{
		BaseComponent: base,
		text:             text,
		state:            ButtonStateNormal,
		bgColor:          &Color{R: 59, G: 130, B: 246, A: 255}, // 蓝色
		textColor:        &Color{R: 255, G: 255, B: 255, A: 255}, // 白色
		border:           &BorderStyle{Width: 1, Color: Color{R: 37, G: 99, B: 235, A: 255}},
		padding:          &Insets{Top: 8, Right: 16, Bottom: 8, Left: 16},
	}

	// 配置Yoga节点
	button.configureYogaNode()

	// 设置事件处理器
	button.setupEventHandlers()

	return button
}

// configureYogaNode 配置Yoga节点（内部方法）
func (b *Button) configureYogaNode() {
	node := b.GetYogaNode()
	if node == nil {
		return
	}

	// 设置默认样式
	node.SetWidth(120)  // 默认宽度
	node.SetHeight(40)  // 默认高度
	node.SetJustifyContent(yoga.JustifyCenter)
	node.SetAlignItems(yoga.AlignCenter)
}

// setupEventHandlers 设置事件处理器（使用新的精简API）
func (b *Button) setupEventHandlers() {
	b.OnHover(func(hovered bool) {
		if !b.IsEnabled() {
			return
		}

		oldState := b.state
		if hovered {
			b.state = ButtonStateHovered
		} else {
			b.state = ButtonStateNormal
		}

		if oldState != b.state {
			b.updateColors()
			b.InvalidateLayout() // 触发重新渲染
		}
	})
}

// SetText 设置按钮文本
func (b *Button) SetText(text string) {
	if b.text != text {
		b.text = text
		b.InvalidateLayout()
	}
}

// GetText 获取按钮文本
func (b *Button) GetText() string {
	return b.text
}

// SetState 设置按钮状态
func (b *Button) SetState(state ButtonState) {
	if b.state != state {
		b.state = state
		b.updateColors()
		b.InvalidateLayout()
	}
}

// GetState 获取按钮状态
func (b *Button) GetState() ButtonState {
	return b.state
}


// updateColors 根据状态更新颜色
func (b *Button) updateColors() {
	switch b.state {
	case ButtonStateNormal:
		b.bgColor = &Color{R: 59, G: 130, B: 246, A: 255}
		b.textColor = &Color{R: 255, G: 255, B: 255, A: 255}
		b.border = &BorderStyle{Width: 1, Color: Color{R: 37, G: 99, B: 235, A: 255}}
	case ButtonStateHovered:
		b.bgColor = &Color{R: 37, G: 99, B: 235, A: 255}
		b.textColor = &Color{R: 255, G: 255, B: 255, A: 255}
		b.border = &BorderStyle{Width: 1, Color: Color{R: 29, G: 78, B: 216, A: 255}}
	case ButtonStatePressed:
		b.bgColor = &Color{R: 29, G: 78, B: 216, A: 255}
		b.textColor = &Color{R: 255, G: 255, B: 255, A: 255}
		b.border = &BorderStyle{Width: 1, Color: Color{R: 30, G: 64, B: 175, A: 255}}
	case ButtonStateDisabled:
		b.bgColor = &Color{R: 229, G: 231, B: 235, A: 255}
		b.textColor = &Color{R: 156, G: 163, B: 175, A: 255}
		b.border = &BorderStyle{Width: 1, Color: Color{R: 209, G: 213, B: 219, A: 255}}
	}
}

// Render 渲染按钮
func (b *Button) Render(canvas *thorvg.Canvas, region image.Rectangle) {
	// 调用基类渲染
	b.BaseComponent.Render(canvas, region)

	// 渲染按钮背景
	b.renderBackground(canvas, region)

	// 渲染按钮边框
	b.renderBorder(canvas, region)

	// 渲染按钮文本
	b.renderText(canvas, region)
}

// renderBackground 渲染背景
func (b *Button) renderBackground(canvas *thorvg.Canvas, region image.Rectangle) {
	bounds := b.GetGlobalBounds().Intersect(region)
	if bounds.Empty() {
		return
	}

	bg := thorvg.NewShape()
	if bg != nil {
		bg.AppendRect(
			float32(bounds.Min.X), float32(bounds.Min.Y),
			float32(bounds.Dx()), float32(bounds.Dy()),
			0, 0, false,
		)

		// 设置背景色
		bg.SetFillColor(b.bgColor.R, b.bgColor.G, b.bgColor.B, b.bgColor.A)

		canvas.Push(&bg.Paint)
		bg.Release()
	}
}

// renderBorder 渲染边框
func (b *Button) renderBorder(canvas *thorvg.Canvas, region image.Rectangle) {
	if b.border.Width <= 0 {
		return
	}

	bounds := b.GetGlobalBounds().Intersect(region)
	if bounds.Empty() {
		return
	}

	border := thorvg.NewShape()
	if border != nil {
		border.AppendRect(
			float32(bounds.Min.X), float32(bounds.Min.Y),
			float32(bounds.Dx()), float32(bounds.Dy()),
			0, 0, false,
		)

		// 设置边框样式（简化版本，ThorVG可能不直接支持边框）
		// border.SetStroke(b.border.Color.R, b.border.Color.G, b.border.Color.B, b.border.Color.A)
		// border.SetStrokeWidth(b.border.Width)

		canvas.Push(&border.Paint)
		border.Release()
	}
}

// renderText 渲染文本
func (b *Button) renderText(canvas *thorvg.Canvas, region image.Rectangle) {
	if b.text == "" {
		return
	}

	bounds := b.GetGlobalBounds().Intersect(region)
	if bounds.Empty() {
		return
	}

	// 计算文本位置（居中）
	textX := float32(bounds.Min.X) + float32(bounds.Dx())/2
	textY := float32(bounds.Min.Y) + float32(bounds.Dy())/2

	// TODO: 实际的文本渲染需要字体管理器
	// 暂时跳过文本渲染，因为ThorVG可能不直接支持文本API
	// 后续需要集成字体管理器和文本布局引擎
	_ = textX // 避免未使用变量警告
	_ = textY
	_ = b.text
}

// HandleEvent 处理事件
func (b *Button) HandleEvent(event goui.Event) bool {
	if !b.IsEnabled() {
		return false
	}

	switch e := event.(type) {
	case goui.MouseEvent:
		if e.Down {
			b.SetState(ButtonStatePressed)
		} else {
			if b.state == ButtonStatePressed {
				b.SetState(ButtonStateHovered)
				// 触发点击事件
				if handlers := b.GetEventHandlers(); handlers != nil && handlers.OnClick != nil {
					handlers.OnClick(e)
				}
			}
		}
		return true
	}

	return b.BaseComponent.HandleEvent(event)
}

// SetEnabled 设置启用状态
func (b *Button) SetEnabled(enabled bool) {
	oldEnabled := b.IsEnabled()
	b.BaseComponent.SetEnabled(enabled)

	if oldEnabled != enabled {
		if enabled {
			b.SetState(ButtonStateNormal)
		} else {
			b.SetState(ButtonStateDisabled)
		}
	}
}

// GetEventHandlers 获取事件处理器
func (b *Button) GetEventHandlers() *goui.EventHandlers {
	return b.BaseComponent.GetEventHandlers()
}

// SetBackgroundColor 设置背景色
func (b *Button) SetBackgroundColor(color *Color) {
	b.bgColor = color
	b.InvalidateLayout()
}

// SetTextColor 设置文本颜色
func (b *Button) SetTextColor(color *Color) {
	b.textColor = color
	b.InvalidateLayout()
}

// SetBorder 设置边框样式
func (b *Button) SetBorder(border *BorderStyle) {
	b.border = border
	b.InvalidateLayout()
}

// SetPadding 设置内边距
func (b *Button) SetPadding(padding *Insets) {
	b.padding = padding
	b.InvalidateLayout()
}