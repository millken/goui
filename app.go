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

		// 渲染帧（在Render内部进行脏检查）
		canvas, err := a.renderer.BeginFrame(w, h)
		if err != nil {
			fmt.Printf("渲染器 BeginFrame 错误: window{%d,%d} %v\n", w, h, err)
			continue
		}

		// 如果窗口大小改变了，标记根组件为脏
		if windowResized && a.root != nil {
			a.root.MarkDirty(true)
		}

		// 渲染帧（Render方法内部会检查是否需要渲染）
		rendered := a.root.Paint(canvas)

		if rendered {
			// 渲染一帧
			a.renderer.EndFrame()
			a.renderer.Present()

			// 交换缓冲区
			a.window.SwapBuffers()
			frameCount++
		} else {
			skippedFrames++
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
	if a.root != nil {

	}

	fmt.Println("应用程序已关闭")
}
