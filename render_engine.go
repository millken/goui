package goui

import (
	"image"
	"sync"
	"time"

	"github.com/millken/go/thorvg"
)

// RenderEngine 渲染引擎（借鉴TpGUI异步渲染设计）
type RenderEngine struct {
	canvas       *thorvg.Canvas
	dirtyMgr     *DirtyRectManager
	updateQueue  chan UpdateCommand
	renderLayers map[RenderLayer][]Renderable
	mu           sync.RWMutex
	running      bool
	stopChan     chan struct{}

	// 性能统计
	renderCount int64
	lastRender  time.Time
	fps         float64
}

// NewRenderEngine 创建渲染引擎
func NewRenderEngine(dirtyRectsLimit int) *RenderEngine {
	return &RenderEngine{
		dirtyMgr:     NewDirtyRectManager(dirtyRectsLimit),
		updateQueue:  make(chan UpdateCommand, 1000),
		renderLayers: make(map[RenderLayer][]Renderable),
		stopChan:     make(chan struct{}),
	}
}

// Start 启动渲染引擎
func (re *RenderEngine) Start(canvas *thorvg.Canvas) {
	re.canvas = canvas
	re.running = true

	// 启动异步渲染处理
	go re.processUpdates()
}

// Stop 停止渲染引擎
func (re *RenderEngine) Stop() {
	re.running = false
	close(re.stopChan)
}

// ScheduleUpdate 调度更新（异步，非阻塞）
func (re *RenderEngine) ScheduleUpdate(cmd UpdateCommand) {
	if !re.running {
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

// processUpdates 异步处理更新命令
func (re *RenderEngine) processUpdates() {
	ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS
	defer ticker.Stop()

	batchCommands := make([]UpdateCommand, 0, 50)

	for {
		select {
		case <-re.stopChan:
			return

		case cmd := <-re.updateQueue:
			batchCommands = append(batchCommands, cmd)

			// 批量处理或定时处理
			if len(batchCommands) >= 50 {
				re.processBatch(batchCommands)
				batchCommands = batchCommands[:0]
			}

		case <-ticker.C:
			if len(batchCommands) > 0 {
				re.processBatch(batchCommands)
				batchCommands = batchCommands[:0]
			}
			re.renderFrame()
		}
	}
}

// processBatch 批量处理更新命令
func (re *RenderEngine) processBatch(commands []UpdateCommand) {
	re.mu.Lock()
	defer re.mu.Unlock()

	// 按层级和优先级排序
	sorted := make([]UpdateCommand, len(commands))
	copy(sorted, commands)

	// 简单排序：优先级高的先处理
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Priority > sorted[j].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 处理命令，合并脏区域
	for _, cmd := range sorted {
		if cmd.Target == nil {
			continue
		}

		// 根据渲染模式处理
		switch cmd.Mode {
		case RenderModeFull:
			re.dirtyMgr.AddDirtyRegion(image.Rect(0, 0, 9999, 9999))
		case RenderModePartial:
			re.dirtyMgr.AddDirtyRegion(cmd.Region)
		case RenderModeOnlyBlit:
			// 仅位图传输，不需要脏区域
		}
	}
}

// renderFrame 渲染一帧
func (re *RenderEngine) renderFrame() {
	dirtyRegions := re.dirtyMgr.GetDirtyRegions()
	if len(dirtyRegions) == 0 {
		return
	}

	re.mu.RLock()
	defer re.mu.RUnlock()

	// 按层级顺序渲染
	for layer := LayerBackground; layer <= LayerTopMost; layer++ {
		renderables := re.renderLayers[layer]
		for _, renderable := range renderables {
			for _, region := range dirtyRegions {
				if renderable.ShouldRender(region) {
					renderable.Render(re.canvas, region)
				}
			}
		}
	}

	// 更新性能统计
	now := time.Now()
	if !re.lastRender.IsZero() {
		duration := now.Sub(re.lastRender)
		if duration > 0 {
			re.fps = 1000.0 / duration.Seconds()
		}
	}
	re.lastRender = now
	re.renderCount++
}

// GetStats 获取渲染统计
func (re *RenderEngine) GetStats() (int64, float64) {
	return re.renderCount, re.fps
}