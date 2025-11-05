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
		// 默认配置：开启 VSync，目标帧率 60
		config: &AppConfig{
			TargetFPS:   60,
			EnableVSync: true,
		},
	}

	// 应用默认/传入配置：若后端支持 VSync，按配置设置
	if app.config != nil {
		window.SetVSync(app.config.EnableVSync)
	}

	// 注册事件管理器到根组件
	base, ok := root.(component.RootComponent)
	if !ok {
		return nil, fmt.Errorf("根组件必须是 RootComponent 类型")
	}
	// 注入窗口引用，供组件内部在交互时调节光标等
	base.SetWindow(window)
	// 先将 EM 挂到组件树并刷新监听器，再启动事件循环，避免事件先于监听器注册被处理
	base.SetEventManager(app.eventM)

	if err := app.eventM.Start(); err != nil {
		return nil, err
	}
	return app, nil
}

// SetConfig 设置应用配置并立即应用与平台相关的选项（如 VSync）
func (a *App) SetConfig(cfg *AppConfig) {
	a.mu.Lock()
	a.config = cfg
	a.mu.Unlock()
	if cfg != nil && a.window != nil {
		a.window.SetVSync(cfg.EnableVSync)
	}
}

// Run 运行应用程序（集成事件系统）
func (a *App) Run() {
	runtime.LockOSThread()

	// 如果启用统计，显示性能信息
	// if a.config.EnableStats {
	// 	go a.showStats()
	// }

	// 渲染统计
	frameCount := 0
	skippedFrames := 0
	statsTicker := time.NewTicker(3 * time.Second)
	defer statsTicker.Stop()

	// 跟踪窗口大小变化
	lastWindowSize := [2]int{0, 0}

	for a.window.IsOpen() {
		frameStart := time.Now()

		// 处理平台事件：先快速排空队列，减少 hover 等连续事件的滞后
		hadEvents := false
		if fetcher, ok := a.window.(interface{ FetchEvents(int) []event.Event }); ok {
			events := fetcher.FetchEvents(0) // 0 表示尽可能多地取出
			for _, ev := range events {
				if ev != nil {
					_ = a.eventM.PostEvent(ev)
					hadEvents = true
				}
			}
		} else {
			// 回退到非阻塞单个获取
			ev := a.window.WaitEvents(0)
			if ev != nil {
				_ = a.eventM.PostEvent(ev)
				hadEvents = true
			}
		}

		// 获取当前窗口大小
		w, h := a.window.GetSize()
		currentWindowSize := [2]int{w, h}

		// 检查窗口大小是否改变
		windowResized := lastWindowSize != currentWindowSize
		lastWindowSize = currentWindowSize

		// 如果窗口大小改变了，先重置渲染场景，避免旧 Canvas 被销毁后仍被使用
		if windowResized && a.root != nil {
			_ = a.renderer.ResetScene(w, h)
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
			layoutNeeded := windowResized
			if checker, ok := a.root.(interface{ CheckLayoutDirtyRecursive() bool }); ok {
				if checker.CheckLayoutDirtyRecursive() {
					layoutNeeded = true
				}
			}
			if layoutNeeded {
				a.root.Layout(float32(w), float32(h), yoga.DirectionLTR)
				if clearer, ok := a.root.(interface{ ClearLayoutDirtyRecursive() }); ok {
					clearer.ClearLayoutDirtyRecursive()
				}
			}
		}

		// 直接让组件树决定是否需要渲染
		rendered := a.root.Paint(canvas)

		if rendered {
			// 渲染一帧
			if err := a.renderer.EndFrame(); err != nil {
				fmt.Printf("EndFrame 错误: %v\n", err)
			}

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
				if err := a.renderer.EndFrame(); err != nil {
					fmt.Printf("EndFrame 错误: %v\n", err)
				}
				a.window.SwapBuffers()
				frameCount++
			} else {
				skippedFrames++
			}
		}

		// 帧率控制
		frameTime := time.Since(frameStart)
		// 如果未设置 TargetFPS，默认 60fps；有事件或刚渲染过一帧不强制睡眠，提升跟手性
		targetFPS := 60
		if a.config != nil && a.config.TargetFPS > 0 {
			targetFPS = a.config.TargetFPS
		}
		if targetFPS > 0 && !hadEvents && !rendered {
			budget := time.Second / time.Duration(targetFPS)
			if frameTime < budget {
				time.Sleep(budget - frameTime)
			}
		}

		// 定期输出统计信息
		select {
		case <-statsTicker.C:
			if a.config == nil || a.config.EnableStats {
				totalFrames := frameCount + skippedFrames
				if totalFrames > 0 {
					skipRate := float64(skippedFrames) / float64(totalFrames) * 100
					// 输出事件丢弃统计（如有）
					var dropped uint64
					if a.eventM != nil {
						dropped = a.eventM.DroppedCount()
					}
					if dropped > 0 {
						fmt.Printf("渲染优化统计: 总帧数=%d, 渲染=%d, 跳过=%d (跳过率: %.1f%%), 丢弃事件=%d\n",
							totalFrames, frameCount, skippedFrames, skipRate, dropped)
					} else {
						fmt.Printf("渲染优化统计: 总帧数=%d, 渲染=%d, 跳过=%d (跳过率: %.1f%%)\n",
							totalFrames, frameCount, skippedFrames, skipRate)
					}
				}
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

// Close 关闭应用程序
func (a *App) Close() {
	// 停止运行
	a.mu.Lock()
	a.running = false
	a.mu.Unlock()

	// 停止事件管理器
	if a.eventM != nil {
		a.eventM.Stop()
	}
	// 释放渲染器资源
	if a.renderer != nil {
		a.renderer.Cleanup()
	}
	// 释放 ThorVG 全局资源
	thorvg.Term()
}
