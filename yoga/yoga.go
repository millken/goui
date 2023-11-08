package yoga

import "math"

type Size struct {
	Width  float32
	Height float32
}

type CachedMeasurement struct {
	availableWidth    float32
	availableHeight   float32
	widthMeasureMode  MeasureMode
	heightMeasureMode MeasureMode
	computedWidth     float32
	computedHeight    float32
}

var (
	defaultCachedMeasurement = CachedMeasurement{
		availableWidth:    Undefined,
		availableHeight:   Undefined,
		widthMeasureMode:  MeasureModeUndefined,
		heightMeasureMode: MeasureModeUndefined,
		computedWidth:     Undefined,
		computedHeight:    Undefined,
	}
	defaultDimensionValues = [2]float32{
		Undefined,
		Undefined,
	}
	currentGenerationCount uint32
)

func isUndefined[T float32 | float64](value T) bool {
	return math.IsNaN(float64(value))
}

func FloatIsUndefined[T float32 | float64](value T) bool {
	return isUndefined(value)
}

type MeasureFunc func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size
type BaseLineFunc func(node *Node, width, height float32) float32
type PrintFunc func(node *Node)

func computedEdgeValue(edges []Value, edge Edge, defaultValue *Value) *Value {
	if edges[edge].Unit != UnitUndefined {
		return &edges[edge]
	}

	isVertEdge := edge == EdgeTop || edge == EdgeBottom
	if isVertEdge && edges[EdgeVertical].Unit != UnitUndefined {
		return &edges[EdgeVertical]
	}

	isHorizEdge := (edge == EdgeLeft || edge == EdgeRight || edge == EdgeStart || edge == EdgeEnd)
	if isHorizEdge && edges[EdgeHorizontal].Unit != UnitUndefined {
		return &edges[EdgeHorizontal]
	}

	if edges[EdgeAll].Unit != UnitUndefined {
		return &edges[EdgeAll]
	}

	if edge == EdgeStart || edge == EdgeEnd {
		return &ValueUndefined
	}

	return defaultValue
}

func resolveValue(value *Value, parentSize float32) float32 {
	switch value.Unit {
	case UnitUndefined, UnitAuto:
		return Undefined
	case UnitPoint:
		return value.Value
	case UnitPercent:
		return value.Value * parentSize / 100
	}
	return Undefined
}

func resolveValueMargin(value *Value, parentSize float32) float32 {
	if value.Unit == UnitAuto {
		return 0
	}
	return resolveValue(value, parentSize)
}

func ValueEqual(a Value, b Value) bool {
	if a.Unit != b.Unit {
		return false
	}
	if a.Unit == UnitUndefined {
		return true
	}
	return Absf(a.Value-b.Value) < 0.0001
}

func FloatsEqual(a, b float32) bool {
	return (a - b) < 0.0001
}

var (
	leading  = [4]Edge{EdgeTop, EdgeBottom, EdgeLeft, EdgeRight}
	trailing = [4]Edge{EdgeBottom, EdgeTop, EdgeRight, EdgeLeft}
	pos      = [4]Edge{EdgeTop, EdgeBottom, EdgeLeft, EdgeRight}
	dim      = [4]Dimension{DimensionHeight, DimensionHeight, DimensionWidth, DimensionWidth}
)

func init() {
	leading[FlexDirectionColumn] = EdgeTop
	leading[FlexDirectionColumnReverse] = EdgeBottom
	leading[FlexDirectionRow] = EdgeLeft
	leading[FlexDirectionRowReverse] = EdgeRight

	trailing[FlexDirectionColumn] = EdgeBottom
	trailing[FlexDirectionColumnReverse] = EdgeTop
	trailing[FlexDirectionRow] = EdgeRight
	trailing[FlexDirectionRowReverse] = EdgeLeft

	pos[FlexDirectionColumn] = EdgeTop
	pos[FlexDirectionColumnReverse] = EdgeBottom
	pos[FlexDirectionRow] = EdgeLeft
	pos[FlexDirectionRowReverse] = EdgeRight

	dim[FlexDirectionColumn] = DimensionHeight
	dim[FlexDirectionColumnReverse] = DimensionHeight
	dim[FlexDirectionRow] = DimensionWidth
	dim[FlexDirectionRowReverse] = DimensionWidth
}

func flexDirectionIsRow(flexDirection FlexDirection) bool {
	return flexDirection == FlexDirectionRow || flexDirection == FlexDirectionRowReverse
}

