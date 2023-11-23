package yoga

import "sync/atomic"

type Node struct {
	Style              Style
	Layout             *Layout
	lineIndex          uint32
	Owner              *Node
	Children           []*Node
	NextChild          *Node
	Config             *Config
	NodeType           NodeType
	Measure            MeasureFunc
	BaseLine           BaseLineFunc
	Print              PrintFunc
	Context            interface{}
	IsDirty            bool
	hasNewLayout       bool
	resolvedDimensions [2]*Value
}

var (
	defaultNode = Node{
		Style:              DefaultStyle,
		Layout:             DefaultLayout,
		Owner:              nil,
		Children:           nil,
		hasNewLayout:       true,
		IsDirty:            false,
		NodeType:           NodeTypeDefault,
		resolvedDimensions: [2]*Value{&ValueUndefined, &ValueUndefined},
	}
)

// NewNode creates new node
func NewNode() *Node {
	return NewNodeWithConfig(&defaultConfig)
}

// NewNodeWithConfig creates new node with config
func NewNodeWithConfig(config *Config) *Node {
	node := &Node{}
	*node = defaultNode
	if config.UseWebDefaults {
		node.Style.FlexDirection = FlexDirectionRow
		node.Style.AlignContent = AlignStretch
	}
	node.Config = config
	return node
}

func (n *Node) Reset() {
	cfg := n.Config
	*n = defaultNode
	if cfg.UseWebDefaults {
		n.Style.FlexDirection = FlexDirectionRow
		n.Style.AlignContent = AlignStretch
	}
	n.Config = cfg
}

