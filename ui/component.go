package ui

type ComponentType int

const (
	ComponentTypeLabel ComponentType = iota
	ComponentTypeLayout
)

type Component interface {
	Position() Rect
	SetPosition(Rect)
}

type BaseComponent struct {
	Component
	Type         ComponentType
	lastPosition Rect
}

func (c BaseComponent) Position() Rect {
	return c.lastPosition
}

func (c *BaseComponent) SetPosition(r Rect) {
	c.lastPosition = r
}
