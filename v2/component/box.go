// Package component contains basic component implementations for v2.
package component

import v2 "github.com/millken/goui/v2"

var _ v2.Component = (*BoxComponent)(nil)

type BoxComponent struct {
	rect     v2.Rect
	children []v2.Component
	shown    bool
	color    v2.Color
	onClick  func()
}

// NewBoxComponent creates a new BoxComponent.
func NewBoxComponent() *BoxComponent {
	return &BoxComponent{color: v2.Color{R: 52, G: 152, B: 219, A: 255}}
}

// NewBoxWithColor creates a new BoxComponent with specified color.
func NewBoxWithColor(c v2.Color) *BoxComponent { return &BoxComponent{color: c} }

// SetBounds sets the bounds of the box (helper for examples).
func (b *BoxComponent) SetBounds(x, y, w, h int) {
	b.rect = v2.Rect{X: float64(x), Y: float64(y), Width: float64(w), Height: float64(h)}
}

// SetPosition sets the position of the box.
func (b *BoxComponent) SetPosition(pos v2.Point) { b.rect.X, b.rect.Y = pos.X, pos.Y }

// SetSize sets the size of the box.
func (b *BoxComponent) SetSize(size v2.Size) { b.rect.Width, b.rect.Height = size.Width, size.Height }

// GetBounds returns the box bounds.
func (b *BoxComponent) GetBounds() v2.Rect { return b.rect }

// Paint draws the box and then its children (if any).
func (b *BoxComponent) Paint(c v2.Canvas) {
	// 自身
	c.DrawRect(b.rect, b.color)
	// 子组件
	if len(b.children) == 0 {
		return
	}
	// 创建一个偏移 Canvas，使子组件以相对坐标绘制
	oc := &offsetCanvas{base: c, dx: b.rect.X, dy: b.rect.Y}
	for _, ch := range b.children {
		if ch == nil {
			continue
		}
		if ch.IsVisible() {
			ch.Paint(oc)
		} else {
			// 如果组件未实现可见性，仍然调用 Paint（兼容旧实现）
			ch.Paint(oc)
		}
	}
}

// offsetCanvas 把相对坐标偏移到父坐标再转发到底层 Canvas
type offsetCanvas struct {
	base   v2.Canvas
	dx, dy float64
}

func (oc *offsetCanvas) DrawRect(rect v2.Rect, color v2.Color) {
	rect.X += oc.dx
	rect.Y += oc.dy
	oc.base.DrawRect(rect, color)
}

// HandleEvent handles input event.
func (b *BoxComponent) HandleEvent(e v2.Event) bool {
	// 先对 children 进行 hit-test（从上到下，后加入的优先）
	if len(b.children) > 0 {
		for i := len(b.children) - 1; i >= 0; i-- {
			ch := b.children[i]
			if ch == nil || !ch.IsVisible() {
				continue
			}
			switch ev := e.(type) {
			case v2.MouseEvent:
				// 将事件坐标转换为子组件本地坐标（子组件以父左上角为原点）
				local := v2.MouseEvent{X: ev.X - b.rect.X, Y: ev.Y - b.rect.Y, Button: ev.Button, Down: ev.Down}
				if ch.HandleEvent(local) {
					return true
				}
			case v2.KeyEvent:
				if ch.HandleEvent(ev) {
					return true
				}
			default:
				// forward other events as-is
				if ch.HandleEvent(e) {
					return true
				}
			}
		}
	}

	// 若没有子组件处理，则由自身处理
	switch ev := e.(type) {
	case v2.MouseEvent:
		if ev.Button == v2.MouseLeft && ev.Down {
			if ev.X >= b.rect.X && ev.X <= b.rect.X+b.rect.Width && ev.Y >= b.rect.Y && ev.Y <= b.rect.Y+b.rect.Height {
				if b.onClick != nil {
					b.onClick()
					return true
				}
				// fallback: 切换颜色（蓝色 <-> 红色）
				if b.color.R == 52 && b.color.G == 152 {
					b.color = v2.Color{R: 231, G: 76, B: 60, A: 255} // red
				} else {
					b.color = v2.Color{R: 52, G: 152, B: 219, A: 255} // blue
				}
				return true
			}
		}
	case v2.KeyEvent:
		if ev.Down && ev.Key == v2.KeyR {
			b.color = v2.Color{R: 231, G: 76, B: 60, A: 255}
			return true
		}
	}
	return false
}

// SetVisible sets visibility.
func (b *BoxComponent) SetVisible(v bool) { b.shown = v }

// IsVisible reports visibility.
func (b *BoxComponent) IsVisible() bool { return b.shown }

// AddChild adds a child component.
func (b *BoxComponent) AddChild(child v2.Component) { b.children = append(b.children, child) }

// RemoveChild removes a child component by id (no-op demo).
func (b *BoxComponent) RemoveChild(id string) { /* id-based removal not implemented in demo */ }

// GetChild returns the child component (single child demo).
func (b *BoxComponent) GetChild(id string) v2.Component { return nil }

// SetColor allows changing the box color.
func (b *BoxComponent) SetColor(c v2.Color) { b.color = c }

// SetOnClick registers a callback to be invoked when the box is clicked.
func (b *BoxComponent) SetOnClick(cb func()) { b.onClick = cb }
