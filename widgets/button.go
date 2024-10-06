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

func NewButtonWithYogaConfig(config *yoga.Config) *Button {
	return &Button{
		BaseWidget: ui.NewBaseWidget(config),
	}
}
