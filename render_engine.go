package goui

import (
	"context"
	"fmt"
	"image"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/millken/go/thorvg"
)

// RenderEngine 渲染引擎（借鉴TpGUI精髓）
type RenderEngine struct {
	// 核心组件
	canvas       *thorvg.Canvas
	dirtyMgr     *DirtyRectManager
	updateQueue  chan UpdateCommand
	renderLayers map[RenderLayer][]Renderable

	// 事件系统（借鉴TpGUI信号槽机制）
	eventQueue   chan Event
	eventHandlers map[EventType][]func(Event)

	// 多线程渲染（借鉴TpGUI多线程设计）
	renderWorkers []*RenderWorker
	workerCount   int
	renderWg      sync.WaitGroup

	// 渲染控制
	mu           sync.RWMutex
	running      int32  // atomic
	stopChan     chan struct{}
	renderCtx    context.Context
	renderCancel context.CancelFunc

	// 帧率控制（借鉴TpGUI VSync）
	targetFPS    int
	frameInterval time.Duration
	lastFrameTime time.Time

	// 性能统计（增强版）
	renderCount  int64  // atomic
	totalFrames  int64  // atomic
	fps          int32  // atomic
	drawCalls    int64  // atomic
	lastStatsTime time.Time

	// 配置
	config       *RenderConfig
}

// RenderConfig 渲染配置
type RenderConfig struct {
	TargetFPS         int
	EnableVSync       bool
	MaxDirtyRects     int
	WorkerCount       int  // 0表示自动检测
	EnableStats       bool
	EnableMultithread bool
	EventQueueSize    int
	UpdateQueueSize   int
}

// RenderWorker 渲染工作线程（借鉴TpGUI多线程渲染）
type RenderWorker struct {
	id         int
	canvas     *thorvg.Canvas
	taskChan   chan RenderTask
	engine     *RenderEngine
	running    int32  // atomic
}

// RenderTask 渲染任务
type RenderTask struct {
	Renderable Renderable
	Region     image.Rectangle
	Layer      RenderLayer
	Priority   int
	Callback   func()
}

// NewRenderEngine 创建渲染引擎
func NewRenderEngine(dirtyRectsLimit int) *RenderEngine {
	return NewRenderEngineWithConfig(&RenderConfig{
		TargetFPS:         60,
		EnableVSync:       true,
		MaxDirtyRects:     dirtyRectsLimit,
		WorkerCount:       0, // 自动检测
		EnableStats:       true,
		EnableMultithread: true,
		EventQueueSize:    1000,
		UpdateQueueSize:   1000,
	})
}

