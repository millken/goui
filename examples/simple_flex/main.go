package main

import (
	"fmt"
	"log"
	"log/slog"
	"runtime"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/goui/color"
	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	"github.com/millken/goui/window/glfw"
)

// YogaBox 是一个使用 Yoga 计算几何并用 thorvg 绘制的简单矩形组件
// 仅用于演示 Flex 布局在实际组件上的效果
type YogaBox struct {
	*component.BaseComponent
	bg color.Color
}

func NewYogaBox(bg color.Color) *YogaBox {
	return &YogaBox{
		BaseComponent: component.NewBaseComponent(),
		bg:            bg,
	}
}

// Paint 使用 Yoga 计算的布局绘制背景矩形
func (b *YogaBox) Paint(vg *thorvg.Canvas) bool {
	// 仅在需要重绘（脏）时绘制；避免每帧都返回渲染
	if !b.IsDirty() {
		return false
	}
	if b.bg != color.TRANSPARENT {
		background := thorvg.NewShape()
		if background != nil {
			l := b.ComputedLayout()
			rct := b.Bounds()
			ax, ay := float32(rct.Min.X), float32(rct.Min.Y)
			background.AppendRect(ax, ay, l.Width, l.Height, 0, 0, false)
			r, g, bl, a := b.bg.RGBA()
			background.SetFillColor(r, g, bl, a)
			vg.Push(&background.Paint)
			background.Release()
		}
	}
	b.ClearDirty()
	return true
}

func main() {
	// 确保 GLFW/GL 在同一 OS 线程
	runtime.LockOSThread()
	slog.SetDefault(slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelDebug})))

	// 创建窗口
	win, err := glfw.NewWindow("GOUI Flex 示例", 800, 600)
	if err != nil {
		log.Fatalf("创建窗口失败: %v", err)
	}
	defer win.Close()
	_ = win.Open()

	// 根容器：垂直方向的 Flex 布局
	root := component.Column("root")
	// 使用新的便捷 API：内边距与行/列间距
	root.SetAlign(component.AlignStretch)
	root.SetPadding(12)
	root.SetRowGap(8)
	root.SetColumnGap(8)
	// 不固定根容器的宽高，交由 App 每帧以窗口尺寸驱动 Yoga.Layout

	// 三个 YogaBox，设置 flexGrow 让它们沿主轴（竖向）分配空间
	b1 := NewYogaBox(color.RED)
	b1.GetYogaNode().SetFlexGrow(1)
	b1.GetYogaNode().SetMinHeight(40)
	b1.OnClick(func(e event.MouseEvent) { fmt.Println("点击：红色") })

	b2 := NewYogaBox(color.GREEN)
	b2.GetYogaNode().SetFlexGrow(2) // 更大的占比
	b2.GetYogaNode().SetMinHeight(40)
	b2.OnClick(func(e event.MouseEvent) { fmt.Println("点击：绿色") })

	b3 := NewYogaBox(color.BLUE)
	b3.GetYogaNode().SetFlexGrow(1)
	b3.GetYogaNode().SetMinHeight(40)
	b3.OnClick(func(e event.MouseEvent) { fmt.Println("点击：蓝色") })

	root.AddChild(b1)
	root.AddChild(b2)
	root.AddChild(b3)

	// 再加一些盒子来演示 Wrap 与 AlignContent
	wrap := component.Row("wrap")
	wrap.SetWrap(component.WrapWrap)
	wrap.SetAlign(component.AlignStart)
	wrap.SetAlignContent(component.AlignStart)
	wrap.SetRowGap(8)
	wrap.SetColumnGap(8)
	for i := 0; i < 10; i++ {
		box := NewYogaBox(color.NewColor(uint8(30*i), uint8(255-20*i), uint8(100+10*i), 255))
		box.GetYogaNode().SetMinWidth(80)
		box.GetYogaNode().SetMinHeight(40)
		wrap.AddChild(box)
	}
	root.AddChild(wrap)

	// 创建并运行应用
	app, err := goui.NewApp(win, root)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}
	fmt.Println("启动 GOUI Flex 示例… 点击不同区域查看命中与布局效果")
	app.Run()
}
