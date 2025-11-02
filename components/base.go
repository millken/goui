package components

import (
	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/yoga"
)

// Base 提供 Component 的通用实现：Yoga 节点、可见性、ID、子节点、事件等
type Base struct {
	id       string
	visible  bool
	node     *yoga.Node
	parent   *Base
	children []goui.Component

	style  goui.Style
	events goui.ComponentEvents
}

// NewBase 创建一个带 Yoga 节点的基础组件
func NewBase() Base {
	n := yoga.NewNode()
	return Base{node: n, visible: true}
}

// ----- 布局 -----
func (b *Base) Layout(width, height float32, direction yoga.Direction) {
	if b.node == nil {
		return
	}
	// 根节点才会触发计算，其它节点随着树一并计算
	b.node.CalculateLayout(width, height, direction)
}

// GetNode 返回 Yoga 节点
func (b *Base) GetNode() *yoga.Node { return b.node }

// GetComputedLayout 返回 Yoga 计算结果（由 Yoga 提供的布局结构体）
func (b *Base) GetComputedLayout() yoga.Layout {
	if b.node == nil {
		return yoga.Layout{}
	}
	return b.node.GetComputedLayout()
}

// ----- 可见性 / ID -----
func (b *Base) SetVisible(visible bool) { b.visible = visible }
func (b *Base) IsVisible() bool         { return b.visible }

func (b *Base) SetID(id string) { b.id = id }
func (b *Base) GetID() string   { return b.id }

// ----- 子节点管理 -----
func (b *Base) AddChild(child goui.Component) {
	b.children = append(b.children, child)
	// 将 Yoga 节点也插入到树里
	if b.node != nil && child != nil && child.GetNode() != nil {
		b.node.InsertChild(child.GetNode(), uint32(len(b.children)-1))
	}
	// 设置父引用（若是 Base 派生）
	if cb, ok := child.(*Base); ok {
		cb.parent = b
	} else if fb, ok := child.(interface{ SetParent(*Base) }); ok {
		fb.SetParent(b)
	}
}

func (b *Base) RemoveChild(child goui.Component) {
	for i, c := range b.children {
		if c == child {
			// 从 Yoga 树删除
			if b.node != nil && child.GetNode() != nil {
				b.node.RemoveChild(child.GetNode())
			}
			// 断开父引用
			if cb, ok := child.(*Base); ok {
				if cb.parent == b {
					cb.parent = nil
				}
			}
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
}

func (b *Base) GetChildren() []goui.Component { return b.children }

// ----- 生命周期 -----
func (b *Base) Mount() {
	if b.events.OnMount != nil {
		b.events.OnMount()
	}
}
func (b *Base) Unmount() {
	if b.events.OnUnmount != nil {
		b.events.OnUnmount()
	}
}
func (b *Base) Update() {
	if b.events.OnUpdate != nil {
		b.events.OnUpdate()
	}
}
func (b *Base) Destroy() {}

// ----- 事件 -----
func (b *Base) HandleEvent(event goui.Event) bool {
	// 默认不处理，交给具体组件覆盖
	return false
}

func (b *Base) SetOnClick(handler func())     { b.events.OnClick = handler }
func (b *Base) SetOnHover(handler func(bool)) { b.events.OnHover = handler }

// Paint 默认不绘制，由具体组件实现。这里保留空实现，保证递归不报错。
func (b *Base) Paint(canvas *thorvg.Canvas) {}

// SetParent 供派生组件在 AddChild 时设置父引用（非导出，仅框架内部使用）
func (b *Base) SetParent(p *Base) { b.parent = p }

// absoluteOffset 通过向上累加 Yoga 计算布局，得到绝对坐标
func (b *Base) absoluteOffset() (float32, float32) {
	var x, y float32
	cur := b
	for cur != nil {
		l := cur.GetComputedLayout()
		x += l.Left
		y += l.Top
		cur = cur.parent
	}
	return x, y
}
