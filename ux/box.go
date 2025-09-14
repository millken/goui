package ux

import "image/color"

// Box 容器（用于背景、占位）
type Box struct {
	style Style
}

func (b Box) Render(ctx *Context) {
	bg := ctx.Style.BackgroundColor
	if bg == nil {
		bg = color.White
	}
	println("Render Box BG at", ctx.X, ctx.Y, "Size:", ctx.Width, "x", ctx.Height)
}

func (b Box) getStyle() Style {
	return b.style
}

func (b Box) WithStyle(mods ...StyleMod) UIElement {
	b.style = applyStyleMods(b.style, mods...)
	return b
}

func NewBox() Box {
	return Box{}
}
