package v2

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type WindowManager struct {
	mu      sync.RWMutex
	windows []Window
	byID    map[string]Window
	roots   map[Window]Component
	events  chan WindowEvent
	idCtr   uint64
	dirty   map[Window]bool
	pump    func()
}

func NewWindowManager() *WindowManager {
	return &WindowManager{
		windows: make([]Window, 0, 4),
		byID:    make(map[string]Window),
		roots:   make(map[Window]Component),
		events:  make(chan WindowEvent, 128),
		dirty:   make(map[Window]bool),
	}
}

// NextID returns a unique window ID managed by this WindowManager.
func (wm *WindowManager) NextID() string { return fmt.Sprintf("w-%d", atomic.AddUint64(&wm.idCtr, 1)) }

// AddWindow registers a window (no id mapping).
func (wm *WindowManager) AddWindow(win Window) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.windows = append(wm.windows, win)
	wm.dirty[win] = true
}

// AddWindowWithID registers a window with a specific id.
func (wm *WindowManager) AddWindowWithID(id string, win Window) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.byID[id] = win
	wm.windows = append(wm.windows, win)
	wm.dirty[win] = true
	fmt.Println("WindowManager: AddWindowWithID id=", id)
}

// RemoveWindow unregisters a window and emits a WindowClosed event.
func (wm *WindowManager) RemoveWindow(win Window) {
	wm.mu.Lock()
	for i, w := range wm.windows {
		if w == win {
			wm.windows = append(wm.windows[:i], wm.windows[i+1:]...)
			break
		}
	}
	// remove from id map
	for id, w := range wm.byID {
		if w == win {
			delete(wm.byID, id)
			break
		}
	}
	delete(wm.roots, win)
	delete(wm.dirty, win)
	wm.mu.Unlock()
	wm.Emit(WindowEvent{Type: WindowClosed, Window: win})
	fmt.Println("WindowManager: RemoveWindow, remaining=", len(wm.windows))
}

// ListWindows returns a snapshot of windows.
func (wm *WindowManager) ListWindows() []Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	out := make([]Window, len(wm.windows))
	copy(out, wm.windows)
	return out
}

// GetWindow returns a window by id.
func (wm *WindowManager) GetWindow(id string) Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.byID[id]
}

// BindRoot binds a root component to the window.
func (wm *WindowManager) BindRoot(w Window, c Component) {
	wm.mu.Lock()
	wm.roots[w] = c
	wm.dirty[w] = true
	wm.mu.Unlock()
}

// GetRoot returns the root component bound to the window.
func (wm *WindowManager) GetRoot(w Window) Component {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.roots[w]
}

// Events returns the window event channel.
func (wm *WindowManager) Events() <-chan WindowEvent { return wm.events }

// Emit sends a window event non-blocking.
func (wm *WindowManager) Emit(ev WindowEvent) {
	select {
	case wm.events <- ev:
	default:
	}
}

// MarkDirty marks a window as needing redraw.
func (wm *WindowManager) MarkDirty(w Window) {
	wm.mu.Lock()
	wm.dirty[w] = true
	wm.mu.Unlock()
}

// PopDirtyWindows returns and clears the set of dirty windows.
func (wm *WindowManager) PopDirtyWindows() []Window {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if len(wm.dirty) == 0 {
		return nil
	}
	out := make([]Window, 0, len(wm.dirty))
	for w := range wm.dirty {
		out = append(out, w)
		delete(wm.dirty, w)
	}
	return out
}

// SetPump registers a platform event pump (e.g., glfw.PollEvents) to be executed by App main loop.
func (wm *WindowManager) SetPump(p func()) {
	wm.mu.Lock()
	wm.pump = p
	wm.mu.Unlock()
}

// CallPump invokes the registered pump if present.
func (wm *WindowManager) CallPump() {
	wm.mu.RLock()
	p := wm.pump
	wm.mu.RUnlock()
	if p != nil {
		p()
	}
}
