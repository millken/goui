package goui

type Panel struct {
	layout *Layout
}

func NewPanel() *Panel {
	return &Panel{}
}

func (p *Panel) SetLayout(layout *Layout) {
	p.layout = layout
}

func (p *Panel) Layout() *Layout {
	return p.layout
}

// PerformLayout is called by the parent container to perform the layout of the panel.
func (p *Panel) PerformLayout() {
}
