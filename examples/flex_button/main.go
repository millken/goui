package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/goui/color"
	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	"github.com/millken/goui/window/glfw"
)

func main() {
	runtime.LockOSThread()
	if err := thorvg.LoadFont("/home/millken/Downloads/OPPO Sans 4.0.ttf"); err != nil {
		panic(err)
	}
	win, err := glfw.NewWindow("GOUI Flex+Button 示例", 800, 600)
	if err != nil {
		log.Fatalf("创建窗口失败: %v", err)
	}
	defer win.Close()
	_ = win.Open()

	// 根容器：垂直 Flex，交叉轴拉伸
	root := component.Column("root")
	root.SetAlign(component.AlignStretch)
	root.SetGap(8)

	// 标题按钮（只是展示文本）
	title := component.NewButton()
	title.SetText("欢迎使用 GOUI Flex + Button")
	title.SetBackgroundColor(color.NewColor(245, 245, 245, 255))
	title.SetAlignSelf(component.AlignCenter)

	// 点击按钮
	clickBtn := component.NewButton()
	clickBtn.SetText("点我试试")
	clickBtn.OnClick(func(e event.MouseEvent) {
		clickBtn.SetText(fmt.Sprintf("点击位置: (%.0f, %.0f)", e.X, e.Y))
	})
	clickBtn.SetAlignSelf(component.AlignCenter)

	// 彩色按钮
	redBtn := component.NewButton()
	redBtn.SetText("红色")
	redBtn.SetBackgroundColor(color.NewColor(255, 230, 230, 255))
	redBtn.SetAlignSelf(component.AlignCenter)

	greenBtn := component.NewButton()
	greenBtn.SetText("绿色")
	greenBtn.SetBackgroundColor(color.NewColor(230, 255, 230, 255))
	greenBtn.SetAlignSelf(component.AlignCenter)

	blueBtn := component.NewButton()
	blueBtn.SetText("蓝色")
	blueBtn.SetBackgroundColor(color.NewColor(230, 230, 255, 255))
	blueBtn.SetAlignSelf(component.AlignCenter)

	// 保持按钮为内容高度，不在主轴方向扩展
	title.SetFlexGrow(0)
	clickBtn.SetFlexGrow(0)
	redBtn.SetFlexGrow(0)
	greenBtn.SetFlexGrow(0)
	blueBtn.SetFlexGrow(0)

	root.AddChild(title)
	root.AddChild(clickBtn)
	root.AddChild(redBtn)
	root.AddChild(greenBtn)
	root.AddChild(blueBtn)

	app, err := goui.NewApp(win, root)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	fmt.Println("启动 GOUI Flex+Button 示例… 点击按钮体验事件与布局；缩放窗口查看自适应效果")
	app.Run()
}
