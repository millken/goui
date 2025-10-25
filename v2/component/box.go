// Package component contains basic component implementations for v2.
package component

import v2 "github.com/millken/goui/v2"

var _ v2.Component = (*BoxComponent)(nil)

type BoxComponent struct {
	rect  v2.Rect
	child v2.Component
	shown bool
}

// NewBoxComponent creates a new BoxComponent.
func NewBoxComponent() *BoxComponent { return &BoxComponent{} }

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

// Paint draws the box as a solid colored rectangle for demo.
func (b *BoxComponent) Paint(c v2.Canvas) {
	c.DrawRect(b.rect, v2.Color{R: 52, G: 152, B: 219, A: 255}) // 蓝色
}

// HandleEvent handles input event.
func (b *BoxComponent) HandleEvent(e v2.Event) bool { return false }

// SetVisible sets visibility.
func (b *BoxComponent) SetVisible(v bool) { b.shown = v }

// IsVisible reports visibility.
func (b *BoxComponent) IsVisible() bool { return b.shown }

// AddChild adds a child component.
func (b *BoxComponent) AddChild(child v2.Component) { b.child = child }

// RemoveChild removes a child component by id (no-op demo).
func (b *BoxComponent) RemoveChild(id string) { b.child = nil }

// GetChild returns the child component (single child demo).
func (b *BoxComponent) GetChild(id string) v2.Component { return b.child }
