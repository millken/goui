// Package component provides UI component implementations for GOUI.
package component

import (
	"sync"

	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

type BaseComponent struct {
	mu       sync.RWMutex
	yoga     *yoga.Node
	children []Component
	dirty    bool // 脏标记
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{
		yoga:     yoga.NewNode(),
		children: []Component{},
		dirty:    true, // 初始状态为需要渲染
	}
}

func (b *BaseComponent) AddChild(c Component) {
	if b.children == nil {
		b.children = []Component{}
	}
	b.children = append(b.children, c)
	b.MarkDirty(false) // 添加子组件时标记为脏
}

func (b *BaseComponent) Children() []Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.children
}

func (b *BaseComponent) Render(vg *thorvg.Canvas) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 递归渲染子组件
	childRendered := false
	for _, child := range b.children {
		if child.Render(vg) {
			childRendered = true
		}
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

// CheckDirtyRecursive 递归检查是否有任何子组件需要重绘（保留向后兼容）
func (b *BaseComponent) CheckDirtyRecursive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}
