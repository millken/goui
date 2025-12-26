package goui

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	"github.com/millken/goui/renderer"
	gouiwindow "github.com/millken/goui/window"
	"github.com/millken/thorvg"
)

// mockWin implements window.Window minimal surface for session loop tests.
type mockWin struct {
	w, h     int
	open     atomic.Bool
	events   chan event.Event
	vsyncSet atomic.Bool
}

func newMockWin(w, h int) *mockWin {
	m := &mockWin{w: w, h: h, events: make(chan event.Event, 8)}
	m.open.Store(true)
	return m
}

func (m *mockWin) Open() error                            { return nil }
func (m *mockWin) Close() error                           { m.open.Store(false); return nil }
func (m *mockWin) SetTitle(string)                        {}
func (m *mockWin) SetSize(w, h int)                       { m.w, m.h = w, h }
func (m *mockWin) GetSize() (int, int)                    { return m.w, m.h }
func (m *mockWin) SetResizable(bool)                      {}
func (m *mockWin) Fullscreen()                            {}
func (m *mockWin) SetMouseVisible(bool)                   {}
func (m *mockWin) SetCursor(_ gouiwindow.Cursor)          {}
func (m *mockWin) SetVSync(enabled bool)                  { m.vsyncSet.Store(enabled) }
func (m *mockWin) MakeContextCurrent()                    {}
func (m *mockWin) IsOpen() bool                           { return m.open.Load() }
func (m *mockWin) SwapBuffers()                           {}
func (m *mockWin) PollEvents()                            {}
func (m *mockWin) SetClipboardText(text string) error     { return nil }
func (m *mockWin) GetClipboardText() (string, error)      { return "", nil }
func (m *mockWin) WaitEvents(_ time.Duration) event.Event { return nil }
func (m *mockWin) FetchEvents(max int) []event.Event {
	var out []event.Event
	for {
		select {
		case ev := <-m.events:
			out = append(out, ev)
			if max > 0 && len(out) >= max {
				return out
			}
		default:
			return out
		}
	}
}

// fakeRenderer tracks call order to assert session loop behavior without GL.
type fakeRenderer struct {
	began int32
	ended int32
	reset int32
}

func (f *fakeRenderer) Cleanup() {}
func (f *fakeRenderer) BeginFrame(w, h int) (*thorvg.Canvas, error) {
	atomic.AddInt32(&f.began, 1)
	return nil, nil
}
func (f *fakeRenderer) EndFrame() error           { atomic.AddInt32(&f.ended, 1); return nil }
func (f *fakeRenderer) ResetScene(w, h int) error { atomic.AddInt32(&f.reset, 1); return nil }

// minimal component that always paints once after layout
type paintOnce struct{ *component.BaseComponent }

func newPaintOnce() *paintOnce                   { return &paintOnce{BaseComponent: component.NewBaseComponent()} }
func (p *paintOnce) Paint(_ *thorvg.Canvas) bool { return true }

func TestSessionLoop_EventLayoutRender(t *testing.T) {
	// arrange
	app := NewApplication()
	// inject fake renderer factory
	fr := &fakeRenderer{}
	app.SetRendererFactory(func() (renderer.Renderer, error) { return fr, nil })

	win := newMockWin(640, 480)
	root := newPaintOnce()
	sess, err := app.NewSessionWithAppConfig(win, root, &AppConfig{TargetFPS: 0})
	if err != nil {
		t.Fatalf("NewSession error: %v", err)
	}
	_ = sess

	// push a mouse move event to trigger OnHover path (even if no handler)
	win.events <- event.NewMouseEvent(5, 5, false, 0)
	// close soon to let one iteration run
	go func() {
		time.Sleep(10 * time.Millisecond)
		_ = win.Close()
	}()

	// act: run a single iteration by calling Run in a goroutine and stop quickly
	done := make(chan struct{})
	go func() { app.Run(); close(done) }()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("session loop did not complete")
	}

	// assert: renderer methods called at least once in expected sequence context
	if atomic.LoadInt32(&fr.began) == 0 {
		t.Errorf("expected BeginFrame to be called")
	}
	if atomic.LoadInt32(&fr.ended) == 0 {
		t.Errorf("expected EndFrame to be called")
	}
}
