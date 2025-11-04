package event

import "time"

//go:generate stringer -type=Type,ModifierKey -output=event_string.go
type Type uint8

const (
	EventUnknown Type = iota
	EventMouse
	EventKey
	EventTouch
	EventWindow
	EventFocus
	EventScroll
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
}

func (e *eventBase) Type() Type {
	return EventUnknown
}

func (e *eventBase) Timestamp() time.Time {
	return e.timestamp
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

func (e *MouseEvent) Type() Type { return EventMouse }

// WindowEvent 窗口事件
type WindowEvent struct {
	Width   int
	Height  int
	Resized bool
	eventBase
}

func (e *WindowEvent) Type() Type { return EventWindow }

// KeyEvent 键盘事件
type KeyEvent struct {
	Key      KeyCode
	Scancode int
	Down     bool
	Repeat   bool
	Mods     ModifierKey
	eventBase
}

func (e *KeyEvent) Type() Type { return EventKey }

// TouchEvent 触摸事件
type TouchEvent struct {
	ID       int
	X        float32
	Y        float32
	Down     bool
	Pressure float32
	eventBase
}

func (e *TouchEvent) Type() Type { return EventTouch }

// FocusEvent 焦点事件
type FocusEvent struct {
	Gained bool
	eventBase
}

func (e *FocusEvent) Type() Type { return EventFocus }

// ScrollEvent 滚动事件
type ScrollEvent struct {
	X      float32
	Y      float32
	DeltaX float32
	DeltaY float32
	eventBase
}

func (e *ScrollEvent) Type() Type { return EventScroll }

// Event 接口定义所有事件的通用方法
type Event interface {
	Type() Type
	Timestamp() time.Time
}

// EventListener 事件监听器接口
type EventListener interface {
	HandleEvent(event Event) bool
}

// NewMouseEvent 创建新的鼠标事件
func NewMouseEvent(x, y float32, down bool, count int) *MouseEvent {
	return &MouseEvent{
		X:     x,
		Y:     y,
		Down:  down,
		Count: count,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}

// NewKeyEvent 创建新的键盘事件
func NewKeyEvent(key KeyCode, scancode int, down, repeat bool, mods ModifierKey) *KeyEvent {
	return &KeyEvent{
		Key:      key,
		Scancode: scancode,
		Down:     down,
		Repeat:   repeat,
		Mods:     mods,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}

// NewWindowEvent 创建新的窗口事件
func NewWindowEvent(width, height int, resized bool) *WindowEvent {
	return &WindowEvent{
		Width:   width,
		Height:  height,
		Resized: resized,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}

// NewFocusEvent 创建新的焦点事件
func NewFocusEvent(gained bool) *FocusEvent {
	return &FocusEvent{
		Gained: gained,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}

// NewScrollEvent 创建新的滚动事件
func NewScrollEvent(x, y, deltaX, deltaY float32) *ScrollEvent {
	return &ScrollEvent{
		X:      x,
		Y:      y,
		DeltaX: deltaX,
		DeltaY: deltaY,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}

// NewTouchEvent 创建新的触摸事件
func NewTouchEvent(id int, x, y float32, down bool, pressure float32) *TouchEvent {
	return &TouchEvent{
		ID:       id,
		X:        x,
		Y:        y,
		Down:     down,
		Pressure: pressure,
		eventBase: eventBase{
			timestamp: time.Now(),
		},
	}
}
