package main

import (
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	// 初始化 SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// 初始化 SDL_ttf
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	// 创建窗口和渲染器
	window, err := sdl.CreateWindow("Test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// 创建矩形
	rect := sdl.Rect{0, 0, 200, 100}

	// 绘制并填充圆角矩形
	gfx.RoundedBoxColor(renderer, int32(rect.X), int32(rect.Y), int32(rect.X+rect.W), int32(rect.Y+rect.H), 10, sdl.Color{255, 0, 0, 255})

	// 创建字体
	font, err := ttf.OpenFont("/Users/millken/golang/ui/ui/eg/DejaVuSans.ttf", 20)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	// 创建表面和纹理
	surface, err := font.RenderUTF8Blended("Button", sdl.Color{255, 255, 255, 255})
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, _ := renderer.CreateTextureFromSurface(surface)
	defer texture.Destroy()

	// 将文本渲染到圆角矩形上
	renderer.Copy(texture, nil, &rect)

	// 更新屏幕
	renderer.Present()

	// 等待用户关闭窗口
	for event := sdl.WaitEvent(); event != nil; event = sdl.WaitEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
			return
		}
	}
}
