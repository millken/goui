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

	// 输入状态
	hovered bool
	pressed bool

	// 脏标记
	dirty bool
}

// NewBase 创建一个带 Yoga 节点的基础组件
func NewBase() Base {
	n := yoga.NewNode()
	return Base{node: n, visible: true, dirty: true}
}

// Layout 触发布局计算
func (b *Base) Layout(width, height float32, direction yoga.Direction) {
	if b.node == nil {
		return
	}
	b.node.CalculateLayout(width, height, direction)
}

// Paint 默认不绘制，由具体组件实现。这里保留空实现，保证递归不报错。
func (b *Base) Paint(canvas *thorvg.Canvas) {}

// GetNode 返回 Yoga 节点
func (b *Base) GetNode() *yoga.Node { return b.node }

// GetComputedLayout 返回 Yoga 计算结果
func (b *Base) GetComputedLayout() yoga.Layout {
	if b.node == nil {
		return yoga.Layout{}
	}
	return b.node.GetComputedLayout()
}

// SetVisible 设置组件可见性
func (b *Base) SetVisible(visible bool) { b.visible = visible }

// IsVisible 返回组件是否可见
func (b *Base) IsVisible() bool { return b.visible }

// SetID 设置组件 ID
func (b *Base) SetID(id string) { b.id = id }

// GetID 返回组件 ID
func (b *Base) GetID() string { return b.id }

// AddChild 添加子组件，并同步 Yoga 树
func (b *Base) AddChild(child goui.Component) {
	b.children = append(b.children, child)
	if b.node != nil && child != nil && child.GetNode() != nil {
		b.node.InsertChild(child.GetNode(), uint32(len(b.children)-1))
	}
	if sp, ok := child.(interface{ SetParent(*Base) }); ok {
		sp.SetParent(b)
	}
}

// RemoveChild 移除子组件，并同步 Yoga 树
func (b *Base) RemoveChild(child goui.Component) {
	for i, c := range b.children {
		if c == child {
			if b.node != nil && child.GetNode() != nil {
				b.node.RemoveChild(child.GetNode())
			}
			if sp, ok := child.(interface{ SetParent(*Base) }); ok {
				sp.SetParent(nil)
			}
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
}

// GetChildren 返回子组件列表
func (b *Base) GetChildren() []goui.Component { return b.children }

// Mount 触发挂载回调
func (b *Base) Mount() {
	if b.events.OnMount != nil {
		b.events.OnMount()
	}
}

// Unmount 触发卸载回调
func (b *Base) Unmount() {
	if b.events.OnUnmount != nil {
		b.events.OnUnmount()
	}
}

// Update 触发更新回调
func (b *Base) Update() {
	if b.events.OnUpdate != nil {
		b.events.OnUpdate()
	}
}

// Destroy 供资源销毁
func (b *Base) Destroy() {}

// HandleEvent 默认处理 Hover/Click，返回是否拦截
func (b *Base) HandleEvent(event goui.Event) bool {
	switch e := event.(type) {
	case goui.MouseEvent:
		inside := b.pointInside(e.X, e.Y)
		if inside != b.hovered {
			b.hovered = inside
			// 悬停状态变化时标记为脏
			b.MarkDirty()
			if b.events.OnHover != nil {
				b.events.OnHover(b.hovered)
			}
		}
		if e.Down {
			if inside {
				b.pressed = true
				return true
			}
		} else {
			if b.pressed {
				if inside && b.events.OnClick != nil {
					b.events.OnClick()
				}
				b.pressed = false
				return true
			}
		}
		return false
	case goui.ScrollEvent:
		// 可以添加滚轮处理逻辑，例如滚动容器
		return false
	case goui.FocusEvent:
		// 可以添加焦点处理逻辑
		return false
	case goui.KeyEvent:
		// 可以添加键盘处理逻辑，例如快捷键
		return false
	default:
		return false
	}
}

// SetOnClick 设置点击回调
func (b *Base) SetOnClick(handler func()) { b.events.OnClick = handler }

// SetOnHover 设置悬停回调
func (b *Base) SetOnHover(handler func(bool)) { b.events.OnHover = handler }

// SetParent 设置父组件（内部使用）
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

// pointInside 判断窗口坐标 (x,y) 是否在该组件边界内
func (b *Base) pointInside(x, y float32) bool {
	l := b.GetComputedLayout()
	ax, ay := b.absoluteOffset()
	left, top := ax, ay
	right, bottom := left+l.Width, top+l.Height
	return x >= left && x <= right && y >= top && y <= bottom
}

// SetStyle 应用并存储样式到 Yoga 节点
func (b *Base) SetStyle(s goui.Style) {
	b.style = s
	if n := b.GetNode(); n != nil {
		goui.ApplyStyle(n, s)
	}
	// 样式变更时标记为脏
	b.MarkDirty()
}

// GetStyle 返回已存储的样式（快照）
func (b *Base) GetStyle() goui.Style { return b.style }

// MarkDirty 标记组件为脏，需要重新渲染
func (b *Base) MarkDirty() {
	b.dirty = true
	// 递归标记父组件也需要重新布局
	if b.parent != nil {
		b.parent.MarkDirty()
	}
}

// IsDirty 检查组件是否需要重新渲染
func (b *Base) IsDirty() bool { return b.dirty }

// ClearDirty 清除脏标记
func (b *Base) ClearDirty() { b.dirty = false }
