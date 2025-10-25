package component

import v2 "github.com/millken/goui/v2"

var _ v2.Component = (*BoxComponent)(nil)

// BoxComponent is a simple container that arranges its children in a box layout.
type BoxComponent struct {
	BaseComponent
	// Additional fields for box layout (e.g., orientation, spacing) can be added here.
}

// NewBoxComponent creates a new BoxComponent.
func NewBoxComponent() *BoxComponent {
	return &BoxComponent{}
}

func (c *BoxComponent) Paint(canvas v2.Canvas) {
	// Paint the box component and its children.
	// Layout logic would go here.
	c.BaseComponent.Draw(canvas)
}
