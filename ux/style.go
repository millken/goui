package ux

import "image/color"

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

func Margin(top, right, bottom, left float32) StyleMod {
	return func(s *Style) { s.Margin = [4]float32{top, right, bottom, left} }
}

func Background(c color.Color) StyleMod {
	return func(s *Style) { s.BackgroundColor = c }
}

func Flex(f float32) StyleMod {
	return func(s *Style) { s.FlexGrow = f }
}
