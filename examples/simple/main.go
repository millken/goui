package main

import (
	"fmt"
	"log"

	"github.com/millken/goui"
	"github.com/millken/goui/color"
	"github.com/millken/goui/component"
	"github.com/millken/goui/window/glfw"
)

func main() {
	// 创建GLFW窗口
	window, err := glfw.NewWindow("GOUI 简单演示", 800, 600)
	if err != nil {
		log.Fatalf("创建窗口失败: %v", err)
	}
	defer window.Close()

	// 创建根组件（垂直布局）
	root := component.NewBox(0, 0, 800, 600)
	root.SetBackgroundColor(color.GREEN)
	box1 := component.NewBox(100, 100, 100, 100)
	box1.SetBackgroundColor(color.RED)
	root.AddChild(box1)

	// 创建应用
	app, err := goui.NewApp(window, root)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 运行应用
	fmt.Println("启动GOUI简单演示...")
	app.Run()
}
