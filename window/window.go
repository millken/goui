// Package window 抽象平台窗口接口，供不同平台实现（例如 GLFW）
package window

import (
	"time"

	"github.com/millken/goui/event"
)

type Window interface {
	Open() error
	Close() error
	SetTitle(title string)
	SetSize(width, height int)
	GetSize() (width, height int)
	SetResizable(resizable bool)
	// SetDecorated(decorated bool)
	Fullscreen()
	SetMouseVisible(visible bool)
	IsOpen() bool

	SwapBuffers()
	PollEvents()
	WaitEvents(timeout time.Duration) event.Event
	// FetchEvents 尽可能在不阻塞的情况下批量获取事件，最多返回 max 个。
	// 如果 max<=0 则返回当前队列中所有可用事件。
	FetchEvents(max int) []event.Event
}
