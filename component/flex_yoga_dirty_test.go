package component

import (
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/yoga"
)

// This test probes whether Yoga marks nodes dirty internally when style changes
// like SetFlexDirection are applied. If the binding exposes IsDirty(), we assert it.
// Otherwise we skip the introspection and only assert behavioral layout changes.
func TestYogaInternalDirty_OnFlexDirectionChange(t *testing.T) {
	r := assert.New(t)

	// Build a small tree: root -> flex (direction: row) with two fixed-size children
	root := NewBaseComponent("root")
	root.SetWidth(200)
	root.SetHeight(100)

	flex := NewFlex("flex")
	c1 := NewButton("c1")
	c1.SetWidth(40)
	c1.SetMinHeight(20)
	c2 := NewButton("c2")
	c2.SetWidth(40)
	c2.SetMinHeight(20)
	flex.AddChild(c1)
	flex.AddChild(c2)
	root.AddChild(flex)

	// initial layout (row)
	root.Layout(200, 100, yoga.DirectionLTR)
	if clearer, ok := any(root).(interface{ ClearLayoutDirtyRecursive() }); ok {
		clearer.ClearLayoutDirtyRecursive()
	}

	x1Row := c1.Bounds().Min.X
	y1Row := c1.Bounds().Min.Y
	x2Row := c2.Bounds().Min.X
	y2Row := c2.Bounds().Min.Y

	// change to column without touching custom layoutDirty
	flex.SetDirection(DirectionColumn)

	// If yoga.Node exposes IsDirty(), validate it's true after style change
	if yd, ok := any(flex.GetYogaNode()).(interface{ IsDirty() bool }); ok {
		r.True(yd.IsDirty())
	} else {
		t.Log("yoga.Node.IsDirty() not exposed; skipping direct dirty flag assertion")
	}

	// Recompute layout and verify positions changed as expected (stacked vertically)
	root.Layout(200, 100, yoga.DirectionLTR)
	if yd, ok := any(flex.GetYogaNode()).(interface{ IsDirty() bool }); ok {
		// After layout, yoga should be clean again
		r.False(yd.IsDirty())
	}

	x1Col := c1.Bounds().Min.X
	y1Col := c1.Bounds().Min.Y
	x2Col := c2.Bounds().Min.X
	y2Col := c2.Bounds().Min.Y

	// In row: y are same; x2 > x1. In column: x are same; y2 > y1.
	r.Equal(y1Row, y2Row)
	r.True(x2Row > x1Row)

	r.Equal(x1Col, x2Col)
	r.True(y2Col > y1Col)
}
