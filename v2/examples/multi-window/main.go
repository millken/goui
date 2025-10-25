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
	box2 := component.NewBoxWithColor(v2.Color{R: 231, G: 76, B: 60, A: 255})
	box2.SetBounds(50, 50, 100, 100)
	rootComp.AddChild(box2)
	win1.SetRoot(rootComp)
	win1.Show()
	w2, err2 := impl.NewGLFWWindow("Second Window", 400, 300)
	if err2 != nil {
		return err2
	}
	if err := w2.Init(ctx); err != nil {
		return err
	}
	w2.SetOnClose(func() {
		w2.Hide()
	})
	// 简单根组件
	r2 := component.NewBoxComponent()
	r2.SetBounds(10, 10, 200, 150)
	w2.SetRoot(r2)
	rootComp.SetOnClick(func() {

		w2.Show()
	})

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
