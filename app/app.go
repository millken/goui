package app

import (
	"log"
	"runtime"
	"sync/atomic"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/millken/goui/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	// All init functions are run on the startup thread. Calling LockOSThread from an init function will cause the main
	// function to be invoked on that thread.
	runtime.LockOSThread()
}

var (
	sdlInited atomic.Bool
)

// Start the application. This function does NOT return. While some calls may be safe to make, it should be assumed no
// calls into unison can be made prior to Start() being called unless explicitly stated otherwise.
func Run(app ui.App) {
	//init sdl
	if err := initSDL(); err != nil {
		log.Fatal(err)
	}
	//init app
	if err := app.Create(); err != nil {
		log.Fatal(err)
	}
	defer app.Destroy()
	ui.ProcessEvents()

}

func initSDL() error {
	if sdlInited.Load() {
		return nil
	}
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)

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
		return errors.Wrap(err, "sdl init error")
	}
	sdlInited.Store(true)
	return gl.Init()
}
