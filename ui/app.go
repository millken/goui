package ui

import (
	"log"
	"math"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/millken/goui/event"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

type App interface {
	Create(*Settings) error
	Start() error

	Tick(elaspedTime time.Duration, syncChan chan<- struct{})
	Pause()
	Resume()
	Destroy()
	Stop()
}

func init() {
	// All init functions are run on the startup thread. Calling LockOSThread from an init function will cause the main
	// function to be invoked on that thread.
	runtime.LockOSThread()
}

var (
	sdlInited atomic.Bool
)
var defaultAppSettings = Settings{
	Width:  800,
	Height: 600,
}

// Start the application. This function does NOT return. While some calls may be safe to make, it should be assumed no
// calls into unison can be made prior to Start() being called unless explicitly stated otherwise.
func Run(app App) {
	//init sdl
	if err := initSDL(); err != nil {
		log.Fatal(err)
	}
	settings := defaultAppSettings
	eng := NewEngine(app)

	if err := app.Create(&settings); err != nil {
		log.Fatal(err)
	}
	defer app.Destroy()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	syncChan := make(chan struct{})
	elapsedTpsTime := time.Duration(0)
	go func() {
		for !eng.IsStopped() {
			if !eng.IsPaused() {
				now := time.Now()
				app.Tick(elapsedTpsTime, syncChan)
				elapsedTpsTime = time.Since(now)
			}
		}
	}()

	elapsedFpsTime := time.Duration(0)
	for !eng.IsStopped() {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch ev.GetType() {
			case sdl.QUIT:
				eng.Stop()
			case sdl.WINDOWEVENT:
				t := ev.(sdl.WindowEvent)
				switch t.Event {
				case sdl.WINDOWEVENT_FOCUS_GAINED:
					app.Resume()
					eng.SetPaused(false)
				case sdl.WINDOWEVENT_FOCUS_LOST:
					eng.SetPaused(true)
					app.Pause()
				case sdl.WINDOWEVENT_RESIZED:
					settings.Width = int(t.Data1)
				}
			case sdl.MOUSEBUTTONDOWN:
				if settings.EventMask&MouseButtonEventEnabled != 0 {
					t := ev.(sdl.MouseButtonEvent)
					button := event.ButtonUndefine
					switch t.Button {
					case sdl.ButtonLeft:
						button = event.ButtonLeft
					case sdl.ButtonMiddle:
						button = event.ButtonMiddle
					case sdl.ButtonRight:
						button = event.ButtonRight
					}
					eng.Publish(event.Event{
						Kind:   event.KindMouse,
						Type:   event.Type(t.Type),
						Button: button,
					})
				}
			case sdl.MOUSEMOTION:
				if settings.EventMask&MouseMotionEventEnabled != 0 {
					t := ev.(sdl.MouseMotionEvent)
					eng.Publish(event.Event{
						Kind:   event.KindMouse,
						Type:   event.TypeMove,
						X:      t.X,
						Y:      t.Y,
						Button: event.ButtonUndefine,
					})
				}
			case sdl.MOUSEWHEEL:
				if settings.EventMask&ScrollEventEnabled != 0 {
					t := ev.(sdl.MouseWheelEvent)
					x := float64(t.X)
					y := float64(t.Y)
					if x != 0 {
						x = x / math.Abs(x)
					}
					if y != 0 {
						y = y / math.Abs(y)
					}
					eng.Publish(event.Event{
						Kind: event.KindMouse,
						Type: event.Type(t.Type),
						X:    int32(x),
						Y:    int32(y),
					})
				}
			case sdl.KEYDOWN:
				if settings.EventMask&KeyEventEnabled != 0 {
					t := ev.(sdl.KeyboardEvent)
					keyName := sdl.GetKeyName(t.Keysym.Sym)
					eng.Publish(event.Event{
						Kind:    event.KindKeyboard,
						Type:    event.Type(t.Type),
						KeyCode: event.KeyMap[keyName],
						KeyName: keyName,
					})
				}
			}
		}
		if !eng.IsPaused() {
			now := time.Now()
			app.Tick(elapsedFpsTime, syncChan)
			elapsedFpsTime = time.Since(now)
		}
	}
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
