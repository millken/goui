package ux

import (
	"cmp"
)

type Box struct {
	Children []Widget
	Style    Style
}

func (b *Box) Mount(ctx *Context) {
	// 1. 临时保存原始样式
	originalStyle := ctx.Style

	// --- 计算可用宽度 ---
	parentAvailableWidth := ctx.Width // 从 Context 传入的容器宽度
	if parentAvailableWidth == 0 {
		parentAvailableWidth = 800 // 可选：全局默认宽度，如屏幕宽
	}

	// 考虑父级 X 位置和 Margin，计算当前可用空间
	// 注意：ctx.X 是当前流式布局的 X 偏移
	remainingWidth := parentAvailableWidth - ctx.X

	// --- 确定 Box 宽度 ---
	width := ctx.Style.Width
	if width == 0 {
		// 🔹 宽度未设置 → 占满父级剩余宽度（减去左右 margin）
		width = remainingWidth - ctx.Style.Margin[1] - ctx.Style.Margin[3]
		if width < 0 {
			width = 0
		}
	}

	resolvedWidth := cmp.Or(b.Style.Width, originalStyle.Width, ctx.Width)
	if resolvedWidth == 0 {
		resolvedWidth = 20
	}

	// 应用 Box 自身的 Style 到当前 Context 的 Style（宽度使用解析后的值）
	ctx.Style.Width = resolvedWidth
	ctx.Style.Height = b.Style.Height
	ctx.Style.Margin = b.Style.Margin
	ctx.Style.Padding = b.Style.Padding
	ctx.Style.ZIndex = b.Style.ZIndex

	// width := resolvedWidth
	height := cmp.Or(ctx.Style.Height, 20)

	// 4. 计算位置（考虑 Margin）
	x := ctx.X + ctx.Style.Margin[3] // left
	y := ctx.Y + ctx.Style.Margin[0] // top

	// 5. 确定背景色
	bgColor := cmp.Or(ctx.Style.BackgroundColor, colorGray)

	// 6. 添加背景渲染命令
	if bgColor != nil {
		ctx.AddCmd(RenderCmd{
			Type:   CmdFill,
			Rect:   Rect{x, y, width, height},
			Color:  bgColor,
			ZIndex: ctx.Style.ZIndex, // 使用样式中的 ZIndex
		})
	}

	// 如果有子组件，递归挂载它们（考虑 Padding）
	if len(b.Children) > 0 {
		savedX, savedY := ctx.X, ctx.Y

		// 移动到 Box 内部（加上 Padding）
		ctx.X = x + ctx.Style.Padding[3] // left padding
		ctx.Y = y + ctx.Style.Padding[0] // top padding

		for _, child := range b.Children {
			child.Mount(ctx)
		}

		// 恢复位置（用于流式布局）
		ctx.X = savedX
		ctx.Y = savedY
	}

	// 8. 更新流式布局位置（仅 Box 本身）
	ctx.X += width + ctx.Style.Margin[1] + ctx.Style.Margin[3]
	ctx.Y += 0 // Box 不影响 Y（除非是 Block 布局）

	// 9. 恢复原始样式
	ctx.Style = originalStyle
}

func NewBox(children ...Widget) *Box {
	return &Box{
		Children: children,
		Style:    Style{},
	}
}

func (b *Box) WithStyle(styleFunc func(*Style)) *Box {
	styleFunc(&b.Style)
	return b
}
