package component

import "github.com/millken/yoga"

// Direction 布局方向 - 用户友好的方向定义
type Direction int

const (
	DirectionHorizontal Direction = iota
	DirectionVertical
	DirectionRow
	DirectionColumn
)

// Align 对齐方式 - 用户友好的对齐定义
type Align int

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
	AlignStretch
	AlignAuto
)

// Justify 主轴对齐 - 用户友好的主轴对齐定义
type Justify int

const (
	JustifyStart Justify = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Wrap 换行方式 - 用户友好的换行定义
type Wrap int

const (
	WrapNoWrap Wrap = iota
	WrapWrap
	WrapWrapReverse
)

// 内部转换函数 - 将用户友好的枚举转换为 Yoga 常量
func (d Direction) toYoga() yoga.FlexDirection {
	switch d {
	case DirectionHorizontal, DirectionRow:
		return yoga.FlexDirectionRow
	case DirectionVertical, DirectionColumn:
		return yoga.FlexDirectionColumn
	default:
		return yoga.FlexDirectionRow
	}
}

func (a Align) toYoga() yoga.Align {
	switch a {
	case AlignStart:
		return yoga.AlignFlexStart
	case AlignCenter:
		return yoga.AlignCenter
	case AlignEnd:
		return yoga.AlignFlexEnd
	case AlignStretch:
		return yoga.AlignStretch
	default:
		return yoga.AlignFlexStart
	}
}

func (j Justify) toYoga() yoga.Justify {
	switch j {
	case JustifyStart:
		return yoga.JustifyFlexStart
	case JustifyCenter:
		return yoga.JustifyCenter
	case JustifyEnd:
		return yoga.JustifyFlexEnd
	case JustifySpaceBetween:
		return yoga.JustifySpaceBetween
	case JustifySpaceAround:
		return yoga.JustifySpaceAround
	case JustifySpaceEvenly:
		return yoga.JustifySpaceEvenly
	default:
		return yoga.JustifyFlexStart
	}
}

func (w Wrap) toYoga() yoga.Wrap {
	switch w {
	case WrapNoWrap:
		return yoga.WrapNoWrap
	case WrapWrap:
		return yoga.WrapWrap
	case WrapWrapReverse:
		return yoga.WrapWrapReverse
	default:
		return yoga.WrapNoWrap
	}
}
