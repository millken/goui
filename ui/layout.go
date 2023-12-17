package ui

import (
	"log/slog"

	"github.com/millken/yoga"
)

// https://codesandbox.io/p/sandbox/canvas-flexbox-2e9sp?file=%2Fpackage.json
type Layer struct {
	layer    *Layer
	children []Widget
	*BaseWidget
}

func NewLayer(yConfig *yoga.Config) *Layer {
	return &Layer{
		BaseWidget: NewBaseWidget(yConfig),
	}
}

func (l *Layer) SetLayer(layer *Layer) {
	l.layer = layer
}

func (l *Layer) Layer() *Layer {
	return l.layer
}

func (l *Layer) InsertWidget(widget Widget, index uint32) {
	widget.SetLayer(l)
	l.children = append(l.children, widget)
	l.node.InsertChild(widget.Node(), index)
}

func (l *Layer) HandleSetLayout() {
	l.BaseWidget.SetLayout()
	for _, child := range l.children {
		child.HandleSetLayout()
	}
}

func (l *Layer) SetLayout() {
	for _, child := range l.children {
		child.HandleSetLayout()
	}
	slog.Info("set layout", "layer", l.layer)
	yoga.CalculateLayout(l.node, yoga.Undefined, yoga.Undefined, yoga.DirectionLTR)
}

func (l *Layer) Render(canvas *Canvas, rect Rect) {
	l.node.StyleSetWidth(rect.Width)
	l.node.StyleSetHeight(rect.Height)
	slog.Info("draw layer", "rect", rect)
	l.Draw(canvas)
}

func (l *Layer) Draw(canvas *Canvas) {
	if l.layer == nil {
		l.SetLayout()
		for _, child := range l.children {
			child.Draw(canvas)
		}
	} else {
		l.CalcLayout()
		// canvas.Save()
		paint := l.Style().Background().ToPaint(Fill)
		canvas.DrawRect(l.Position(), paint)
		// canvas.Restore()
	}
}