// NewRenderEngineWithConfig 使用配置创建渲染引擎
func NewRenderEngineWithConfig(config *RenderConfig) *RenderEngine {
	if config == nil {
		config = &RenderConfig{
			TargetFPS:         60,
			EnableVSync:       true,
			MaxDirtyRects:     10,
			WorkerCount:       0,
			EnableStats:       true,
			EnableMultithread: true,
			EventQueueSize:    1000,
			UpdateQueueSize:   1000,
		}
	}

	// 自动检测线程数
	if config.WorkerCount == 0 {
		if config.EnableMultithread {
			config.WorkerCount = max(1, runtime.NumCPU()/2)
		} else {
			config.WorkerCount = 1
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	engine := &RenderEngine{
		dirtyMgr:       NewDirtyRectManager(config.MaxDirtyRects),
		updateQueue:    make(chan UpdateCommand, config.UpdateQueueSize),
		renderLayers:   make(map[RenderLayer][]Renderable),
		eventQueue:     make(chan Event, config.EventQueueSize),
		eventHandlers:  make(map[EventType][]func(Event)),
		stopChan:       make(chan struct{}),
		renderCtx:      ctx,
		renderCancel:   cancel,
		targetFPS:      config.TargetFPS,
		frameInterval:  time.Second / time.Duration(config.TargetFPS),
		workerCount:    config.WorkerCount,
		config:         config,
		lastStatsTime:  time.Now(),
	}

	// 初始化ThorVG引擎
	if err := thorvg.Init(uint(config.WorkerCount)); err != nil {
		// 降级到单线程模式
		config.WorkerCount = 1
		config.EnableMultithread = false
		fmt.Printf("Warning: ThorVG初始化失败，降级到单线程模式: %v\n", err)
	}

	// 创建渲染工作线程
	if config.EnableMultithread {
		engine.renderWorkers = make([]*RenderWorker, config.WorkerCount)
		for i := 0; i < config.WorkerCount; i++ {
			// 暂时使用nil，稍后在Start时创建
			engine.renderWorkers[i] = &RenderWorker{
				id:       i,
				canvas:   nil,
				taskChan: make(chan RenderTask, 50),
				engine:   engine,
			}
		}
	}

	return engine
}

// Start 启动渲染引擎
func (re *RenderEngine) Start(canvas *thorvg.Canvas) {
	if !atomic.CompareAndSwapInt32(&re.running, 0, 1) {
		return // 已经在运行
	}

	re.canvas = canvas
	re.lastFrameTime = time.Now()

	// 暂时禁用多线程渲染，先让基本功能工作
	/*
	if re.config.EnableMultithread {
		for _, worker := range re.renderWorkers {
			// 创建SwCanvas用于工作线程
			if worker.canvas == nil {
				worker.canvas = thorvg.NewSwCanvas(thorvg.EngineOptionNone)
			}
			atomic.StoreInt32(&worker.running, 1)
			re.renderWg.Add(1)
			go worker.run()
		}
	}
	*/

	// 启动主渲染循环（借鉴TpGUI主循环）
	go re.processUpdates()

	// 启动事件处理循环（借鉴TpGUI信号槽机制）
	go re.processEvents()

	// 启动统计更新循环
	if re.config.EnableStats {
		go re.updateStats()
	}
}

// Stop 停止渲染引擎
func (re *RenderEngine) Stop() {
	if !atomic.CompareAndSwapInt32(&re.running, 1, 0) {
		return // 已经停止
	}

	// 取消上下文
	re.renderCancel()

	// 停止工作线程
	if re.config.EnableMultithread {
		for _, worker := range re.renderWorkers {
			atomic.StoreInt32(&worker.running, 0)
			close(worker.taskChan)
		}
	}

	// 等待所有工作线程结束
	re.renderWg.Wait()

	// 关闭队列
	close(re.stopChan)
	close(re.updateQueue)
	close(re.eventQueue)

	// 清理ThorVG资源
	thorvg.Term()
}

// ScheduleUpdate 调度更新（异步，非阻塞）
func (re *RenderEngine) ScheduleUpdate(cmd UpdateCommand) {
	if atomic.LoadInt32(&re.running) == 0 {
		return
	}

	cmd.Timestamp = time.Now()

	select {
	case re.updateQueue <- cmd:
	default:
		// 队列满时丢弃低优先级更新
		if cmd.Priority > 0 {
			select {
			case re.updateQueue <- cmd:
			default:
			}
		}
	}
}

// RegisterRenderable 注册渲染对象
func (re *RenderEngine) RegisterRenderable(layer RenderLayer, renderable Renderable) {
	re.mu.Lock()
	defer re.mu.Unlock()

	re.renderLayers[layer] = append(re.renderLayers[layer], renderable)
}

// UnregisterRenderable 注销渲染对象
func (re *RenderEngine) UnregisterRenderable(layer RenderLayer, renderable Renderable) {
	re.mu.Lock()
	defer re.mu.Unlock()

	layerItems := re.renderLayers[layer]
	for i, item := range layerItems {
		if item == renderable {
			re.renderLayers[layer] = append(layerItems[:i], layerItems[i+1:]...)
			break
		}
	}
}

// processUpdates 异步处理更新命令（借鉴TpGUI主循环设计）
func (re *RenderEngine) processUpdates() {
	// 使用高精度定时器实现VSync
	frameTicker := time.NewTicker(re.frameInterval)
	defer frameTicker.Stop()

	batchCommands := make([]UpdateCommand, 0, 100)

	for {
		select {
		case <-re.renderCtx.Done():
			return

		// 处理更新命令（类似TpGUI的updateTasks_）
		case cmd := <-re.updateQueue:
			batchCommands = append(batchCommands, cmd)

			// 批量处理，避免频繁的小批量更新
			if len(batchCommands) >= 100 {
				re.processBatchUpdates(batchCommands)
				batchCommands = batchCommands[:0]
			}

		// 帧率控制渲染（借鉴TpGUI的16ms间隔）
		case <-frameTicker.C:
			frameStart := time.Now()

			// 1. 处理剩余的批量更新
			if len(batchCommands) > 0 {
				re.processBatchUpdates(batchCommands)
				batchCommands = batchCommands[:0]
			}

			// 2. 渲染一帧
			re.renderFrame()

			// 3. VSync控制
			if re.config.EnableVSync {
				re.waitForVSync(frameStart)
			}
		}
	}
}

// processBatchUpdates 批量处理更新命令（借鉴TpGUI合并策略）
func (re *RenderEngine) processBatchUpdates(commands []UpdateCommand) {
	if len(commands) == 0 {
		return
	}

	re.mu.Lock()
	defer re.mu.Unlock()

	// 按组件分组，合并同一组件的多个更新（TpGUI精髓）
	componentUpdates := make(map[Component][]UpdateCommand)
	for _, cmd := range commands {
		if cmd.Target != nil {
			componentUpdates[cmd.Target] = append(componentUpdates[cmd.Target], cmd)
		}
	}

	// 处理每个组件的更新
	for component, updates := range componentUpdates {
		// 合并更新区域
		mergedRegion := image.Rectangle{}
		hasFullUpdate := false

		for _, update := range updates {
			if update.Mode == RenderModeFull {
				hasFullUpdate = true
				break
			}

			if mergedRegion.Empty() {
				mergedRegion = update.Region
			} else {
				mergedRegion = mergedRegion.Union(update.Region)
			}
		}

		// 应用合并后的更新
		if hasFullUpdate {
			re.dirtyMgr.AddDirtyRegion(image.Rect(0, 0, 9999, 9999))
		} else if !mergedRegion.Empty() {
			re.dirtyMgr.AddDirtyRegion(mergedRegion)
		}

		// 标记组件需要重新渲染
		if baseComp, ok := component.(*BaseComponent); ok {
			baseComp.markNeedsRender()
		}
	}
}

// processEvents 处理事件队列（借鉴TpGUI信号槽机制）
func (re *RenderEngine) processEvents() {
	for {
		select {
		case <-re.renderCtx.Done():
			return
		case event := <-re.eventQueue:
			re.dispatchEvent(event)
		}
	}
}

// dispatchEvent 分发事件
func (re *RenderEngine) dispatchEvent(event Event) {
	re.mu.RLock()
	handlers := re.eventHandlers[event.Type()]
	re.mu.RUnlock()

	// 调用所有注册的事件处理器
	for _, handler := range handlers {
		go handler(event) // 异步处理，避免阻塞
	}
}

// waitForVSync 等待垂直同步（借鉴TpGUI VSync控制）
func (re *RenderEngine) waitForVSync(frameStart time.Time) {
	elapsed := time.Since(frameStart)
	if elapsed < re.frameInterval {
		time.Sleep(re.frameInterval - elapsed)
	}
}

// renderFrame 渲染一帧（增强版，支持多线程）
func (re *RenderEngine) renderFrame() {
	frameStart := time.Now()
	dirtyRegions := re.dirtyMgr.GetDirtyRegions()
	if len(dirtyRegions) == 0 {
		return
	}

	atomic.AddInt64(&re.totalFrames, 1)

	re.mu.RLock()
	renderLayers := make(map[RenderLayer][]Renderable)
	for k, v := range re.renderLayers {
		renderLayers[k] = append([]Renderable{}, v...)
	}
	re.mu.RUnlock()

	// 多线程渲染或单线程渲染
	if re.config.EnableMultithread && len(re.renderWorkers) > 0 {
		re.renderFrameMultithread(dirtyRegions, renderLayers)
	} else {
		re.renderFrameSingleThread(dirtyRegions, renderLayers)
	}

	// 同步ThorVG画布
	re.canvas.Sync()

	// 更新性能统计
	if re.config.EnableStats {
		frameTime := time.Since(frameStart)
		re.updateFrameStats(frameTime)
	}
}

// renderFrameMultithread 多线程渲染帧
func (re *RenderEngine) renderFrameMultithread(dirtyRegions []image.Rectangle, renderLayers map[RenderLayer][]Renderable) {
	// 创建渲染任务
	tasks := make([]RenderTask, 0, len(dirtyRegions)*len(renderLayers))

	// 按层级和区域创建任务
	for layer, renderables := range renderLayers {
		for _, renderable := range renderables {
			for _, region := range dirtyRegions {
				if renderable.ShouldRender(region) {
					task := RenderTask{
						Renderable: renderable,
						Region:     region,
						Layer:      layer,
						Priority:   int(layer),
					}
					tasks = append(tasks, task)
				}
			}
		}
	}

	// 分发任务到工作线程（负载均衡）
	for i, task := range tasks {
		workerIndex := i % re.workerCount
		select {
		case re.renderWorkers[workerIndex].taskChan <- task:
		default:
			// 队列满时在主线程渲染
			re.renderTask(task)
		}
	}

	// 等待所有任务完成
	re.renderWg.Wait()
}

// renderFrameSingleThread 单线程渲染帧
func (re *RenderEngine) renderFrameSingleThread(dirtyRegions []image.Rectangle, renderLayers map[RenderLayer][]Renderable) {
	// 按层级顺序渲染
	for layer := LayerBackground; layer <= LayerTopMost; layer++ {
		renderables := renderLayers[layer]
		for _, renderable := range renderables {
			for _, region := range dirtyRegions {
				if renderable.ShouldRender(region) {
					renderable.Render(re.canvas, region)
					atomic.AddInt64(&re.drawCalls, 1)
				}
			}
		}
	}
}

// renderTask 渲染单个任务
func (re *RenderEngine) renderTask(task RenderTask) {
	if task.Renderable.ShouldRender(task.Region) {
		task.Renderable.Render(re.canvas, task.Region)
		atomic.AddInt64(&re.drawCalls, 1)
	}

	if task.Callback != nil {
		task.Callback()
	}
}

// updateFrameStats 更新帧统计
func (re *RenderEngine) updateFrameStats(frameTime time.Duration) {
	// 计算FPS
	if time.Since(re.lastStatsTime) >= time.Second {
		frames := atomic.SwapInt64(&re.totalFrames, 0)
		atomic.StoreInt32(&re.fps, int32(frames))
		re.lastStatsTime = time.Now()
	}
}

// GetStats 获取渲染统计
func (re *RenderEngine) GetStats() (int64, float64) {
	return re.renderCount, float64(atomic.LoadInt32(&re.fps))
}

// GetDetailedStats 获取详细统计
func (re *RenderEngine) GetDetailedStats() RenderStats {
	return RenderStats{
		FrameCount: atomic.LoadInt64(&re.totalFrames),
		FPS:        float32(atomic.LoadInt32(&re.fps)),
		DrawCalls:  atomic.LoadInt64(&re.drawCalls),
		DirtyRects: re.dirtyMgr.GetDirtyRectCount(),
	}
}

// PostEvent 投递事件
func (re *RenderEngine) PostEvent(event Event) {
	if atomic.LoadInt32(&re.running) == 1 {
		select {
		case re.eventQueue <- event:
		default:
			// 事件队列满时丢弃
		}
	}
}

// RegisterEventHandler 注册事件处理器
func (re *RenderEngine) RegisterEventHandler(eventType EventType, handler func(Event)) {
	re.mu.Lock()
	defer re.mu.Unlock()
	re.eventHandlers[eventType] = append(re.eventHandlers[eventType], handler)
}

// IsRunning 检查引擎是否运行
func (re *RenderEngine) IsRunning() bool {
	return atomic.LoadInt32(&re.running) == 1
}

// RenderWorker.run 渲染工作线程主循环
func (w *RenderWorker) run() {
	defer w.engine.renderWg.Done()

	for {
		select {
		case <-w.engine.renderCtx.Done():
			return
		case task := <-w.taskChan:
			w.engine.renderTask(task)
		}
	}
}

// markNeedsRender 标记需要渲染（内部方法）
func (b *BaseComponent) markNeedsRender() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.needsRender = true
}

// updateStats 更新统计信息
func (re *RenderEngine) updateStats() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-re.renderCtx.Done():
			return
		case <-ticker.C:
			// 这里可以实现更复杂的统计逻辑
			if re.config.EnableStats {
				// 清理过期的统计数据
			}
		}
	}
}