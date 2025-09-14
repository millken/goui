package ux

import "image/color"

// UIElement 所有 UI 元素的接口
type UIElement interface {
	Render(*Context)
}

// Style 样式
type Style struct {
	Data            any // 可存储 fontSize, borderRadius 等
	Color           color.Color
	Width           float32
	Height          float32
	Padding         [4]float32 // top, right, bottom, left
	Margin          [4]float32
	BackgroundColor color.Color
	BorderRadius    float32
	FlexGrow        float32
	ZIndex          int16 // 组件建议的层级
}
