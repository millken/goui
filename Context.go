package goui

import (
	"time"

	"github.com/millken/goui/font"
)

type Context struct {
	Now       time.Time
	Delta     time.Duration
	frameNo   uint64
	winEvents []WindowEvent
	wm        *WindowManager
	sched     *Scheduler
	fm        *font.FontManager
}

// NewContext 创建新的渲染上下文
func NewContext() *Context { return &Context{} }

// GetScheduler returns the scheduler associated with this context.
func (c *Context) GetScheduler() *Scheduler {
	if c.sched == nil {
		c.sched = NewScheduler()
	}
	return c.sched
}

// PushWindowEvent 收集窗口事件（由 App 主循环调用）
func (c *Context) PushWindowEvent(e WindowEvent) { c.winEvents = append(c.winEvents, e) }

// WindowEvents 返回并清空当前累计的窗口事件
func (c *Context) WindowEvents() []WindowEvent {
	out := c.winEvents
	c.winEvents = c.winEvents[:0]
	return out
}

func (c *Context) tick(now time.Time) {
	if c.Now.IsZero() {
		c.Delta = 0
	} else {
		c.Delta = now.Sub(c.Now)
	}
	c.Now = now
	c.frameNo++
}

// SetWindowManager attaches a WindowManager to the context.
func (c *Context) SetWindowManager(wm *WindowManager) { c.wm = wm }

// GetWindowManager returns the attached WindowManager.
func (c *Context) GetWindowManager() *WindowManager { return c.wm }

// GetFontManager returns or creates the Context's FontManager
func (c *Context) GetFontManager() *font.FontManager {
	if c.fm == nil {
		c.fm = font.NewFontManager()
	}
	return c.fm
}
