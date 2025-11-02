package goui

import (
	"context"
	"time"

	"github.com/millken/goui/font"
)

type Context struct {
	context.Context // 嵌入标准 context
	cancel          context.CancelFunc

	Now       time.Time
	Delta     time.Duration
	frameNo   uint64
	winEvents []WindowEvent
	fm        *font.FontManager
}

// NewContext 创建新的渲染上下文
func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		Context: ctx,
		cancel:  cancel,
	}
}

// Shutdown 优雅关闭上下文
func (c *Context) Shutdown() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Done 返回关闭通道
func (c *Context) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err 返回上下文错误
func (c *Context) Err() error {
	return c.Context.Err()
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

// GetFontManager returns or creates the Context's FontManager
func (c *Context) GetFontManager() *font.FontManager {
	if c.fm == nil {
		c.fm = font.NewFontManager()
	}
	return c.fm
}
