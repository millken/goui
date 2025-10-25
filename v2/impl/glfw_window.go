package impl

import (
	"fmt"
	"sync"
	"time"

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

	loopOnce    sync.Once
	ctrl        chan controlOp
	removedOnce sync.Once
}

type controlOp int

const (
	ctrlShow controlOp = iota
	ctrlHide
	ctrlClose
)

func NewGLFWWindow(title string, width, height int) (*GLFWWindow, error) {
	if err := ensureGLFW(); err != nil {
		return nil, fmt.Errorf("glfw init failed: %w", err)
	}
	startPoller()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("create window failed: %w", err)
	}
	w := &GLFWWindow{
		win:  win,
		id:   "", // set in Init when ctx is available
		ctrl: make(chan controlOp, 8),
	}

	// 当 framebuffer 尺寸变化（窗口缩放或 DPI 变化）时，标记窗口为脏以触发重绘
	win.SetFramebufferSizeCallback(func(_ *glfw.Window, fbw, fbh int) {
		if w == nil {
			return
		}
		if w.ctx != nil {
			if wm := w.ctx.GetWindowManager(); wm != nil {
				wm.MarkDirty(w)
			}
		}
	})
	w.startLoop()
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
	if w.win != nil {
		w.win.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			if w.ctx == nil {
				return
			}
			x, y := win.GetCursorPos()
			me := v2.MouseEvent{X: x, Y: y}
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
				// mark dirty so cursor driven interactions may force repaint if required
				wm.MarkDirty(w)
				me := v2.MouseEvent{X: xpos, Y: ypos, Button: 0, Down: false}
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
	}
	return nil
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
		}
		wm.Emit(v2.WindowEvent{Type: v2.WindowShown, Window: w})
		wm.MarkDirty(w)
	}
	if w.win != nil {
		w.win.Show()
	}
	select {
	case w.ctrl <- ctrlShow:
	default:
	}
}

func (w *GLFWWindow) Hide() {
	w.visible = false
	if w.win != nil {
		w.win.Hide()
	}
	select {
	case w.ctrl <- ctrlHide:
	default:
	}
	if wm := w.ctx.GetWindowManager(); wm != nil {
		wm.Emit(v2.WindowEvent{Type: v2.WindowHidden, Window: w})
		wm.MarkDirty(w)
	}
}

func (w *GLFWWindow) Close() {
	if w.win != nil {
		w.win.SetShouldClose(true)
	}
	select {
	case w.ctrl <- ctrlClose:
	default:
	}
	w.removeFromWM()
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
	w.removedOnce.Do(func() {
		if wm := w.ctx.GetWindowManager(); wm != nil {
			wm.RemoveWindow(w)
		}
	})
}

func (w *GLFWWindow) startLoop() {
	w.loopOnce.Do(func() {
		go func() {
			for {
				select {
				case op := <-w.ctrl:
					switch op {
					case ctrlShow:
						// platform show handled in Show()
					case ctrlHide:
						// platform hide handled in Hide()
					case ctrlClose:
						if w.win != nil {
							w.win.Destroy()
							w.visible = false
							w.win = nil
							w.removeFromWM()
						}
						return
					}
				case <-time.After(time.Second / 60):
					// check close flag periodically
					if w.win != nil && w.win.ShouldClose() {
						w.win.Destroy()
						w.visible = false
						w.win = nil
						w.removeFromWM()
						return
					}
				}
			}
		}()
	})
}
