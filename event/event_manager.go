// Package event 提供事件类型与事件管理器。
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
	listeners  map[string]map[Type][]EventListener
	filters    []EventFilter
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
		config:    config,
		listeners: make(map[string]map[Type][]EventListener),
		// 长度为0，仅预留容量，避免 append 与 len 判断冲突
		filters:    make([]EventFilter, 0, config.MaxFilters),
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
		// slog.Debug("Processing event", "type", event.Type(), "event", event)
		// 处理事件
		em.handleEvent(event)
	}
}

func (em *EventManager) handleEvent(event Event) {
	// 应用过滤器
	if !em.applyFilters(event) {
		return // 事件被过滤
	}
	// 根据事件类型分发到所有已注册的监听器
	em.mu.RLock()
	var batch []EventListener
	for _, byType := range em.listeners {
		if ls, ok := byType[event.Type()]; ok && len(ls) > 0 {
			// 复制引用到本地切片，避免长时间持锁
			batch = append(batch, ls...)
		}
	}
	// slog.Debug("Dispatching event", "type", event.Type(), "em.listeners", len(em.listeners), "listeners", len(batch))
	em.mu.RUnlock()
	for _, l := range batch {
		if l.HandleEvent(event) {
			// 消费后可按需停止传播；这里选择停止，避免重复触发
			break
		}
	}
}

// applyFilters 应用事件过滤器
func (em *EventManager) applyFilters(event Event) bool {
	em.mu.RLock()
	defer em.mu.RUnlock()
	if len(em.filters) == 0 {
		return true
	}
	filters := make([]EventFilter, len(em.filters))
	copy(filters, em.filters)

	for _, filter := range filters {
		if !filter.FilterEvent(event) {
			return false
		}
	}
	return true
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

// RegisterListener 注册事件监听器
func (em *EventManager) RegisterListener(componentID string, eventType Type, listener EventListener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.listeners[componentID]; !exists {
		em.listeners[componentID] = make(map[Type][]EventListener)
	}
	em.listeners[componentID][eventType] = append(em.listeners[componentID][eventType], listener)
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
