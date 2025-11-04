package component

import (
	"image"
	"strconv"
	"sync/atomic"

	"github.com/millken/go/thorvg"
	"github.com/millken/yoga"
)

var componentIDCounter int64 = 0

func GenerateID() string {
	return "component:" + strconv.FormatInt(atomic.AddInt64(&componentIDCounter, 1), 10)
}

type Component interface {
	SetID(string)
	ID() string

	SetParent(Component)
	Parent() Component
	AddChild(Component)
	// RemoveChild(Component)

	IsVisible() bool
	SetVisible(bool)

	// Layout 计算布局
	Layout(width, height float32, direction yoga.Direction)
	// ComputedLayout 获取计算后的布局信息
	ComputedLayout() yoga.Layout
	//todo:remove
	Bounds() image.Rectangle

	Children() []Component

	Paint(*thorvg.Canvas) bool // 返回是否实际进行了渲染
	//MarkDirty 标记组件为脏，如果递归标记子组件则传入 true
	MarkDirty(bool)
}
