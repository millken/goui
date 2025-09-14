package ux

import "sort"

// Context 渲染上下文（可扩展）
type Context struct {
	X, Y     float32
	Width    float32
	Height   float32
	Commands []RenderCmd
	NextZ    int16 //自动递增的 Z 值
	Padding  [4]float32
	Margin   [4]float32
	Style    Style
}

func (ctx *Context) AddCmd(cmd RenderCmd) {
	// 如果命令未指定 ZIndex，使用 Style.ZIndex
	if cmd.ZIndex == 0 {
		cmd.ZIndex = ctx.Style.ZIndex
	}
	// 如果仍未设置，使用自动递增
	if cmd.ZIndex == 0 {
		cmd.ZIndex = ctx.NextZ
		ctx.NextZ++
	}
	ctx.Commands = append(ctx.Commands, cmd)
}

// AddWidget 挂载一个 Widget，触发其布局逻辑
func (ctx *Context) AddWidget(w Widget) {
	if w != nil {
		w.Mount(ctx)
	}
}

func (ctx *Context) EndLayout() []RenderCmd {
	commands := ctx.Commands
	sort.SliceStable(commands, func(i, j int) bool {
		return commands[i].ZIndex < commands[j].ZIndex
	})
	return commands
}
func (ctx *Context) Reset() { ctx.Commands = ctx.Commands[:0] }
