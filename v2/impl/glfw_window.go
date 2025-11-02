package impl

import (
	"fmt"
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
	v2 "github.com/millken/goui/v2"
)

var _ v2.Window = (*GLFWWindow)(nil)

var glfwOnce sync.Once
var glfwInitErr error
var pollOnce sync.Once

func ensureGLFW() error {
	glfwOnce.Do(func() { glfwInitErr = glfw.Init() })
	return glfwInitErr
}

func startPoller() {
	// 保留空实现（由 App 主循环统一调用 WM 的 pump）
	pollOnce.Do(func() {})
}

type GLFWWindow struct {
	win     *glfw.Window
	root    v2.Component
	visible bool
	id      string
	ctx     *v2.Context

	// store creation params so we can re-create the native window after it was destroyed
	title  string
	width  int
	height int

	// control goroutine lifecycle removed: GLFW must be called from main thread
	loopMu      sync.Mutex
	loopRunning bool
	onClose     func()
}

// control goroutine removed: all platform calls must be done on main thread

func NewGLFWWindow(title string, width, height int) (*GLFWWindow, error) {
	if err := ensureGLFW(); err != nil {
		return nil, fmt.Errorf("glfw init failed: %w", err)
	}
	startPoller()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// 默认创建为不可见，直到调用 Show()
	glfw.WindowHint(glfw.Visible, glfw.False)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("create window failed: %w", err)
	}
	w := &GLFWWindow{
		win:    win,
		id:     "", // set in Init when ctx is available
		title:  title,
		width:  width,
		height: height,
	}

	// 当 framebuffer 尺寸变化（窗口缩放或 DPI 变化）时，标记窗口为脏以触发重绘
	// attach framebuffer callback and start event loop for this native window
	w.attachFramebufferCallback()
	return w, nil
}

// Init injects app context and registers window id.
func (w *GLFWWindow) Init(ctx *v2.Context) error {
	w.ctx = ctx
	if wm := ctx.GetWindowManager(); wm != nil {
		w.id = wm.NextID()
		// 注册平台事件泵
		wm.SetPump(glfw.PollEvents)
	}

	// 注册输入回调：鼠标按键、鼠标移动与键盘按键
	// ensure callbacks are attached to the native window (if present)
	if w.win != nil {
		w.attachInputCallbacks()
	}
	return nil
}

// attachFramebufferCallback attaches a framebuffer callback to the current native window (if any)
func (w *GLFWWindow) attachFramebufferCallback() {
	if w.win == nil {
		return
	}
	w.win.SetFramebufferSizeCallback(func(_ *glfw.Window, fbw, fbh int) {
		if w == nil {
			return
		}
		if w.ctx != nil {
			if wm := w.ctx.GetWindowManager(); wm != nil {
				wm.MarkDirty(w)
			}
		}
	})
}

// attachInputCallbacks attaches mouse/keyboard/cursor callbacks to the native window
func (w *GLFWWindow) attachInputCallbacks() {
	if w.win == nil {
		return
	}
	w.win.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if w.ctx == nil {
			return
		}
		x, y := win.GetCursorPos()
		me := v2.MouseEvent{X: float32(x), Y: float32(y)}
		switch button {
		case glfw.MouseButtonLeft:
			me.Button = v2.MouseLeft
		case glfw.MouseButtonRight:
			me.Button = v2.MouseRight
		default:
			me.Button = v2.MouseMiddle
		}
		me.Down = (action != glfw.Release)
		if wm := w.ctx.GetWindowManager(); wm != nil {
			wm.Emit(v2.WindowEvent{Type: v2.WindowMouse, Window: w, Event: me})
		}
	})

	w.win.SetCursorPosCallback(func(win *glfw.Window, xpos float64, ypos float64) {
		if w.ctx == nil {
			return
		}
		if wm := w.ctx.GetWindowManager(); wm != nil {
			wm.MarkDirty(w)
			me := v2.MouseEvent{X: float32(xpos), Y: float32(ypos), Button: 0, Down: false}
			wm.Emit(v2.WindowEvent{Type: v2.WindowMouse, Window: w, Event: me})
		}
	})

	w.win.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if w.ctx == nil {
			return
		}
		ke := v2.KeyEvent{Key: mapGlfwKey(key), Mod: int(mods), Down: action != glfw.Release}
		if wm := w.ctx.GetWindowManager(); wm != nil {
			wm.Emit(v2.WindowEvent{Type: v2.WindowKey, Window: w, Event: ke})
		}
	})

	// Close callback: emit WindowClosed so App.Run can handle cleanup
	w.win.SetCloseCallback(func(win *glfw.Window) {
		if w.ctx == nil {
			return
		}
		// if an OnClose handler is registered, invoke it and do not perform default close
		if w.onClose != nil {
			w.onClose()
			return
		}
		if wm := w.ctx.GetWindowManager(); wm != nil {
			wm.Emit(v2.WindowEvent{Type: v2.WindowClosed, Window: w})
		}
	})
}

