package main

import (
	"log"
	"runtime"

	"gitee.com/millken/goui/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// https://chromium.googlesource.com/skia/+/chrome/m53/example/SkiaSDLExample.cpp
func main() {
	runtime.LockOSThread()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
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
	//defer sdl.Quit()

	dm, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		panic(err)
	}

	windowFlags = sdl.WINDOW_OPENGL
	// 创建一个窗口
	window, err := sdl.CreateWindow("Skia SDL Example", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		dm.W, dm.H, windowFlags)
	if err != nil {
		panic(err)
	}

	//defer window.Destroy()
	sdl.GLSetSwapInterval(1)

	ctx, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	if err := window.GLMakeCurrent(ctx); err != nil {
		panic(err)
	}
	// // 创建一个渲染器
	// renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	// if err != nil {
	// 	panic(err)
	// }
	// rendererInfo, err := renderer.GetInfo()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Using SDL renderer %s", rendererInfo.Name)
	// defer renderer.Destroy()

	// // 创建一个纹理
	// texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 800, 600)
	// if err != nil {
	// 	panic(err)
	// }
	// defer texture.Destroy()
	win, err := ui.NewWindow()
	if err != nil {
		panic(err)
	}
	win.SetSDLWindow(window)

	win.BackingScale()
	// // 创建一个Bitmap
	// bitmap := skia.NewBitmap()
	// bitmap.AllocN32Pixels(800, 600)

	// // 创建一个Canvas
	// canvas := skia.NewCanvas(bitmap)

	// // 创建一个Paint
	// paint := skia.NewPaint()
	// paint.SetColor(skia.ColorSetARGB(255, 255, 0, 0)) // 设置颜色为红色

	// // 在Canvas上绘制一个矩形
	// canvas.DrawRect(skia.NewRectFromLTRB(10, 10, 90, 90), paint)

	// // 将Bitmap的数据复制到纹理
	// pixels := bitmap.GetPixels()
	// texture.Update(nil, pixels, 800*4)

	// // 将纹理渲染到窗口
	// renderer.Copy(texture, nil, nil)
	// renderer.Present()

	// glContext, err := window.GLCreateContext()
	// if err != nil {
	// 	panic(err)
	// }
	// if err := window.GLMakeCurrent(glContext); err != nil {
	// 	panic(err)
	// }

	w, h := window.GLGetDrawableSize()
	log.Println("w,h:", w, h)
	// surface, err := window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// surface.FillRect(nil, 0)

	// rect := sdl.Rect{0, 0, 200, 200}
	// colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	// pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	// surface.FillRect(&rect, pixel)
	// sdl.WaitEvent()
	win.ReDraw()
	//window.UpdateSurface()
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
