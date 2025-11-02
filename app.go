package goui

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/millken/go/rgfw"
	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

// App 负责主窗口、渲染、事件分发
type App struct {
	Win      *rgfw.Window
	Renderer Renderer
	Root     Component
	Ctx      *Context
	width    int
	height   int
}

// NewApp 创建一个基本 App，内部初始化 window/OpenGL/ThorVG
func NewApp(title string, w, h int, root Component) (*App, error) {
	// ThorVG init 由外部也可调用，这里确保一次
	if err := thorvg.Init(0); err != nil {
		return nil, fmt.Errorf("thorvg init: %w", err)
	}

	rgfw.SetClassName("Goui")
	win := rgfw.CreateWindow(title, rgfw.Rect{X: 100, Y: 100, W: int32(w), H: int32(h)}, rgfw.WindowCenter|rgfw.WindowAllowDND|rgfw.WindowOpenGL)
	win.MakeCurrentOpenGL()
	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("opengl init: %w", err)
	}
	win.SwapIntervalOpenGL(1)

	gr, err := NewGLRenderer(win, uint32(w), uint32(h))
	if err != nil {
		return nil, err
	}

	app := &App{
		Win:      win,
		Renderer: gr,
		Root:     root,
		Ctx:      NewContext(),
		width:    w,
		height:   h,
	}
	return app, nil
}

// Run 进入主循环
func (a *App) Run() {
	start := time.Now()
	a.Root.Mount()
	for !a.Win.ShouldClose() {
		// 事件轮询（后续可转发到组件树）
		ev := rgfw.GetEvent()
		for a.Win.CheckEvent(ev) {
			// TODO: 转换为 goui.Event 分发
			_ = ev
			rgfw.ReleaseEvent(ev)
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

// Close 释放资源
func (a *App) Close() {
	if a.Renderer != nil {
		a.Renderer.Cleanup()
	}
	if a.Win != nil {
		a.Win.Close()
	}
	if a.Ctx != nil {
		a.Ctx.Shutdown()
	}
	thorvg.Term()
}
