package v2

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
	X, Y float64
}

type Size struct {
	Width, Height float64
}

type Rect struct {
	X, Y          float64
	Width, Height float64
}

type Event interface {
	Type() EventType
}

type EventType int

// MouseButton represents a mouse button.
type MouseButton int

// WindowEventType 表示窗口级事件类型
type WindowEventType int

const (
	WindowShown WindowEventType = iota
	WindowHidden
	WindowClosed
)

// WindowEvent 是由 App 主循环产生的窗口事件
type WindowEvent struct {
	Type   WindowEventType
	Window Window
}
