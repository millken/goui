package component

import "github.com/millken/go/thorvg"

type Component interface {
	AddChild(Component)
	// RemoveChild(Component)

	Children() []Component

	Render(*thorvg.Canvas) bool // 返回是否实际进行了渲染
	//MarkDirty 标记组件为脏，如果递归标记子组件则传入 true
	MarkDirty(bool)
}
