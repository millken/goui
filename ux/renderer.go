package ux

// Renderer 是渲染后端的通用接口，便于扩展不同实现（CPU 软渲染、GPU 等）。
type Renderer interface {
	// Resize 在窗口/画布尺寸变化时调用（可选）。
	Resize(width, height int)
	// ProcessRenderCommands 处理并渲染一批命令。
	ProcessRenderCommands(commands []RenderCmd) error
	// Cleanup 释放渲染资源。
	Cleanup()
	//WriteToPNG(filename string) error
	WriteToPNG(filename string) error
}
