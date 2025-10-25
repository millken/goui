package v2

// Component UI组件接口
type Component interface {
	// 标识
	ID() string // 组件唯一标识

	// 布局
	SetPosition(position Point) // 设置位置
	SetSize(size Size)          // 设置尺寸
	GetBounds() Rect            // 获取边界矩形

	// 渲染
	Draw(canvas Canvas) // 绘制组件

	// 事件处理
	HandleEvent(event Event) bool // 处理事件，返回是否已处理

	// 状态
	SetVisible(visible bool) // 设置可见性
	IsVisible() bool         // 是否可见

	// 层次结构
	AddChild(child Component)     // 添加子组件
	RemoveChild(id string)        // 移除子组件
	GetChild(id string) Component // 获取子组件
}
