package component

import (
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/yoga"
)

func TestLayoutDirty_Cycle(t *testing.T) {
	r := assert.New(t)
	root := NewBaseComponent("root")
	root.SetWidth(300)
	root.SetHeight(200)

	btn := NewButton("b1")
	btn.SetWidth(80)
	btn.SetMinHeight(30)
	root.AddChild(btn)

	// initial layout
	root.Layout(300, 200, yoga.DirectionLTR)
	// in app loop, layoutDirty is cleared after layout; mirror that here
	if clearer, ok := any(root).(interface{ ClearLayoutDirtyRecursive() }); ok {
		clearer.ClearLayoutDirtyRecursive()
	}

	// initially should be clean
	if checker, ok := any(root).(interface{ CheckLayoutDirtyRecursive() bool }); ok {
		r.False(checker.CheckLayoutDirtyRecursive())
	}

	// change button intrinsic size -> marks layout dirty
	btn.SetText("Hello")
	if checker, ok := any(root).(interface{ CheckLayoutDirtyRecursive() bool }); ok {
		r.True(checker.CheckLayoutDirtyRecursive())
	}

	// do a layout pass and then clear recursively
	root.Layout(300, 200, yoga.DirectionLTR)
	if clearer, ok := any(root).(interface{ ClearLayoutDirtyRecursive() }); ok {
		clearer.ClearLayoutDirtyRecursive()
	}
	if checker, ok := any(root).(interface{ CheckLayoutDirtyRecursive() bool }); ok {
		r.False(checker.CheckLayoutDirtyRecursive())
	}
}
