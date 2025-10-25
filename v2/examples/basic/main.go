package main

import (
	"fmt"

	v2 "github.com/millken/goui/v2"
	"github.com/millken/goui/v2/impl"
)

type app struct{}

func (a *app) Create() error {
	win1 := impl.NewGLFWWindow("First Window", 800, 600)
	win1.Show()
	fmt.Println("App Created")
	return nil
}
func (a *app) Destroy() error {
	fmt.Println("App Destroyed")
	return nil
}
func main() {
	settings := v2.NewSettings()
	settings.SetRenderer(impl.NewOpenGLRenderer()) // 假设有一个 OpenGL 渲染器实现
	v2.Run(settings, v2.NewApp(app{}))
}
