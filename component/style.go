package component

import "github.com/millken/yoga"

// Style 定义组件样式
type Style struct {
	// 布局尺寸（支持 px/percent/auto）
	Width     *Dimension
	Height    *Dimension
	MinWidth  *Dimension
	MinHeight *Dimension
	MaxWidth  *Dimension
	MaxHeight *Dimension

	// Flexbox属性
	Flex           *float32
	FlexGrow       *float32
	FlexShrink     *float32
	FlexBasis      *Dimension // auto/px/%
	FlexDirection  *yoga.FlexDirection
	JustifyContent *yoga.Justify
	AlignItems     *yoga.Align
	AlignSelf      *yoga.Align
	AlignContent   *yoga.Align
	FlexWrap       *yoga.Wrap

	// 位置与偏移
	Position *yoga.PositionType
	Left     *Dimension
	Top      *Dimension
	Right    *Dimension
	Bottom   *Dimension

	// 边距 / 内边距 / 边框（支持简写与按边设置）
	Margin  *Edges
	Padding *Edges
	Border  *Edges // 这里仅表示边框宽度，颜色见下

	BorderRadius *Dimension // 简化为统一圆角，后续可扩展为四角
	BorderColor  *string

	// 文本属性
	Color    *string
	FontSize *float32
	Font     *string

	// 背景属性
	BackgroundColor *string

	// 其他属性
	AspectRatio *float32
	Overflow    *yoga.Overflow
	Display     *yoga.Display
	Gap         *Dimension // 行列间距（容器）

	// 交互属性
	Opacity       *float32
	Cursor        *string
	Transition    *string
	PointerEvents *string

	// CSS属性（用于调试和未来扩展）
	ZIndex    *int
	BoxShadow *string
	Transform *string
}

// Unit 表示数值单位
type Unit uint8

const (
	UnitUndefined Unit = iota // 未设置
	UnitAuto                  // 自适应
	UnitPoint                 // 物理像素/pt
	UnitPercent               // 百分比（0..100）
)

// Dimension 带单位的尺寸
type Dimension struct {
	Value float32
	Unit  Unit
}

// Edges 边集（支持 CSS 样式的简写与方向）
type Edges struct {
	All        *Dimension
	Horizontal *Dimension
	Vertical   *Dimension
	Left       *Dimension
	Top        *Dimension
	Right      *Dimension
	Bottom     *Dimension
}

// Dimension helpers
func Px(v float32) *Dimension  { return &Dimension{Value: v, Unit: UnitPoint} }
func Pc(v float32) *Dimension  { return &Dimension{Value: v, Unit: UnitPercent} }
func AutoDim() *Dimension      { return &Dimension{Unit: UnitAuto} }
func UndefinedDim() *Dimension { return &Dimension{Unit: UnitUndefined} }
