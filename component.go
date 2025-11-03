package goui

import (
	"image"
	"sync"

	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)


// Component 精简的组件接口
type Component interface {
	// 基础属性
	ID() string
	Children() []Component
	Parent() Component

	// 几何计算
	Bounds() image.Rectangle

	// 渲染相关
	Renderable

	// 生命周期
	Update()

	// 事件处理
	OnClick(func(MouseEvent))
	OnHover(func(bool))

	// 状态管理
	SetVisible(bool)
	IsVisible() bool
	SetEnabled(bool)
	IsEnabled() bool

	// 组件树操作
	AddChild(Component)
	RemoveChild(Component)
	FindChild(string) Component
	setParent(Component)
}

// EventHandlers 事件处理器集合
type EventHandlers struct {
	OnClick    func(MouseEvent)
	OnHover    func(bool)
	OnFocus    func()
	OnBlur     func()
	OnKeyPress func(KeyEvent)
	OnScroll   func(ScrollEvent)
	OnMount    func()
	OnUnmount  func()
	OnUpdate   func()
}

// BaseComponent 基础组件实现
type BaseComponent struct {
	id        string
	visible   bool
	enabled   bool
	node      *yoga.Node
	parent    Component
	children  []Component

	// 渲染相关
	layer       RenderLayer
	needsLayout bool
	needsRender bool

	// 事件处理
	events *EventHandlers

	// 同步锁
	mu sync.RWMutex

	// 渲染缓存
	canvas *thorvg.Canvas
}

// NewBaseComponent 创建基础组件
func NewBaseComponent(id string) *BaseComponent {
	node := yoga.NewNode()
	return &BaseComponent{
		id:          id,
		visible:     true,
		enabled:     true,
		node:        node,
		layer:       LayerContent,
		needsLayout: true,
		needsRender: true,
		events:      &EventHandlers{},
	}
}

// 精简接口实现
func (b *BaseComponent) ID() string { return b.id }
func (b *BaseComponent) Children() []Component { return b.children }
func (b *BaseComponent) Parent() Component { return b.parent }

// 内部方法（不暴露在接口中，但可选择性使用）
func (b *BaseComponent) SetID(id string) { b.id = id }

// 内部方法（不暴露给用户，但子组件可访问）
func (b *BaseComponent) getNode() *yoga.Node { return b.node }

// GetYogaNode 公开方法供子组件使用（内部实现细节）
func (b *BaseComponent) GetYogaNode() *yoga.Node { return b.node }

// 几何计算方法（精简接口）
func (b *BaseComponent) Bounds() image.Rectangle {
	if !b.visible {
		return image.Rectangle{}
	}
	return b.getGlobalBounds()
}

// 内部方法
func (b *BaseComponent) GetBounds() image.Rectangle {
	return b.Bounds()
}

func (b *BaseComponent) GetLayer() RenderLayer { return b.layer }

func (b *BaseComponent) ShouldRender(region image.Rectangle) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.visible || !b.needsRender {
		return false
	}

	bounds := b.GetGlobalBounds()
	return bounds.Intersect(region).Empty() == false
}

func (b *BaseComponent) Render(canvas *thorvg.Canvas, region image.Rectangle) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.visible {
		return
	}

	// 创建渲染上下文
	if b.canvas == nil {
		b.canvas = canvas
	}

	// 子类实现具体渲染逻辑
	b.onRender(canvas, region)

	// 递归渲染子组件
	for _, child := range b.children {
		if child.ShouldRender(region) {
			child.Render(canvas, region)
		}
	}

	b.needsRender = false
}

func (b *BaseComponent) onRender(canvas *thorvg.Canvas, region image.Rectangle) {
	// 默认空实现，由子类重写
}

// 内部布局方法（不暴露给用户）
func (b *BaseComponent) layout(width, height float32, direction yoga.Direction) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.node == nil {
		return
	}

	b.node.CalculateLayout(width, height, direction)
	b.needsLayout = false

	// 递归布局子组件
	for _, child := range b.children {
		if baseChild, ok := child.(*BaseComponent); ok {
			childLayout := baseChild.node.GetComputedLayout()
			baseChild.layout(childLayout.Width, childLayout.Height, direction)
		}
	}
}

func (b *BaseComponent) invalidateLayout() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.needsLayout = true
	b.needsRender = true

	// 向上传播布局失效
	if b.parent != nil {
		if parentBase, ok := b.parent.(*BaseComponent); ok {
			parentBase.invalidateLayout()
		}
	}
}

// 向后兼容方法
func (b *BaseComponent) Layout(width, height float32, direction yoga.Direction) {
	b.layout(width, height, direction)
}

// InvalidateLayout 公开的布局失效方法（供子组件调用）
func (b *BaseComponent) InvalidateLayout() {
	b.invalidateLayout()
}

// 精简的生命周期方法（仅保留Update在接口中）
func (b *BaseComponent) Update() {
	if b.events.OnUpdate != nil {
		b.events.OnUpdate()
	}

	// 标记需要重新渲染
	b.mu.Lock()
	b.needsRender = true
	b.mu.Unlock()

	// 递归更新子组件
	for _, child := range b.children {
		child.Update()
	}
}


// 内部方法（不暴露在接口中）
func (b *BaseComponent) mount() {
	if b.events.OnMount != nil {
		b.events.OnMount()
	}

	// 递归挂载子组件
	for _, child := range b.children {
		if baseChild, ok := child.(*BaseComponent); ok {
			baseChild.mount()
		}
	}
}

