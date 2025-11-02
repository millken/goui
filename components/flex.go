package components

import (
	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/yoga"
)

var _ goui.Component = (*Flex)(nil)

type Flex struct {
	Base
}

func NewFlex() *Flex {
	f := &Flex{Base: NewBase()}
	if n := f.GetNode(); n != nil {
		// Best-effort defaults; use common Yoga Go API names
		n.SetDisplay(yoga.DisplayFlex)
		n.SetFlexDirection(yoga.FlexDirectionRow)
		n.SetAlignItems(yoga.AlignFlexStart)
		n.SetJustifyContent(yoga.JustifyFlexStart)
	}
	return f
}

// Paint 递归绘制子组件
func (f *Flex) Paint(canvas *thorvg.Canvas) {
	if !f.IsVisible() {
		return
	}
	for _, c := range f.GetChildren() {
		c.Paint(canvas)
	}
}
