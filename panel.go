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
