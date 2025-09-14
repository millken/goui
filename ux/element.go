package ux

import "image/color"

// UIElement 所有 UI 元素的接口
type UIElement interface {
	Render(*Context)
}

// Context 渲染上下文（可扩展）
type Context struct {
	Style  Style
	X, Y   float32
	Width  float32
	Height float32
}

// Style 样式
type Style struct {
	Width           float32
	Height          float32
	Padding         [4]float32 // top, right, bottom, left
	Margin          [4]float32
	BackgroundColor color.Color
	BorderRadius    float32
	FlexGrow        float32
}
