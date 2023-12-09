package ui

type Image interface {
	PixelSize() (width, height int32)
	Draw(canvas Canvas)
}
