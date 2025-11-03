package goui

// Color represents an RGBA color.
type Color struct {
	R, G, B, A uint8
}

type Font struct {
	Name string
	Size int
}

type Image interface {
	Width() int
	Height() int
}

type Point struct {
	X, Y float32
}

type Size struct {
	Width, Height float32
}

type Rect struct {
	X, Y          float32
	Width, Height float32
}

// MouseButton represents a mouse button.
type MouseButton int

const (
	MouseLeft MouseButton = iota + 1
	MouseRight
	MouseMiddle
)
