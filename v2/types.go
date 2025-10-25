package v2

// 基本类型
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

// 事件类型
type MouseButton int
