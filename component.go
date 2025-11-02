package goui

import "github.com/millken/yoga"

// Component 定义所有UI组件的基础接口
type Component interface {
	// GetNode 返回yoga节点
	GetNode() *yoga.Node

	// Layout 计算布局
	Layout(width, height float32, direction yoga.Direction)
	// 渲染
	Paint(canvas Canvas) // 绘制组件（原 Draw 改名，更贴近通用术语）

	// GetComputedLayout 获取计算后的布局信息
	GetComputedLayout() yoga.Layout

	// SetVisible 设置可见性
	SetVisible(visible bool)
	IsVisible() bool

	// SetID 设置组件ID
	SetID(id string)
	GetID() string

	// AddChild 添加子组件
	AddChild(child Component)
	RemoveChild(child Component)
	GetChildren() []Component

	HandleEvent(event Event) bool // 处理事件，返回是否已处理

	// 生命周期方法
	Mount()
	Unmount()
	Update()
	Destroy()

	// 事件处理
	SetOnClick(handler func())
	SetOnHover(handler func(bool))
}

// Style 定义组件样式
type Style struct {
	// 布局属性
	Width     *float32
	Height    *float32
	MinWidth  *float32
	MinHeight *float32
	MaxWidth  *float32
	MaxHeight *float32

	// Flexbox属性
	Flex           *float32
	FlexGrow       *float32
	FlexShrink     *float32
	FlexBasis      *float32
	FlexDirection  *yoga.FlexDirection
	JustifyContent *yoga.Justify
	AlignItems     *yoga.Align
	AlignSelf      *yoga.Align
	AlignContent   *yoga.Align
	FlexWrap       *yoga.Wrap

	// 位置属性
	Position *yoga.PositionType
	Left     *float32
	Top      *float32
	Right    *float32
	Bottom   *float32

	// 边距属性
	Margin       *float32
	MarginLeft   *float32
	MarginTop    *float32
	MarginRight  *float32
	MarginBottom *float32

	// 内边距属性
	Padding       *float32
	PaddingLeft   *float32
	PaddingTop    *float32
	PaddingRight  *float32
	PaddingBottom *float32

	// 边框属性
	Border       *float32
	BorderLeft   *float32
	BorderTop    *float32
	BorderRight  *float32
	BorderBottom *float32
	BorderRadius *float32
	BorderColor  *string

	// 文本属性
	Color    *string
	FontSize *float32
	Font     *string

	// 背景属性
	BackgroundColor *string

	// 其他属性
	AspectRatio *float32
	Overflow    *yoga.Overflow
	Display     *yoga.Display
	Gap         *float32

	// 交互属性
	Opacity       *float32
	Cursor        *string
	Transition    *string
	PointerEvents *string

	// CSS属性（用于调试和未来扩展）
	ZIndex    *int
	BoxShadow *string
	Transform *string
}

// Props 定义组件属性接口
type Props interface {
	GetStyle() Style
	GetChildren() []Component
}

// EventHandler 定义事件处理器类型
type EventHandler func()

// ComponentEvents 组件事件
type ComponentEvents struct {
	OnClick   EventHandler
	OnHover   func(bool)
	OnFocus   func()
	OnBlur    func()
	OnMount   func()
	OnUnmount func()
	OnUpdate  func()
}
