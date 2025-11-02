package components

import (
	"github.com/millken/goui"
	"github.com/millken/yoga"
)

type Base struct {
}

func (b *Base) Layout(width, height float32, direction yoga.Direction) {
}

func (b *Base) Paint(canvas goui.Canvas) {
}

func (b *Base) GetNode() *yoga.Node {
	return nil
}

func (b *Base) GetComputedLayout() yoga.Layout {
	return yoga.Layout{}
}

func (b *Base) SetVisible(visible bool) {
}

func (b *Base) IsVisible() bool {
	return true
}

func (b *Base) SetID(id string) {
}

func (b *Base) GetID() string {
	return ""
}

func (b *Base) AddChild(child goui.Component) {
}

func (b *Base) RemoveChild(child goui.Component) {
}

func (b *Base) GetChildren() []goui.Component {
	return nil
}

func (b *Base) HandleEvent(event goui.Event) bool {
	return false
}

func (b *Base) Mount() {
}

func (b *Base) Unmount() {
}
func (b *Base) Update() {
}
func (b *Base) Destroy() {
}

func (b *Base) SetOnClick(handler func()) {
}

func (b *Base) SetOnHover(handler func(bool)) {
}
