package widgets

import (
	"github.com/millken/goui/ui"
	"github.com/millken/yoga"
)

type Button struct {
	title string
	*ui.BaseWidget
}

func NewButton() *Button {
	return &Button{
		BaseWidget: ui.NewBaseWidget(yoga.ConfigNew()),
	}
}

func (b *Button) SetTitle(title string) {
	b.title = title
}

func (b *Button) OnPress() {

}

func NewButtonWithYogaConfig(config *yoga.Config) *Button {
	return &Button{
		BaseWidget: ui.NewBaseWidget(config),
	}
}
