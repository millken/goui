package goui

import (
	"time"
)

// EventType 事件类型
type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeMouse
	EventTypeKey
	EventTypeTouch
	EventTypeWindow
	EventTypeFocus
	EventTypeScroll
)

// Event 基础事件接口
type Event interface {
	Type() EventType
	Timestamp() time.Time
	Target() Component
	SetTarget(Component)
}

// MouseEvent 鼠标事件
type MouseEvent struct {
	X     float32
	Y     float32
	Down  bool    // true表示按下，false表示释放
	Count int     // 点击次数（用于双击等）
	Delta float32 // 滚轮增量
	eventBase
}

// KeyEvent 键盘事件
type KeyEvent struct {
	Key      KeyCode
	Scancode int
	Down     bool
	Repeat   bool
	Mods     ModifierKey
	eventBase
}

// TouchEvent 触摸事件
type TouchEvent struct {
	ID       int
	X        float32
	Y        float32
	Down     bool
	Pressure  float32
	eventBase
}

// WindowEvent 窗口事件
type WindowEvent struct {
	Width  int
	Height int
	Resized bool
	eventBase
}

// FocusEvent 焦点事件
type FocusEvent struct {
	Gained bool
	eventBase
}

// ScrollEvent 滚动事件
type ScrollEvent struct {
	X     float32
	Y     float32
	DeltaX float32
	DeltaY float32
	eventBase
}

// 键码定义
type KeyCode int

const (
	KeyUnknown KeyCode = iota
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
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
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2

	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
)

// 修饰键
type ModifierKey int

const (
	ModShift ModifierKey = 1 << iota
	ModControl
	ModAlt
	ModSuper
	ModCapsLock
	ModNumLock
)

// 事件基类
type eventBase struct {
	timestamp time.Time
	target    Component
}

func (e *eventBase) Type() EventType {
	return EventTypeUnknown
}

func (e *eventBase) Timestamp() time.Time {
	return e.timestamp
}

func (e *eventBase) Target() Component {
	return e.target
}

func (e *eventBase) SetTarget(target Component) {
	e.target = target
}

// MouseEvent 类型实现
func (e *MouseEvent) Type() EventType { return EventTypeMouse }
func (e *KeyEvent) Type() EventType    { return EventTypeKey }
func (e *TouchEvent) Type() EventType  { return EventTypeTouch }
func (e *WindowEvent) Type() EventType { return EventTypeWindow }
func (e *FocusEvent) Type() EventType  { return EventTypeFocus }
func (e *ScrollEvent) Type() EventType { return EventTypeScroll }

// EventListener 事件监听器
type EventListener interface {
	HandleEvent(event Event) bool
}

// EventFilter 事件过滤器
type EventFilter interface {
	FilterEvent(event Event) bool
}