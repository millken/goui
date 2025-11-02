package windows

import "github.com/millken/goui"

var _ goui.Window = (*RgfwWindow)(nil)

type RgfwWindow struct {
}

func NewRgfwWindow() *RgfwWindow {
	return &RgfwWindow{}
}

func (w *RgfwWindow) Init(ctx *goui.Context) error {
	return nil
}

func (w *RgfwWindow) Show() {

}

func (w *RgfwWindow) Hide() {

}

func (w *RgfwWindow) Close() {

}

func (w *RgfwWindow) SetTitle(title string) {

}

func (w *RgfwWindow) SetSize(width, height int) {

}

func (w *RgfwWindow) SetRoot(root goui.Component) {

}

func (w *RgfwWindow) PreRender(r goui.Renderer) {

}

func (w *RgfwWindow) PostRenderer(r goui.Renderer) {

}

func (w *RgfwWindow) RequestRedraw() {

}

func (w *RgfwWindow) SetOnClose(f func()) {

}

func (w *RgfwWindow) IsVisible() bool {
	return false
}

func (w *RgfwWindow) GetSize() (width, height int) {
	return 0, 0
}
