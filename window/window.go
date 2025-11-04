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
}
