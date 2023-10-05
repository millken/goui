package main

import (
	"runtime"

	"gitee.com/millken/goui/sdl2/rmlui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	runtime.LockOSThread()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		800, 600, sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	rendererInfo, err := renderer.GetInfo()
	if err != nil {
		panic(err)
	}
	println(rendererInfo.Name)
	rmlui.SetSystemInterface(rmlui.NewSystemInterface(window))
	rmlui.SetRenderInterface(rmlui.NewRenterInterface(renderer))

	if !rmlui.Initialise() {
		panic("Failed to initialise RmlUi")
	}

	// Create the main RmlUi context.
	// uictx := rmlui.CreateContext("main", rmlui.NewVector2Int_2(800, 600))
	// if uictx == nil {
	// 	println("Failed to create RmlUi context")
	// 	rmlui.Shutdown()

	// 	return
	// }
	// document := uictx.LoadDocumentFromMemory("<template><h1>Hello World!</h1></template>")
	// if document == nil {
	// 	println("Failed to load document")
	// 	rmlui.Shutdown()
	// 	return
	// }
	// document.Show()

	// running := true

	// for running {
	// 	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	// 		switch event.(type) {
	// 		case *sdl.QuitEvent:
	// 			println("Quit")
	// 			running = false
	// 			break
	// 		}
	// 	}
	// }
}
