package goui

import "github.com/richardwilkes/toolbox/xmath/geom"

type LabelFlag int

const (
	LabelFlagNone LabelFlag = iota
	LabelFlagSingle
	LabelFlagMulti
)

type Label struct {
	component
	text    string
	flag    LabelFlag
	size    geom.Size[float32]
	align   Align
	onClick *Listener
}

func NewLabel(text string) *Label {
	return &Label{
		text: text,
		flag: LabelFlagSingle,
	}
}

func NewMultilineLabel(text string) *Label {
	return &Label{
		text: text,
		flag: LabelFlagMulti,
	}
}

func (l *Label) SetText(text string) {
	l.text = text
}

func (l *Label) Text() string {
	return l.text
}