func flexDirectionIsColumn(flexDirection FlexDirection) bool {
	return flexDirection == FlexDirectionColumn || flexDirection == FlexDirectionColumnReverse
}

func nodeAlignItem(node *Node, child *Node) Align {
	align := child.Style.AlignSelf
	if child.Style.AlignSelf == AlignAuto {
		align = node.Style.AlignItems
	}
	if align == AlignBaseline && flexDirectionIsColumn(node.Style.FlexDirection) {
		return AlignFlexStart
	}
	return align
}

func nodeResolveDirection(node *Node, parentDirection Direction) Direction {
	if node.Style.Direction == DirectionInherit {
		if parentDirection > DirectionInherit {
			return parentDirection
		}
		return DirectionLTR
	}
	return node.Style.Direction
}

func nodeIsStyleDimDefined(node *Node, axis FlexDirection, ownerSize float32) bool {
	rDim := node.getResolvedDimension(int(dim[axis]))
	isUndefined := FloatIsUndefined(rDim.Value)
	return !(rDim.Unit == UnitAuto || rDim.Unit == UnitUndefined ||
		(rDim.Unit == UnitPoint &&
			isUndefined && rDim.Value < 0.0) ||
		(rDim.Unit == UnitPercent &&
			isUndefined && (rDim.Value < 0.0 || FloatIsUndefined(ownerSize))))
}

// layoutNodeInternal
func layoutNodeInternal(node *Node, availableWidth, availableHeight float32, ownerDirection Direction, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, performLayout bool, reason LayoutPassReason, config *Config, layoutMarkerData *LayoutData, layoutContext interface{}, depth uint32, generationCount uint32) {
	layout := node.Layout
	depth++
	needToVisitNode := (node.IsNodeDirty() && layout.generationCount != generationCount) || layout.lastOwnerDirection != ownerDirection
	if needToVisitNode {
		// Invalidate the cached results.
		layout.nextCachedMeasurementsIndex = 0
		layout.cachedLayout.availableWidth = -1
		layout.cachedLayout.availableHeight = -1
		layout.cachedLayout.widthMeasureMode = MeasureModeUndefined
		layout.cachedLayout.heightMeasureMode = MeasureModeUndefined
		layout.cachedLayout.computedWidth = -1
		layout.cachedLayout.computedHeight = -1
	}
	var cachedResults *CachedMeasurement
	// Determine whether the results are already cached. We maintain a separate
	// cache for layouts and measurements. A layout operation modifies the
	// positions and dimensions for nodes in the subtree. The algorithm assumes
	// that each node gets laid out a maximum of one time per tree layout, but
	// multiple measurements may be required to resolve all of the flex
	// dimensions. We handle nodes with measure functions specially here because
	// they are the most expensive to measure, so it's worth avoiding redundant
	// measurements if at all possible.
	if node.Measure != nil {
		marginAxisRow := node.getMarginForAxis(FlexDirectionRow, ownerWidth).GetValue()
		marginAxisColumn := node.getMarginForAxis(FlexDirectionColumn, ownerWidth).GetValue()

		// First, try to use the layout cache.
		if NodeCanUseCachedMeasurement(widthMeasureMode,
			availableWidth,
			heightMeasureMode,
			availableHeight,
			layout.cachedLayout.widthMeasureMode,
			layout.cachedLayout.availableWidth,
			layout.cachedLayout.heightMeasureMode,
			layout.cachedLayout.availableHeight,
			layout.cachedLayout.computedWidth,
			layout.cachedLayout.computedHeight,
			marginAxisRow,
			marginAxisColumn,
			config) {
			cachedResults = &layout.cachedLayout
		} else {
			// Try to use the measurement cache.
			for i := 0; i < int(layout.nextCachedMeasurementsIndex); i++ {
				if NodeCanUseCachedMeasurement(widthMeasureMode,
					availableWidth,
					heightMeasureMode,
					availableHeight,
					layout.cachedMeasurements[i].widthMeasureMode,
					layout.cachedMeasurements[i].availableWidth,
					layout.cachedMeasurements[i].heightMeasureMode,
					layout.cachedMeasurements[i].availableHeight,
					layout.cachedMeasurements[i].computedWidth,
					layout.cachedMeasurements[i].computedHeight,
					marginAxisRow,
					marginAxisColumn,
					config) {
					cachedResults = &layout.cachedMeasurements[i]
					break
				}
			}
		}
	} else if performLayout {
		if FloatsEqual(layout.cachedLayout.availableWidth, availableWidth) && FloatsEqual(layout.cachedLayout.availableHeight, availableHeight) && layout.cachedLayout.widthMeasureMode == widthMeasureMode && layout.cachedLayout.heightMeasureMode == heightMeasureMode {
			cachedResults = &layout.cachedLayout
		}
	} else {
		for i := 0; i < int(layout.nextCachedMeasurementsIndex); i++ {
			if FloatsEqual(layout.cachedMeasurements[i].availableWidth, availableWidth) && FloatsEqual(layout.cachedMeasurements[i].availableHeight, availableHeight) && layout.cachedMeasurements[i].widthMeasureMode == widthMeasureMode && layout.cachedMeasurements[i].heightMeasureMode == heightMeasureMode {
				cachedResults = &layout.cachedMeasurements[i]
				break
			}
		}
	}
	if !needToVisitNode && cachedResults != nil {
		layout.measuredDimensions[DimensionWidth] = cachedResults.computedWidth
		layout.measuredDimensions[DimensionHeight] = cachedResults.computedHeight
		if performLayout {
			layoutMarkerData.cachedLayouts += 1
		} else {
			layoutMarkerData.cachedMeasures += 1
		}
	} else {
		NodelayoutImpl(node, availableWidth, availableHeight, ownerDirection, widthMeasureMode, heightMeasureMode, ownerWidth, ownerHeight, performLayout, config, layoutMarkerData, layoutContext, depth, generationCount)
	}
}

