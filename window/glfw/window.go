// Package glfw provides GLFW window implementation for the goui framework.
// It handles mouse, keyboard, window, focus, and scroll events using the GLFW library.
package glfw

import (
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/millken/goui/event"
	gouiwin "github.com/millken/goui/window"
)

var initOnce sync.Once

type Window struct {
	native    *glfw.Window
	lastEvent event.Event
	events    chan event.Event
	// cursor cache
	cursors map[gouiwin.Cursor]*glfw.Cursor
}

func NewWindow(title string, width, height int) (*Window, error) {
	initOnce.Do(func() {
		if err := glfw.Init(); err != nil {
			panic(err)
		}
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		// glfw.WindowHint(glfw.Resizable, glfw.True)
	})
	native, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, err
	}
	native.MakeContextCurrent()
	win := &Window{
		native:  native,
		events:  make(chan event.Event, 256),
		cursors: make(map[gouiwin.Cursor]*glfw.Cursor),
	}
	// 设置事件回调
	win.setupCallbacks()
	return win, nil
}

func (w *Window) Open() error {
	w.native.Show()
	return nil
}

func (w *Window) Close() error {
	// 销毁缓存的光标
	for k, c := range w.cursors {
		if c != nil {
			c.Destroy()
		}
		delete(w.cursors, k)
	}
	w.native.Destroy()
	return nil
}

func (w *Window) SetTitle(title string) {
	w.native.SetTitle(title)
}

func (w *Window) SetSize(width, height int) {
	w.native.SetSize(width, height)
}

func (w *Window) GetSize() (width, height int) {
	return w.native.GetSize()
}

func (w *Window) SetResizable(resizable bool) {
	if resizable {
		w.native.SetAttrib(glfw.Resizable, glfw.True)
	} else {
		w.native.SetAttrib(glfw.Resizable, glfw.False)
	}
}

// SetDecorated 设置窗口装饰
func (w *Window) SetDecorated(decorated bool) {
	if decorated {
		w.native.SetAttrib(glfw.Decorated, glfw.True)
	} else {
		w.native.SetAttrib(glfw.Decorated, glfw.False)
	}
}

func (w *Window) Fullscreen() {
	// 简化处理，实际实现可能需要更多逻辑
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	w.native.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
}

