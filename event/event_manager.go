package event

import (
	"errors"
	"sync"
	"sync/atomic"
)

// 错误定义
var (
	ErrEventManagerNotRunning = errors.New("event manager is not running")
	ErrEventQueueFull         = errors.New("event queue is full")
	ErrTooManyFilters         = errors.New("too many event filters")
)

type EventManager struct {
	// 事件队列
	eventQueue chan Event
	// 状态管理
	running int32
	mu      sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) Start() error {
	// 启动事件处理循环
	if !atomic.CompareAndSwapInt32(&em.running, 0, 1) {
		return nil // 已经运行
	}
	return nil
}

func (em *EventManager) Stop() {
	if !atomic.CompareAndSwapInt32(&em.running, 1, 0) {
		return // 已经停止
	}
}

// PostEvent 发送普通事件
func (em *EventManager) PostEvent(event Event) error {
	if atomic.LoadInt32(&em.running) == 0 {
		return ErrEventManagerNotRunning
	}

	select {
	case em.eventQueue <- event:
		return nil
	default:
		return ErrEventQueueFull
	}
}
