package main

import (
	"gitee.com/millken/goui/sdl2/rmlui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

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
	rmlui.SetRenderInterface(rmlui.NewRenterInterface(renderer))

	if rmlui.Initialise() == false {
		panic("Failed to initialise RmlUi")
	}

	// Create the main RmlUi context.
	uictx := rmlui.CreateContext("main", rmlui.NewVector2Int_2(800, 600))
	if uictx == nil {
		println("Failed to create RmlUi context")
		rmlui.Shutdown()

		return
	}
	document := uictx.LoadDocumentFromMemory("<template><h1>Hello World!</h1></template>")
	if document == nil {
		println("Failed to load document")
		rmlui.Shutdown()
		return
	}
	document.Show()
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
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
