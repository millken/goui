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
	pollOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Second / 60)
			defer ticker.Stop()
			for range ticker.C {
				glfw.PollEvents()
			}
		}()
	})
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
	w.startLoop()
	return w, nil
}

// Init injects app context and registers window id.
func (w *GLFWWindow) Init(ctx *v2.Context) error {
	w.ctx = ctx
	if wm := ctx.GetWindowManager(); wm != nil {
		w.id = wm.NextID()
	}
	return nil
}

func (w *GLFWWindow) Show() {
	w.visible = true
	if wm := w.ctx.GetWindowManager(); wm != nil {
		// 避免重复注册：如果 id 为空则在 Init 已生成；若未在 WM 中，注册一次
		// 这里简单地调用 AddWindowWithID，WM 内部可容忍重复或后续可改为查询
		wm.AddWindowWithID(w.id, w)
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
