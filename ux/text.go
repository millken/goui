package ux

// Text 文本元素
type Text struct {
	content string
	style   Style
}

func (t Text) Render(ctx *Context) {
	// 模拟渲染
	println("Render Text:", t.content, "at", ctx.X, ctx.Y, "Size:", ctx.Width, "x", ctx.Height)
}

func (t Text) getStyle() Style {
	return t.style
}

// WithStyle 返回带样式的 Text
func (t Text) WithStyle(mods ...StyleMod) UIElement {
	t.style = applyStyleMods(t.style, mods...)
	return t
}

// NewText 创建新文本
func NewText(s string) Text {
	return Text{content: s}
}
