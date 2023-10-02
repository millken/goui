package goui

type UIType int

const (
	UITypeLabel UIType = iota
	UITypeButton
)

type Align int

const (
	AlignLeft Align = iota
	AlignTop
	AlignCenter
	AlignRight
	AlignBottom
	AlignJustify
)
