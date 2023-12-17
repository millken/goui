package ui

import (
	"log/slog"

	"github.com/millken/yoga"
)

type Widget interface {
	SetLayer(layer *Layer)
	HandleSetLayout()
	Layer() *Layer
	Node() *yoga.Node
	Draw(canvas *Canvas)
}

type BaseWidget struct {
	parent   *Layer
	node     *yoga.Node
	style    *Style
	position Rect
}

func undefinedPosition() Rect {
	return NewRect(yoga.Undefined, yoga.Undefined, yoga.Undefined, yoga.Undefined)
}

func NewBaseWidget(yConfig *yoga.Config) *BaseWidget {
	return &BaseWidget{
		node:     yoga.NewNodeWithConfig(yConfig),
		style:    NewStyle(),
		position: undefinedPosition(),
	}
}

func (w *BaseWidget) Node() *yoga.Node {
	return w.node
}

func (w *BaseWidget) Style() *Style {
	return w.style
}

func (w *BaseWidget) SetParent(layer *Layer) {
	w.parent = layer
}

func (w *BaseWidget) SetPosition(x, y, width, height float32) {
	w.position = NewRect(x, y, width, height)
}

func (w *BaseWidget) Position() Rect {
	return w.position
}

// setting yoga layout
func (w *BaseWidget) SetLayout() {
	w.node.StyleSetDisplay(w.style.Display())
	if yoga.IsDefined(w.position.Width) {
		w.node.StyleSetWidth(w.position.Width)
	}
	if yoga.IsDefined(w.position.Height) {
		w.node.StyleSetHeight(w.position.Height)
	}
}

func (w *BaseWidget) CalcLayout() {
	paddingLeft := w.node.LayoutPadding(yoga.EdgeLeft)
	paddingTop := w.node.LayoutPadding(yoga.EdgeTop)
	marginLeft := w.node.LayoutMargin(yoga.EdgeLeft)
	marginTop := w.node.LayoutMargin(yoga.EdgeTop)
	left, top, width, height := w.node.LayoutLeft(), w.node.LayoutTop(), w.node.LayoutWidth(), w.node.LayoutHeight()
	w.position = NewRect(left+marginLeft+paddingLeft, top+marginTop+paddingTop, width, height)
	slog.Info("CalcLayout", "left", left, "top", top, "width", width, "height", height)
}
