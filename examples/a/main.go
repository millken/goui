package main

import (
	"runtime"

	"github.com/millken/goui"
	"github.com/millken/goui/components"
)

func main() {
	runtime.LockOSThread()

	// 构建组件树：根 Flex 下放三个 Box
	root := components.NewFlex()
	red := components.NewBox(220, 60, 60, 255)
	green := components.NewBox(60, 220, 120, 255)
	blue := components.NewBox(60, 120, 220, 255)

	// 简单设置子项尺寸（用 Yoga 节点的 Style 设置）
	if n := red.GetNode(); n != nil {
		n.SetWidth(120)
		n.SetHeight(120)
	}
	if n := green.GetNode(); n != nil {
		n.SetWidth(180)
		n.SetHeight(80)
	}
	if n := blue.GetNode(); n != nil {
		n.SetWidth(100)
		n.SetHeight(160)
	}

	root.AddChild(red)
	root.AddChild(green)
	root.AddChild(blue)

	app, err := goui.NewApp("Goui Demo", 800, 600, root)
	if err != nil {
		panic(err)
	}
	defer app.Close()
	app.Run()
}
