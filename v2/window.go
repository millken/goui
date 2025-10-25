package v2

// Window 抽象了顶层窗口/视图容器的最小契约。
// 职责：
// - 驱动事件循环（Start/Close/RequestRedraw）
// - 持有一个 Renderer，并在需要时请求其绘制
// - 管理根组件的挂载与尺寸变化
// 注意：Window 本身不关心具体渲染细节，交给 Renderer/Canvas。
// Window 表示应用程序窗口
type Window interface {
	// 初始化：注入应用 Context，便于窗口访问 WindowManager 等上下文资源
	Init(ctx *Context) error
	// 窗口管理
	Show()                     // 显示窗口
	Hide()                     // 隐藏窗口
	Close()                    // 关闭窗口
	SetTitle(title string)     // 设置窗口标题
	SetSize(width, height int) // 设置窗口尺寸
	// 组件树
	SetRoot(root Component) // 设置根组件
	// 渲染管理（hooks 由 Window 调度）
	PreRender(r Renderer)    // 渲染前钩子（如进场动画）
	PostRenderer(r Renderer) // 渲染后钩子（如退场、清理）
	RequestRedraw()          // 请求下一帧重绘

	// 注册关闭回调：当用户点击窗口的关闭按钮时会触发此回调
	// 若回调处理了关闭（例如仅隐藏窗口），窗口实现应尊重此行为而不进行默认销毁。
	SetOnClose(func())

	// 状态获取
	IsVisible() bool              // 窗口是否可见
	GetSize() (width, height int) // 获取窗口尺寸
}
