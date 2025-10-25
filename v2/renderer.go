package v2

// Renderer 负责底层图形渲染
type Renderer interface {
	// 初始化与清理
	Initialize(window Window) error // 初始化渲染器
	Cleanup()                       // 清理资源

	// 渲染控制
	BeginFrame() // 开始渲染帧
	EndFrame()   // 结束渲染帧
	Present()    // 呈现到屏幕

	// 状态设置
	SetViewport(x, y, width, height int) // 设置视口
}
