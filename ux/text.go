package ux

import (
	"cmp"
)

// text.go
type Text string

func (t Text) Mount(ctx *Context) {
	style := ctx.Style

	textStyle := DefaultTextStyle

	// 从 Style.Data 中提取 TextStyle
	if ts, ok := ctx.Style.Data.(TextStyle); ok {
		textStyle = ts
	} else if tsPtr, ok := ctx.Style.Data.(*TextStyle); ok {
		textStyle = *tsPtr
	}
	width := cmp.Or(style.Width, textStyle.FontSize*float32(len(t))*0.6)

	height := cmp.Or(style.Height, textStyle.FontSize*textStyle.LineHeight)

	x := ctx.X + style.Margin[3]
	y := ctx.Y + style.Margin[0]

	// 使用 TextStyle 的颜色优先级高于通用 Color
	finalColor := cmp.Or(textStyle.Color, style.Color, colorBlack)

	ctx.AddCmd(RenderCmd{
		Type:   CmdText,
		Rect:   Rect{x, y, width, height},
		Text:   string(t),
		Color:  finalColor,
		Data:   textStyle, // 传递 TextStyle 以供渲染使用
		ZIndex: style.ZIndex,
	})

	// 更新流式布局位置
	ctx.X += width + ctx.Margin[1] + ctx.Margin[3]
}

// 便捷构造函数
func NewText(s string, mods ...StyleMod) Widget {
	return WithStyle(Text(s), mods...)
}

func WithTextStyle(mods ...func(*TextStyle)) StyleMod {
	return func(s *Style) {
		ts := getDefaultTextStyle(s)
		for _, mod := range mods {
			mod(&ts)
		}
		s.Data = ts // 存入 Data，供 Text 使用
	}
}
