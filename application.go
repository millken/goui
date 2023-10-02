package goui

import (
	"context"
)

type Application interface {
	Startup(_ context.Context) error
	Shutdown(_ context.Context) error
	//Window returns the root window of the application.
	Window() Window
}
