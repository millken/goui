package main

import (
	"runtime"

	"gitee.com/millken/goui/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// https://chromium.googlesource.com/skia/+/chrome/m53/example/SkiaSDLExample.cpp
func main() {
	runtime.LockOSThread()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	var windowFlags uint32
	if ui.IsAndroid() {
		windowFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE |
			sdl.WINDOW_BORDERLESS | sdl.WINDOW_FULLSCREEN_DESKTOP |
			sdl.WINDOW_ALLOW_HIGHDPI
	} else {
		// For all other clients we use the core profile and operate in a window
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
		windowFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	}
	kStencilBits := 8
	//kMsaaSampleCount := 4
	sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 0)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, kStencilBits)
	sdl.GLSetAttribute(sdl.GL_ACCELERATED_VISUAL, 1)

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	dm, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Skia SDL Example", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		dm.W, dm.H, windowFlags)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
	window.UpdateSurface()
	running := true
	for running {
		event := sdl.WaitEvent()
		switch event.(type) {
		case *sdl.QuitEvent:
			running = false
		default:
			//rmlui.SDLInputEventHandler(uictx, event)
		}

	}
}
