package glfw

import (
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/millken/goui/event"
)

var initOnce sync.Once

type Window struct {
	native *glfw.Window
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
	return &Window{
		native: native,
	}, nil
}

func (w *Window) Open() error {
	w.native.Show()
	return nil
}

func (w *Window) Close() error {
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

func (w *Window) IsOpen() bool {
	return !w.native.ShouldClose()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) WaitEvents(timeout time.Duration) *event.Event {
	glfw.WaitEventsTimeout(timeout.Seconds())
	return nil
}

func (w *Window) SwapBuffers() {
	w.native.SwapBuffers()
}
