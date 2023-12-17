package ui

type Border interface {
	// Insets returns the insets describing the space the border occupies on each side.
	Insets() Insets
	// Draw the border into rect.
	Draw(canvas *Canvas, rect Rect)
}
