package main

import (
	"fmt"
	"runtime"

	v2 "github.com/millken/goui/v2"
	"github.com/millken/goui/v2/component"
	"github.com/millken/goui/v2/impl"
)

type app struct{}

func (a *app) Create(ctx *v2.Context) error {
	win1, err := impl.NewGLFWWindow("First Window", 800, 600)
	if err != nil {
		return err
	}
	if err := win1.Init(ctx); err != nil {
		return err
	}
	rootComp := component.NewBoxComponent()
	// 在窗口内的某个坐标绘制一个矩形（示例：x=100,y=80,width=220,height=160）
	rootComp.SetBounds(100, 80, 220, 160)
	// 内层子组件：初始为蓝色，点击切换为白色
	box2 := component.NewBoxComponent()
	box2.SetBounds(20, 20, 180, 120) // 相对于 root 的坐标（示例中未实现局部坐标变换，示意）
	box2.SetOnClick(func() {
		// 切换颜色
		if box2.GetBounds().Width >= 0 {
			// 直接切换颜色示例
			box2.SetColor(v2.Color{R: 255, G: 255, B: 255, A: 255})
		}
	})
	rootComp.AddChild(box2)
	// 添加一个 Label 演示文本渲染
	lbl := component.NewLabel("Hello Goui - Label")
	lbl.SetBounds(10, 10, 200, 24)
	// set label color to white so it is visible on default background
	lbl.SetColor(v2.Color{R: 0, G: 0, B: 255, A: 255})
	rootComp.AddChild(lbl)
	win1.SetRoot(rootComp)
	win1.Show()
	fmt.Println("App Created")
	return nil
}
func (a *app) Destroy() error {
	fmt.Println("App Destroyed")
	return nil
}
func main() {
	// OpenGL/GLFW 需要在同一 OS 线程内进行所有 GL 调用
	runtime.LockOSThread()
	settings := v2.NewSettings()
	settings.SetRenderer(impl.NewOpenGLRenderer()) // 假设有一个 OpenGL 渲染器实现
	v2.Run(settings, v2.NewApp(app{}))
}
