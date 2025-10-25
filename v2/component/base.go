package component

import v2 "github.com/millken/goui/v2"

// BaseComponent 提供了 v2.Component 的基础实现（空实现/最小状态）。
type BaseComponent struct {
	position v2.Point
	size     v2.Size
	visible  bool
	children map[string]v2.Component
}

func (c *BaseComponent) SetPosition(position v2.Point) {
	c.position = position
}

func (c *BaseComponent) SetSize(size v2.Size) {
	c.size = size
}

// SetBounds 以整数形式设置位置与尺寸，便于示例/适配器调用。
func (c *BaseComponent) SetBounds(x, y, width, height int) {
	c.position = v2.Point{X: float64(x), Y: float64(y)}
	c.size = v2.Size{Width: float64(width), Height: float64(height)}
}

func (c *BaseComponent) GetBounds() v2.Rect {
	return v2.Rect{X: c.position.X, Y: c.position.Y, Width: c.size.Width, Height: c.size.Height}
}

// Paint 默认不进行绘制；子类型按需覆写。
func (c *BaseComponent) Paint(canvas v2.Canvas) {}

// HandleEvent 默认不处理事件。
func (c *BaseComponent) HandleEvent(event v2.Event) bool { return false }

func (c *BaseComponent) SetVisible(visible bool) { c.visible = visible }

func (c *BaseComponent) IsVisible() bool { return c.visible }

func (c *BaseComponent) AddChild(child v2.Component) {
	if c.children == nil {
		c.children = make(map[string]v2.Component)
	}
	// 这里没有 ID 概念，存一个临时 key；真实实现应由子类提供 ID。
	c.children["child"] = child
}

func (c *BaseComponent) RemoveChild(id string) {
	if c.children == nil {
		return
	}
	delete(c.children, id)
}

func (c *BaseComponent) GetChild(id string) v2.Component {
	if c.children == nil {
		return nil
	}
	return c.children[id]
}
