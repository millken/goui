package ux

type Column struct {
	Children []UIElement
	style    Style
}

func (c Column) Render(ctx *Context) {
	style := c.style
	if ctx.Style.Width > 0 {
		style.Width = ctx.Style.Width
	}
	if ctx.Style.Height > 0 {
		style.Height = ctx.Style.Height
	}

	// 计算可用空间（扣除 padding）
	availWidth := style.Width - style.Padding[1] - style.Padding[3]
	availHeight := style.Height - style.Padding[0] - style.Padding[2]

	y := ctx.Y + style.Padding[0] // 起始 Y 偏移（上 padding）

	// === 第一次遍历：收集信息 ===
	var totalFlex float32
	var usedHeight float32

	for _, child := range c.Children {
		if childStyle, ok := getStyle(child); ok {
			height := childStyle.Height

			if height > 0 {
				usedHeight += height + childStyle.Margin[0] + childStyle.Margin[2]
			} else if childStyle.FlexGrow > 0 {
				totalFlex += childStyle.FlexGrow
			} else {
				// 默认高度 + margin
				usedHeight += 40 + childStyle.Margin[0] + childStyle.Margin[2]
			}
		}
	}

	// 可用于 Flex 分配的空间
	remainingHeight := availHeight - usedHeight

	// === 第二次遍历：布局每个子元素 ===
	for _, child := range c.Children {
		childCtx := &Context{
			X:      ctx.X + style.Padding[3],
			Y:      y,
			Width:  availWidth,
			Height: 40, // 临时值，下面会覆盖
			Style:  Style{},
		}

		if childStyle, ok := getStyle(child); ok {
			height := childStyle.Height

			// 如果是 Flex 元素，分配剩余空间
			if height == 0 && childStyle.FlexGrow > 0 && totalFlex > 0 {
				height = (childStyle.FlexGrow / totalFlex) * remainingHeight
			} else if height == 0 {
				height = 40 // 默认非 flex 高度
			}

			childCtx.Height = height
			childCtx.Style = childStyle

			// 渲染子元素
			child.Render(childCtx)

			// 更新 y 位置（加上高度和 margin bottom）
			y += height + childStyle.Margin[0] + childStyle.Margin[2]
		}
	}
}

func (c Column) getStyle() Style {
	return c.style
}

func (c Column) WithStyle(mods ...StyleMod) UIElement {
	c.style = applyStyleMods(c.style, mods...)
	return c
}

func NewColumn(children ...UIElement) Column {
	return Column{Children: children}
}
