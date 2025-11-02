package goui

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// 应用错误类型
var (
	ErrNoWindows     = errors.New("no active windows")
	ErrRendererInit  = errors.New("renderer initialization failed")
	ErrContextClosed = errors.New("context closed")
	ErrShouldExit    = errors.New("application should exit")
	ErrInvalidConfig = errors.New("invalid configuration")
)

type App interface {
	OnStart(ctx *Context) error
	OnStop(ctx *Context) error
}

func NewApp[T any](t T) *T {
	return &t
}

// createTicker 根据配置创建定时器
func createTicker(settings *Settings) *time.Ticker {
	if settings.FrameRate <= 0 {
		return time.NewTicker(time.Second / 60) // 默认60 FPS
	}
	return time.NewTicker(time.Second / time.Duration(settings.FrameRate))
}

// handleError 处理错误，优先使用自定义错误处理器
func handleError(settings *Settings, err error) {
	if err == nil {
		return
	}

	if settings.ErrorHandler != nil {
		settings.ErrorHandler(err)
	} else {
		// 默认错误处理
		log.Printf("Application error: %v", err)
	}
}

// shouldExit 检查是否应该退出应用
func shouldExit(wm *WindowManager) bool {
	windows := wm.ListWindows()
	if len(windows) == 0 {
		fmt.Println("App.Run: no windows left, exiting")
		return true
	}

	// 检查是否所有窗口都不可见
	anyVisible := false
	for _, w := range windows {
		if w != nil && w.IsVisible() {
			anyVisible = true
			break
		}
	}

	if !anyVisible {
		fmt.Println("App.Run: all windows hidden, exiting")
		return true
	}

	return false
}

// initializeApp 初始化应用
func initializeApp(settings *Settings) (*Context, error) {
	if settings == nil {
		return nil, ErrInvalidConfig
	}

	ctx := NewContext()
	ctx.WM = NewWindowManager()

	return ctx, nil
}

// startApp 启动应用
func startApp(ctx *Context, app App) error {
	// 设置主线程
	if sched := ctx.GetScheduler(); sched != nil {
		sched.SetMain()
	}

	// 启动应用
	if err := app.OnStart(ctx); err != nil {
		return fmt.Errorf("app start failed: %w", err)
	}

	return nil
}

// cleanupApp 清理应用资源
func cleanupApp(ctx *Context) {
	if ctx != nil {
		ctx.Shutdown()
	}
}

// handleWindowEvent 处理单个窗口事件
func handleWindowEvent(ctx *Context, ev WindowEvent) {

	// 处理具体的事件类型
	if ev.Event != nil {
		// 处理鼠标/键盘等输入事件
		if root := ctx.WM.GetRoot(ev.Window); root != nil {
			handled := root.HandleEvent(ev.Event)
			if handled {
				ctx.WM.MarkDirty(ev.Window)
			}
		}
		return
	}

	// 处理窗口级生命周期事件
	switch ev.Type {
	case WindowClosed:
		// 确保在主线程上调用窗口关闭
		if ev.Window != nil {
			ev.Window.Close()
		}
	}
}

// processEvents 处理所有待处理的窗口事件
func processEvents(ctx *Context) error {
	evs := ctx.WindowEvents()
	for _, e := range evs {
		handleWindowEvent(ctx, e)
	}

	return nil
}

// renderWindow 渲染单个窗口
func renderWindow(w Window, renderer Renderer, ctx *Context) error {
	if renderer == nil {
		return ErrRendererInit
	}

	// 初始化渲染器
	if err := renderer.Initialize(w); err != nil {
		return fmt.Errorf("renderer initialize failed: %w", err)
	}

	// 渲染前钩子
	w.PreRender(renderer)

	// 开始帧渲染
	width, height := w.GetSize()
	canvas, err := renderer.BeginFrame(width, height)
	if err != nil {
		return fmt.Errorf("begin frame failed: %w", err)
	}

	// 渲染组件树
	if canvas != nil {
		if root := ctx.WM.GetRoot(w); root != nil {
			root.Paint(canvas)
		}
	}

	// 结束帧渲染
	renderer.EndFrame()
	renderer.Present()

	// 渲染后钩子
	w.PostRenderer(renderer)

	// 清理资源
	renderer.Cleanup()

	return nil
}

// renderDirtyWindows 渲染所有需要重绘的窗口
func renderDirtyWindows(ctx *Context, renderer Renderer, settings *Settings) error {

	dirtyWindows := ctx.WM.PopDirtyWindows()
	if len(dirtyWindows) == 0 {
		return nil
	}

	// 串行渲染（简单实现）
	for _, window := range dirtyWindows {
		if err := renderWindow(window, renderer, ctx); err != nil {
			handleError(settings, err)
		}
	}

	return nil
}

// processFrame 处理一帧的逻辑
func processFrame(ctx *Context, settings *Settings) error {
	// 更新时间
	now := time.Now()
	ctx.tick(now)

	// 驱动平台事件循环

	ctx.WM.CallPump()

	// 执行主线程任务
	if sched := ctx.GetScheduler(); sched != nil {
		sched.Poll()
	}

	// 处理事件
	if err := processEvents(ctx); err != nil {
		return err
	}

	// 检查是否应该退出
	if shouldExit(ctx.WM) {
		return ErrShouldExit
	}

	// 渲染脏窗口
	if err := renderDirtyWindows(ctx, settings.Renderer, settings); err != nil {
		return err
	}

	return nil
}

// runMainLoop 主事件循环
func runMainLoop(ctx *Context, settings *Settings) error {
	ticker := createTicker(settings)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrContextClosed
		case ev := <-ctx.WM.Events():
			// 收集窗口事件
			ctx.PushWindowEvent(ev)
		case <-ticker.C:
			// 处理一帧
			if err := processFrame(ctx, settings); err != nil {
				if err == ErrShouldExit {
					return nil
				}
				if err == ErrContextClosed {
					return err
				}
				// 处理其他错误但继续运行
				handleError(settings, err)
			}
		}
	}
}

// Run 运行应用（重构后的版本）
func Run(settings *Settings, app App) error {
	// 初始化
	ctx, err := initializeApp(settings)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	defer cleanupApp(ctx)

	// 启动应用
	if err := startApp(ctx, app); err != nil {
		return err
	}

	// 注册关闭处理
	defer func() {
		if stopErr := app.OnStop(ctx); stopErr != nil {
			handleError(settings, fmt.Errorf("app stop failed: %w", stopErr))
		}
	}()

	// 运行主循环
	return runMainLoop(ctx, settings)
}