// mapGlfwKey maps glfw key codes to backend-agnostic v2 key constants.
func mapGlfwKey(k glfw.Key) int {
	switch k {
	case glfw.KeySpace:
		return v2.KeySpace
	case glfw.KeyA:
		return v2.KeyA
	case glfw.KeyB:
		return v2.KeyB
	case glfw.KeyC:
		return v2.KeyC
	case glfw.KeyD:
		return v2.KeyD
	case glfw.KeyE:
		return v2.KeyE
	case glfw.KeyF:
		return v2.KeyF
	case glfw.KeyG:
		return v2.KeyG
	case glfw.KeyH:
		return v2.KeyH
	case glfw.KeyI:
		return v2.KeyI
	case glfw.KeyJ:
		return v2.KeyJ
	case glfw.KeyK:
		return v2.KeyK
	case glfw.KeyL:
		return v2.KeyL
	case glfw.KeyM:
		return v2.KeyM
	case glfw.KeyN:
		return v2.KeyN
	case glfw.KeyO:
		return v2.KeyO
	case glfw.KeyP:
		return v2.KeyP
	case glfw.KeyQ:
		return v2.KeyQ
	case glfw.KeyR:
		return v2.KeyR
	case glfw.KeyS:
		return v2.KeyS
	case glfw.KeyT:
		return v2.KeyT
	case glfw.KeyU:
		return v2.KeyU
	case glfw.KeyV:
		return v2.KeyV
	case glfw.KeyW:
		return v2.KeyW
	case glfw.KeyX:
		return v2.KeyX
	case glfw.KeyY:
		return v2.KeyY
	case glfw.KeyZ:
		return v2.KeyZ
	case glfw.KeyEnter:
		return v2.KeyEnter
	case glfw.KeyEscape:
		return v2.KeyEsc
	default:
		return v2.KeyUnknown
	}
}

func (w *GLFWWindow) Show() {
	w.visible = true
	if wm := w.ctx.GetWindowManager(); wm != nil {
		// 确保只有在未注册时才添加到 WindowManager，防止重复注册导致无法正确移除
		if w.id == "" {
			w.id = wm.NextID()
		}
		if existing := wm.GetWindow(w.id); existing == nil {
			wm.AddWindowWithID(w.id, w)
			fmt.Println("GLFWWindow.Show: added window id=", w.id, "title=", w.title)
		} else {
			fmt.Println("GLFWWindow.Show: window already registered id=", w.id)
		}
		wm.Emit(v2.WindowEvent{Type: v2.WindowShown, Window: w})
		wm.MarkDirty(w)
	}
	// if native window was destroyed earlier, recreate it
	if w.win == nil {
		// reapply hints to ensure the recreated window gets the expected GL context
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		// keep created window hidden until we explicitly show it
		glfw.WindowHint(glfw.Visible, glfw.False)

		// recreate native window on main thread
		if w.ctx != nil && w.ctx.GetScheduler() != nil {
			// synchronously create the window on main thread and wait for result
			w.ctx.GetScheduler().RunWithResult(func() any {
				nw, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
				if err == nil {
					w.win = nw
					w.win.MakeContextCurrent()
					w.attachFramebufferCallback()
					w.attachInputCallbacks()
				}
				return nil
			})
		}
	}
	// show native window on main thread
	if w.win != nil {
		if w.ctx != nil && w.ctx.GetScheduler() != nil {
			ww := w.win
			w.ctx.GetScheduler().RunWithResult(func() any { ww.Show(); return nil })
		} else {
			w.win.Show()
		}
		fmt.Println("GLFWWindow.Show: shown native window id=", w.id)
	}
}

func (w *GLFWWindow) Hide() {
	w.visible = false
	if w.win != nil {
		w.win.Hide()
	}
	// no background goroutine; hide is synchronous
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.Emit(v2.WindowEvent{Type: v2.WindowHidden, Window: w})
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) Close() {
	if w.win != nil {
		ww := w.win
		if w.ctx != nil && w.ctx.GetScheduler() != nil {
			// synchronously destroy on main thread
			w.ctx.GetScheduler().RunWithResult(func() any {
				ww.SetShouldClose(true)
				ww.Destroy()
				return nil
			})
		} else {
			ww.SetShouldClose(true)
			ww.Destroy()
		}
		w.visible = false
		w.win = nil
	}
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.RemoveWindow(w)
	}
}

func (w *GLFWWindow) SetTitle(title string) {
	if w.win != nil {
		w.win.SetTitle(title)
	}
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) SetSize(width, height int) {
	if w.win != nil {
		w.win.SetSize(width, height)
	}
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) SetRoot(root v2.Component) {
	w.root = root
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.BindRoot(w, root)
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) IsVisible() bool { return w.visible }

func (w *GLFWWindow) GetSize() (width, height int) {
	if w.win == nil {
		return 0, 0
	}
	wW, wH := w.win.GetSize()
	return int(wW), int(wH)
}

// PreRender is called before the renderer is initialized.
func (w *GLFWWindow) PreRender(r v2.Renderer) {}

// PostRenderer is called after the renderer is cleaned up.
func (w *GLFWWindow) PostRenderer(r v2.Renderer) {}

func (w *GLFWWindow) RequestRedraw() {
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) removeFromWM() {
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.RemoveWindow(w)
	}
}

// SetOnClose registers a callback invoked when the user requests to close the window
// (e.g., clicking the native close button). If the callback is provided it is
// invoked instead of performing the default destroy behavior, allowing the
// caller to hide the window or perform graceful shutdown.
func (w *GLFWWindow) SetOnClose(cb func()) { w.onClose = cb }

// background loop removed: GLFW calls must be performed on the main thread in App.Run
