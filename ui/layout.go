package ui

import "github.com/millken/yoga"

type Layout struct {
	yNode *yoga.Node
}

func NewLayout() *Layout {
	return &Layout{
		yNode: yoga.NewNode(),
	}
}

func (node *Layout) InsertWidget(widget Widget, index uint32) {
	node.yNode.InsertChild(widget.Node(), index)
}
