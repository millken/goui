package goui

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	"github.com/millken/goui/renderer"
	"github.com/millken/goui/window"
	"github.com/millken/yoga"
)

// App 应用程序主结构
type App struct {
	root   component.Component
	config *AppConfig

	eventM   *event.EventManager
	window   window.Window
	renderer renderer.Renderer

	// 运行状态
	running bool
	mu      sync.RWMutex
}

// AppConfig 应用配置
type AppConfig struct {
	TargetFPS     int
	EnableVSync   bool
	EnableStats   bool
	MaxDirtyRects int
	WorkerCount   int

	Resizable    bool
	Fullscreen   bool
	MouseVisible bool
}

// NewApp 创建应用程序
func NewApp(window window.Window, root component.Component) (*App, error) {
	root.SetRoot(root) // 设置根组件的引用
	if err := thorvg.Init(1); err != nil {
		return nil, err
	}
	r, err := renderer.NewGLRenderer()
	if err != nil {
		return nil, err
	}
	app := &App{
		eventM:   event.NewEventManager(),
		root:     root,
		window:   window,
		renderer: r,
	}

	// 注册事件管理器到根组件
	base, ok := root.(component.RootComponent)
	if !ok {
		return nil, fmt.Errorf("根组件必须是 RootComponent 类型")
	}
	// 先将 EM 挂到组件树并刷新监听器，再启动事件循环，避免事件先于监听器注册被处理
	base.SetEventManager(app.eventM)

	if err := app.eventM.Start(); err != nil {
		return nil, err
	}
	return app, nil
}

// Run 运行应用程序（集成事件系统）
func (a *App) Run() {
	runtime.LockOSThread()
	// 画布会自动使用窗口尺寸

	// // 启动系统
	// a.renderEngine.Start(a.canvas)

	fmt.Println("渲染引擎和事件系统已启动，应用程序正在运行...")

	// 如果启用统计，显示性能信息
	// if a.config.EnableStats {
	// 	go a.showStats()
	// }

	fmt.Println("按窗口关闭按钮或Ctrl+C退出应用程序")

	// 渲染统计
	frameCount := 0
	skippedFrames := 0
	statsTicker := time.NewTicker(3 * time.Second)
	defer statsTicker.Stop()

	// 跟踪窗口大小变化
	lastWindowSize := [2]int{0, 0}

	for a.window.IsOpen() {
		frameStart := time.Now()

		// 处理平台事件
		event := a.window.WaitEvents(time.Millisecond * 16) // 大约60 FPS
		if event != nil {
			a.eventM.PostEvent(event)
		}

		// 获取当前窗口大小
		w, h := a.window.GetSize()
		currentWindowSize := [2]int{w, h}

		// 检查窗口大小是否改变
		windowResized := lastWindowSize != currentWindowSize
		lastWindowSize = currentWindowSize

		// 如果窗口大小改变了，先重置渲染场景，避免旧 Canvas 被销毁后仍被使用
		if windowResized && a.root != nil {
			if rs, ok := a.renderer.(interface{ ResetScene() error }); ok {
				_ = rs.ResetScene()
			}
			// 标记根组件为脏以触发一次完整重绘
			a.root.MarkDirty(true)
		}

		// 渲染帧（在Render内部进行脏检查）。注意必须在 ResetScene 之后获取 Canvas
		canvas, err := a.renderer.BeginFrame(w, h)
		if err != nil {
			fmt.Printf("渲染器 BeginFrame 错误: window{%d,%d} %v\n", w, h, err)
			continue
		}

		// 在渲染前计算布局（以窗口尺寸为根容器尺寸）
		if a.root != nil {
			a.root.Layout(float32(w), float32(h), yoga.DirectionLTR)
		}

		// 直接让组件树决定是否需要渲染
		rendered := a.root.Paint(canvas)

		if rendered {
			// 渲染一帧
			a.renderer.EndFrame()
			a.renderer.Present()

			// 交换缓冲区
			a.window.SwapBuffers()
			frameCount++
		} else {
			// 即使没有组件标记脏，窗口缩放等也可能需要刷新背景
			if windowResized {
				// 没有组件绘制，但窗口大小变了：压入背景并绘制一次，避免残影
				bg := thorvg.NewShape()
				if bg != nil {
					bg.AppendRect(0, 0, float32(w), float32(h), 0, 0, false)
					bg.SetFillColor(229, 229, 229, 255)
					canvas.Push(&bg.Paint)
					bg.Release()
				}
				a.renderer.EndFrame()
				a.renderer.Present()
				a.window.SwapBuffers()
				frameCount++
			} else {
				skippedFrames++
			}
		}

		// 帧率控制
		frameTime := time.Since(frameStart)
		if frameTime < 16*time.Millisecond { // 约60 FPS
			time.Sleep(16*time.Millisecond - frameTime)
		}

		// 定期输出统计信息
		select {
		case <-statsTicker.C:
			totalFrames := frameCount + skippedFrames
			if totalFrames > 0 {
				skipRate := float64(skippedFrames) / float64(totalFrames) * 100
				fmt.Printf("渲染优化统计: 总帧数=%d, 渲染=%d, 跳过=%d (跳过率: %.1f%%)\n",
					totalFrames, frameCount, skippedFrames, skipRate)
			}
		default:
		}
	}
}

// isRunning 检查应用程序是否运行中
func (a *App) isRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// registerComponentsRecursive 递归注册组件到事件管理器
// func (a *App) registerComponentsRecursive(component Component) {
// 	a.eventManager.RegisterComponent(component)

// 	for _, child := range component.Children() {
// 		a.registerComponentsRecursive(child)
// 	}
// }

// Close 关闭应用程序
func (a *App) Close() {
	fmt.Println("正在关闭应用程序...")

	// 停止运行
	a.mu.Lock()
	a.running = false
	a.mu.Unlock()

	// 停止事件管理器
	// if a.eventManager != nil {
	// 	a.eventManager.Stop()
	// }

	// 卸载根组件
	// 卸载逻辑留空，后续按需实现

	fmt.Println("应用程序已关闭")
}
