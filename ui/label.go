package ui

import "github.com/millken/yoga"

type Label struct {
	yNode *yoga.Node
	text  string
}

func NewLabel(text string) *Label {
	node := yoga.NewNode()
	return &Label{
		yNode: node,
		text:  text,
	}
}

func (label *Label) Text() string {
	return label.text
}

func (label *Label) SetText(text string) {
	label.text = text
}

func (label *Label) Node() *yoga.Node {
	return label.yNode
}
