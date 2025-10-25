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

const (
	EventMouse EventType = iota + 1
	EventKey
)

// MouseButton represents a mouse button.
type MouseButton int

const (
	MouseLeft MouseButton = iota + 1
	MouseRight
	MouseMiddle
)

// WindowEventType 表示窗口级事件类型
type WindowEventType int

const (
	WindowShown WindowEventType = iota
	WindowHidden
	WindowClosed
	WindowMouse
	WindowKey
)

// WindowEvent 是由 App 主循环产生的窗口事件
type WindowEvent struct {
	Type   WindowEventType
	Window Window
	Event  Event // optional payload (MouseEvent/KeyEvent)
}

// MouseEvent 传递鼠标信息（窗口坐标，按钮，按下或释放）
type MouseEvent struct {
	X, Y   float64
	Button MouseButton
	Down   bool
}

func (m MouseEvent) Type() EventType { return EventMouse }

// KeyEvent 传递键盘信息（使用 GLFW key 常量整数值）
type KeyEvent struct {
	Key  int
	Mod  int
	Down bool
}

func (k KeyEvent) Type() EventType { return EventKey }
