package ux

import "image/color"

type CmdType int

const (
	CmdFill CmdType = iota
	CmdText
	CmdImage
)

func (c CmdType) String() string {
	switch c {
	case CmdFill:
		return "Fill"
	case CmdText:
		return "Text"
	case CmdImage:
		return "Image"
	default:
		return "Unknown"
	}
}

// Rect represents a rectangle with position and size.
type Rect struct{ X, Y, Width, Height float32 }
type Size struct{ Width, Height float32 }

// RenderCmd 渲染命令（可被 GPU/Canvas/HTML 后端消费）
type RenderCmd struct {
	Data   any
	Type   CmdType
	Rect   Rect
	Color  color.Color
	Text   string
	ZIndex int16 // 最终渲染层级
}
