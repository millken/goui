package ux

// 通用样式装饰器
type StyledWidget struct {
	Widget Widget
	Mods   []StyleMod
}

func (s *StyledWidget) Mount(ctx *Context) {
	// 临时保存原样式
	original := ctx.Style

	// 应用样式
	for _, mod := range s.Mods {
		mod(&ctx.Style)
	}

	// 挂载内部组件
	s.Widget.Mount(ctx)

	// 恢复样式（避免污染后续组件）
	ctx.Style = original
}

// WithStyle 创建一个带样式的 Widget
func WithStyle(widget Widget, mods ...StyleMod) Widget {
	return &StyledWidget{
		Widget: widget,
		Mods:   mods,
	}
}
