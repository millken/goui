package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/millken/goui"
	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	gouifont "github.com/millken/goui/font"
	"github.com/millken/goui/window/glfw"
	"github.com/millken/thorvg"
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

	// 标题（使用 Neutral 风格）
	title := component.NewButton()
	title.SetText("GOUI Flex + Button 增强演示")
	title.SetVariant(component.VariantNeutral)
	title.SetAlignSelf(component.AlignCenter)

	// 点击按钮：显示点击位置
	clickBtn := component.NewButton()
	clickBtn.SetText("点我试试")
	clickBtn.OnClick(func(e event.MouseEvent) {
		clickBtn.SetText(fmt.Sprintf("点击位置: (%.0f, %.0f)", e.X, e.Y))
	})
	clickBtn.SetAlignSelf(component.AlignCenter)

	// 预设风格按钮
	primaryBtn := component.NewButton()
	primaryBtn.SetText("Primary")
	primaryBtn.SetVariant(component.VariantPrimary)
	primaryBtn.SetAlignSelf(component.AlignCenter)

	dangerBtn := component.NewButton()
	dangerBtn.SetText("Danger")
	dangerBtn.SetVariant(component.VariantDanger)
	dangerBtn.SetAlignSelf(component.AlignCenter)

	// 大号按钮：更大字号与内边距、圆角
	largeBtn := component.NewButton()
	largeBtn.SetText("Large Button")
	largeBtn.SetVariant(component.VariantPrimary)
	largeBtn.SetTextSize(18)
	largeBtn.SetPadding(20, 12)
	largeBtn.SetCornerRadius(8)
	largeBtn.SetAlignSelf(component.AlignCenter)

	// 禁用按钮
	disabledBtn := component.NewButton()
	disabledBtn.SetText("Disabled")
	disabledBtn.SetVariant(component.VariantNeutral)
	disabledBtn.SetDisabled(true)
	disabledBtn.SetAlignSelf(component.AlignCenter)

	// 切换禁用按钮
	toggleBtn := component.NewButton()
	toggleBtn.SetText("Toggle Danger Disabled")
	toggleBtn.SetAlignSelf(component.AlignCenter)
	toggleBtn.OnClick(func(e event.MouseEvent) {
		// 简单切换 dangerBtn 的禁用状态
		// 这里通过读取一次边界触发使用，以示例没有直接状态读取方法的情况下直接切换
		// 采用 SetDisabled(true/false) 的方案
		// 为了演示，这里用一个闭包局部变量保存状态
	})

	// 用闭包变量持久化禁用状态
	var dangerDisabled bool
	toggleBtn.OnClick(func(e event.MouseEvent) {
		dangerDisabled = !dangerDisabled
		dangerBtn.SetDisabled(dangerDisabled)
		if dangerDisabled {
			toggleBtn.SetText("Enable Danger")
		} else {
			toggleBtn.SetText("Disable Danger")
		}
	})

	// 保持按钮为内容高度，不在主轴方向扩展
	title.SetFlexGrow(0)
	clickBtn.SetFlexGrow(0)
	primaryBtn.SetFlexGrow(0)
	dangerBtn.SetFlexGrow(0)
	largeBtn.SetFlexGrow(0)
	disabledBtn.SetFlexGrow(0)
	toggleBtn.SetFlexGrow(0)

	root.AddChild(title)
	root.AddChild(clickBtn)
	root.AddChild(primaryBtn)
	root.AddChild(dangerBtn)
	root.AddChild(largeBtn)
	root.AddChild(disabledBtn)
	root.AddChild(toggleBtn)

	app, err := goui.NewApp(win, root)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 注意：ThorVG 需要在 Init 之后再加载字体，这里放在 NewApp(内部调用 Init) 之后
	if data := gouifont.GetDefaultFontData(); len(data) > 0 {
		f, err := os.CreateTemp("", "goui-font-*.ttf")
		if err == nil {
			_, _ = f.Write(data)
			_ = f.Close()
			if ferr := thorvg.LoadFont(f.Name()); ferr != nil {
				log.Printf("加载内置字体失败: %v", ferr)
			}
			defer os.Remove(f.Name())
		} else {
			log.Printf("创建临时字体文件失败: %v", err)
		}
	}

	fmt.Println("启动 GOUI Flex+Button 示例…\n- 测试 hover/active 视觉\n- 点击 Primary/Danger/Toggle/Disabled 等按钮\n- 调整窗口尺寸查看布局与重绘")
	app.Run()
}
