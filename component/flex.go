package component

import "github.com/millken/yoga"

// Flex 弹性布局容器
type Flex struct {
	*BaseComponent
	direction Direction
	wrap      Wrap
	justify   Justify
	align     Align
	alignSelf Align
	gap       float32
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
		alignSelf:     AlignAuto,
		gap:           0,
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
	node.SetAlignSelf(f.alignSelf.toYoga())

	if f.gap > 0 {
		node.SetGap(yoga.GutterAll, f.gap) // 根据实际Yoga版本调整
	}
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
