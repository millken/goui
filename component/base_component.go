package component

import (
	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

type BaseComponent struct {
	yoga     *yoga.Node
	children []Component
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{
		yoga:     yoga.NewNode(),
		children: []Component{},
	}
}

func (b *BaseComponent) AddChild(c Component) {
	if b.children == nil {
		b.children = []Component{}
	}
	b.children = append(b.children, c)
}

func (b *BaseComponent) Children() []Component {
	return b.children
}

func (b *BaseComponent) Render(vg *thorvg.Canvas) {
	for _, child := range b.children {
		child.Render(vg)
	}
}
