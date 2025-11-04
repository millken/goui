package component

import (
	"strconv"
	"sync/atomic"

	"github.com/millken/go/thorvg"
)

var componentIDCounter int64 = 0

func GenerateID() string {
	return "component:" + strconv.FormatInt(atomic.AddInt64(&componentIDCounter, 1), 10)
}

type Component interface {
	SetID(string)
	ID() string
	AddChild(Component)
	// RemoveChild(Component)

	Children() []Component

	Paint(*thorvg.Canvas) bool // 返回是否实际进行了渲染
	//MarkDirty 标记组件为脏，如果递归标记子组件则传入 true
	MarkDirty(bool)
}
