package component

import (
	"github.com/millken/go/thorvg"
	gouicolor "github.com/millken/goui/color"
	"github.com/millken/yoga"
)

// Flex 弹性布局容器
type Flex struct {
	*BaseComponent
	direction    Direction
	wrap         Wrap
	justify      Justify
	align        Align
	alignContent Align
	alignSelf    Align
	gap          float32
	rowGap       float32
	columnGap    float32
	// padding (in px)
	padLeft   float32
	padTop    float32
	padRight  float32
	padBottom float32
	// debug drawing
	debug      bool
	debugColor gouicolor.Color
}

// 确保Flex实现了Component接口
var _ Component = (*Flex)(nil)

func NewFlex(id ...string) *Flex {
	base := NewBaseComponent(id...)
	flex := &Flex{
		BaseComponent: base,
		direction:     DirectionRow,
		wrap:          WrapNoWrap,
		justify:       JustifyStart,
		align:         AlignStart,
		alignContent:  AlignStart,
		alignSelf:     AlignAuto,
		gap:           0,
		rowGap:        0,
		columnGap:     0,
		padLeft:       0,
		padTop:        0,
		padRight:      0,
		padBottom:     0,
		debug:         false,
		debugColor:    gouicolor.NewColor(0, 0, 0, 20),
	}
	// 配置Yoga节点
	flex.configureYogaNode()
	return flex
}

// configureYogaNode 配置Yoga节点（内部方法）
func (f *Flex) configureYogaNode() {
	node := f.GetYogaNode()
	if node == nil {
		return
	}

	node.SetDisplay(yoga.DisplayFlex)
	node.SetFlexDirection(f.direction.toYoga())
	node.SetFlexWrap(f.wrap.toYoga())
	node.SetJustifyContent(f.justify.toYoga())
	node.SetAlignItems(f.align.toYoga())
	node.SetAlignContent(f.alignContent.toYoga())
	node.SetAlignSelf(f.alignSelf.toYoga())

	// 间距：优先使用行/列独立 gap；否则使用统一 gap
	if f.rowGap > 0 || f.columnGap > 0 {
		node.SetGap(yoga.GutterRow, f.rowGap)
		node.SetGap(yoga.GutterColumn, f.columnGap)
	} else if f.gap > 0 {
		node.SetGap(yoga.GutterAll, f.gap)
	}

	// 内边距
	if f.padLeft != 0 {
		node.SetPadding(yoga.EdgeLeft, f.padLeft)
	}
	if f.padTop != 0 {
		node.SetPadding(yoga.EdgeTop, f.padTop)
	}
	if f.padRight != 0 {
		node.SetPadding(yoga.EdgeRight, f.padRight)
	}
	if f.padBottom != 0 {
		node.SetPadding(yoga.EdgeBottom, f.padBottom)
	}
}

// Paint 按 Yoga 布局渲染子组件
func (f *Flex) Paint(vg *thorvg.Canvas) bool {
	rendered := false
	// 仅在启用调试时绘制背景辅助
	if f.debug && f.IsDirty() {
		l := f.ComputedLayout()
		if bg := thorvg.NewShape(); bg != nil {
			bg.AppendRect(l.Left, l.Top, l.Width, l.Height, 0, 0, false)
			r, g, b, a := f.debugColor.RGBA()
			bg.SetFillColor(r, g, b, a)
			vg.Push(&bg.Paint)
			bg.Release()
			rendered = true
		}
		f.ClearDirty()
	}
	// 递归渲染子组件，聚合它们的渲染结果
	for _, child := range f.Children() {
		if child.Paint(vg) {
			rendered = true
		}
	}
	return rendered
}

// SetDirection 设置主轴方向
func (f *Flex) SetDirection(direction Direction) {
	if f.direction != direction {
		f.direction = direction
		f.GetYogaNode().SetFlexDirection(direction.toYoga())
		f.MarkDirty(false)
	}
}

// SetWrap 设置换行方式
func (f *Flex) SetWrap(wrap Wrap) {
	if f.wrap != wrap {
		f.wrap = wrap
		f.GetYogaNode().SetFlexWrap(wrap.toYoga())
		f.MarkDirty(false)
	}
}

// SetJustify 设置主轴对齐方式
func (f *Flex) SetJustify(justify Justify) {
	if f.justify != justify {
		f.justify = justify
		f.GetYogaNode().SetJustifyContent(justify.toYoga())
		f.MarkDirty(false)
	}
}

