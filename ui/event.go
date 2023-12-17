package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	eventHandlers = make(map[Widget]func())
)

func ProcessEvents() {
	running := true
	for running {
		event := sdl.WaitEvent()
		switch event.(type) {
		case *sdl.QuitEvent:
			running = false
		default:
			dispatchEvent(event)
		}
	}
}

var mouseDownTime uint32

const longPressThreshold = 500

func dispatchEvent(event sdl.Event) {
	t := event.GetType()
	switch t {
	case sdl.MOUSEMOTION:
		e := event.(*sdl.MouseMotionEvent)
		if e.State == sdl.BUTTON_LEFT {
			fmt.Println("mouse drag", e.X, e.Y)
		}
	case sdl.MOUSEBUTTONDOWN:
		e := event.(*sdl.MouseButtonEvent)
		mouseDownTime = e.Timestamp
	case sdl.MOUSEBUTTONUP:
		e := event.(*sdl.MouseButtonEvent)
		if e.Timestamp-mouseDownTime > longPressThreshold {
			fmt.Println("mouse long press", e.X, e.Y, e.Timestamp-mouseDownTime)
		} else if e.Timestamp-mouseDownTime < 200 {
			fmt.Println("mouse click", e.X, e.Y)
		}
		mouseDownTime = 0
	case sdl.KEYDOWN:
		e := event.(*sdl.KeyboardEvent)
		fmt.Println("key down", e.Keysym.Sym)
	}

	redrawSet := LoadRedrawWindows()
	if len(redrawSet) > 0 {
		for _, wnd := range redrawSet {
			if wnd.IsVisible() {
				DeleteWindowFromRedraw(wnd)
				wnd.ReDraw()
			}
		}
	}
}
