//go:build !mobile && !js
// +build !mobile,!js

package goui

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// GLFWPlatformSimple 简化的GLFW平台实现
type GLFWPlatformSimple struct {
	window     *glfw.Window
	config     *PlatformConfig
	initTime   time.Time
}

// GLFWPlatformSimpleFactory 简化的GLFW平台工厂
type GLFWPlatformSimpleFactory struct{}

// NewGLFWPlatformSimpleFactory 创建简化的GLFW平台工厂
func NewGLFWPlatformSimpleFactory() *GLFWPlatformSimpleFactory {
	return &GLFWPlatformSimpleFactory{}
}

// CreatePlatform 创建GLFW平台实例
func (f *GLFWPlatformSimpleFactory) CreatePlatform(config *PlatformConfig) (Platform, error) {
	platform := &GLFWPlatformSimple{
		config:   config,
		initTime: time.Now(),
	}

	err := platform.Initialize()
	if err != nil {
		return nil, err
	}

	return platform, nil
}

// GetSupportedBackends 获取支持的后端
func (f *GLFWPlatformSimpleFactory) GetSupportedBackends() []BackendType {
	return []BackendType{BackendGLFW}
}

// IsBackendSupported 检查后端是否支持
func (f *GLFWPlatformSimpleFactory) IsBackendSupported(backend BackendType) bool {
	return backend == BackendGLFW
}

// Initialize 初始化GLFW
func (p *GLFWPlatformSimple) Initialize() error {
	if err := glfw.Init(); err != nil {
		return NewError("failed to initialize GLFW: " + err.Error())
	}

	// 简化窗口提示设置
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.True)

	return nil
}

// Shutdown 关闭GLFW
func (p *GLFWPlatformSimple) Shutdown() error {
	if p.window != nil {
		p.window.Destroy()
		p.window = nil
	}
	glfw.Terminate()
	return nil
}

// CreateWindow 创建窗口
func (p *GLFWPlatformSimple) CreateWindow(config *PlatformConfig) (WindowHandle, error) {
	window, err := glfw.CreateWindow(config.Width, config.Height, config.Title, nil, nil)
	if err != nil {
		return nil, NewError("failed to create GLFW window: " + err.Error())
	}

	p.window = window
	window.MakeContextCurrent()

	// 设置VSync
	if config.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	window.Show()

	return &GLFWWindowHandleSimple{window: window}, nil
}

// DestroyWindow 销毁窗口
func (p *GLFWPlatformSimple) DestroyWindow(window WindowHandle) error {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.Destroy()
		p.window = nil
	}
	return nil
}

// ShowWindow 显示窗口
func (p *GLFWPlatformSimple) ShowWindow(window WindowHandle) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.Show()
	}
}

// HideWindow 隐藏窗口
func (p *GLFWPlatformSimple) HideWindow(window WindowHandle) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.Hide()
	}
}

// GetWindowSize 获取窗口大小
func (p *GLFWPlatformSimple) GetWindowSize(window WindowHandle) (int, int) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		return handle.window.GetSize()
	}
	return 0, 0
}

// SetWindowSize 设置窗口大小
func (p *GLFWPlatformSimple) SetWindowSize(window WindowHandle, width, height int) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.SetSize(width, height)
	}
}

// GetWindowPosition 获取窗口位置
func (p *GLFWPlatformSimple) GetWindowPosition(window WindowHandle) (int, int) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		return handle.window.GetPos()
	}
	return 0, 0
}

// SetWindowPosition 设置窗口位置
func (p *GLFWPlatformSimple) SetWindowPosition(window WindowHandle, x, y int) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.SetPos(x, y)
	}
}

// GetFramebufferSize 获取帧缓冲大小
func (p *GLFWPlatformSimple) GetFramebufferSize(window WindowHandle) (int, int) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		return handle.window.GetFramebufferSize()
	}
	return 0, 0
}

// PollEvents 轮询事件
func (p *GLFWPlatformSimple) PollEvents() Event {
	glfw.PollEvents()
	// 简化的事件处理 - 暂时返回nil，稍后实现具体事件
	return nil
}

