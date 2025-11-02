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

	// 使用组件 Style API（避免直接访问 Yoga 节点）
	red.SetStyle(goui.Style{
		Width:       goui.Px(120),
		Height:      goui.Px(120),
		Margin:      goui.EdgesAll(goui.Px(6)),
		Border:      goui.EdgesAll(goui.Px(4)),
		BorderColor: strPtr("#333333"),
	})
	green.SetStyle(goui.Style{
		Width:      goui.Px(180),
		Height:     goui.Px(80),
		Margin:     goui.EdgesAll(goui.Px(6)),
		FlexGrow:   float32Ptr(1), // 演示占剩余空间（与蓝色对比）
		FlexShrink: float32Ptr(0),
	})
	blue.SetStyle(goui.Style{
		Width:       goui.Pc(30), // 30% 宽度（相对父）
		Height:      goui.Px(160),
		Margin:      goui.EdgesAll(goui.Px(6)),
		Border:      goui.EdgesAll(goui.Px(2)),
		BorderColor: strPtr("rgba(30,60,120,0.8)"),
	})

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

func float32Ptr(v float32) *float32 { return &v }
func strPtr(s string) *string       { return &s }