func (b *BaseComponent) unmount() {
	if b.events.OnUnmount != nil {
		b.events.OnUnmount()
	}

	// 递归卸载子组件
	for _, child := range b.children {
		if baseChild, ok := child.(*BaseComponent); ok {
			baseChild.unmount()
		}
	}
}

func (b *BaseComponent) destroy() {
	// 清理资源
	b.canvas = nil

	// 递归销毁子组件
	for _, child := range b.children {
		if baseChild, ok := child.(*BaseComponent); ok {
			baseChild.destroy()
		}
	}
}

// 向后兼容方法
func (b *BaseComponent) Mount() {
	b.mount()
}

func (b *BaseComponent) Unmount() {
	b.unmount()
}

func (b *BaseComponent) Destroy() {
	b.destroy()
}

// 精简的事件处理API（接口方法）
func (b *BaseComponent) OnClick(handler func(MouseEvent)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events.OnClick = handler
}

func (b *BaseComponent) OnHover(handler func(bool)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events.OnHover = handler
}

// 内部事件处理
func (b *BaseComponent) HandleEvent(event Event) bool {
	// 默认实现处理基本事件
	switch e := event.(type) {
	case MouseEvent:
		if b.events.OnClick != nil && e.Down {
			b.events.OnClick(e)
			return true
		}
	case KeyEvent:
		if b.events.OnKeyPress != nil {
			b.events.OnKeyPress(e)
			return true
		}
	}
	return false
}

// 向后兼容方法
func (b *BaseComponent) SetEventHandlers(handlers *EventHandlers) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = handlers
}

func (b *BaseComponent) GetEventHandlers() *EventHandlers {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.events
}

// 状态管理
func (b *BaseComponent) SetVisible(visible bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.visible != visible {
		b.visible = visible
		b.needsRender = true
	}
}

func (b *BaseComponent) IsVisible() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.visible
}

func (b *BaseComponent) SetEnabled(enabled bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.enabled != enabled {
		b.enabled = enabled
		b.needsRender = true
	}
}

func (b *BaseComponent) IsEnabled() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.enabled
}

// 便利方法：显示/隐藏
func (b *BaseComponent) Show() {
	b.SetVisible(true)
}

func (b *BaseComponent) Hide() {
	b.SetVisible(false)
}

// 便利方法：启用/禁用
func (b *BaseComponent) Enable() {
	b.SetEnabled(true)
}

func (b *BaseComponent) Disable() {
	b.SetEnabled(false)
}

// 便利方法：切换状态
func (b *BaseComponent) ToggleVisible() {
	b.SetVisible(!b.IsVisible())
}

func (b *BaseComponent) ToggleEnabled() {
	b.SetEnabled(!b.IsEnabled())
}

// 组件树操作（接口方法）
func (b *BaseComponent) AddChild(child Component) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.children = append(b.children, child)
	child.setParent(b)

	// 同步到Yoga节点
	if b.node != nil {
		if childBase, ok := child.(*BaseComponent); ok && childBase.node != nil {
			b.node.InsertChild(childBase.node, uint32(len(b.children)-1))
		}
	}

	// 自动挂载子组件（如果父组件已挂载）
	if b.parent != nil {
		if childBase, ok := child.(*BaseComponent); ok {
			childBase.mount()
		}
	}

	// 标记需要重新布局
	b.needsLayout = true
	b.needsRender = true
}

func (b *BaseComponent) RemoveChild(child Component) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			child.setParent(nil)

			// 自动卸载子组件
			if childBase, ok := child.(*BaseComponent); ok {
				childBase.unmount()
			}

			// 从Yoga节点移除
			if b.node != nil {
				if childBase, ok := child.(*BaseComponent); ok && childBase.node != nil {
					b.node.RemoveChild(childBase.node)
				}
			}

			b.needsLayout = true
			b.needsRender = true
			break
		}
	}
}

func (b *BaseComponent) SetParent(parent Component) {
	b.setParent(parent)
}

func (b *BaseComponent) setParent(parent Component) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.parent = parent
}

func (b *BaseComponent) FindChild(id string) Component {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, child := range b.children {
		if child.ID() == id {
			return child
		}
		if result := child.FindChild(id); result != nil {
			return result
		}
	}
	return nil
}

// 向后兼容方法
func (b *BaseComponent) GetChildByID(id string) Component {
	return b.FindChild(id)
}

// 内部边界计算方法
func (b *BaseComponent) getLocalBounds() image.Rectangle {
	if b.node == nil {
		return image.Rectangle{}
	}

	layout := b.node.GetComputedLayout()
	return image.Rect(
		int(layout.Left), int(layout.Top),
		int(layout.Left+layout.Width), int(layout.Top+layout.Height),
	)
}

func (b *BaseComponent) getGlobalBounds() image.Rectangle {
	localBounds := b.getLocalBounds()
	if b.parent == nil {
		return localBounds
	}

	parentBounds := b.parent.Bounds()
	return image.Rect(
		localBounds.Min.X+parentBounds.Min.X,
		localBounds.Min.Y+parentBounds.Min.Y,
		localBounds.Max.X+parentBounds.Min.X,
		localBounds.Max.Y+parentBounds.Min.Y,
	)
}

// 向后兼容方法
func (b *BaseComponent) GetLocalBounds() image.Rectangle {
	return b.getLocalBounds()
}

func (b *BaseComponent) GetGlobalBounds() image.Rectangle {
	return b.getGlobalBounds()
}