package goui

// // EventManager 事件管理器
// type EventManager struct {
// 	// 事件队列
// 	eventQueue    chan Event
// 	priorityQueue chan PriorityEvent

// 	// 事件监听器
// 	listeners map[EventType][]EventListener
// 	filters   []EventFilter

// 	// 事件分发
// 	dispatchers map[EventType]EventDispatcher

// 	// 组件事件处理
// 	componentMap map[string]component.Component // ID -> Component
// 	hitTester    HitTester

// 	// 状态管理
// 	running int32
// 	mu      sync.RWMutex

// 	// 配置
// 	config *EventManagerConfig
// }

// // EventManagerConfig 事件管理器配置
// type EventManagerConfig struct {
// 	QueueSize       int
// 	PrioritySize    int
// 	MaxFilters      int
// 	EnableProfiling bool
// }

// // PriorityEvent 带优先级的事件
// type PriorityEvent struct {
// 	Event    Event
// 	Priority int
// 	Time     time.Time
// }

// // EventDispatcher 事件分发器接口
// type EventDispatcher interface {
// 	CanDispatch(event Event) bool
// 	Dispatch(event Event, target component.Component) bool
// }

// // HitTester 命中测试接口
// type HitTester interface {
// 	HitTest(x, y float32) component.Component
// }

// // NewEventManager 创建事件管理器
// func NewEventManager(platform Platform) *EventManager {
// 	config := &EventManagerConfig{
// 		QueueSize:       1000,
// 		PrioritySize:    100,
// 		MaxFilters:      10,
// 		EnableProfiling: false,
// 	}

// 	return NewEventManagerWithConfig(platform, config)
// }

// // NewEventManagerWithConfig 使用配置创建事件管理器
// func NewEventManagerWithConfig(platform Platform, config *EventManagerConfig) *EventManager {
// 	em := &EventManager{
// 		eventQueue:    make(chan Event, config.QueueSize),
// 		priorityQueue: make(chan PriorityEvent, config.PrioritySize),
// 		listeners:     make(map[EventType][]EventListener),
// 		filters:       make([]EventFilter, 0, config.MaxFilters),
// 		dispatchers:   make(map[EventType]EventDispatcher),
// 		componentMap:  make(map[string]Component),
// 		platform:      platform,
// 		config:        config,
// 	}

// 	// 注册默认分发器
// 	em.registerDefaultDispatchers()

// 	return em
// }

// // Start 启动事件管理器
// func (em *EventManager) Start() error {
// 	if !atomic.CompareAndSwapInt32(&em.running, 0, 1) {
// 		return nil // 已经运行
// 	}

// 	// 启动事件处理循环
// 	go em.eventLoop()
// 	go em.priorityEventLoop()

// 	return nil
// }

// // Stop 停止事件管理器
// func (em *EventManager) Stop() {
// 	if !atomic.CompareAndSwapInt32(&em.running, 1, 0) {
// 		return // 已经停止
// 	}

// 	// 关闭队列
// 	close(em.eventQueue)
// 	close(em.priorityQueue)
// }

// // PostEvent 发送普通事件
// func (em *EventManager) PostEvent(event Event) error {
// 	if atomic.LoadInt32(&em.running) == 0 {
// 		return ErrEventManagerNotRunning
// 	}

// 	select {
// 	case em.eventQueue <- event:
// 		return nil
// 	default:
// 		return ErrEventQueueFull
// 	}
// }

// // PostPriorityEvent 发送优先级事件
// func (em *EventManager) PostPriorityEvent(event Event, priority int) error {
// 	if atomic.LoadInt32(&em.running) == 0 {
// 		return ErrEventManagerNotRunning
// 	}

// 	pe := PriorityEvent{
// 		Event:    event,
// 		Priority: priority,
// 		Time:     time.Now(),
// 	}

// 	select {
// 	case em.priorityQueue <- pe:
// 		return nil
// 	default:
// 		return ErrEventQueueFull
// 	}
// }

// // RegisterListener 注册事件监听器
// func (em *EventManager) RegisterListener(eventType EventType, listener EventListener) {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	em.listeners[eventType] = append(em.listeners[eventType], listener)
// }

// // UnregisterListener 注销事件监听器
// func (em *EventManager) UnregisterListener(eventType EventType, listener EventListener) {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	listeners := em.listeners[eventType]
// 	for i, l := range listeners {
// 		if l == listener {
// 			em.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
// 			break
// 		}
// 	}
// }

// // RegisterFilter 注册事件过滤器
// func (em *EventManager) RegisterFilter(filter EventFilter) error {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	if len(em.filters) >= em.config.MaxFilters {
// 		return ErrTooManyFilters
// 	}

// 	em.filters = append(em.filters, filter)
// 	return nil
// }

// // RegisterComponent 注册组件（用于事件分发）
// func (em *EventManager) RegisterComponent(component Component) {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	em.componentMap[component.ID()] = component
// }

// // UnregisterComponent 注销组件
// func (em *EventManager) UnregisterComponent(component Component) {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	delete(em.componentMap, component.ID())
// }

// // SetHitTester 设置命中测试器
// func (em *EventManager) SetHitTester(tester HitTester) {
// 	em.mu.Lock()
// 	defer em.mu.Unlock()

// 	em.hitTester = tester
// }

// // eventLoop 事件主循环
// func (em *EventManager) eventLoop() {
// 	for {
// 		select {
// 		case event, ok := <-em.eventQueue:
// 			if !ok {
// 				return // 队列已关闭
// 			}
// 			em.processEvent(event)
// 		}
// 	}
// }

