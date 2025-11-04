// Package component provides UI component implementations for GOUI.
package component

import (
	"fmt"
	"image"
	"log/slog"
	"sync"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui/event"
	"github.com/millken/yoga"
)

type BaseComponent struct {
	root     Component
	id       string
	mu       sync.RWMutex
	yoga     *yoga.Node
	children []Component
	//events
	eventManager     *event.EventManager
	pendingListeners map[event.Type][]event.EventListener

	parent  Component
	visible bool // 可见性标记
	dirty   bool // 脏标记
}

func NewBaseComponent(id ...string) *BaseComponent {
	base := &BaseComponent{
		id: func() string {
			if len(id) > 0 {
				return id[0]
			} else {
				return GenerateID()
			}
		}(),
		yoga:             yoga.NewNode(),
		children:         []Component{},
		parent:           nil,
		visible:          true, // 默认可见
		dirty:            true, // 初始状态为需要渲染
		pendingListeners: make(map[event.Type][]event.EventListener),
	}
	return base
}

func (b *BaseComponent) SetRoot(c Component) {
	b.mu.Lock()
	b.root = c
	// 复制 children 与 pending 快照
	children := append([]Component(nil), b.children...)
	pending := make(map[event.Type][]event.EventListener, len(b.pendingListeners))
	for et, ls := range b.pendingListeners {
		if len(ls) == 0 {
			continue
		}
		cp := make([]event.EventListener, len(ls))
		copy(cp, ls)
		pending[et] = cp
	}
	// 清空缓存，避免重复
	b.pendingListeners = make(map[event.Type][]event.EventListener)
	b.mu.Unlock()

	// 若能从 root 上拿到 EventManager，则立即注册缓存的监听器
	if rc, ok := c.(RootComponent); ok && rc != nil {
		if em := rc.GetEventManager(); em != nil {
			for et, ls := range pending {
				for _, l := range ls {
					em.RegisterListener(b.ID(), et, l)
				}
			}
			// 清空已注册成功的本地快照
			pending = nil
		}
	}
	// 如果仍无法注册（例如 EM 未就绪），将快照放回缓存
	if len(pending) > 0 {
		b.mu.Lock()
		for et, ls := range pending {
			b.pendingListeners[et] = append(b.pendingListeners[et], ls...)
		}
		b.mu.Unlock()
	}

	// 将根组件的引用向下传播，避免已添加的子组件 Root 仍为 nil
	for _, child := range children {
		child.SetRoot(c)
	}
}

func (b *BaseComponent) Root() Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.root
}

func (b *BaseComponent) SetEventManager(em *event.EventManager) {
	// 设置自身 EM，并尽量刷新缓存，再向下传播
	b.mu.Lock()
	b.eventManager = em
	// 快照 children 与 pending，避免持锁调用外部
	children := append([]Component(nil), b.children...)
	pending := make(map[event.Type][]event.EventListener, len(b.pendingListeners))
	for et, ls := range b.pendingListeners {
		if len(ls) == 0 {
			continue
		}
		cp := make([]event.EventListener, len(ls))
		copy(cp, ls)
		pending[et] = cp
	}
	// 清空缓存，避免重复
	b.pendingListeners = make(map[event.Type][]event.EventListener)
	b.mu.Unlock()

	if em != nil {
		for et, ls := range pending {
			for _, l := range ls {
				em.RegisterListener(b.ID(), et, l)
			}
		}
	} else if len(pending) > 0 {
		// EM 仍为空，把快照放回缓存（极少见）
		b.mu.Lock()
		for et, ls := range pending {
			b.pendingListeners[et] = append(b.pendingListeners[et], ls...)
		}
		b.mu.Unlock()
	}

	// 将 EM 向下传播，促使子组件也能刷新缓存
	for _, child := range children {
		child.SetEventManager(em)
	}
}

func (b *BaseComponent) GetEventManager() *event.EventManager {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.eventManager
}

func (b *BaseComponent) SetID(id string) {
	b.id = id
}

func (b *BaseComponent) ID() string {
	return b.id
}

func (b *BaseComponent) IsVisible() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.visible
}

func (b *BaseComponent) SetVisible(visible bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.visible != visible {
		b.visible = visible
		b.MarkDirty(false) // 可见性改变时标记为脏
	}
}