func (w *Window) SetMouseVisible(visible bool) {
	if visible {
		w.native.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.native.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

func (w *Window) SetVSync(enabled bool) {
	if enabled {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

func (w *Window) SetCursor(c gouiwin.Cursor) {
	// lazily create glfw cursors
	get := func(shape gouiwin.Cursor) *glfw.Cursor {
		if cur, ok := w.cursors[shape]; ok && cur != nil {
			return cur
		}
		var std glfw.StandardCursor
		switch shape {
		case gouiwin.CursorPointer:
			std = glfw.HandCursor
		case gouiwin.CursorText:
			std = glfw.IBeamCursor
		case gouiwin.CursorCrosshair:
			std = glfw.CrosshairCursor
		case gouiwin.CursorEWResize:
			std = glfw.HResizeCursor
		case gouiwin.CursorNSResize:
			std = glfw.VResizeCursor
		default:
			std = glfw.ArrowCursor
		}
		cur := glfw.CreateStandardCursor(std)
		w.cursors[shape] = cur
		return cur
	}
	if w.native != nil {
		w.native.SetCursor(get(c))
	}
}

func (w *Window) IsOpen() bool {
	return !w.native.ShouldClose()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) WaitEvents(timeout time.Duration) event.Event {
	// 当 timeout 为 0 时，不阻塞等待，直接轮询一次以触发回调并尽快返回
	if timeout <= 0 {
		glfw.PollEvents()
	} else {
		glfw.WaitEventsTimeout(timeout.Seconds())
	}
	select {
	case ev := <-w.events:
		return ev
	default:
		// 兼容旧路径：返回并清空 lastEvent（如果有）
		ev := w.lastEvent
		w.lastEvent = nil
		return ev
	}
}

// FetchEvents 尽可能在不阻塞的情况下批量获取事件，最多返回 max 个。
// 当 max<=0 时，返回当前可用的所有事件。该方法还会对连续的鼠标移动事件进行简单合并，减少冗余处理。
func (w *Window) FetchEvents(max int) []event.Event {
	// 快速轮询以触发回调，将系统事件推送到队列
	glfw.PollEvents()

	// 从通道拉取事件
	var out []event.Event
	pull := func() (event.Event, bool) {
		select {
		case ev := <-w.events:
			return ev, true
		default:
			if w.lastEvent != nil {
				ev := w.lastEvent
				w.lastEvent = nil
				return ev, true
			}
			return nil, false
		}
	}

	limit := max
	if max <= 0 {
		limit = 1<<31 - 1
	}

	// 简易合并：如果连续读取到 MouseEvent（移动）且仅位置变化，则只保留最后一个
	var pendingMove *event.MouseEvent
	for len(out) < limit {
		ev, ok := pull()
		if !ok {
			break
		}
		if me, ok := ev.(*event.MouseEvent); ok {
			if !me.Down && me.Count == 0 {
				// 移动事件，合并
				pendingMove = me
				continue
			}
			// 非移动或按键相关事件，先冲刷 pendingMove
			if pendingMove != nil {
				out = append(out, pendingMove)
				pendingMove = nil
			}
			out = append(out, ev)
			continue
		}
		// 其他事件，冲刷 pendingMove 后追加
		if pendingMove != nil {
			out = append(out, pendingMove)
			pendingMove = nil
		}
		out = append(out, ev)
	}
	if pendingMove != nil && len(out) < limit {
		out = append(out, pendingMove)
	}
	return out
}

// GetLastEvent 获取最后一个事件
func (w *Window) GetLastEvent() event.Event {
	return w.lastEvent
}

// ClearLastEvent 清除最后一个事件
func (w *Window) ClearLastEvent() {
	w.lastEvent = nil
}

// GetKeyState 获取按键状态
func (w *Window) GetKeyState(key event.KeyCode) bool {
	// 将event.KeyCode转换为glfw.Key的映射
	glfwKey := w.eventKeyCodeToGlfwKey(key)
	return w.native.GetKey(glfwKey) == glfw.Press
}

// GetMouseButtonState 获取鼠标按键状态
func (w *Window) GetMouseButtonState(button int) bool {
	var glfwButton glfw.MouseButton
	switch button {
	case 0:
		glfwButton = glfw.MouseButtonLeft
	case 1:
		glfwButton = glfw.MouseButtonRight
	case 2:
		glfwButton = glfw.MouseButtonMiddle
	default:
		return false
	}
	return w.native.GetMouseButton(glfwButton) == glfw.Press
}

// GetMousePosition 获取鼠标位置
func (w *Window) GetMousePosition() (float32, float32) {
	x, y := w.native.GetCursorPos()
	return float32(x), float32(y)
}

// SetMousePosition 设置鼠标位置
func (w *Window) SetMousePosition(x, y float32) {
	w.native.SetCursorPos(float64(x), float64(y))
}

// eventKeyCodeToGlfwKey 将event.KeyCode转换为glfw.Key
func (w *Window) eventKeyCodeToGlfwKey(key event.KeyCode) glfw.Key {
	// 反向查找映射表
	for glfwKey, keyCode := range glfwKeyToEventKeyCodeMap {
		if keyCode == key {
			return glfwKey
		}
	}
	return glfw.KeyUnknown
}

func (w *Window) SwapBuffers() {
	w.native.SwapBuffers()
}

func (w *Window) setupCallbacks() {
	// 鼠标按钮回调
	w.native.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		xpos, ypos := win.GetCursorPos()
		mouseEvent := event.NewMouseEvent(
			float32(xpos),
			float32(ypos),
			action == glfw.Press,
			1, // TODO: 实现点击计数逻辑
		)
		select { // 尝试写入事件队列，满了则回退到 lastEvent
		case w.events <- mouseEvent:
		default:
			w.lastEvent = mouseEvent
		}
	})

	// 鼠标移动回调
	w.native.SetCursorPosCallback(func(win *glfw.Window, xpos, ypos float64) {
		mouseEvent := event.NewMouseEvent(
			float32(xpos),
			float32(ypos),
			false, // 鼠标移动时不表示按下状态
			0,
		)
		select {
		case w.events <- mouseEvent:
		default:
			w.lastEvent = mouseEvent
		}
	})

	// 键盘回调
	w.native.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		keyEvent := event.NewKeyEvent(
			w.glfwKeyToKeyCode(key),
			scancode,
			action == glfw.Press || action == glfw.Repeat,
			action == glfw.Repeat,
			w.glfwModsToModifierKey(mods),
		)
		select {
		case w.events <- keyEvent:
		default:
			w.lastEvent = keyEvent
		}
	})

	// 窗口大小调整回调
	w.native.SetFramebufferSizeCallback(func(win *glfw.Window, width, height int) {
		windowEvent := event.NewWindowEvent(width, height, true)
		select {
		case w.events <- windowEvent:
		default:
			w.lastEvent = windowEvent
		}
	})

	// 滚轮回调
	w.native.SetScrollCallback(func(win *glfw.Window, xoff, yoff float64) {
		xpos, ypos := win.GetCursorPos()
		scrollEvent := event.NewScrollEvent(
			float32(xpos),
			float32(ypos),
			float32(xoff),
			float32(yoff),
		)
		select {
		case w.events <- scrollEvent:
		default:
			w.lastEvent = scrollEvent
		}
	})

	// 窗口焦点回调
	w.native.SetFocusCallback(func(win *glfw.Window, focused bool) {
		focusEvent := event.NewFocusEvent(focused)
		select {
		case w.events <- focusEvent:
		default:
			w.lastEvent = focusEvent
		}
	})
}

