package ux

import "image/color"

var (
	colorWhite  color.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colorBlack  color.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	colorBlue   color.Color = color.RGBA{R: 100, G: 100, B: 255, A: 255}
	colorRed    color.Color = color.RGBA{R: 255, G: 100, B: 100, A: 255}
	colorGreen  color.Color = color.RGBA{R: 100, G: 255, B: 100, A: 255}
	colorYellow color.Color = color.RGBA{R: 255, G: 255, B: 100, A: 255}
	colorGray   color.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
)

// styledElement 接口：用于提取元素的样式
type styledElement interface {
	getStyle() Style
}

// getStyle 从 UIElement 中提取 Style（如果支持）
func getStyle(e UIElement) (Style, bool) {
	if se, ok := e.(styledElement); ok {
		return se.getStyle(), true
	}
	return Style{}, false
}

// StyleMod 样式修饰函数
type StyleMod func(*Style)

func applyStyleMods(base Style, mods ...StyleMod) Style {
	for _, mod := range mods {
		mod(&base)
	}
	return base
}

// 常用样式函数
func Width(w float32) StyleMod {
	return func(s *Style) { s.Width = w }
}

func Height(h float32) StyleMod {
	return func(s *Style) { s.Height = h }
}

func Padding(all float32) StyleMod {
	return func(s *Style) { s.Padding = [4]float32{all, all, all, all} }
}

func Margin(all float32) StyleMod {
	return func(s *Style) { s.Margin = [4]float32{all, all, all, all} }
}

func MarginTop(v float32) StyleMod {
	return func(s *Style) { s.Margin[0] = v }
}
func MarginRight(v float32) StyleMod {
	return func(s *Style) { s.Margin[1] = v }
}
func MarginBottom(v float32) StyleMod {
	return func(s *Style) { s.Margin[2] = v }
}
func MarginLeft(v float32) StyleMod {
	return func(s *Style) { s.Margin[3] = v }
}

func Background(c color.Color) StyleMod {
	return func(s *Style) { s.BackgroundColor = c }
}

func Flex(f float32) StyleMod {
	return func(s *Style) { s.FlexGrow = f }
}
func ZIndex(z int16) StyleMod {
	return func(s *Style) {
		s.ZIndex = z
	}
}
func getColor(c, def color.Color) color.Color {
	if c != nil {
		return c
	}
	return def
}
