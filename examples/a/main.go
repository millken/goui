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

	// 辅助函数：创建 Box Style
	createBoxStyle := func(bgColor string, width *goui.Dimension, height *goui.Dimension, borderWidth float32, borderColor string, flexGrow *float32) goui.Style {
		style := goui.Style{
			Width:  width,
			Height: height,
			Margin: goui.EdgesAll(goui.Px(6)),
		}
		if bgColor != "" {
			style.BackgroundColor = &bgColor
		}
		if borderWidth > 0 {
			style.Border = goui.EdgesAll(goui.Px(borderWidth))
			if borderColor != "" {
				style.BorderColor = &borderColor
			}
		}
		if flexGrow != nil {
			style.FlexGrow = flexGrow
			style.FlexShrink = float32Ptr(0)
		}
		return style
	}

	// 设置初始 Style 和事件
	red.SetStyle(createBoxStyle("#ff0000", goui.Px(120), goui.Px(120), 4, "#333333", nil))
	red.SetOnHover(func(hovered bool) {
		if hovered {
			red.SetStyle(createBoxStyle("#ffcccc", goui.Px(120), goui.Px(120), 4, "#333333", nil))
		} else {
			red.SetStyle(createBoxStyle("#ff0000", goui.Px(120), goui.Px(120), 4, "#333333", nil))
		}
	})

	green.SetStyle(createBoxStyle("#00ff00", goui.Px(180), goui.Px(80), 0, "", float32Ptr(1)))
	green.SetOnHover(func(hovered bool) {
		if hovered {
			green.SetStyle(createBoxStyle("#ccffcc", goui.Px(180), goui.Px(80), 0, "", float32Ptr(1)))
		} else {
			green.SetStyle(createBoxStyle("#00ff00", goui.Px(180), goui.Px(80), 0, "", float32Ptr(1)))
		}
	})

	blue.SetStyle(createBoxStyle("#0000ff", goui.Pc(30), goui.Px(160), 2, "rgba(30,60,120,0.8)", nil))
	blue.SetOnHover(func(hovered bool) {
		if hovered {
			blue.SetStyle(createBoxStyle("#ccccff", goui.Pc(30), goui.Px(160), 2, "rgba(30,60,120,0.8)", nil))
		} else {
			blue.SetStyle(createBoxStyle("#0000ff", goui.Pc(30), goui.Px(160), 2, "rgba(30,60,120,0.8)", nil))
		}
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
