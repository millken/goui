package goui

import (
	"image"
	"sync"
	"time"

	"github.com/millken/go/thorvg"
)

// RenderMode 渲染模式
type RenderMode int

const (
	RenderModePartial RenderMode = iota // 局部渲染
	RenderModeFull                     // 全量渲染
	RenderModeOnlyBlit                 // 仅位图传输
)

// UpdateCommand 更新命令（借鉴TpGUI设计）
type UpdateCommand struct {
	Target    Component
	Region    image.Rectangle // 更新区域
	Mode      RenderMode     // 渲染模式
	Timestamp time.Time      // 时间戳，用于合并优化
	Priority  int            // 优先级
}

// RenderLayer 渲染层级（借鉴TpGUI分层设计）
type RenderLayer int

const (
	LayerBackground RenderLayer = iota
	LayerContent
	LayerOverlay
	LayerTopMost
)

// Renderable 渲染接口
type Renderable interface {
	GetBounds() image.Rectangle
	GetLayer() RenderLayer
	ShouldRender(region image.Rectangle) bool
	Render(canvas *thorvg.Canvas, region image.Rectangle)
}

// DirtyRectManager 脏区域管理器（核心优化机制）
type DirtyRectManager struct {
	mu       sync.RWMutex
	rects    []image.Rectangle
	maxRects int
}

// RenderStats 渲染统计
type RenderStats struct {
	FrameCount int64
	FPS        float32
	DrawCalls  int64
	DirtyRects int
}

func NewDirtyRectManager(maxRects int) *DirtyRectManager {
	return &DirtyRectManager{
		rects:    make([]image.Rectangle, 0),
		maxRects: maxRects,
	}
}

// AddDirtyRegion 添加脏区域
func (d *DirtyRectManager) AddDirtyRegion(rect image.Rectangle) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if rect.Empty() {
		return
	}

	// 尝试与现有区域合并
	merged := false
	for i, existing := range d.rects {
		if mergedRect := unionIfClose(existing, rect); !mergedRect.Empty() {
			d.rects[i] = mergedRect
			merged = true
			break
		}
	}

	if !merged {
		d.rects = append(d.rects, rect)
	}

	// 限制区域数量，过多的区域合并为全屏
	if len(d.rects) > d.maxRects {
		d.rects = []image.Rectangle{image.Rect(0, 0, 9999, 9999)}
	}
}

// GetDirtyRegions 获取并清除脏区域
func (d *DirtyRectManager) GetDirtyRegions() []image.Rectangle {
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make([]image.Rectangle, len(d.rects))
	copy(result, d.rects)
	d.rects = d.rects[:0] // 清空
	return result
}

// GetDirtyRectCount 获取脏区域数量
func (d *DirtyRectManager) GetDirtyRectCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.rects)
}

// unionIfClose 如果两个矩形接近则合并
func unionIfClose(a, b image.Rectangle) image.Rectangle {
	const threshold = 10 // 阈值距离

	// 如果两个矩形在阈值范围内，则合并
	if (abs(a.Min.X-b.Min.X) <= threshold && abs(a.Min.Y-b.Min.Y) <= threshold) ||
		(a.Intersect(b).Dx() > threshold || a.Intersect(b).Dy() > threshold) {
		return a.Union(b)
	}
	return image.Rectangle{}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}