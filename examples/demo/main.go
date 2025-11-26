package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/millken/go/rgfw"
	"github.com/millken/goui/font"
	"github.com/millken/thorvg"
)

var icon = [4 * 3 * 3]uint8{
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0x00,
	0xFF, 0xFF, 0x00, 0xFF,
	0xFF, 0xFF, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
}

func main() {
	runtime.LockOSThread()

	// 初始化 ThorVG 引擎
	if err := thorvg.Init(0); err != nil {
		panic(err)
	}
	defer thorvg.Term()

	if err := thorvg.LoadFont("/home/millken/Downloads/OPPO Sans 4.0.ttf"); err != nil {
		panic(err)
	}
	rgfw.SetClassName("RGFW Simple Animation Test")
	win := rgfw.CreateWindow("ThorVG Simple Animation", rgfw.Rect{X: 100, Y: 100, W: 800, H: 600}, rgfw.WindowCenter|rgfw.WindowAllowDND|rgfw.WindowOpenGL)

	win.MakeCurrentOpenGL()

	// 初始化 OpenGL 绑定
	if err := gl.Init(); err != nil {
		panic(err)
	}

	win.SetExitKey(rgfw.KeyEscape)
	win.SwapIntervalOpenGL(1)
	fmt.Printf("OpenGL Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))

	// 创建 ThorVG OpenGL 画布
	canvas := thorvg.NewGlCanvas()
	if canvas == nil {
		panic("Failed to create ThorVG OpenGL canvas")
	}
	defer canvas.Destroy()

	width, height := uint32(800), uint32(600)

	err := canvas.GlCanvasSetTarget(win.GetOpenGLContext(), 0, width, height, thorvg.ColorspaceABGR8888S)
	if err != nil {
		panic(err)
	}
	startTime := time.Now()
	frames := uint32(0)

	fmt.Println("简单动画测试 - 按 ESC 退出")
	fmt.Println("如果看到蓝色矩形在旋转，说明基础渲染正常")
	measureSize := font.MeasureText("简单动画测试 - 按 ESC 退出", font.WithFontSize(14))
	fmt.Printf("文本测量示例: 文本宽度=%f 高度=%f\n", measureSize.Width, measureSize.Height)
	for !win.ShouldClose() {
		var event = rgfw.GetEvent()
		for win.CheckEvent(event) {
			switch event.Type {
			case rgfw.EventQuit:
				break
			}
			rgfw.ReleaseEvent(event)
		}

		if win.IsKeyRelease(rgfw.KeyEscape) {
			break
		}
		// 简单的动画：每帧创建一个矩形
		elapsed := time.Since(startTime).Seconds()
		rotation := float32(elapsed * 60) // 每秒60度

		// 创建背景矩形（清除上一帧）
		background := thorvg.NewShape()
		if background != nil {
			background.AppendRect(0, 0, float32(width), float32(height), 0, 0, false)
			background.SetFillColor(229, 229, 229, 255) // 灰色背景
			canvas.Push(&background.Paint)
			background.Release()
		}

		// 创建动画矩形
		rect := thorvg.NewShape()
		if rect != nil {
			// 在屏幕中心创建一个矩形
			// centerX, centerY := float32(width)/2, float32(height)/2
			rect.AppendRect(0, 0, float32(measureSize.Width), float32(measureSize.Height), 0, 0, false)
			rect.SetFillColor(0, 100, 155, 200) // 半透明蓝色

			// 应用旋转变换
			// rect.Translate(centerX, centerY)
			canvas.Push(&rect.Paint)
			rect.Release()
		}

		//写文本
		text := thorvg.NewText()
		if text != nil {
			text.SetText("简单动画测试 - 按 ESC 退出")
			text.SetSize(14.0)
			text.SetColor(30, 60, 120)
			// text.SetAlign(0.5, 0.0)
			text.SetFont("OPPO Sans 4.0")
			text.Translate(0, 0)
			canvas.Push(&text.Paint)
			text.Release()
		}

		// 渲染 ThorVG 画布
		if err := canvas.Update(); err != nil {
			fmt.Printf("Canvas Update 失败: %v\n", err)
		}
		if err := canvas.Draw(true); err != nil {
			fmt.Printf("Canvas Draw 失败: %v\n", err)
		}
		if err := canvas.Sync(); err != nil {
			fmt.Printf("Canvas Sync 失败: %v\n", err)
		}

		win.SwapBuffersOpenGL()
		frames++

		// 每60帧显示一次状态
		if frames%60 == 0 {
			fmt.Printf("帧数: %d, 角度: %.1f°\n", frames, rotation)
		}

		time.Sleep(16 * time.Millisecond) // 约60 FPS
	}

	fmt.Printf("总帧数: %d\n", frames)
	win.Close()
}