func (b *BaseComponent) GetYogaNode() *yoga.Node {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.yoga
}

func (b *BaseComponent) Layout(width, height float32, direction yoga.Direction) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.yoga.CalculateLayout(width, height, direction)
}

func (b *BaseComponent) ComputedLayout() yoga.Layout {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.yoga.GetComputedLayout()
}

func (b *BaseComponent) AddChild(c Component) {
	b.children = append(b.children, c)
	// 若当前 Root 已设置则传递；否则先用自身作为临时 Root，后续 SetRoot 会向下传播
	if r := b.Root(); r != nil {
		c.SetRoot(r)
	} else {
		c.SetRoot(b)
	}
	c.SetParent(b)
	b.MarkDirty(false) // 添加子组件时标记为脏
}

func (b *BaseComponent) Children() []Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.children
}

func (b *BaseComponent) Paint(vg *thorvg.Canvas) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 递归渲染子组件
	childRendered := false
	for _, child := range b.children {
		if child.Paint(vg) {
			childRendered = true
		}
	}

	//如果不可见，跳过自己渲染
	if !b.visible {
		return childRendered
	}

	// 如果自己需要渲染，清理脏状态
	if b.dirty {
		b.dirty = false
		return true
	}

	return childRendered
}

// MarkDirty 标记组件为脏（需要重绘）
func (b *BaseComponent) MarkDirty(recursive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = true
	if recursive {
		for _, child := range b.children {
			child.MarkDirty(true)
		}
	}
}

// IsDirty 检查组件是否需要重绘
func (b *BaseComponent) IsDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

// SetParent 设置父组件
func (b *BaseComponent) SetParent(c Component) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.parent = c
}

// Parent 获取父组件
func (b *BaseComponent) Parent() Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.parent
}

// Bounds 返回组件的边界矩形
func (b *BaseComponent) Bounds() image.Rectangle {
	l := b.ComputedLayout()
	ax, ay := b.absoluteOffset()
	left, top := ax, ay
	right, bottom := left+l.Width, top+l.Height
	return image.Rectangle{
		Min: image.Point{X: int(left), Y: int(top)},
		Max: image.Point{X: int(right), Y: int(bottom)},
	}
}

// CheckDirtyRecursive 递归检查是否有任何子组件需要重绘（保留向后兼容）
func (b *BaseComponent) CheckDirtyRecursive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

// absoluteOffset 通过向上累加 Yoga 计算布局，得到绝对坐标
func (b *BaseComponent) absoluteOffset() (float32, float32) {
	var x, y float32
	var cur Component = b
	for cur != nil {
		l := cur.ComputedLayout()
		x += l.Left
		y += l.Top
		cur = cur.Parent()
	}
	return x, y
}

// PointInside 判断窗口坐标 (x,y) 是否在该组件边界内
func (b *BaseComponent) PointInside(x, y float32) bool {
	l := b.ComputedLayout()
	ax, ay := b.absoluteOffset()
	left, top := ax, ay
	right, bottom := left+l.Width, top+l.Height
	return x >= left && x <= right && y >= top && y <= bottom
}

func (b *BaseComponent) RegisterEventListener(eventType event.Type, listener event.EventListener) {
	slog.Debug("Registering event listener", "componentID", b.ID(), "eventType", eventType)

	root := b.Root()
	if root == nil {
		// Root 尚未就绪，暂存监听器，待 SetRoot/EM 可用时注册
		b.mu.Lock()
		b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
		b.mu.Unlock()
		return
	}
	if rootComp, ok := root.(RootComponent); ok {
		em := rootComp.GetEventManager()
		if em != nil {
			em.RegisterListener(b.ID(), eventType, listener)
			return
		}
		// EM 未就绪，暂存
		b.mu.Lock()
		b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
		b.mu.Unlock()
		return
	}
	// root 非 RootComponent，保底暂存
	b.mu.Lock()
	b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
	b.mu.Unlock()
}

func (b *BaseComponent) OnClick(handler func(event event.MouseEvent)) {
	b.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		if mouseEvent, ok := e.(*event.MouseEvent); ok && mouseEvent.Down {
			handler(*mouseEvent)
			return true
		}
		return false
	}))
}

func (b *BaseComponent) OnHover(handler func(hover bool)) {
	b.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		fmt.Println("ss")
		return true
	}))
}
