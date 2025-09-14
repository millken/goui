package ux_test

import (
	"fmt"
	"testing"

	"github.com/millken/goui/ux"
)

func TestAA(t *testing.T) {
	ctx := &ux.Context{X: 0, Y: 0, Width: 400, Height: 300}
	ctx.AddWidget(ux.NewBox().WithStyle(func(s *ux.Style) {
		s.Width = 200
		s.Height = 200
		s.Padding = [4]float32{10, 10, 10, 10}
		// s.ZIndex = 5
	}))
	ctx.AddWidget(ux.NewText("Normal Text", ux.ZIndex(0)))
	cmds := ctx.EndLayout()
	for _, cmd := range cmds {
		fmt.Printf("Z%d: %s | %s | %+v\n", cmd.ZIndex, cmd.Type, cmd.Text, cmd.Rect)
	}
}