// // priorityEventLoop 优先级事件循环
// func (em *EventManager) priorityEventLoop() {
// 	for {
// 		select {
// 		case pe, ok := <-em.priorityQueue:
// 			if !ok {
// 				return // 队列已关闭
// 			}
// 			em.processPriorityEvent(pe)
// 		}
// 	}
// }

// // processEvent 处理普通事件
// func (em *EventManager) processEvent(event Event) {
// 	// 应用过滤器
// 	if !em.applyFilters(event) {
// 		return // 事件被过滤
// 	}

// 	// 确定事件目标
// 	target := em.determineEventTarget(event)
// 	if target == nil {
// 		return
// 	}

// 	event.SetTarget(target)

// 	// 使用分发器分发事件
// 	em.dispatchEvent(event, target)
// }

// // processPriorityEvent 处理优先级事件
// func (em *EventManager) processPriorityEvent(pe PriorityEvent) {
// 	// 优先级事件直接处理，不经过过滤器
// 	target := em.determineEventTarget(pe.Event)
// 	if target == nil {
// 		return
// 	}

// 	pe.Event.SetTarget(target)
// 	em.dispatchEvent(pe.Event, target)
// }

// // applyFilters 应用事件过滤器
// func (em *EventManager) applyFilters(event Event) bool {
// 	em.mu.RLock()
// 	filters := make([]EventFilter, len(em.filters))
// 	copy(filters, em.filters)
// 	em.mu.RUnlock()

// 	for _, filter := range filters {
// 		if !filter.FilterEvent(event) {
// 			return false
// 		}
// 	}
// 	return true
// }

// // determineEventTarget 确定事件目标
// func (em *EventManager) determineEventTarget(event Event) Component {
// 	// 如果事件已经有目标，直接返回
// 	if event.Target() != nil {
// 		return event.Target()
// 	}

// 	// 根据事件类型确定目标
// 	switch e := event.(type) {
// 	case *MouseEvent:
// 		return em.hitTest(e.X, e.Y)
// 	case *TouchEvent:
// 		return em.hitTest(e.X, e.Y)
// 	case *KeyEvent:
// 		// 键盘事件通常发送给当前焦点组件
// 		return em.getFocusedComponent()
// 	default:
// 		// 其他事件发送给根组件
// 		return em.getRootComponent()
// 	}
// }

// // dispatchEvent 分发事件
// func (em *EventManager) dispatchEvent(event Event, target Component) {
// 	eventType := event.Type()

// 	// 使用专用分发器
// 	if dispatcher, exists := em.dispatchers[eventType]; exists {
// 		if dispatcher.CanDispatch(event) {
// 			if dispatcher.Dispatch(event, target) {
// 				return // 事件被处理
// 			}
// 		}
// 	}

// 	// 默认分发：冒泡阶段
// 	em.bubbleEvent(event, target)

// 	// 通知监听器
// 	em.notifyListeners(event)
// }

// // bubbleEvent 事件冒泡
// func (em *EventManager) bubbleEvent(event Event, target Component) {
// 	current := target
// 	for current != nil {
// 		if handler, ok := current.(EventListener); ok {
// 			if handler.HandleEvent(event) {
// 				break // 事件被处理，停止冒泡
// 			}
// 		}
// 		current = current.Parent()
// 	}
// }

// // notifyListeners 通知事件监听器
// func (em *EventManager) notifyListeners(event Event) {
// 	em.mu.RLock()
// 	listeners := make([]EventListener, 0)
// 	if ls, exists := em.listeners[event.Type()]; exists {
// 		listeners = append(listeners, ls...)
// 	}
// 	em.mu.RUnlock()

// 	for _, listener := range listeners {
// 		go listener.HandleEvent(event) // 异步通知
// 	}
// }

// // hitTest 命中测试
// func (em *EventManager) hitTest(x, y float32) Component {
// 	em.mu.RLock()
// 	tester := em.hitTester
// 	em.mu.RUnlock()

// 	if tester != nil {
// 		return tester.HitTest(x, y)
// 	}

// 	// 默认实现：遍历所有组件
// 	em.mu.RLock()
// 	defer em.mu.RUnlock()

// 	for _, component := range em.componentMap {
// 		bounds := component.Bounds()
// 		point := image.Pt(int(x), int(y))
// 		if point.X >= bounds.Min.X && point.X < bounds.Max.X && point.Y >= bounds.Min.Y && point.Y < bounds.Max.Y {
// 			return component
// 		}
// 	}
// 	return nil
// }

// // getFocusedComponent 获取焦点组件
// func (em *EventManager) getFocusedComponent() Component {
// 	// TODO: 实现焦点管理
// 	return nil
// }

// // getRootComponent 获取根组件
// func (em *EventManager) getRootComponent() Component {
// 	em.mu.RLock()
// 	defer em.mu.RUnlock()

// 	for _, component := range em.componentMap {
// 		if component.Parent() == nil {
// 			return component
// 		}
// 	}
// 	return nil
// }

// // registerDefaultDispatchers 注册默认分发器
// func (em *EventManager) registerDefaultDispatchers() {
// 	// 鼠标事件分发器
// 	em.dispatchers[EventTypeMouse] = &MouseEventDispatcher{}

// 	// 键盘事件分发器
// 	em.dispatchers[EventTypeKey] = &KeyEventDispatcher{}

// 	// 窗口事件分发器
// 	em.dispatchers[EventTypeWindow] = &WindowEventDispatcher{}
// }

// // 错误定义
// var (
// 	ErrEventManagerNotRunning = NewError("event manager is not running")
// 	ErrEventQueueFull         = NewError("event queue is full")
// 	ErrTooManyFilters         = NewError("too many event filters")
// )
