package app

import (
	"log"
	"log/slog"
	"sync/atomic"

	"gitee.com/millken/goui/ui"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	redrawSet = make(map[*ui.Window]struct{})
	sdlInited atomic.Bool
)

type startupOption struct { // This exists just to prevent arbitrary functions from being passed to application startup.
}

// StartupOption holds an option for application startup.
type StartupOption func(startupOption) error

// Start the application. This function does NOT return. While some calls may be safe to make, it should be assumed no
// calls into unison can be made prior to Start() being called unless explicitly stated otherwise.
func Start(options ...StartupOption) {
	//init sdl
	for _, option := range options {
		if err := option(startupOption{}); err != nil {
			log.Fatal(err)
		}
	}
	sdlInited.Store(true)
	//init app
	sdl.WaitEvent()
	slog.Debug("redrawSet", len(redrawSet))
	if len(redrawSet) > 0 {
		set := redrawSet
		redrawSet = make(map[*ui.Window]struct{})
		for wnd := range set {
			if wnd.IsVisible() {
				wnd.ReDraw()
			} else {
				redrawSet[wnd] = struct{}{}
			}
		}
	}
}
