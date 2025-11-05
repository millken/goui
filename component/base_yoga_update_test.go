package component

import (
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/yoga"
)

func TestYogaUpdate_MarksDirtyAndLayoutDirty(t *testing.T) {
	r := assert.New(t)

	b := NewBaseComponent("b1")
	// BaseComponent 初始 dirty=true，这里先清理，方便观察 YogaUpdate 的效果
	b.ClearDirty()
	if cl, ok := any(b).(interface{ ClearLayoutDirty() }); ok {
		cl.ClearLayoutDirty()
	}

	// 通过 YogaUpdate 直接更新 Yoga 节点样式
	b.YogaUpdate(func(n *yoga.Node) {
		n.SetWidth(123)
	})

	r.True(b.IsDirty())
	if il, ok := any(b).(interface{ IsLayoutDirty() bool }); ok {
		r.True(il.IsLayoutDirty())
	}

	// 执行一次布局并清理布局脏标记（模拟 app loop）
	b.Layout(200, 100, yoga.DirectionLTR)
	if cl, ok := any(b).(interface{ ClearLayoutDirty() }); ok {
		cl.ClearLayoutDirty()
	}
	if il, ok := any(b).(interface{ IsLayoutDirty() bool }); ok {
		r.False(il.IsLayoutDirty())
	}
}