// NodelayoutImpl
func NodelayoutImpl(node *Node, availableWidth, availableHeight float32, ownerDirection Direction, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, performLayout bool, config *Config, layoutMarkerData *LayoutData, layoutContext interface{}, depth uint32, generationCount uint32, reason LayoutPassReason) {
	if performLayout {
		layoutMarkerData.layouts += 1
	} else {
		layoutMarkerData.measures += 1
	}
	direction := node.resolveDirection(ownerDirection)
	node.SetLayoutDirection(direction)
	flexRowDirection := ResolveFlexDirection(FlexDirectionRow, direction)
	flexColumnDirection := ResolveFlexDirection(FlexDirectionColumn, direction)
	var startEdge, endEdge Edge
	if direction == DirectionLTR {
		startEdge = EdgeLeft
		endEdge = EdgeRight
	} else {
		startEdge = EdgeRight
		endEdge = EdgeLeft
	}
	marginRowLeading := node.getLeadingMargin(flexRowDirection, ownerWidth).GetValue()
	node.SetLayoutMargin(startEdge, marginRowLeading)
	marginRowTrailing := node.getTrailingMargin(flexRowDirection, ownerWidth).GetValue()
	node.SetLayoutMargin(endEdge, marginRowTrailing)
	marginColumnLeading := node.getLeadingMargin(flexColumnDirection, ownerWidth).GetValue()
	node.SetLayoutMargin(EdgeTop, marginColumnLeading)
	marginColumnTrailing := node.getTrailingMargin(flexColumnDirection, ownerWidth).GetValue()
	node.SetLayoutMargin(EdgeBottom, marginColumnTrailing)

	marginAxisRow := marginRowLeading + marginRowTrailing
	marginAxisColumn := marginColumnLeading + marginColumnTrailing
	node.SetLayoutBorder(startEdge, node.getLeadingBorder(flexRowDirection))
	node.SetLayoutBorder(endEdge, node.getTrailingBorder(flexRowDirection))
	node.SetLayoutBorder(EdgeTop, node.getLeadingBorder(flexColumnDirection))
	node.SetLayoutBorder(EdgeBottom, node.getTrailingBorder(flexColumnDirection))

	node.SetLayoutPadding(startEdge, node.getLeadingPadding(flexRowDirection, ownerWidth).GetValue())
	node.SetLayoutPadding(endEdge, node.getTrailingPadding(flexRowDirection, ownerWidth).GetValue())
	node.SetLayoutPadding(EdgeTop, node.getLeadingPadding(flexColumnDirection, ownerWidth).GetValue())
	node.SetLayoutPadding(EdgeBottom, node.getTrailingPadding(flexColumnDirection, ownerWidth).GetValue())

	if node.Measure != nil {
		NodelayoutImplWithMeasureFunc(node, availableWidth, availableHeight, ownerDirection, widthMeasureMode, heightMeasureMode, ownerWidth, ownerHeight, performLayout, config, layoutMarkerData, layoutContext, depth, generationCount, reason)
		return
	}
}