// GLFW键码到事件键码的映射表
var glfwKeyToEventKeyCodeMap = map[glfw.Key]event.KeyCode{
	glfw.KeyA:            event.KeyCodeA,
	glfw.KeyB:            event.KeyCodeB,
	glfw.KeyC:            event.KeyCodeC,
	glfw.KeyD:            event.KeyCodeD,
	glfw.KeyE:            event.KeyCodeE,
	glfw.KeyF:            event.KeyCodeF,
	glfw.KeyG:            event.KeyCodeG,
	glfw.KeyH:            event.KeyCodeH,
	glfw.KeyI:            event.KeyCodeI,
	glfw.KeyJ:            event.KeyCodeJ,
	glfw.KeyK:            event.KeyCodeK,
	glfw.KeyL:            event.KeyCodeL,
	glfw.KeyM:            event.KeyCodeM,
	glfw.KeyN:            event.KeyCodeN,
	glfw.KeyO:            event.KeyCodeO,
	glfw.KeyP:            event.KeyCodeP,
	glfw.KeyQ:            event.KeyCodeQ,
	glfw.KeyR:            event.KeyCodeR,
	glfw.KeyS:            event.KeyCodeS,
	glfw.KeyT:            event.KeyCodeT,
	glfw.KeyU:            event.KeyCodeU,
	glfw.KeyV:            event.KeyCodeV,
	glfw.KeyW:            event.KeyCodeW,
	glfw.KeyX:            event.KeyCodeX,
	glfw.KeyY:            event.KeyCodeY,
	glfw.KeyZ:            event.KeyCodeZ,
	glfw.Key0:            event.KeyCode0,
	glfw.Key1:            event.KeyCode1,
	glfw.Key2:            event.KeyCode2,
	glfw.Key3:            event.KeyCode3,
	glfw.Key4:            event.KeyCode4,
	glfw.Key5:            event.KeyCode5,
	glfw.Key6:            event.KeyCode6,
	glfw.Key7:            event.KeyCode7,
	glfw.Key8:            event.KeyCode8,
	glfw.Key9:            event.KeyCode9,
	glfw.KeyEnter:        event.KeyCodeReturnEnter,
	glfw.KeyTab:          event.KeyCodeTab,
	glfw.KeySpace:        event.KeyCodeSpacebar,
	glfw.KeyMinus:        event.KeyCodeHyphenMinus,
	glfw.KeyEqual:        event.KeyCodeEqualSign,
	glfw.KeyLeftBracket:  event.KeyCodeLeftSquareBracket,
	glfw.KeyRightBracket: event.KeyCodeRightSquareBracket,
	glfw.KeyBackslash:    event.KeyCodeBackslash,
	glfw.KeySemicolon:    event.KeyCodeSemicolon,
	glfw.KeyApostrophe:   event.KeyCodeApostrophe,
	glfw.KeyGraveAccent:  event.KeyCodeGraveAccent,
	glfw.KeyComma:        event.KeyCodeComma,
	glfw.KeyPeriod:       event.KeyCodeFullStop,
	glfw.KeySlash:        event.KeyCodeSlash,
	glfw.KeyEscape:       event.KeyCodeEscape,
	glfw.KeyCapsLock:     event.KeyCodeCapsLock,
	glfw.KeyBackspace:    event.KeyCodeDeleteBackspace,
	glfw.KeyPause:        event.KeyCodePause,
	glfw.KeyInsert:       event.KeyCodeInsert,
	glfw.KeyHome:         event.KeyCodeHome,
	glfw.KeyPageUp:       event.KeyCodePageUp,
	glfw.KeyDelete:       event.KeyCodeDeleteForward,
	glfw.KeyEnd:          event.KeyCodeEnd,
	glfw.KeyPageDown:     event.KeyCodePageDown,
	glfw.KeyRight:        event.KeyCodeRightArrow,
	glfw.KeyLeft:         event.KeyCodeLeftArrow,
	glfw.KeyDown:         event.KeyCodeDownArrow,
	glfw.KeyUp:           event.KeyCodeUpArrow,
	glfw.KeyF1:           event.KeyCodeF1,
	glfw.KeyF2:           event.KeyCodeF2,
	glfw.KeyF3:           event.KeyCodeF3,
	glfw.KeyF4:           event.KeyCodeF4,
	glfw.KeyF5:           event.KeyCodeF5,
	glfw.KeyF6:           event.KeyCodeF6,
	glfw.KeyF7:           event.KeyCodeF7,
	glfw.KeyF8:           event.KeyCodeF8,
	glfw.KeyF9:           event.KeyCodeF9,
	glfw.KeyF10:          event.KeyCodeF10,
	glfw.KeyF11:          event.KeyCodeF11,
	glfw.KeyF12:          event.KeyCodeF12,
	glfw.KeyLeftShift:    event.KeyCodeLeftShift,
	glfw.KeyLeftControl:  event.KeyCodeLeftControl,
	glfw.KeyLeftAlt:      event.KeyCodeLeftAlt,
	glfw.KeyLeftSuper:    event.KeyCodeLeftGUI,
	glfw.KeyRightShift:   event.KeyCodeRightShift,
	glfw.KeyRightControl: event.KeyCodeRightControl,
	glfw.KeyRightAlt:     event.KeyCodeRightAlt,
	glfw.KeyRightSuper:   event.KeyCodeRightGUI,
}

func (w *Window) glfwKeyToKeyCode(key glfw.Key) event.KeyCode {
	if keyCode, exists := glfwKeyToEventKeyCodeMap[key]; exists {
		return keyCode
	}
	return event.KeyCodeUnknown
}

func (w *Window) glfwModsToModifierKey(mods glfw.ModifierKey) event.ModifierKey {
	var result event.ModifierKey
	if mods&glfw.ModShift != 0 {
		result |= event.ModShift
	}
	if mods&glfw.ModControl != 0 {
		result |= event.ModControl
	}
	if mods&glfw.ModAlt != 0 {
		result |= event.ModAlt
	}
	if mods&glfw.ModSuper != 0 {
		result |= event.ModSuper
	}
	if mods&glfw.ModCapsLock != 0 {
		result |= event.ModCapsLock
	}
	if mods&glfw.ModNumLock != 0 {
		result |= event.ModNumLock
	}
	return result
}
