package ux_test

import (
	"image/color"
	"testing"

	"github.com/millken/goui/ux"
)

var (
	colorBlue   = color.RGBA{R: 100, G: 100, B: 255, A: 255}
	colorRed    = color.RGBA{R: 255, G: 100, B: 100, A: 255}
	colorGreen  = color.RGBA{R: 100, G: 255, B: 100, A: 255}
	colorYellow = color.RGBA{R: 255, G: 255, B: 100, A: 255}
)

func TestAA(t *testing.T) {
	ui := ux.NewColumn(
		ux.NewText("Header").WithStyle(
			ux.Width(200),
			ux.Height(40),
			ux.Background(color.RGBA{R: 255, G: 100, B: 100, A: 255}),
			ux.Margin(10, 10, 10, 10),
		),
		ux.NewRow(
			ux.NewText("Left").WithStyle(ux.Width(80), ux.Background(color.Gray{Y: 200})),
			ux.NewBox().WithStyle(ux.Flex(1), ux.Background(colorBlue)),
			ux.NewText("Right").WithStyle(ux.Width(60)),
		).WithStyle(ux.Width(300), ux.Height(60), ux.Margin(0, 0, 10, 0)),
		ux.NewText("Footer").WithStyle(ux.Height(30), ux.Background(colorGreen)),
	).WithStyle(ux.Width(400), ux.Height(300))

	ctx := &ux.Context{X: 0, Y: 0, Width: 400, Height: 300}
	ui.Render(ctx)
}
