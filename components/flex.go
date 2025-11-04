package components

import (
	"image"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui"
	"github.com/millken/goui/component"
	"github.com/millken/yoga"
)

// Flex 弹性布局容器
type Flex struct {
	*goui.BaseComponent
	direction Direction
	wrap      Wrap
	justify   Justify
	align     Align
	alignSelf Align
	gap       float32
}

// 确保Flex实现了Component接口
var _ component.Component = (*Flex)(nil)

// NewFlex 创建弹性布局容器
func NewFlex(id string) *Flex {
	base := goui.NewBaseComponent(id)
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

	// 设置间距（如果Yoga版本支持gap）
	if f.gap > 0 {
		// 注意：部分Yoga版本可能不支持gap属性
		// node.SetGap(f.gap) // 根据实际Yoga版本调整
	}
}

// SetDirection 设置主轴方向
func (f *Flex) SetDirection(direction Direction) {
	if f.direction != direction {
		f.direction = direction
		f.GetYogaNode().SetFlexDirection(direction.toYoga())
		f.InvalidateLayout()
	}
}

// SetWrap 设置换行
func (f *Flex) SetWrap(wrap Wrap) {
	if f.wrap != wrap {
		f.wrap = wrap
		f.GetYogaNode().SetFlexWrap(wrap.toYoga())
		f.InvalidateLayout()
	}
}

// SetJustify 设置主轴对齐
func (f *Flex) SetJustify(justify Justify) {
	if f.justify != justify {
		f.justify = justify
		f.GetYogaNode().SetJustifyContent(justify.toYoga())
		f.InvalidateLayout()
	}
}

// SetAlign 设置交叉轴对齐
func (f *Flex) SetAlign(align Align) {
	if f.align != align {
		f.align = align
		f.GetYogaNode().SetAlignItems(align.toYoga())
		f.InvalidateLayout()
	}
}

// SetAlignSelf 设置自身对齐
func (f *Flex) SetAlignSelf(alignSelf Align) {
	if f.alignSelf != alignSelf {
		f.alignSelf = alignSelf
		f.GetYogaNode().SetAlignSelf(alignSelf.toYoga())
		f.InvalidateLayout()
	}
}

// SetGap 设置间距
func (f *Flex) SetGap(gap float32) {
	if f.gap != gap {
		f.gap = gap
		// f.GetNode().SetGap(gap) // 根据Yoga版本支持情况
		f.InvalidateLayout()
	}
}

// Render 渲染容器（可选择性绘制背景等）
func (f *Flex) Render(canvas *thorvg.Canvas) {

}

// renderBackground 渲染背景
func (f *Flex) renderBackground(canvas *thorvg.Canvas, region image.Rectangle) {
	// 获取背景色（从样式或属性中获取）
	bgColor := f.getBackgroundColor()
	if bgColor == nil {
		return
	}

	bounds := f.Bounds().Intersect(region)
	if bounds.Empty() {
		return
	}

	bg := thorvg.NewShape()
	if bg != nil {
		bg.AppendRect(
			float32(bounds.Min.X), float32(bounds.Min.Y),
			float32(bounds.Dx()), float32(bounds.Dy()),
			0, 0, false,
		)
		// 应用背景色
		// bg.SetFillColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
		canvas.Push(&bg.Paint)
		bg.Release()
	}
}

// renderBorder 渲染边框
func (f *Flex) renderBorder(canvas *thorvg.Canvas, region image.Rectangle) {
	// 获取边框样式
	borderWidth := f.getBorderWidth()
	borderColor := f.getBorderColor()
	if borderWidth <= 0 || borderColor == nil {
		return
	}

	bounds := f.Bounds().Intersect(region)
	if bounds.Empty() {
		return
	}

	border := thorvg.NewShape()
	if border != nil {
		border.AppendRect(
			float32(bounds.Min.X), float32(bounds.Min.Y),
			float32(bounds.Dx()), float32(bounds.Dy()),
			0, 0, false,
		)
		// 应用边框样式
		// border.SetStroke(borderColor.R, borderColor.G, borderColor.B, borderColor.A)
		// border.SetStrokeWidth(borderWidth)
		canvas.Push(&border.Paint)
		border.Release()
	}
}

// getBackgroundColor 获取背景色（示例方法）
func (f *Flex) getBackgroundColor() *struct{ R, G, B, A uint8 } {
	// TODO: 从样式系统获取背景色
	return nil
}

// getBorderColor 获取边框颜色（示例方法）
func (f *Flex) getBorderColor() *struct{ R, G, B, A uint8 } {
	// TODO: 从样式系统获取边框颜色
	return nil
}

// getBorderWidth 获取边框宽度（示例方法）
func (f *Flex) getBorderWidth() float32 {
	// TODO: 从样式系统获取边框宽度
	return 0
}

// AddChild 添加子组件（重写以添加特定的布局逻辑）
func (f *Flex) AddChild(child goui.Component) {
	f.BaseComponent.AddChild(child)

	// 可以为子组件设置特定的Flex属性
	// 例如：根据子组件类型设置默认的flex属性
	f.configureChildFlex(child)
}

// configureChildFlex 配置子组件的Flex属性
func (f *Flex) configureChildFlex(child goui.Component) {
	// 根据子组件类型设置默认的flex属性
	// 暂时简化处理，避免循环依赖
	// 后续可以通过接口检查来优化
	if childBase, ok := child.(*goui.BaseComponent); ok {
		childBase.GetYogaNode().SetFlexGrow(0)
	}
}

// 其他便利方法...

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
func Center(id string) *Flex {
	flex := NewFlex(id)
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