// SetAlign 设置交叉轴对齐方式
func (f *Flex) SetAlign(align Align) {
	if f.align != align {
		f.align = align
		f.GetYogaNode().SetAlignItems(align.toYoga())
		f.MarkDirty(false)
	}
}

// SetAlignSelf 设置自身对齐方式
func (f *Flex) SetAlignSelf(alignSelf Align) {
	if f.alignSelf != alignSelf {
		f.alignSelf = alignSelf
		f.GetYogaNode().SetAlignSelf(alignSelf.toYoga())
		f.MarkDirty(false)
	}
}

// SetGap 设置主轴/交叉轴间距（行列间隔），单位为像素
func (f *Flex) SetGap(gap float32) {
	if f.gap != gap {
		f.gap = gap
		f.rowGap, f.columnGap = 0, 0
		f.GetYogaNode().SetGap(yoga.GutterAll, gap)
		f.MarkDirty(false)
	}
}

// SetRowGap 设置行间距（换行后各行之间的垂直间距）
func (f *Flex) SetRowGap(gap float32) {
	if f.rowGap != gap {
		f.rowGap = gap
		f.gap = 0
		f.GetYogaNode().SetGap(yoga.GutterRow, gap)
		f.MarkDirty(false)
	}
}

// SetColumnGap 设置列间距（主轴方向相邻子项间距）
func (f *Flex) SetColumnGap(gap float32) {
	if f.columnGap != gap {
		f.columnGap = gap
		f.gap = 0
		f.GetYogaNode().SetGap(yoga.GutterColumn, gap)
		f.MarkDirty(false)
	}
}

// SetPadding 设置四边相同的内边距
func (f *Flex) SetPadding(all float32) {
	f.SetPaddingEdges(all, all, all, all)
}

// SetPaddingXY 设置水平/垂直内边距
func (f *Flex) SetPaddingXY(horizontal, vertical float32) {
	f.SetPaddingEdges(horizontal, vertical, horizontal, vertical)
}

// SetPaddingEdges 设置分别的内边距：左、上、右、下
func (f *Flex) SetPaddingEdges(left, top, right, bottom float32) {
	if f.padLeft == left && f.padTop == top && f.padRight == right && f.padBottom == bottom {
		return
	}
	f.padLeft, f.padTop, f.padRight, f.padBottom = left, top, right, bottom
	node := f.GetYogaNode()
	node.SetPadding(yoga.EdgeLeft, left)
	node.SetPadding(yoga.EdgeTop, top)
	node.SetPadding(yoga.EdgeRight, right)
	node.SetPadding(yoga.EdgeBottom, bottom)
	f.MarkDirty(false)
}

// SetAlignContent 设置多行内容在交叉轴上的对齐（当 Wrap 启用并产生多行时生效）
func (f *Flex) SetAlignContent(align Align) {
	if f.alignContent != align {
		f.alignContent = align
		f.GetYogaNode().SetAlignContent(align.toYoga())
		f.MarkDirty(false)
	}
}

// EnableDebug 开关调试背景绘制
func (f *Flex) EnableDebug(enable bool) {
	if f.debug != enable {
		f.debug = enable
		f.MarkDirty(false)
	}
}

// SetDebugColor 设置调试背景颜色
func (f *Flex) SetDebugColor(c gouicolor.Color) {
	f.debugColor = c
	if f.debug {
		f.MarkDirty(false)
	}
}

// Row 创建水平布局
func Row(id string) *Flex {
	flex := NewFlex(id)
	flex.SetDirection(DirectionRow)
	return flex
}

// Column 创建垂直布局
func Column(id string) *Flex {
	flex := NewFlex(id)
	flex.SetDirection(DirectionColumn)
	return flex
}

// Center 创建居中布局
func Center(id ...string) *Flex {
	flex := NewFlex(id...)
	flex.SetJustify(JustifyCenter)
	flex.SetAlign(AlignCenter)
	return flex
}

// SpaceBetween 创建两端对齐布局
func SpaceBetween(id string) *Flex {
	flex := NewFlex(id)
	flex.SetJustify(JustifySpaceBetween)
	return flex
}

// SpaceAround 创建环绕对齐布局
func SpaceAround(id string) *Flex {
	flex := NewFlex(id)
	flex.SetJustify(JustifySpaceAround)
	return flex
}

// SpaceEvenly 创建均匀对齐布局
func SpaceEvenly(id string) *Flex {
	flex := NewFlex(id)
	flex.SetJustify(JustifySpaceEvenly)
	return flex
}

// Bounds 使用 BaseComponent 的 Yoga 计算结果
// 若后续需要为 Flex 提供内边距/边框，可在这里叠加样式计算
