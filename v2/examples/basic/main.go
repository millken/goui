package main

import (
	"fmt"

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
	rootComp.SetBounds(0, 0, 800, 600)
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
	settings := v2.NewSettings()
	settings.SetRenderer(impl.NewOpenGLRenderer()) // 假设有一个 OpenGL 渲染器实现
	v2.Run(settings, v2.NewApp(app{}))
}
