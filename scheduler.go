package goui

import (
	"runtime"
	"sync"
)

// Scheduler allows queuing tasks to be executed on the main thread (polled by App.Run).
type Scheduler struct {
	tasks  chan func()
	mu     sync.RWMutex
	isMain bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{tasks: make(chan func(), 256)}
}

// Run queues a task to be executed by the main thread.
func (s *Scheduler) Run(f func()) {
	if f == nil {
		return
	}
	s.tasks <- f
}

// RunWithResult queues a task and waits for its result.
func (s *Scheduler) RunWithResult(f func() any) any {
	if f == nil {
		return nil
	}
	// If caller is main thread, execute directly to avoid deadlock
	s.mu.RLock()
	main := s.isMain
	s.mu.RUnlock()
	if main {
		return f()
	}
	ch := make(chan any, 1)
	s.tasks <- func() { ch <- f() }
	return <-ch
}

// Poll executes all pending tasks without blocking.
func (s *Scheduler) Poll() {
	for {
		select {
		case t := <-s.tasks:
			if t != nil {
				t()
			}
		default:
			return
		}
	}
}

// EnsureMainThread optionally locks the OS thread for the caller; kept for compatibility.
func (s *Scheduler) EnsureMainThread() { runtime.LockOSThread() }

// SetMain marks the scheduler as running on the main thread. Call this from the main loop goroutine.
func (s *Scheduler) SetMain() {
	s.mu.Lock()
	s.isMain = true
	s.mu.Unlock()
}