// NodeWithMeasureFuncSetMeasuredDimensions
func NodeWithMeasureFuncSetMeasuredDimensions(node *Node, availableWidth, availableHeight float32, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, layoutMarkerData *LayoutData, layoutContext interface{}, reason LayoutPassReason) {
	if widthMeasureMode == MeasureModeUndefined {
		availableWidth = Undefined
	}
	if heightMeasureMode == MeasureModeUndefined {
		availableHeight = Undefined
	}
	padding := &node.Layout.Padding
	border := &node.Layout.Border
	paddingAndBorderAxisRow := padding[EdgeLeft] + padding[EdgeRight] + border[EdgeLeft] + border[EdgeRight]
	paddingAndBorderAxisColumn := padding[EdgeTop] + padding[EdgeBottom] + border[EdgeTop] + border[EdgeBottom]
	// We want to make sure we don't call measure with negative size
	var innerWidth, innerHeight float32
	if FloatIsUndefined(availableWidth) {
		innerWidth = availableWidth
	} else {
		innerWidth = FloatMax(0, availableWidth-paddingAndBorderAxisRow)
	}
	if FloatIsUndefined(availableHeight) {
		innerHeight = availableHeight
	} else {
		innerHeight = FloatMax(0, availableHeight-paddingAndBorderAxisColumn)
	}
	if widthMeasureMode == MeasureModeExactly &&
		heightMeasureMode == MeasureModeExactly {
		// Don't bother sizing the text if both dimensions are already defined.
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionRow, availableWidth, ownerWidth, ownerWidth), int(DimensionWidth))
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionColumn, availableHeight, ownerHeight, ownerWidth), int(DimensionHeight))
	} else {
		measuredSize := node.Measure(node, innerWidth, widthMeasureMode, innerHeight, heightMeasureMode)
		layoutMarkerData.measureCallbacks += 1
		layoutMarkerData.measureCallbackReasonsCount[reason] += 1
		var value float32
		if widthMeasureMode == MeasureModeUndefined || widthMeasureMode == MeasureModeAtMost {
			value = measuredSize.Width + paddingAndBorderAxisRow
		} else {
			value = availableWidth
		}
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionRow, value, ownerWidth, ownerWidth), int(DimensionWidth))
		if heightMeasureMode == MeasureModeUndefined || heightMeasureMode == MeasureModeAtMost {
			value = measuredSize.Height + paddingAndBorderAxisColumn
		} else {
			value = availableHeight
		}
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionColumn, value, ownerHeight, ownerWidth), int(DimensionHeight))
	}
}

// NodeBoundAxis
func NodeBoundAxis(node *Node, axis FlexDirection, value, axisSize, widthSize float32) float32 {
	return FloatMax(NodeBoundAxisWithinMinAndMax(node, axis, NewFloatOptional(value), axisSize).GetValue(), NodePaddingAndBorderForAxis(node, axis, widthSize))
}

// NodePaddingAndBorderForAxis
func NodePaddingAndBorderForAxis(node *Node, axis FlexDirection, widthSize float32) float32 {
	return node.getLeadingPaddingAndBorder(axis, widthSize).Add(node.getTrailingPaddingAndBorder(axis, widthSize)).GetValue()
}

// NodeBoundAxisWithinMinAndMax
func NodeBoundAxisWithinMinAndMax(node *Node, axis FlexDirection, value FloatOptional, axisSize float32) FloatOptional {
	var min, max FloatOptional
	if FlexDirectionIsColumn(axis) {
		min = ResolveValue(node.Style.MinDimensions[DimensionHeight], axisSize)
		max = ResolveValue(node.Style.MaxDimensions[DimensionHeight], axisSize)
	} else if FlexDirectionIsRow(axis) {
		min = ResolveValue(node.Style.MinDimensions[DimensionWidth], axisSize)
		max = ResolveValue(node.Style.MaxDimensions[DimensionWidth], axisSize)
	}
	if max.GetValue() >= 0 && value.GetValue() > max.GetValue() {
		return max
	}
	if min.GetValue() >= 0 && value.GetValue() < min.GetValue() {
		return min
	}
	return value
}

