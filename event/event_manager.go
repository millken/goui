package event

import (
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

// 错误定义
var (
	ErrEventManagerNotRunning = errors.New("event manager is not running")
	ErrEventQueueFull         = errors.New("event queue is full")
	ErrTooManyFilters         = errors.New("too many event filters")
)

// EventManagerConfig 事件管理器配置
type EventManagerConfig struct {
	QueueSize       int
	PrioritySize    int
	MaxFilters      int
	EnableProfiling bool
}

type EventManager struct {
	config *EventManagerConfig
	// 事件队列
	eventQueue chan Event

	filters []EventFilter
	// 状态管理
	running int32
	mu      sync.RWMutex
}

func NewEventManager() *EventManager {
	config := &EventManagerConfig{
		QueueSize:       1000,
		PrioritySize:    100,
		MaxFilters:      100,
		EnableProfiling: false,
	}
	return &EventManager{
		filters:    make([]EventFilter, config.MaxFilters),
		eventQueue: make(chan Event, config.QueueSize), // 设置合理的缓冲区大小
	}
}

func (em *EventManager) Start() error {
	// 启动事件处理循环
	if !atomic.CompareAndSwapInt32(&em.running, 0, 1) {
		return nil // 已经运行
	}
	// 启动事件处理循环
	go em.eventLoop()
	return nil
}

func (em *EventManager) eventLoop() {
	for event := range em.eventQueue {
		slog.Debug("Processing event", "type", event.Type(), "event", event)
		// 处理事件
		em.handleEvent(event)
	}
}

func (em *EventManager) handleEvent(event Event) {
}

func (em *EventManager) Stop() {
	if !atomic.CompareAndSwapInt32(&em.running, 1, 0) {
		return // 已经停止
	}
	close(em.eventQueue) // 关闭 channel 以退出 eventLoop
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

// RegisterFilter 注册事件过滤器
func (em *EventManager) RegisterFilter(filter EventFilter) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.filters) >= em.config.MaxFilters {
		return ErrTooManyFilters
	}

	em.filters = append(em.filters, filter)
	return nil
}
