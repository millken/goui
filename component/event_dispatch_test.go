package component

import (
	"testing"
	"time"

	"github.com/dnsoa/go/assert"
	"github.com/millken/goui/event"
	"github.com/millken/yoga"
)

// waitUntil waits up to n*1ms until cond returns true.
func waitUntil(t *testing.T, n int, cond func() bool) {
	t.Helper()
	for i := 0; i < n; i++ {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("condition not met after ~%d ms", n*5)
}

func newRootWithEM(w, h float32) (*BaseComponent, *event.EventManager) {
	root := NewBaseComponent("root")
	root.SetRoot(root)
	root.SetWidth(w)
	root.SetHeight(h)
	em := event.NewEventManager()
	root.SetEventManager(em)
	_ = em.Start()
	// compute initial layout
	root.Layout(w, h, yoga.DirectionLTR)
	if clearer, ok := any(root).(interface{ ClearLayoutDirtyRecursive() }); ok {
		clearer.ClearLayoutDirtyRecursive()
	}
	return root, em
}

func TestEventOrder_SameComponentRegistrationOrder(t *testing.T) {
	r := assert.New(t)
	root, em := newRootWithEM(200, 100)

	order := make([]string, 0, 2)
	root.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool { order = append(order, "l1"); return false }))
	root.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool { order = append(order, "l2"); return false }))

	_ = em.PostEvent(&event.MouseEvent{X: 10, Y: 10, Down: true, Count: 1})
	waitUntil(t, 100, func() bool { return len(order) == 2 })

	r.Equal(2, len(order))
	r.Equal("l1", order[0])
	r.Equal("l2", order[1])
}

func TestEventConsume_StopsFurtherListeners(t *testing.T) {
	r := assert.New(t)
	root, em := newRootWithEM(200, 100)

	called2 := 0
	root.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool { return true }))
	root.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool { called2++; return false }))

	_ = em.PostEvent(&event.MouseEvent{X: 10, Y: 10, Down: true, Count: 1})
	// give some time for dispatch
	time.Sleep(5 * time.Millisecond)

	r.Equal(0, called2)
}

func TestHoverDoesNotSwallowClick(t *testing.T) {
	r := assert.New(t)
	root, em := newRootWithEM(200, 100)

	// child component with real hover & click handlers
	btn := NewButton("btn")
	btn.SetWidth(80)
	btn.SetMinHeight(30)
	root.AddChild(btn)
	// recompute layout for child bounds
	root.Layout(200, 100, yoga.DirectionLTR)

	clickCount := 0
	btn.OnClick(func(me event.MouseEvent) { clickCount++ })

	// move inside -> only hover
	rb := btn.Bounds()
	insideX := float32(rb.Min.X + 1)
	insideY := float32(rb.Min.Y + 1)
	_ = em.PostEvent(&event.MouseEvent{X: insideX, Y: insideY, Count: 0})
	// ensure click not triggered by hover move
	time.Sleep(10 * time.Millisecond)
	r.Equal(0, clickCount)

	// click inside -> hover doesn't swallow; click counts
	_ = em.PostEvent(&event.MouseEvent{X: insideX, Y: insideY, Down: true, Count: 1})
	waitUntil(t, 100, func() bool { return clickCount == 1 })
	r.Equal(1, clickCount)
}