// WaitEvents 等待事件
func (p *GLFWPlatformSimple) WaitEvents(timeout time.Duration) Event {
	if timeout > 0 {
		glfw.WaitEventsTimeout(timeout.Seconds())
	} else {
		glfw.WaitEvents()
	}
	return nil
}

// ShouldClose 检查窗口是否应该关闭
func (p *GLFWPlatformSimple) ShouldClose(window WindowHandle) bool {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		return handle.window.ShouldClose()
	}
	return true
}

// SetShouldClose 设置窗口关闭标志
func (p *GLFWPlatformSimple) SetShouldClose(window WindowHandle, shouldClose bool) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.SetShouldClose(shouldClose)
	}
}

// GetKeyState 获取按键状态
func (p *GLFWPlatformSimple) GetKeyState(key KeyCode) bool {
	if p.window != nil {
		// 简化实现
		return false
	}
	return false
}

// GetMouseButton 获取鼠标按钮状态
func (p *GLFWPlatformSimple) GetMouseButton(button int) bool {
	if p.window != nil {
		action := p.window.GetMouseButton(glfw.MouseButtonLeft)
		return action == glfw.Press
	}
	return false
}

// GetMousePosition 获取鼠标位置
func (p *GLFWPlatformSimple) GetMousePosition() (float32, float32) {
	if p.window != nil {
		x, y := p.window.GetCursorPos()
		return float32(x), float32(y)
	}
	return 0, 0
}

// SetMousePosition 设置鼠标位置
func (p *GLFWPlatformSimple) SetMousePosition(x, y float32) {
	if p.window != nil {
		p.window.SetCursorPos(float64(x), float64(y))
	}
}

// SetMouseVisible 设置鼠标可见性
func (p *GLFWPlatformSimple) SetMouseVisible(visible bool) {
	if p.window != nil {
		if visible {
			p.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		} else {
			p.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		}
	}
}

// GetTime 获取时间
func (p *GLFWPlatformSimple) GetTime() float64 {
	return time.Since(p.initTime).Seconds()
}

// SetTime 设置时间
func (p *GLFWPlatformSimple) SetTime(t float64) {
	p.initTime = time.Now().Add(-time.Duration(t * float64(time.Second)))
}

// GetPlatformType 获取平台类型
func (p *GLFWPlatformSimple) GetPlatformType() PlatformType {
	return PlatformDesktop
}

// GetBackendType 获取后端类型
func (p *GLFWPlatformSimple) GetBackendType() BackendType {
	return BackendGLFW
}

// GetClipboard 获取剪贴板内容
func (p *GLFWPlatformSimple) GetClipboard() string {
	if p.window != nil {
		return p.window.GetClipboardString()
	}
	return ""
}

// SetClipboard 设置剪贴板内容
func (p *GLFWPlatformSimple) SetClipboard(text string) {
	if p.window != nil {
		p.window.SetClipboardString(text)
	}
}

// MakeContextCurrent 设置当前OpenGL上下文
func (p *GLFWPlatformSimple) MakeContextCurrent(window WindowHandle) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.MakeContextCurrent()
	}
}

// SwapBuffers 交换缓冲区
func (p *GLFWPlatformSimple) SwapBuffers(window WindowHandle) {
	if handle, ok := window.(*GLFWWindowHandleSimple); ok {
		handle.window.SwapBuffers()
	}
}

// GetProcAddress 获取函数地址
func (p *GLFWPlatformSimple) GetProcAddress(name string) interface{} {
	return glfw.GetProcAddress(name)
}

// GLFWWindowHandleSimple 简化的GLFW窗口句柄
type GLFWWindowHandleSimple struct {
	window *glfw.Window
}

func (h *GLFWWindowHandleSimple) IsValid() bool {
	return h.window != nil
}

func (h *GLFWWindowHandleSimple) GetNativeHandle() interface{} {
	return h.window
}

// 初始化简化的GLFW后端
func init() {
	GlobalPlatformManager.RegisterFactory(BackendGLFW, NewGLFWPlatformSimpleFactory())
}