func NodeCanUseCachedMeasurement(widthMode MeasureMode, width float32, heightMode MeasureMode, height float32, lastWidthMode MeasureMode, lastWidth float32, lastHeightMode MeasureMode, lastHeight float32, lastComputedWidth, lastComputedHeight, marginRow, marginColumn float32, config *Config) bool {
	if (!FloatIsUndefined(lastComputedHeight) && lastComputedHeight < 0) ||
		(!FloatIsUndefined(lastComputedWidth) && lastComputedWidth < 0) {
		return false
	}
	useRoundedComparison := config.PointScaleFactor != 0
	var effectiveWidth, effectiveHeight, effectiveLastWidth, effectiveLastHeight float32
	if useRoundedComparison {
		effectiveWidth = RoundValueToPixelGrid(width, config.PointScaleFactor, false, false)
		effectiveHeight = RoundValueToPixelGrid(height, config.PointScaleFactor, false, false)
		effectiveLastWidth = RoundValueToPixelGrid(lastWidth, config.PointScaleFactor, false, false)
		effectiveLastHeight = RoundValueToPixelGrid(lastHeight, config.PointScaleFactor, false, false)
	} else {
		effectiveWidth = width
		effectiveHeight = height
		effectiveLastWidth = lastWidth
		effectiveLastHeight = lastHeight
	}
	hasSameWidthSpec := lastWidthMode == widthMode && FloatsEqual(effectiveLastWidth, effectiveWidth)
	hasSameHeightSpec := lastHeightMode == heightMode && FloatsEqual(effectiveLastHeight, effectiveHeight)
	widthIsCompatible := hasSameWidthSpec || MeasureModeSizeIsExactAndMatchesOldMeasuredSize(widthMode, width-marginRow, lastComputedWidth) || MeasureModeOldSizeIsUnspecifiedAndStillFits(widthMode, width-marginRow, lastWidthMode, lastComputedWidth) || MeasureModeNewMeasureSizeIsStricterAndStillValid(widthMode,
		width-marginRow,
		lastWidthMode,
		lastWidth,
		lastComputedWidth)
	heightIsCompatible := hasSameHeightSpec || MeasureModeSizeIsExactAndMatchesOldMeasuredSize(heightMode, height-marginColumn, lastComputedHeight) || MeasureModeOldSizeIsUnspecifiedAndStillFits(heightMode, height-marginColumn, lastHeightMode, lastComputedHeight) || MeasureModeNewMeasureSizeIsStricterAndStillValid(heightMode, height-marginColumn, lastHeightMode, lastHeight, lastComputedHeight)
	return widthIsCompatible && heightIsCompatible
}

// RoundValueToPixelGrid
func RoundValueToPixelGrid(value float32, pointScaleFactor float32, forceCeil bool, forceFloor bool) float32 {
	scaledValue := value * pointScaleFactor
	fractial := fmodf(scaledValue, 1)
	if fractial < 0 {
		fractial = 1 + fractial
	}
	if FloatsEqual(fractial, 0) {
		scaledValue = scaledValue - fractial
	} else if FloatsEqual(fractial, 1) {
		scaledValue = scaledValue - fractial + 1
	} else if forceCeil {
		scaledValue = scaledValue - fractial + 1
	} else if forceFloor {
		scaledValue = scaledValue - fractial
	} else {
		if !FloatIsUndefined(fractial) && fractial >= 0.5 {
			scaledValue = scaledValue - fractial + 1
		} else {
			scaledValue = scaledValue - fractial
		}
	}
	if FloatIsUndefined(scaledValue) || FloatIsUndefined(pointScaleFactor) {
		return Undefined
	}
	return scaledValue / pointScaleFactor
}

// MeasureModeSizeIsExactAndMatchesOldMeasuredSize
func MeasureModeSizeIsExactAndMatchesOldMeasuredSize(sizeMode MeasureMode, size, lastComputedSize float32) bool {
	return sizeMode == MeasureModeExactly && FloatsEqual(size, lastComputedSize)
}

// MeasureModeOldSizeIsUnspecifiedAndStillFits
func MeasureModeOldSizeIsUnspecifiedAndStillFits(sizeMode MeasureMode, size float32, lastSizeMode MeasureMode, lastComputedSize float32) bool {
	return sizeMode == MeasureModeAtMost && lastSizeMode == MeasureModeUndefined && (size >= lastComputedSize || FloatsEqual(size, lastComputedSize))
}

// MeasureModeNewMeasureSizeIsStricterAndStillValid
func MeasureModeNewMeasureSizeIsStricterAndStillValid(sizeMode MeasureMode, size float32, lastSizeMode MeasureMode, lastSize, lastComputedSize float32) bool {
	return lastSizeMode == MeasureModeAtMost && sizeMode == MeasureModeAtMost && !FloatIsUndefined(lastSize) && !FloatIsUndefined(size) && !FloatIsUndefined(lastComputedSize) && lastSize > size && (lastComputedSize <= size || FloatsEqual(size, lastComputedSize))
}
