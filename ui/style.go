package ui

import "github.com/millken/yoga"

type Style struct {
	background Color
	color      Color
	display    yoga.Display
}

func NewStyle() *Style {
	return &Style{}
}

func (s *Style) Background() Color {
	return s.background
}

func (s *Style) SetBackground(color string) {
	s.background = MustColorDecode(color)
}

func (s *Style) Display() yoga.Display {
	return s.display
}

func (s *Style) SetDisplay(display yoga.Display) {
	s.display = display
}

func (s *Style) Color() Color {
	return s.color
}

func (s *Style) SetColor(color string) {
	s.color = MustColorDecode(color)
}
