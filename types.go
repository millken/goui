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

type Event interface {
	Type() EventType
}

type EventType int

const (
	EventMouse EventType = iota + 1
	EventKey
	EventScroll
	EventFocus
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
	WindowResize
	WindowScroll
	WindowFocus
)

// WindowEvent 是由 App 主循环产生的窗口事件
type WindowEvent struct {
	Type  WindowEventType
	Event Event // optional payload (MouseEvent/KeyEvent)
}

// MouseEvent 传递鼠标信息（窗口坐标，按钮，按下或释放）
type MouseEvent struct {
	X, Y   float32
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

// ScrollEvent 传递鼠标滚轮信息
type ScrollEvent struct {
	X, Y float32 // 滚轮偏移量
}

func (s ScrollEvent) Type() EventType { return EventScroll }

// FocusEvent 传递窗口焦点信息
type FocusEvent struct {
	Focused bool
}

func (f FocusEvent) Type() EventType { return EventFocus }

// Key constants (generic, backend-agnostic).
// Backends should map their native key codes to these generic constants.
const (
	KeyUnknown = iota
	KeySpace
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyEnter
	KeyEsc
)
