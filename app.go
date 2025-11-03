package goui

import (
	"fmt"
	"sync"
	"time"

	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

// App 应用程序主结构
type App struct {
	root         Component
	renderEngine *RenderEngine
	canvas       *thorvg.Canvas
	config       *AppConfig

	// 事件系统
	eventManager *EventManager
	platform     Platform
	window       WindowHandle

	// 运行状态
	running      bool
	mu           sync.RWMutex
}

// AppConfig 应用配置
type AppConfig struct {
	Title          string
	Width          int
	Height         int
	TargetFPS      int
	EnableVSync    bool
	EnableStats    bool
	MaxDirtyRects  int
	WorkerCount    int

	// 平台相关
	Backend        BackendType
	Resizable      bool
	Decorated      bool
	Fullscreen     bool
	MouseVisible   bool
}

// NewApp 创建应用程序
func NewApp(title string, width, height int, root Component) (*App, error) {
	config := &AppConfig{
		Title:         title,
		Width:         width,
		Height:        height,
		TargetFPS:     60,
		EnableVSync:   true,
		EnableStats:   true,
		MaxDirtyRects: 10,
		WorkerCount:   0, // 自动检测
		Backend:       BackendGLFW,
		Resizable:     true,
		Decorated:     true,
		Fullscreen:    false,
		MouseVisible:  true,
	}
	return NewAppWithConfig(root, config)
}

// NewAppWithConfig 使用配置创建应用程序
func NewAppWithConfig(root Component, config *AppConfig) (*App, error) {
	if root == nil {
		return nil, fmt.Errorf("root component cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	app := &App{
		root:   root,
		config: config,
	}

	// 创建平台实例
	platformConfig := &PlatformConfig{
		Title:        config.Title,
		Width:        config.Width,
		Height:       config.Height,
		Resizable:    config.Resizable,
		Decorated:    config.Decorated,
		Fullscreen:   config.Fullscreen,
		VSync:        config.EnableVSync,
		Backend:      config.Backend,
		MouseVisible: config.MouseVisible,
	}

	platform, err := GlobalPlatformManager.CreatePlatform(platformConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform: %w", err)
	}
	app.platform = platform

	// 创建事件管理器
	app.eventManager = NewEventManager(platform)

	// 创建渲染引擎
	renderConfig := &RenderConfig{
		TargetFPS:         config.TargetFPS,
		EnableVSync:       config.EnableVSync,
		MaxDirtyRects:     config.MaxDirtyRects,
		WorkerCount:       config.WorkerCount,
		EnableStats:       config.EnableStats,
		EnableMultithread: config.WorkerCount != 1,
	}
	app.renderEngine = NewRenderEngineWithConfig(renderConfig)

	return app, nil
}

// Run 运行应用程序（集成事件系统）
func (a *App) Run() {
	fmt.Printf("启动 %s - GOUI 应用程序\n", a.config.Title)
	fmt.Printf("分辨率: %dx%d, 目标FPS: %d\n", a.config.Width, a.config.Height, a.config.TargetFPS)

	// 创建窗口
	window, err := a.platform.CreateWindow(&PlatformConfig{
		Title:        a.config.Title,
		Width:        a.config.Width,
		Height:       a.config.Height,
		Resizable:    a.config.Resizable,
		Decorated:    a.config.Decorated,
		Fullscreen:   a.config.Fullscreen,
		VSync:        a.config.EnableVSync,
		Backend:      a.config.Backend,
		MouseVisible: a.config.MouseVisible,
	})
	if err != nil {
		fmt.Printf("错误：无法创建窗口: %v\n", err)
		return
	}
	a.window = window
	defer a.platform.DestroyWindow(window)

	// 初始化根组件
	if rootBase, ok := a.root.(*BaseComponent); ok {
		rootBase.mount()

		// 执行初始布局
		rootBase.layout(float32(a.config.Width), float32(a.config.Height), yoga.DirectionLTR)

		// 注册根组件到渲染引擎
		a.renderEngine.RegisterRenderable(LayerContent, rootBase)

		// 注册组件到事件管理器
		a.eventManager.RegisterComponent(rootBase)
		a.registerComponentsRecursive(rootBase)
	}

	// 创建ThorVG画布（使用软件画布，不需要OpenGL上下文）
	a.canvas = thorvg.NewSwCanvas(thorvg.EngineOptionNone)
	if a.canvas == nil {
		fmt.Println("错误：无法创建ThorVG画布")
		return
	}
	defer a.canvas.Destroy()

	// 画布会自动使用窗口尺寸

	// 启动系统
	a.renderEngine.Start(a.canvas)
	a.eventManager.Start()

	fmt.Println("渲染引擎和事件系统已启动，应用程序正在运行...")

	// 如果启用统计，显示性能信息
	if a.config.EnableStats {
		go a.showStats()
	}

	// 主事件循环
	a.mu.Lock()
	a.running = true
	a.mu.Unlock()

	fmt.Println("按窗口关闭按钮或Ctrl+C退出应用程序")

	for a.isRunning() && !a.platform.ShouldClose(window) {
		// 处理平台事件
		event := a.platform.PollEvents()
		if event != nil {
			a.eventManager.PostEvent(event)
		}

		// 触发根组件更新
		if rootBase, ok := a.root.(*BaseComponent); ok {
			rootBase.Update()
		}

		// 渲染一帧
		a.renderFrame()

		// 交换缓冲区
		a.platform.SwapBuffers(window)
	}
}

// isRunning 检查应用程序是否运行中
func (a *App) isRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// registerComponentsRecursive 递归注册组件到事件管理器
func (a *App) registerComponentsRecursive(component Component) {
	a.eventManager.RegisterComponent(component)

	for _, child := range component.Children() {
		a.registerComponentsRecursive(child)
	}
}

// renderFrame 渲染一帧
func (a *App) renderFrame() {
	// 强制触发一次渲染更新，确保有绘制调用
	if a.root != nil {
		if rootBase, ok := a.root.(*BaseComponent); ok {
			// 标记需要重新渲染
			rootBase.markNeedsRender()
		}
	}
}

// Close 关闭应用程序
func (a *App) Close() {
	fmt.Println("正在关闭应用程序...")

	// 停止运行
	a.mu.Lock()
	a.running = false
	a.mu.Unlock()

	// 停止事件管理器
	if a.eventManager != nil {
		a.eventManager.Stop()
	}

	// 停止渲染引擎
	if a.renderEngine != nil {
		a.renderEngine.Stop()
	}

	// 卸载根组件
	if a.root != nil {
		if rootBase, ok := a.root.(*BaseComponent); ok {
			rootBase.unmount()
			rootBase.destroy()
		}
	}

	// 关闭平台
	if a.platform != nil {
		a.platform.Shutdown()
	}

	fmt.Println("应用程序已关闭")
}

// showStats 显示性能统计
func (a *App) showStats() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for a.renderEngine.IsRunning() {
		select {
		case <-ticker.C:
			stats := a.renderEngine.GetDetailedStats()
			fmt.Printf("渲染统计 - 帧数: %d, FPS: %.1f, 绘制调用: %d, 脏区域: %d\n",
				stats.FrameCount, stats.FPS, stats.DrawCalls, stats.DirtyRects)
		}
	}
}