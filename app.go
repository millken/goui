package goui

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

// App 负责主窗口、渲染、事件分发
type App struct {
	Win      *glfw.Window
	Renderer Renderer
	Root     Component
	Ctx      *Context
	width    int
	height   int
}

// NewApp 创建一个基本 App，内部初始化 window/OpenGL/ThorVG
func NewApp(title string, w, h int, root Component) (*App, error) {
	runtime.LockOSThread()
	// ThorVG init 由外部也可调用，这里确保一次
	if err := thorvg.Init(0); err != nil {
		return nil, fmt.Errorf("thorvg init: %w", err)
	}

	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("glfw init: %w", err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	win, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("create window: %w", err)
	}
	win.MakeContextCurrent()

	if glErr := gl.Init(); glErr != nil {
		win.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("opengl init: %w", glErr)
	}
	glfw.SwapInterval(1)

	fbw, fbh := win.GetFramebufferSize()
	if fbw <= 0 || fbh <= 0 {
		fbw, fbh = w, h
	}
	gr, err := NewGLRenderer(win, uint32(fbw), uint32(fbh))
	if err != nil {
		win.Destroy()
		glfw.Terminate()
		return nil, err
	}

	app := &App{
		Win:      win,
		Renderer: gr,
		Root:     root,
		Ctx:      NewContext(),
		width:    fbw,
		height:   fbh,
	}

	// 回调：帧缓冲尺寸变化
	win.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		if width <= 0 || height <= 0 {
			return
		}
		if gr2, ok := app.Renderer.(*GLRenderer); ok {
			_ = gr2.resize(uint32(width), uint32(height))
		}
		app.width, app.height = width, height
		app.Ctx.PushWindowEvent(WindowEvent{Type: WindowResize})
	})

	// 回调：鼠标
	win.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		x, y := w.GetCursorPos()
		var mb MouseButton
		switch button {
		case glfw.MouseButtonLeft:
			mb = MouseLeft
		case glfw.MouseButtonRight:
			mb = MouseRight
		case glfw.MouseButtonMiddle:
			mb = MouseMiddle
		default:
			mb = 0
		}
		down := action == glfw.Press
		app.InjectMouse(float32(x), float32(y), mb, down)
	})
	win.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		app.Ctx.PushWindowEvent(WindowEvent{Type: WindowMouse, Event: MouseEvent{X: float32(x), Y: float32(y), Button: 0, Down: false}})
	})

	// 回调：键盘
	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		down := action != glfw.Release
		// 简单映射示例（可扩展）
		mapped := int(key)
		app.InjectKey(mapped, int(mods), down)
	})

	// 回调：鼠标滚轮
	win.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		app.Ctx.PushWindowEvent(WindowEvent{Type: WindowScroll, Event: ScrollEvent{X: float32(xoffset), Y: float32(yoffset)}})
	})

	// 回调：窗口焦点
	win.SetFocusCallback(func(w *glfw.Window, focused bool) {
		app.Ctx.PushWindowEvent(WindowEvent{Type: WindowFocus, Event: FocusEvent{Focused: focused}})
	})

	return app, nil
}

// Run 进入主循环
func (a *App) Run() {
	start := time.Now()
	a.Root.Mount()
	for !a.Win.ShouldClose() {
		// 轮询 GLFW 事件
		glfw.PollEvents()

		// 查询窗口尺寸变化（在没有直接尺寸事件映射时，主动拉取）
		// 注意：若后端提供 GetSize/Rect API，可替换以下获取逻辑
		if a.Win != nil {
			fbw, fbh := a.Win.GetFramebufferSize()
			if fbw > 0 && fbh > 0 && (fbw != a.width || fbh != a.height) {
				if gr, ok := a.Renderer.(*GLRenderer); ok {
					_ = gr.resize(uint32(fbw), uint32(fbh))
				}
				a.width, a.height = fbw, fbh
				a.Ctx.PushWindowEvent(WindowEvent{Type: WindowResize})
			}
		}

		// 分发已收集的 WindowEvents 到组件树
		for _, we := range a.Ctx.WindowEvents() {
			if we.Event == nil {
				continue
			}
			_ = dispatchEvent(a.Root, we.Event)
		}

		// 布局
		w, h := a.width, a.height
		a.Root.Layout(float32(w), float32(h), yoga.DirectionLTR)

		// 开始渲染
		canvas, err := a.Renderer.BeginFrame(int(w), int(h))
		if err != nil {
			fmt.Println("BeginFrame error:", err)
			break
		}

		// 清屏背景
		bg := thorvg.NewShape()
		if bg != nil {
			bg.AppendRect(0, 0, float32(w), float32(h), 0, 0, false)
			bg.SetFillColor(245, 245, 245, 255)
			canvas.Push(&bg.Paint)
			bg.Release()
		}

		// 递归绘制
		a.Root.Paint(canvas)

		a.Renderer.EndFrame()
		a.Renderer.Present()

		// 简单限帧
		_ = start
		time.Sleep(16 * time.Millisecond)
	}

	a.Root.Unmount()
}

// InjectMouse 向应用注入鼠标事件（窗口坐标系，像素；用于平台层对接 RGFW/其它后端）
func (a *App) InjectMouse(x, y float32, btn MouseButton, down bool) {
	a.Ctx.PushWindowEvent(WindowEvent{Type: WindowMouse, Event: MouseEvent{X: x, Y: y, Button: btn, Down: down}})
}

// InjectKey 向应用注入按键事件
func (a *App) InjectKey(key int, mods int, down bool) {
	a.Ctx.PushWindowEvent(WindowEvent{Type: WindowKey, Event: KeyEvent{Key: key, Mod: mods, Down: down}})
}

// Close 释放资源
func (a *App) Close() {
	if a.Renderer != nil {
		a.Renderer.Cleanup()
	}
	if a.Win != nil {
		a.Win.Destroy()
	}
	if a.Ctx != nil {
		a.Ctx.Shutdown()
	}
	thorvg.Term()
	glfw.Terminate()
}
