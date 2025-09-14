package ux

import "image/color"

type FontWeight int

const (
	FontWeightNormal FontWeight = iota
	FontWeightBold
	FontWeightLight
)

type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

// TextStyle 专用于文本渲染的样式
type TextStyle struct {
	FontName      string
	FontSize      float32
	FontWeight    FontWeight
	Color         color.Color
	LineHeight    float32 // 倍数，如 1.2
	Alignment     TextAlignment
	LetterSpacing float32 // 字间距
}

// 默认文本样式
var DefaultTextStyle = TextStyle{
	FontName:   "Arial",
	FontSize:   14,
	FontWeight: FontWeightNormal,
	Color:      color.Black,
	LineHeight: 1.2,
	Alignment:  TextAlignLeft,
}

func getDefaultTextStyle(s *Style) TextStyle {
	if ts, ok := s.Data.(TextStyle); ok {
		return ts
	}
	return DefaultTextStyle
}
