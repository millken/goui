package v2

// Window 抽象了顶层窗口/视图容器的最小契约。
// 职责：
// - 驱动事件循环（Start/Close/RequestRedraw）
// - 持有一个 Renderer，并在需要时请求其绘制
// - 管理根组件的挂载与尺寸变化
// 注意：Window 本身不关心具体渲染细节，交给 Renderer/Canvas。
// Window 表示应用程序窗口
type Window interface {
    // 窗口管理
    Show()                          // 显示窗口
    Hide()                          // 隐藏窗口
    Close()                         // 关闭窗口
    SetTitle(title string)          // 设置窗口标题
    SetSize(width, height int)      // 设置窗口尺寸
    
    // 渲染管理
    SetRenderer(renderer Renderer)  // 设置渲染器
    Render()                         // 渲染窗口内容
    
    // 组件管理
    AddComponent(component Component) // 添加组件
    RemoveComponent(id string)        // 移除组件
    GetComponent(id string) Component // 获取组件
    
    // 事件处理
    SetEventHandler(handler EventHandler) // 设置事件处理器
    ProcessEvents()               // 处理窗口事件
    
    // 状态获取
    IsVisible() bool              // 窗口是否可见
    GetSize() (width, height int) // 获取窗口尺寸
}