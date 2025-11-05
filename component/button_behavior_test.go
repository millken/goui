package component

import (
	"testing"
	"time"

	"github.com/dnsoa/go/assert"
	"github.com/millken/goui/event"
	gouiwindow "github.com/millken/goui/window"
	"github.com/millken/yoga"
)

// mockWindow is a minimal mock of window.Window used for behavior tests.
type mockWindow struct {
	w, h       int
	lastCursor gouiwindow.Cursor
}

func (m *mockWindow) Open() error                            { return nil }
func (m *mockWindow) Close() error                           { return nil }
func (m *mockWindow) SetTitle(string)                        {}
func (m *mockWindow) SetSize(w, h int)                       { m.w, m.h = w, h }
func (m *mockWindow) GetSize() (int, int)                    { return m.w, m.h }
func (m *mockWindow) SetResizable(bool)                      {}
func (m *mockWindow) Fullscreen()                            {}
func (m *mockWindow) SetMouseVisible(bool)                   {}
func (m *mockWindow) SetCursor(c gouiwindow.Cursor)          { m.lastCursor = c }
func (m *mockWindow) SetVSync(enabled bool)                  {}
func (m *mockWindow) IsOpen() bool                           { return true }
func (m *mockWindow) SwapBuffers()                           {}
func (m *mockWindow) PollEvents()                            {}
func (m *mockWindow) WaitEvents(_ time.Duration) event.Event { return nil }
func (m *mockWindow) FetchEvents(max int) []event.Event      { return nil }
func (m *mockWindow) MakeContextCurrent()                    {}

func TestButton_HoverAndActive_DoNotSwallowClick(t *testing.T) {
	r := assert.New(t)

	// Root & infra
	root := NewBaseComponent("root")
	// give root a concrete size
	root.SetWidth(800)
	root.SetHeight(600)

	btn := NewButton("btn1")
	btn.SetText("OK")
	// to avoid font measurement dependency for size, set explicit width/minHeight
	btn.SetWidth(80)
	btn.SetMinHeight(36)
	root.AddChild(btn)

	// wire EventManager and mock window via root
	em := event.NewEventManager()
	root.SetEventManager(em)
	mw := &mockWindow{w: 800, h: 600}
	root.SetWindow(mw)
	_ = em.Start()

	// compute initial layout so Bounds() are valid
	root.Layout(800, 600, yoga.DirectionLTR)

	// record click
	clicked := 0
	btn.OnClick(func(me event.MouseEvent) { clicked++ })

	// simulate hover enter -> should change cursor to pointer
	bounds := btn.Bounds()
	enterX := float32(bounds.Min.X + 1)
	enterY := float32(bounds.Min.Y + 1)
	_ = em.PostEvent(&event.MouseEvent{X: enterX, Y: enterY, Count: 0})

	// give a tiny async window for handler (eventLoop is async)
	withRetries(t, 100, func() bool { return mw.lastCursor == gouiwindow.CursorPointer })
	r.Equal(int(gouiwindow.CursorPointer), int(mw.lastCursor))

	// press inside -> pressed state and click recorded on Down
	_ = em.PostEvent(&event.MouseEvent{X: enterX, Y: enterY, Down: true, Count: 1})
	withRetries(t, 100, func() bool { return clicked == 1 })
	r.Equal(1, clicked)

	// release inside -> not strictly necessary for click count here
	_ = em.PostEvent(&event.MouseEvent{X: enterX, Y: enterY, Down: false, Count: 1})

	// hover leave -> cursor back to default
	leaveX := float32(bounds.Max.X + 10)
	leaveY := float32(bounds.Max.Y + 10)
	_ = em.PostEvent(&event.MouseEvent{X: leaveX, Y: leaveY, Count: 0})

	withRetries(t, 100, func() bool { return mw.lastCursor == gouiwindow.CursorDefault })
	r.Equal(int(gouiwindow.CursorDefault), int(mw.lastCursor))
}

// withRetries waits up to n*1ms for cond to become true.
func withRetries(t *testing.T, n int, cond func() bool) {
	t.Helper()
	for i := 0; i < n; i++ {
		if cond() {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	t.Fatalf("condition not met after %d ms", n)
}
