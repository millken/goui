package ux

// Row 水平布局
type Row struct {
	Children []UIElement
	style    Style
}

func (r Row) Render(ctx *Context) {
	// 使用传入的 Context 样式（父级样式）或自身样式
	style := r.style
	if ctx.Style.Width > 0 {
		style.Width = ctx.Style.Width
	}
	if ctx.Style.Height > 0 {
		style.Height = ctx.Style.Height
	}

	availWidth := style.Width - style.Padding[1] - style.Padding[3]
	availHeight := style.Height - style.Padding[0] - style.Padding[2]

	var totalFlex float32
	var nonFlexWidth float32

	// 第一次遍历：计算固定宽度和总 flexGrow
	for _, child := range r.Children {
		if childStyle, ok := getStyle(child); ok {
			width := childStyle.Width
			if width > 0 {
				nonFlexWidth += width + childStyle.Margin[1] + childStyle.Margin[3]
			} else if childStyle.FlexGrow > 0 {
				totalFlex += childStyle.FlexGrow
			} else {
				nonFlexWidth += 60 + childStyle.Margin[1] + childStyle.Margin[3] // 默认宽度
			}
		}
	}

	remaining := availWidth - nonFlexWidth
	x := ctx.X + style.Padding[3] // left padding
	y := ctx.Y + style.Padding[0] // top padding

	// 第二次遍历：分配位置
	for _, child := range r.Children {
		childCtx := &Context{
			X:      x,
			Y:      y,
			Width:  100,
			Height: availHeight,
			Style:  Style{}, // 将 child.style 合并进来
		}

		if childStyle, ok := getStyle(child); ok {
			width := childStyle.Width
			if width == 0 && childStyle.FlexGrow > 0 && totalFlex > 0 {
				width = (childStyle.FlexGrow / totalFlex) * remaining
			} else if width == 0 {
				width = 60
			}
			childCtx.Width = width
			childCtx.Style = childStyle

			// 渲染子元素
			child.Render(childCtx)

			// 更新 x
			x += width + childStyle.Margin[1] + childStyle.Margin[3]
		}
	}
}

func (r Row) getStyle() Style {
	return r.style
}

func (r Row) WithStyle(mods ...StyleMod) UIElement {
	r.style = applyStyleMods(r.style, mods...)
	return r
}

func NewRow(children ...UIElement) Row {
	return Row{Children: children}
}
