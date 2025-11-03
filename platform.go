package goui

import (
	"time"
)

// PlatformType 平台类型
type PlatformType int

const (
	PlatformDesktop PlatformType = iota
	PlatformMobile
	PlatformWeb
)

// BackendType 后端类型
type BackendType int

const (
	BackendGLFW BackendType = iota
	BackendRGFW
	BackendSDL
	BackendAndroid
	BackendIOS
	BackendWebGL
)

// PlatformConfig 平台配置
type PlatformConfig struct {
	Title        string
	Width        int
	Height       int
	Resizable    bool
	Decorated    bool
	Fullscreen   bool
	VSync        bool
	Samples      int
	Backend      BackendType
	MouseVisible bool
}

// WindowHandle 窗口句柄（抽象）
type WindowHandle interface {
	IsValid() bool
	GetNativeHandle() interface{}
}

// Platform 平台抽象接口
type Platform interface {
	// 生命周期管理
	Initialize() error
	Shutdown() error

	// 窗口管理
	CreateWindow(config *PlatformConfig) (WindowHandle, error)
	DestroyWindow(window WindowHandle) error
	ShowWindow(window WindowHandle)
	HideWindow(window WindowHandle)

	// 窗口属性
	GetWindowSize(window WindowHandle) (int, int)
	SetWindowSize(window WindowHandle, width, height int)
	GetWindowPosition(window WindowHandle) (int, int)
	SetWindowPosition(window WindowHandle, x, y int)
	GetFramebufferSize(window WindowHandle) (int, int)

	// 事件处理
	PollEvents() Event
	WaitEvents(timeout time.Duration) Event
	ShouldClose(window WindowHandle) bool
	SetShouldClose(window WindowHandle, shouldClose bool)

	// 输入状态
	GetKeyState(key KeyCode) bool
	GetMouseButton(button int) bool
	GetMousePosition() (float32, float32)
	SetMousePosition(x, y float32)
	SetMouseVisible(visible bool)

	// 时间管理
	GetTime() float64
	SetTime(time float64)

	// 平台信息
	GetPlatformType() PlatformType
	GetBackendType() BackendType
	GetClipboard() string
	SetClipboard(text string)

	// 渲染上下文
	MakeContextCurrent(window WindowHandle)
	SwapBuffers(window WindowHandle)

	// 扩展接口
	GetProcAddress(name string) interface{}
}

// PlatformFactory 平台工厂接口
type PlatformFactory interface {
	CreatePlatform(config *PlatformConfig) (Platform, error)
	GetSupportedBackends() []BackendType
	IsBackendSupported(backend BackendType) bool
}

// PlatformManager 平台管理器
type PlatformManager struct {
	factories map[BackendType]PlatformFactory
	current   Platform
}

// NewPlatformManager 创建平台管理器
func NewPlatformManager() *PlatformManager {
	return &PlatformManager{
		factories: make(map[BackendType]PlatformFactory),
	}
}

// RegisterFactory 注册平台工厂
func (pm *PlatformManager) RegisterFactory(backend BackendType, factory PlatformFactory) {
	pm.factories[backend] = factory
}

// CreatePlatform 创建平台实例
func (pm *PlatformManager) CreatePlatform(config *PlatformConfig) (Platform, error) {
	factory, exists := pm.factories[config.Backend]
	if !exists {
		// 尝试自动选择最佳后端
		for backend, f := range pm.factories {
			if f.IsBackendSupported(backend) {
				factory = f
				config.Backend = backend
				break
			}
		}
		if factory == nil {
			return nil, ErrNoSupportedBackend
		}
	}

	platform, err := factory.CreatePlatform(config)
	if err != nil {
		return nil, err
	}

	pm.current = platform
	return platform, nil
}

// GetCurrentPlatform 获取当前平台
func (pm *PlatformManager) GetCurrentPlatform() Platform {
	return pm.current
}

// GetSupportedBackends 获取支持的后端列表
func (pm *PlatformManager) GetSupportedBackends() []BackendType {
	var backends []BackendType
	for backend, factory := range pm.factories {
		if factory.IsBackendSupported(backend) {
			backends = append(backends, backend)
		}
	}
	return backends
}

// 全局平台管理器实例
var GlobalPlatformManager = NewPlatformManager()

// 错误定义
var (
	ErrNoSupportedBackend = NewError("no supported backend found")
	ErrPlatformInitFailed = NewError("platform initialization failed")
	ErrWindowCreateFailed = NewError("window creation failed")
)

// Error 自定义错误类型
type Error struct {
	message string
}

func NewError(message string) *Error {
	return &Error{message: message}
}

func (e *Error) Error() string {
	return e.message
}