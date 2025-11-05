package goui

import (
	"runtime"
	"sync"

	"github.com/millken/goui/component"
	"github.com/millken/goui/window"
)

type App struct {
	// 兼容旧 API 的封装：内部委托 Application/Session
	appl   *Application
	sess   *Session
	window window.Window
	root   component.Component
	config *AppConfig
	mu     sync.RWMutex
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
	appl := NewApplication()
	sess, err := appl.NewSession(window, root, &AppConfig{TargetFPS: 60, EnableVSync: true})
	if err != nil {
		return nil, err
	}
	return &App{appl: appl, sess: sess, window: window, root: root, config: sess.cfg}, nil
}

// SetConfig 设置应用配置并立即应用与平台相关的选项（如 VSync）
func (a *App) SetConfig(cfg *AppConfig) {
	a.mu.Lock()
	a.config = cfg
	a.mu.Unlock()
	if cfg != nil {
		if a.window != nil {
			a.window.SetVSync(cfg.EnableVSync)
		}
		if a.sess != nil {
			a.sess.cfg = cfg
		}
	}
}

// Run 运行应用程序（集成事件系统）
func (a *App) Run() {
	runtime.LockOSThread()
	if a.appl != nil {
		a.appl.Run()
	}
}

// isRunning 检查应用程序是否运行中
func (a *App) isRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	// 旧 API 保留：运行状态由 Application 管理，这里简单返回 true 当 session/window 存在
	return a.appl != nil && a.sess != nil
}

// Close 关闭应用程序
func (a *App) Close() {
	if a.appl != nil {
		a.appl.Close()
	}
}