func (n *Node) SetWidth(width float32) *Node {
	dim := &n.Style.Dimensions[DimensionWidth]
	if dim.Unit != UnitPoint || dim.Value != width {
		dim.Value = width
		dim.Unit = UnitPoint
		if width == Undefined {
			dim.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) SetWidthPercent(width float32) *Node {
	dim := &n.Style.Dimensions[DimensionWidth]
	if dim.Unit != UnitPercent || dim.Value != width {
		dim.Value = width
		dim.Unit = UnitPercent
		if width == Undefined {
			dim.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) SetWidthAuto() *Node {
	dim := &n.Style.Dimensions[DimensionWidth]
	if dim.Unit != UnitAuto {
		dim.Value = Undefined
		dim.Unit = UnitAuto
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetWidth() Value {
	return n.Style.Dimensions[DimensionWidth]
}

func (n *Node) SetHeight(height float32) *Node {
	dim := &n.Style.Dimensions[DimensionHeight]
	if dim.Unit != UnitPoint || dim.Value != height {
		dim.Value = height
		dim.Unit = UnitPoint
		if height == Undefined {
			dim.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) SetHeightPercent(height float32) *Node {
	dim := &n.Style.Dimensions[DimensionHeight]
	if dim.Unit != UnitPercent || dim.Value != height {
		dim.Value = height
		dim.Unit = UnitPercent
		if height == Undefined {
			dim.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) SetHeightAuto() *Node {
	dim := &n.Style.Dimensions[DimensionHeight]
	if dim.Unit != UnitAuto {
		dim.Value = Undefined
		dim.Unit = UnitAuto
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetHeight() Value {
	return n.Style.Dimensions[DimensionHeight]
}

// MarkDirtyAndPropagate marks node as dirty
func (n *Node) MarkDirtyAndPropagate() {
	if !n.IsDirty {
		n.IsDirty = true
		n.Layout.computedFlexBasis = undefinedFloatOptional
		if n.Owner != nil {
			n.Owner.MarkDirtyAndPropagate()
		}
	}
}

func (n *Node) iterChildrenAfterCloningIfNeeded(callback func(child *Node, layoutContext interface{}), cloneContext interface{}) {
	for i, child := range n.Children {
		if child.Owner != n {
			n.Children[i] = n.Config.cloneNode(child)
			n.Children[i].Owner = n.Owner
		}
		callback(n.Children[i], cloneContext)
	}
}

func (n *Node) IsNodeDirty() bool {
	return n.IsDirty
}

func (n *Node) SetPositionType(positionType PositionType) *Node {
	if n.Style.PositionType != positionType {
		n.Style.PositionType = positionType
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetPositionType() PositionType {
	return n.Style.PositionType
}

func (n *Node) SetPosition(edge Edge, position float32) *Node {
	current := n.Style.Position[edge]
	if current.Unit != UnitPoint || current.Value != position {
		current.Unit = UnitPoint
		current.Value = position
		if position == Undefined {
			current.Unit = UnitUndefined
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) SetPositionPercent(edge Edge, position float32) *Node {
	current := n.Style.Position[edge]
	if current.Unit != UnitPercent || current.Value != position {
		current.Unit = UnitPercent
		current.Value = position
		if position == Undefined {
			current.Unit = UnitUndefined
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetPosition(edge Edge) Value {
	return n.Style.Position[edge]
}

func (n *Node) SetDirection(direction Direction) *Node {
	if n.Style.Direction != direction {
		n.Style.Direction = direction
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetDirection() Direction {
	return n.Style.Direction
}

func (n *Node) SetFlexDirection(flexDirection FlexDirection) *Node {
	if n.Style.FlexDirection != flexDirection {
		n.Style.FlexDirection = flexDirection
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetFlexDirection() FlexDirection {
	return n.Style.FlexDirection
}

func (n *Node) SetJustifyContent(justifyContent Justify) *Node {
	if n.Style.JustifyContent != justifyContent {
		n.Style.JustifyContent = justifyContent
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetJustifyContent() Justify {
	return n.Style.JustifyContent
}

func (n *Node) SetAlignContent(alignContent Align) *Node {
	if n.Style.AlignContent != alignContent {
		n.Style.AlignContent = alignContent
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetAlignContent() Align {
	return n.Style.AlignContent
}

func (n *Node) SetAlignItems(alignItems Align) *Node {
	if n.Style.AlignItems != alignItems {
		n.Style.AlignItems = alignItems
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetAlignItems() Align {
	return n.Style.AlignItems
}

func (n *Node) SetAlignSelf(alignSelf Align) *Node {
	if n.Style.AlignSelf != alignSelf {
		n.Style.AlignSelf = alignSelf
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetAlignSelf() Align {
	return n.Style.AlignSelf
}

func (n *Node) SetFlexWrap(flexWrap Wrap) *Node {
	if n.Style.FlexWrap != flexWrap {
		n.Style.FlexWrap = flexWrap
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetFlexWrap() Wrap {
	return n.Style.FlexWrap
}

func (n *Node) SetOverflow(overflow Overflow) *Node {
	if n.Style.Overflow != overflow {
		n.Style.Overflow = overflow
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetOverflow() Overflow {
	return n.Style.Overflow
}

func (n *Node) SetDisplay(display Display) *Node {
	if n.Style.Display != display {
		n.Style.Display = display
		n.MarkDirtyAndPropagate()
	}
	return n
}

func (n *Node) GetDisplay() Display {
	return n.Style.Display
}

func (n *Node) SetFlex(flex float32) *Node {
	flexOption := NewFloatOptional(flex)
	if n.Style.Flex.Equals(flexOption) {
		return n
	}
	n.Style.Flex = flexOption
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) GetFlex() float32 {
	return n.Style.Flex.GetValue()
}

func (n *Node) SetFlexGrow(flexGrow float32) *Node {
	flexGrowOption := NewFloatOptional(flexGrow)
	if n.Style.FlexGrow.Equals(flexGrowOption) {
		return n
	}
	n.Style.FlexGrow = flexGrowOption
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) GetFlexGrow() float32 {
	return n.Style.FlexGrow.GetValue()
}

func (n *Node) SetFlexShrink(flexShrink float32) *Node {
	flexShrinkOption := NewFloatOptional(flexShrink)
	if n.Style.FlexShrink.Equals(flexShrinkOption) {
		return n
	}
	n.Style.FlexShrink = flexShrinkOption
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) GetFlexShrink() float32 {
	return n.Style.FlexShrink.GetValue()
}

// SetFlexBasis sets flex basis
func (n *Node) SetFlexBasis(flexBasis float32) *Node {
	if n.Style.FlexBasis.Unit != UnitPoint || n.Style.FlexBasis.Value != flexBasis {
		n.Style.FlexBasis.Unit = UnitPoint
		n.Style.FlexBasis.Value = flexBasis
		if flexBasis == Undefined {
			n.Style.FlexBasis.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetFlexBasisPercent sets flex basis percent
func (n *Node) SetFlexBasisPercent(flexBasis float32) *Node {
	if n.Style.FlexBasis.Unit != UnitPercent || n.Style.FlexBasis.Value != flexBasis {
		n.Style.FlexBasis.Unit = UnitPercent
		n.Style.FlexBasis.Value = flexBasis
		if flexBasis == Undefined {
			n.Style.FlexBasis.Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetFlexBasisAuto sets flex basis auto
func (n *Node) SetFlexBasisAuto() *Node {
	if n.Style.FlexBasis.Unit != UnitAuto {
		n.Style.FlexBasis.Unit = UnitAuto
		n.Style.FlexBasis.Value = Undefined
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetFlexBasis returns flex basis
func (n *Node) GetFlexBasis() Value {
	return n.Style.FlexBasis
}

// SetMargin sets margin
func (n *Node) SetMargin(edge Edge, margin float32) *Node {
	if n.Style.Margin[edge].Unit != UnitPoint || n.Style.Margin[edge].Value != margin {
		n.Style.Margin[edge].Unit = UnitPoint
		n.Style.Margin[edge].Value = margin
		if margin == Undefined {
			n.Style.Margin[edge].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMarginPercent sets margin percent
func (n *Node) SetMarginPercent(edge Edge, margin float32) *Node {
	if n.Style.Margin[edge].Unit != UnitPercent || n.Style.Margin[edge].Value != margin {
		n.Style.Margin[edge].Unit = UnitPercent
		n.Style.Margin[edge].Value = margin
		if margin == Undefined {
			n.Style.Margin[edge].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMarginAuto sets margin auto
func (n *Node) SetMarginAuto(edge Edge) *Node {
	if n.Style.Margin[edge].Unit != UnitAuto {
		n.Style.Margin[edge].Unit = UnitAuto
		n.Style.Margin[edge].Value = Undefined
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetMargin returns margin
func (n *Node) GetMargin(edge Edge) Value {
	return n.Style.Margin[edge]
}

// SetPadding sets padding
func (n *Node) SetPadding(edge Edge, padding float32) *Node {
	if n.Style.Padding[edge].Unit != UnitPoint || n.Style.Padding[edge].Value != padding {
		n.Style.Padding[edge].Unit = UnitPoint
		n.Style.Padding[edge].Value = padding
		if padding == Undefined {
			n.Style.Padding[edge].Unit = UnitUndefined
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetPaddingPercent sets padding percent
func (n *Node) SetPaddingPercent(edge Edge, padding float32) *Node {
	if n.Style.Padding[edge].Unit != UnitPercent || n.Style.Padding[edge].Value != padding {
		n.Style.Padding[edge].Unit = UnitPercent
		n.Style.Padding[edge].Value = padding
		if padding == Undefined {
			n.Style.Padding[edge].Unit = UnitUndefined
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetPadding returns padding
func (n *Node) GetPadding(edge Edge) Value {
	return n.Style.Padding[edge]
}

// SetBorder sets border
func (n *Node) SetBorder(edge Edge, border float32) *Node {
	if n.Style.Border[edge].Unit != UnitPoint || n.Style.Border[edge].Value != border {
		n.Style.Border[edge].Unit = UnitPoint
		n.Style.Border[edge].Value = border
		if border == Undefined {
			n.Style.Border[edge].Unit = UnitUndefined
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetBorder returns border
func (n *Node) GetBorder(edge Edge) Value {
	return n.Style.Border[edge]
}

// SetMinWidth sets min width
func (n *Node) SetMinWidth(minWidth float32) *Node {
	if n.Style.MinDimensions[DimensionWidth].Unit != UnitPoint || n.Style.MinDimensions[DimensionWidth].Value != minWidth {
		n.Style.MinDimensions[DimensionWidth].Unit = UnitPoint
		n.Style.MinDimensions[DimensionWidth].Value = minWidth
		if minWidth == Undefined {
			n.Style.MinDimensions[DimensionWidth].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMinWidthPercent sets min width percent
func (n *Node) SetMinWidthPercent(minWidth float32) *Node {
	if n.Style.MinDimensions[DimensionWidth].Unit != UnitPercent || n.Style.MinDimensions[DimensionWidth].Value != minWidth {
		n.Style.MinDimensions[DimensionWidth].Unit = UnitPercent
		n.Style.MinDimensions[DimensionWidth].Value = minWidth
		if minWidth == Undefined {
			n.Style.MinDimensions[DimensionWidth].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetMinWidth returns min width
func (n *Node) GetMinWidth() Value {
	return n.Style.MinDimensions[DimensionWidth]
}

// SetMinHeight sets min height
func (n *Node) SetMinHeight(minHeight float32) *Node {
	if n.Style.MinDimensions[DimensionHeight].Unit != UnitPoint || n.Style.MinDimensions[DimensionHeight].Value != minHeight {
		n.Style.MinDimensions[DimensionHeight].Unit = UnitPoint
		n.Style.MinDimensions[DimensionHeight].Value = minHeight
		if minHeight == Undefined {
			n.Style.MinDimensions[DimensionHeight].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMinHeightPercent sets min height percent
func (n *Node) SetMinHeightPercent(minHeight float32) *Node {
	if n.Style.MinDimensions[DimensionHeight].Unit != UnitPercent || n.Style.MinDimensions[DimensionHeight].Value != minHeight {
		n.Style.MinDimensions[DimensionHeight].Unit = UnitPercent
		n.Style.MinDimensions[DimensionHeight].Value = minHeight
		if minHeight == Undefined {
			n.Style.MinDimensions[DimensionHeight].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetMinHeight returns min height
func (n *Node) GetMinHeight() Value {
	return n.Style.MinDimensions[DimensionHeight]
}

// SetMaxWidth sets max width
func (n *Node) SetMaxWidth(maxWidth float32) *Node {
	if n.Style.MaxDimensions[DimensionWidth].Unit != UnitPoint || n.Style.MaxDimensions[DimensionWidth].Value != maxWidth {
		n.Style.MaxDimensions[DimensionWidth].Unit = UnitPoint
		n.Style.MaxDimensions[DimensionWidth].Value = maxWidth
		if maxWidth == Undefined {
			n.Style.MaxDimensions[DimensionWidth].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMaxWidthPercent sets max width percent
func (n *Node) SetMaxWidthPercent(maxWidth float32) *Node {
	if n.Style.MaxDimensions[DimensionWidth].Unit != UnitPercent || n.Style.MaxDimensions[DimensionWidth].Value != maxWidth {
		n.Style.MaxDimensions[DimensionWidth].Unit = UnitPercent
		n.Style.MaxDimensions[DimensionWidth].Value = maxWidth
		if maxWidth == Undefined {
			n.Style.MaxDimensions[DimensionWidth].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetMaxWidth returns max width
func (n *Node) GetMaxWidth() Value {
	return n.Style.MaxDimensions[DimensionWidth]
}

// SetMaxHeight sets max height
func (n *Node) SetMaxHeight(maxHeight float32) *Node {
	if n.Style.MaxDimensions[DimensionHeight].Unit != UnitPoint || n.Style.MaxDimensions[DimensionHeight].Value != maxHeight {
		n.Style.MaxDimensions[DimensionHeight].Unit = UnitPoint
		n.Style.MaxDimensions[DimensionHeight].Value = maxHeight
		if maxHeight == Undefined {
			n.Style.MaxDimensions[DimensionHeight].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetMaxHeightPercent sets max height percent
func (n *Node) SetMaxHeightPercent(maxHeight float32) *Node {
	if n.Style.MaxDimensions[DimensionHeight].Unit != UnitPercent || n.Style.MaxDimensions[DimensionHeight].Value != maxHeight {
		n.Style.MaxDimensions[DimensionHeight].Unit = UnitPercent
		n.Style.MaxDimensions[DimensionHeight].Value = maxHeight
		if maxHeight == Undefined {
			n.Style.MaxDimensions[DimensionHeight].Unit = UnitAuto
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// GetMaxHeight returns max height
func (n *Node) GetMaxHeight() Value {
	return n.Style.MaxDimensions[DimensionHeight]
}

// SetAspectRatio sets aspect ratio
func (n *Node) SetAspectRatio(aspectRatio float32) *Node {
	aspectRatioOption := NewFloatOptional(aspectRatio)
	if n.Style.AspectRatio.Equals(aspectRatioOption) {
		return n
	}
	n.Style.AspectRatio = aspectRatioOption
	n.MarkDirtyAndPropagate()
	return n
}

// GetAspectRatio returns aspect ratio
func (n *Node) GetAspectRatio() float32 {
	return n.Style.AspectRatio.GetValue()
}

// LayoutLeft returns left
func (n *Node) LayoutLeft() float32 {
	return n.Layout.Position[EdgeLeft]
}

// LayoutTop returns top
func (n *Node) LayoutTop() float32 {
	return n.Layout.Position[EdgeTop]
}

// LayoutRight returns right
func (n *Node) LayoutRight() float32 {
	return n.Layout.Position[EdgeRight]
}

// LayoutBottom returns bottom
func (n *Node) LayoutBottom() float32 {
	return n.Layout.Position[EdgeBottom]
}

// LayoutWidth returns width
func (n *Node) LayoutWidth() float32 {
	return n.Layout.Dimensions[DimensionWidth]
}

// LayoutHeight returns height
func (n *Node) LayoutHeight() float32 {
	return n.Layout.Dimensions[DimensionHeight]
}

// LayoutMargin returns margin
func (n *Node) LayoutMargin(edge Edge) float32 {
	switch edge {
	case EdgeLeft:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Margin[EdgeEnd]
		}
		return n.Layout.Margin[EdgeStart]
	case EdgeRight:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Margin[EdgeStart]
		}
		return n.Layout.Margin[EdgeEnd]
	default:
		return n.Layout.Margin[edge]
	}
}

// LayoutBorder returns border
func (n *Node) LayoutBorder(edge Edge) float32 {
	switch edge {
	case EdgeLeft:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Border[EdgeEnd]
		}
		return n.Layout.Border[EdgeStart]
	case EdgeRight:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Border[EdgeStart]
		}
		return n.Layout.Border[EdgeEnd]
	default:
		return n.Layout.Border[edge]
	}
}

// LayoutPadding returns padding
func (n *Node) LayoutPadding(edge Edge) float32 {
	switch edge {
	case EdgeLeft:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Padding[EdgeEnd]
		}
		return n.Layout.Padding[EdgeStart]
	case EdgeRight:
		if n.Layout.Direction == DirectionRTL {
			return n.Layout.Padding[EdgeStart]
		}
		return n.Layout.Padding[EdgeEnd]
	default:
		return n.Layout.Padding[edge]
	}
}

// InsertChild inserts child
func (n *Node) InsertChild(child *Node, index int) *Node {
	n.Children = append(n.Children[:index], append([]*Node{child}, n.Children[index:]...)...)
	child.Owner = n
	n.MarkDirtyAndPropagate()
	return n
}

// CalculateLayout calculates layout
func (n *Node) CalculateLayout(ownerWidth, ownerHeight float32, ownerDirection Direction) {
	atomic.AddUint32(&currentGenerationCount, 1)
	n.resolveDimension()
	width := Undefined
	widthMeasureMode := MeasureModeUndefined
	maxDimensions := n.Style.MaxDimensions
	if n.isStyleDimDefined(FlexDirectionRow, ownerWidth) {
		width = resolveValue(n.getResolvedDimension(int(dim[FlexDirectionRow])), ownerWidth) + n.getMarginForAxis(FlexDirectionRow, ownerWidth).GetValue()
		widthMeasureMode = MeasureModeExactly
	} else if !ResolveValue(maxDimensions[DimensionWidth], ownerWidth).IsUndefined {
		width = ResolveValue(maxDimensions[DimensionWidth], ownerWidth).GetValue()
		widthMeasureMode = MeasureModeAtMost
	} else {
		width = ownerWidth
		widthMeasureMode = MeasureModeUndefined
	}
	height := Undefined
	heightMeasureMode := MeasureModeUndefined
	if n.isStyleDimDefined(FlexDirectionColumn, ownerHeight) {
		height = resolveValue(n.getResolvedDimension(int(dim[FlexDirectionColumn])), ownerHeight) + n.getMarginForAxis(FlexDirectionColumn, ownerWidth).GetValue()
		heightMeasureMode = MeasureModeExactly
	} else if !ResolveValue(maxDimensions[DimensionHeight], ownerHeight).IsUndefined {
		height = ResolveValue(maxDimensions[DimensionHeight], ownerHeight).GetValue()
		heightMeasureMode = MeasureModeAtMost
	} else {
		height = ownerHeight
		heightMeasureMode = MeasureModeUndefined
	}

}

// isStyleDimDefined
func (n *Node) isStyleDimDefined(axis FlexDirection, ownerSize float32) bool {
	rdim := n.getResolvedDimension(int(dim[axis]))
	isUndefined := FloatIsUndefined(rdim.Value)
	return !(rdim.Unit == UnitAuto ||
		rdim.Unit == UnitUndefined ||
		(rdim.Unit == UnitPoint && !isUndefined && rdim.Value < 0) ||
		(rdim.Unit == UnitPercent && !isUndefined && (rdim.Value < 0 || FloatIsUndefined(ownerSize))))

}

// getConfig
func (n *Node) getConfig() *Config {
	return n.Config
}

// getResolvedDimension
func (n *Node) getResolvedDimension(index int) *Value {
	return n.resolvedDimensions[index]
}

// getResolvedDimensions
func (n *Node) getResolvedDimensions() [2]*Value {
	return n.resolvedDimensions
}

// resolveDimension
func (n *Node) resolveDimension() {
	for dim := DimensionWidth; dim < DimensionCount; dim++ {
		if n.Style.MaxDimensions[dim].Unit != UnitUndefined &&
			ValueEqual(n.Style.MaxDimensions[dim], n.Style.MinDimensions[dim]) {
			n.resolvedDimensions[dim] = &n.Style.MaxDimensions[dim]
		} else {
			n.resolvedDimensions[dim] = &n.Style.Dimensions[dim]
		}
	}
}

func (n *Node) getLeadingPosition(axis FlexDirection, axisSize float32) FloatOptional {
	leadingPosition := ValueZero
	if flexDirectionIsRow(axis) {
		leadingPosition = computeEdgeValueForRow(n.Style.Position, EdgeStart, leading[axis], ValueZero)
	} else {
		leadingPosition = computeEdgeValueForColumn(n.Style.Position, leading[axis], ValueZero)
	}
	return ResolveValue(leadingPosition, axisSize)
}

// getTrailingPosition
func (n *Node) getTrailingPosition(axis FlexDirection, axisSize float32) FloatOptional {
	trailingPosition := ValueZero
	if flexDirectionIsRow(axis) {
		trailingPosition = computeEdgeValueForRow(n.Style.Position, EdgeEnd, trailing[axis], ValueZero)
	} else {
		trailingPosition = computeEdgeValueForColumn(n.Style.Position, trailing[axis], ValueZero)
	}
	return ResolveValue(trailingPosition, axisSize)
}

// isLeadingPositionDefined
func (n *Node) isLeadingPositionDefined(axis FlexDirection) bool {
	leadingPosition := ValueZero
	if flexDirectionIsRow(axis) {
		leadingPosition = computeEdgeValueForRow(n.Style.Position, EdgeStart, leading[axis], ValueUndefined)
	} else {
		leadingPosition = computeEdgeValueForColumn(n.Style.Position, leading[axis], ValueUndefined)
	}
	return !leadingPosition.IsUndefined()
}

// isTrailingPosDefined
func (n *Node) isTrailingPosDefined(axis FlexDirection) bool {
	trailingPosition := ValueZero
	if flexDirectionIsRow(axis) {
		trailingPosition = computeEdgeValueForRow(n.Style.Position, EdgeEnd, trailing[axis], ValueUndefined)
	} else {
		trailingPosition = computeEdgeValueForColumn(n.Style.Position, trailing[axis], ValueUndefined)
	}
	return !trailingPosition.IsUndefined()
}

// getLeadingMargin
func (n *Node) getLeadingMargin(axis FlexDirection, widthSize float32) FloatOptional {
	leadingMargin := ValueZero
	if flexDirectionIsRow(axis) {
		leadingMargin = computeEdgeValueForRow(n.Style.Margin, EdgeStart, leading[axis], ValueZero)
	} else {
		leadingMargin = computeEdgeValueForColumn(n.Style.Margin, leading[axis], ValueZero)
	}
	return ResolveValueMargin(leadingMargin, widthSize)
}

// getTrailingMargin
func (n *Node) getTrailingMargin(axis FlexDirection, widthSize float32) FloatOptional {
	trailingMargin := ValueZero
	if flexDirectionIsRow(axis) {
		trailingMargin = computeEdgeValueForRow(n.Style.Margin, EdgeEnd, trailing[axis], ValueZero)
	} else {
		trailingMargin = computeEdgeValueForColumn(n.Style.Margin, trailing[axis], ValueZero)
	}
	return ResolveValueMargin(trailingMargin, widthSize)
}

// getMarginForAxis
func (n *Node) getMarginForAxis(axis FlexDirection, widthSize float32) FloatOptional {
	return n.getLeadingMargin(axis, widthSize).Add(n.getTrailingMargin(axis, widthSize))
}

// getGapForAxis
func (n *Node) getGapForAxis(axis FlexDirection, widthSize float32) FloatOptional {
	gap := ValueZero
	if flexDirectionIsRow(axis) {
		gap = computeRowGap(n.Style.Gap, ValueZero)
	} else {
		gap = computeColumnGap(n.Style.Gap, ValueZero)
	}
	return ResolveValue(gap, widthSize)
}

// measure
func (n *Node) measure(width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) {
	if n.Measure != nil {
		n.Measure(n, width, widthMode, height, heightMode)
	}
}

// baseline
func (n *Node) baseline(width float32, height float32) float32 {
	return n.BaseLine(n, width, height)
}

// SetMeasureFunc
func (n *Node) SetMeasureFunc(measureFunc MeasureFunc) *Node {
	n.Measure = measureFunc
	return n
}

// ReplaceChild
func (n *Node) ReplaceChild(child *Node, newChild *Node) *Node {
	for i, c := range n.Children {
		if c == child {
			n.Children[i] = newChild
			newChild.Owner = n
			child.Owner = nil
			n.MarkDirtyAndPropagate()
			return n
		}
	}
	return n
}

// SetChild
func (n *Node) SetChild(index int, child *Node) *Node {
	n.Children[index] = child
	n.MarkDirtyAndPropagate()
	return n
}

// SetDirty
func (n *Node) SetDirty(isDirty bool) *Node {
	if n.IsDirty != isDirty {
		n.IsDirty = isDirty
		n.MarkDirtyAndPropagate()
	}
	return n
}

// RemoveChild
func (n *Node) RemoveChild(child *Node) *Node {
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			n.MarkDirtyAndPropagate()
			return n
		}
	}
	return n
}

// setLineIndex
func (n *Node) setLineIndex(lineIndex uint32) {
	n.lineIndex = lineIndex
}

// cloneChildrenIfNeeded
func (n *Node) cloneChildrenIfNeeded() {
	n.iterChildrenAfterCloningIfNeeded(func(child *Node, cloneContext interface{}) {}, nil)
}

// HasErrata
func (n *Node) HasErrata(errata Errata) bool {
	return n.Config.HasErrata(errata)
}

// GetChild
func (n *Node) GetChild(index int) *Node {
	return n.Children[index]
}

// getStyle
func (n *Node) getStyle() *Style {
	return &n.Style
}

// getLineIndex
func (n *Node) getLineIndex() uint32 {
	return n.lineIndex
}

// getLayout
func (n *Node) getLayout() *Layout {
	return n.Layout
}

// getChildren
func (n *Node) getChildren() []*Node {
	return n.Children
}

// SetLayoutDirection
func (n *Node) SetLayoutDirection(direction Direction) *Node {
	if n.Layout.Direction != direction {
		n.Layout.Direction = direction
		n.MarkDirtyAndPropagate()
	}
	return n
}

// SetLayoutMargin
func (n *Node) SetLayoutMargin(edge Edge, margin float32) *Node {
	n.Layout.Margin[edge] = margin
	return n
}

// SetLayoutBorder
func (n *Node) SetLayoutBorder(edge Edge, border float32) *Node {
	n.Layout.Border[edge] = border
	return n
}

// SetLayoutPadding
func (n *Node) SetLayoutPadding(edge Edge, padding float32) *Node {
	n.Layout.Padding[edge] = padding
	return n
}

// SetLayoutLastOwnerDirection
func (n *Node) SetLayoutLastOwnerDirection(direction Direction) *Node {
	n.Layout.lastOwnerDirection = direction
	return n
}

// SetLayoutComputedFlexBasis
func (n *Node) SetLayoutComputedFlexBasis(computedFlexBasis FloatOptional) *Node {
	n.Layout.computedFlexBasis = computedFlexBasis
	return n
}

// SetLayoutPosition
func (n *Node) SetLayoutPosition(edge Edge, position float32) *Node {
	n.Layout.Position[edge] = position
	return n
}

// SetLayoutComputedFlexBasisGeneration
func (n *Node) SetLayoutComputedFlexBasisGeneration(generation uint32) *Node {
	n.Layout.computedFlexBasisGeneration = generation
	return n
}

// SetLayoutMeasuredDimension
func (n *Node) SetLayoutMeasuredDimension(dimension float32, index int) *Node {
	n.Layout.measuredDimensions[index] = dimension
	return n
}

// SetLayoutHadOverflow
func (n *Node) SetLayoutHadOverflow(hadOverflow bool) *Node {
	n.Layout.HadOverflow = hadOverflow
	return n
}

// SetLayoutDimension
func (n *Node) SetLayoutDimension(dimension float32, index int) *Node {
	n.Layout.Dimensions[index] = dimension
	return n
}

// relativePosition
func (n *Node) relativePosition(axis FlexDirection, axisSize float32) FloatOptional {
	if n.isLeadingPositionDefined(axis) {
		return n.getLeadingPosition(axis, axisSize)
	}
	trailingPosition := n.getTrailingPosition(axis, axisSize)
	if !trailingPosition.IsUndefined {
		trailingPosition = trailingPosition.Mult(NewFloatOptional(-1))
	}
	return trailingPosition
}

// SetPosition
func (n *Node) setPosition(direction Direction, mainSize float32, crossSize float32, ownerWidth float32) {
	var directionRespectingRoot Direction
	if n.Owner != nil {
		directionRespectingRoot = direction
	} else {
		directionRespectingRoot = DirectionLTR
	}
	mainAxis := ResolveFlexDirection(n.Style.FlexDirection, directionRespectingRoot)
	crossAxis := FlexDirectionCross(mainAxis, directionRespectingRoot)
	relativePositionMain := n.relativePosition(mainAxis, mainSize)
	relativePositionCross := n.relativePosition(crossAxis, crossSize)

	n.SetLayoutPosition(leading[mainAxis], n.getLeadingMargin(mainAxis, ownerWidth).Add(relativePositionMain).GetValue())
	n.SetLayoutPosition(trailing[mainAxis], n.getTrailingMargin(mainAxis, ownerWidth).Add(relativePositionMain).GetValue())
	n.SetLayoutPosition(leading[crossAxis], n.getLeadingMargin(crossAxis, ownerWidth).Add(relativePositionCross).GetValue())
	n.SetLayoutPosition(trailing[crossAxis], n.getTrailingMargin(crossAxis, ownerWidth).Add(relativePositionCross).GetValue())
}

// marginLeadingValue
func (n *Node) marginLeadingValue(axis FlexDirection) Value {
	if FlexDirectionIsRow(axis) && !n.Style.Margin[EdgeStart].IsUndefined() {
		return n.Style.Margin[EdgeStart]
	}
	return n.Style.Margin[leading[axis]]
}

// marginTrailingValue
func (n *Node) marginTrailingValue(axis FlexDirection) Value {
	if FlexDirectionIsRow(axis) && !n.Style.Margin[EdgeEnd].IsUndefined() {
		return n.Style.Margin[EdgeEnd]
	}
	return n.Style.Margin[trailing[axis]]
}

// resolveFlexBasisPtr
func (n *Node) resolveFlexBasisPtr() Value {
	flexBasis := n.Style.FlexBasis
	if flexBasis.Unit != UnitAuto && flexBasis.Unit != UnitUndefined {
		return flexBasis
	}
	if !n.Style.Flex.IsUndefined && n.Style.Flex.GetValue() > 0 {
		if n.Config.UseWebDefaults {
			return ValueAuto
		}
		return ValueZero
	}
	return ValueAuto
}

// resolveDirection
func (n *Node) resolveDirection(ownerDirection Direction) Direction {
	if n.Style.Direction == DirectionInherit {
		if ownerDirection > DirectionInherit {
			return ownerDirection
		}
		return DirectionLTR
	}
	return n.Style.Direction
}

// clearChildren
func (n *Node) clearChildren() {
	n.Children = nil
	n.MarkDirtyAndPropagate()
}

// markDirtyAndPropagateDownwards
func (n *Node) markDirtyAndPropagateDownwards() {
	n.SetDirty(true)
	for _, child := range n.Children {
		child.markDirtyAndPropagateDownwards()
	}
}

// resolveFlexGrow
func (n *Node) resolveFlexGrow() float32 {
	// Root nodes flexGrow should always be 0
	if n.Owner == nil {
		return 0
	}
	if !n.Style.FlexGrow.IsUndefined {
		return n.Style.FlexGrow.GetValue()
	}
	if !n.Style.Flex.IsUndefined && n.Style.Flex.GetValue() > 0 {
		return n.Style.Flex.GetValue()
	}
	// According to the spec, if the flexGrow property is undefined,
	// its value is 0 (and not 1).
	return 0
}

// resolveFlexShrink
func (n *Node) resolveFlexShrink() float32 {
	if n.Owner == nil {
		return 0
	}
	if !n.Style.FlexShrink.IsUndefined {
		return n.Style.FlexShrink.GetValue()
	}
	if !n.Config.UseWebDefaults && !n.Style.Flex.IsUndefined && n.Style.Flex.GetValue() < 0 {
		return -n.Style.Flex.GetValue()
	}

	if n.Config.UseWebDefaults {
		// The default flexShrink value is 1 in the web spec
		// See https://www.w3.org/TR/css-flexbox-1/#default-flex-shrink
		return 1
	}
	return 0
}

// isNodeFlexible
func (n *Node) isNodeFlexible() bool {
	return (n.Style.PositionType == PositionTypeRelative && (n.resolveFlexGrow() != 0 ||
		n.resolveFlexShrink() != 0))
}

// getLeadingBorder
func (n *Node) getLeadingBorder(axis FlexDirection) float32 {
	var val Value
	if FlexDirectionIsRow(axis) && !n.Style.Border[EdgeStart].IsUndefined() {
		val = computeEdgeValueForRow(n.Style.Border, EdgeStart, leading[axis], ValueZero)
	} else {
		val = computeEdgeValueForColumn(n.Style.Border, leading[axis], ValueZero)
	}
	return fmaxf(val.Value, 0)
}

// getTrailingBorder
func (n *Node) getTrailingBorder(axis FlexDirection) float32 {
	var val Value
	if FlexDirectionIsRow(axis) && !n.Style.Border[EdgeEnd].IsUndefined() {
		val = computeEdgeValueForRow(n.Style.Border, EdgeEnd, trailing[axis], ValueZero)
	} else {
		val = computeEdgeValueForColumn(n.Style.Border, trailing[axis], ValueZero)
	}
	return fmaxf(val.Value, 0)
}

// getChildCount
func (n *Node) getChildCount() int {
	return len(n.Children)
}

// getLeadingPadding
func (n *Node) getLeadingPadding(axis FlexDirection, widthSize float32) FloatOptional {
	var val Value
	if FlexDirectionIsRow(axis) {
		val = computeEdgeValueForRow(n.Style.Padding, EdgeStart, leading[axis], ValueZero)
	} else {
		val = computeEdgeValueForColumn(n.Style.Padding, leading[axis], ValueZero)
	}
	resolveValue := ResolveValue(val, widthSize)
	return FloatOptionalMax(resolveValue, NewFloatOptional(0))
}

// getTrailingPadding
func (n *Node) getTrailingPadding(axis FlexDirection, widthSize float32) FloatOptional {
	var val Value
	if FlexDirectionIsRow(axis) {
		val = computeEdgeValueForRow(n.Style.Padding, EdgeEnd, trailing[axis], ValueZero)
	} else {
		val = computeEdgeValueForColumn(n.Style.Padding, trailing[axis], ValueZero)
	}
	resolveValue := ResolveValue(val, widthSize)
	return FloatOptionalMax(resolveValue, NewFloatOptional(0))
}

// getLeadingPaddingAndBorder
func (n *Node) getLeadingPaddingAndBorder(axis FlexDirection, widthSize float32) FloatOptional {
	return n.getLeadingPadding(axis, widthSize).Add(NewFloatOptional(n.getLeadingBorder(axis)))
}

// getTrailingPaddingAndBorder
func (n *Node) getTrailingPaddingAndBorder(axis FlexDirection, widthSize float32) FloatOptional {
	return n.getTrailingPadding(axis, widthSize).Add(NewFloatOptional(n.getTrailingBorder(axis)))
}

// reset
func (n *Node) reset() {
	n = &Node{}
}

func computeEdgeValueForRow(edges [EdgeCount]Value, rowEdge Edge, edge Edge, defaultValue Value) Value {
	if !edges[rowEdge].IsUndefined() {
		return edges[rowEdge]
	} else if !edges[edge].IsUndefined() {
		return edges[edge]
	} else if !edges[EdgeHorizontal].IsUndefined() {
		return edges[EdgeHorizontal]
	} else if !edges[EdgeAll].IsUndefined() {
		return edges[EdgeAll]
	}
	return defaultValue
}

func computeEdgeValueForColumn(edges [EdgeCount]Value, edge Edge, defaultValue Value) Value {
	if !edges[edge].IsUndefined() {
		return edges[edge]
	} else if !edges[EdgeVertical].IsUndefined() {
		return edges[EdgeVertical]
	} else if !edges[EdgeAll].IsUndefined() {
		return edges[EdgeAll]
	}
	return defaultValue
}

func computeRowGap(gutters [GutterCount]Value, defaultValue Value) Value {
	if !gutters[GutterRow].IsUndefined() {
		return gutters[GutterRow]
	} else if !gutters[GutterAll].IsUndefined() {
		return gutters[GutterAll]
	}
	return defaultValue
}

func computeColumnGap(gutters [GutterCount]Value, defaultValue Value) Value {
	if !gutters[GutterColumn].IsUndefined() {
		return gutters[GutterColumn]
	} else if !gutters[GutterAll].IsUndefined() {
		return gutters[GutterAll]
	}
	return defaultValue
}
