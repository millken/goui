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

func NodeIsStyleDimDefined(node *Node, axis FlexDirection, ownerSize float32) bool {
	rDim := node.getResolvedDimension(int(dim[axis]))
	isUndefined := FloatIsUndefined(rDim.Value)
	return !(rDim.Unit == UnitAuto || rDim.Unit == UnitUndefined ||
		(rDim.Unit == UnitPoint &&
			isUndefined && rDim.Value < 0.0) ||
		(rDim.Unit == UnitPercent &&
			isUndefined && (rDim.Value < 0.0 || FloatIsUndefined(ownerSize))))
}

// layoutNodeInternal
func layoutNodeInternal(node *Node, availableWidth, availableHeight float32, ownerDirection Direction, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, performLayout bool, reason LayoutPassReason, config *Config, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32) {
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
		NodelayoutImpl(node, availableWidth, availableHeight, ownerDirection, widthMeasureMode, heightMeasureMode, ownerWidth, ownerHeight, performLayout, config, layoutMarkerData, layoutContext, depth, generationCount, reason)
	}
}

// NodelayoutImpl
func NodelayoutImpl(node *Node, availableWidth, availableHeight float32, ownerDirection Direction, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, performLayout bool, config *Config, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32, reason LayoutPassReason) {
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
		NodeWithMeasureFuncSetMeasuredDimensions(
			node,
			availableWidth-marginAxisRow,
			availableHeight-marginAxisColumn,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight,
			layoutMarkerData,
			layoutContext,
			reason)
		return
	}

	childCount := node.getChildCount()
	if childCount == 0 {
		NodeEmptyContainerSetMeasuredDimensions(
			node,
			availableWidth-marginAxisRow,
			availableHeight-marginAxisColumn,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight)
		return
	}

	// If we're not being asked to perform a full layout, we can handle a number of common
	// cases here without incurring the cost of the remaining function.
	if !performLayout {
		NodeFixedSizeSetMeasuredDimensions(
			node,
			availableWidth-marginAxisRow,
			availableHeight-marginAxisColumn,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight)
		return
	}
	// At this point we know we're going to perform work. Ensure that each child
	// has a mutable copy.
	node.cloneChildrenIfNeeded()
	// Reset layout flags, as they could have changed.
	node.SetLayoutHadOverflow(false)
	// STEP 1: CALCULATE VALUES FOR REMAINDER OF ALGORITHM
	mainAxis := ResolveFlexDirection(node.Style.FlexDirection, direction)
	crossAxis := FlexDirectionCross(mainAxis, direction)
	isMainAxisRow := flexDirectionIsRow(mainAxis)
	isNodeFlexWrap := node.Style.FlexWrap != WrapNoWrap
	var mainAxisownerSize, crossAxisownerSize float32
	if isMainAxisRow {
		mainAxisownerSize = ownerWidth
		crossAxisownerSize = ownerHeight
	} else {
		mainAxisownerSize = ownerHeight
		crossAxisownerSize = ownerWidth
	}
	paddingAndBorderAxisMain := NodePaddingAndBorderForAxis(node, mainAxis, ownerWidth)
	leadingPaddingAndBorderCross := node.getLeadingPaddingAndBorder(crossAxis, ownerWidth).GetValue()
	trailingPaddingAndBorderCross := node.getTrailingPaddingAndBorder(crossAxis, ownerWidth).GetValue()
	paddingAndBorderAxisCross := leadingPaddingAndBorderCross + trailingPaddingAndBorderCross

	var measureModeMainDim, measureModeCrossDim MeasureMode
	if isMainAxisRow {
		measureModeMainDim = widthMeasureMode
		measureModeCrossDim = heightMeasureMode
	} else {
		measureModeMainDim = heightMeasureMode
		measureModeCrossDim = widthMeasureMode
	}
	var paddingAndBorderAxisRow, paddingAndBorderAxisColumn float32
	if isMainAxisRow {
		paddingAndBorderAxisRow = paddingAndBorderAxisMain
		paddingAndBorderAxisColumn = paddingAndBorderAxisCross
	} else {
		paddingAndBorderAxisRow = paddingAndBorderAxisCross
		paddingAndBorderAxisColumn = paddingAndBorderAxisMain
	}

	// STEP 2: DETERMINE AVAILABLE SIZE IN MAIN AND CROSS DIRECTIONS
	availableInnerWidth := NodeCalculateAvailableInnerDim(
		node,
		DimensionWidth,
		availableWidth-marginAxisRow,
		paddingAndBorderAxisRow,
		ownerWidth)
	availableInnerHeight := NodeCalculateAvailableInnerDim(
		node,
		DimensionHeight,
		availableHeight-marginAxisColumn,
		paddingAndBorderAxisColumn,
		ownerHeight)
	var availableInnerMainDim, availableInnerCrossDim float32
	if isMainAxisRow {
		availableInnerMainDim = availableInnerWidth
		availableInnerCrossDim = availableInnerHeight
	} else {
		availableInnerMainDim = availableInnerHeight
		availableInnerCrossDim = availableInnerWidth
	}

	// STEP 3: DETERMINE FLEX BASIS FOR EACH ITEM
	totalMainDim := NodeComputeFlexBasisForChildren(
		node,
		availableInnerWidth,
		availableInnerHeight,
		widthMeasureMode,
		heightMeasureMode,
		direction,
		mainAxis,
		config,
		performLayout,
		layoutMarkerData,
		layoutContext,
		depth,
		generationCount)
	if childCount > 1 {
		totalMainDim += node.getGapForAxis(mainAxis, availableInnerCrossDim).GetValue() * float32(childCount-1)
	}

	mainAxisOverflows :=
		(measureModeMainDim != MeasureModeUndefined) &&
			totalMainDim > availableInnerMainDim
	if isNodeFlexWrap && mainAxisOverflows &&
		measureModeMainDim == MeasureModeAtMost {
		measureModeMainDim = MeasureModeExactly
	}

	// STEP 4: COLLECT FLEX ITEMS INTO FLEX LINES

	// Indexes of children that represent the first and last items in the line.
	var startOfLineIndex, endOfLineIndex, lineCount uint32

	// Accumulated cross dimensions of all lines so far.
	var float32 totalLineCrossDim = 0

	crossAxisGap := node.getGapForAxis(crossAxis, availableInnerCrossDim).GetValue()

	// Max main dimension of all the lines.
	var maxLineMainDim float32 = 0
	var collectedFlexItemsValues CollectFlexItemsRowValues
	for endOfLineIndex < uint32(childCount) {
		collectedFlexItemsValues = CalculateCollectFlexItemsRowValues(
			node,
			ownerDirection,
			mainAxisownerSize,
			availableInnerWidth,
			availableInnerMainDim,
			startOfLineIndex,
			lineCount)
		endOfLineIndex = collectedFlexItemsValues.endOfLineIndex
		// If we don't need to measure the cross axis, we can skip the entire flex
		// step.
		canSkipFlex := !performLayout && measureModeCrossDim == MeasureModeExactly
		// STEP 5: RESOLVING FLEXIBLE LENGTHS ON MAIN AXIS
		// Calculate the remaining available space that needs to be allocated. If
		// the main dimension size isn't known, it is computed based on the line
		// length, so there's no more space left to distribute.

		sizeBasedOnContent := false
		// If we don't measure with exact main dimension we want to ensure we don't
		// violate min and max
		if measureModeMainDim != MeasureModeExactly {
			minDimensions := node.Style.MinDimensions
			maxDimensions := node.Style.MaxDimensions
			minInnerWidth := ResolveValue(minDimensions[DimensionWidth], ownerWidth).GetValue() - paddingAndBorderAxisRow
			maxInnerWidth := ResolveValue(maxDimensions[DimensionWidth], ownerWidth).GetValue() - paddingAndBorderAxisRow
			minInnerHeight := ResolveValue(minDimensions[DimensionHeight], ownerHeight).GetValue() - paddingAndBorderAxisColumn
			maxInnerHeight := ResolveValue(maxDimensions[DimensionHeight], ownerHeight).GetValue() - paddingAndBorderAxisColumn
			minInnerMainDim := minInnerHeight
			maxInnerMainDim := maxInnerHeight
			if isMainAxisRow {
				minInnerMainDim = minInnerWidth
				maxInnerMainDim = maxInnerWidth
			}
			if !FloatIsUndefined(minInnerMainDim) && collectedFlexItemsValues.sizeConsumedOnCurrentLine < minInnerMainDim {
				availableInnerMainDim = minInnerMainDim
			} else if !FloatIsUndefined(maxInnerMainDim) && collectedFlexItemsValues.sizeConsumedOnCurrentLine > maxInnerMainDim {
				availableInnerMainDim = maxInnerMainDim
			} else {
				useLegacyStretchBehaviour := node.HasErrata(ErrataStretchFlexBasis)
				if !useLegacyStretchBehaviour &&
					((!FloatIsUndefined(
						collectedFlexItemsValues.totalFlexGrowFactors) &&
						collectedFlexItemsValues.totalFlexGrowFactors == 0) ||
						(!FloatIsUndefined(node.resolveFlexGrow()) &&
							node.resolveFlexGrow() == 0)) {
					// If we don't have any children to flex or we can't flex the node
					// itself, space we've used is all space we need. Root node also
					// should be shrunk to minimum
					availableInnerMainDim =
						collectedFlexItemsValues.sizeConsumedOnCurrentLine
				}

				sizeBasedOnContent = !useLegacyStretchBehaviour
			}
		}
		if !sizeBasedOnContent && !FloatIsUndefined(availableInnerMainDim) {
			collectedFlexItemsValues.remainingFreeSpace = availableInnerMainDim -
				collectedFlexItemsValues.sizeConsumedOnCurrentLine
		} else if collectedFlexItemsValues.sizeConsumedOnCurrentLine < 0 {
			// availableInnerMainDim is indefinite which means the node is being sized
			// based on its content. sizeConsumedOnCurrentLine is negative which means
			// the node will allocate 0 points for its content. Consequently,
			// remainingFreeSpace is 0 - sizeConsumedOnCurrentLine.
			collectedFlexItemsValues.remainingFreeSpace =
				-collectedFlexItemsValues.sizeConsumedOnCurrentLine
		}
		if !canSkipFlex {
			ResolveFlexibleLength(
				node,
				&collectedFlexItemsValues,
				mainAxis,
				crossAxis,
				mainAxisownerSize,
				availableInnerMainDim,
				availableInnerCrossDim,
				availableInnerWidth,
				availableInnerHeight,
				mainAxisOverflows,
				measureModeCrossDim,
				performLayout,
				config,
				layoutMarkerData,
				layoutContext,
				depth,
				generationCount)
		}
		node.SetLayoutHadOverflow(
			node.Layout.HadOverflow ||
				(collectedFlexItemsValues.remainingFreeSpace < 0))
   // STEP 6: MAIN-AXIS JUSTIFICATION & CROSS-AXIS SIZE DETERMINATION

    // At this point, all the children have their dimensions set in the main
    // axis. Their dimensions are also set in the cross axis with the exception
    // of items that are aligned "stretch". We need to compute these stretch
    // values and set the final positions.

    JustifyMainAxis(
        node,
        &collectedFlexItemsValues,
        startOfLineIndex,
        mainAxis,
        crossAxis,
        measureModeMainDim,
        measureModeCrossDim,
        mainAxisownerSize,
        ownerWidth,
        availableInnerMainDim,
        availableInnerCrossDim,
        availableInnerWidth,
        performLayout,
        layoutContext);
		containerCrossAxis := availableInnerCrossDim;
		if (measureModeCrossDim == MeasureModeUndefined ||
			measureModeCrossDim == MeasureModeAtMost) {
		  // Compute the cross axis from the max cross dimension of the children.
		  containerCrossAxis =
			  NodeBoundAxis(
				  node,
				  crossAxis,
				  collectedFlexItemsValues.crossDim + paddingAndBorderAxisCross,
				  crossAxisownerSize,
				  ownerWidth) -
			  paddingAndBorderAxisCross;
		}

		    // If there's no flex wrap, the cross dimension is defined by the container.
			if (!isNodeFlexWrap && measureModeCrossDim == MeasureModeExactly) {
				collectedFlexItemsValues.crossDim = availableInnerCrossDim;
			  }
    // Clamp to the min/max size specified on the container.
    collectedFlexItemsValues.crossDim =
        NodeBoundAxis(
            node,
            crossAxis,
            collectedFlexItemsValues.crossDim + paddingAndBorderAxisCross,
            crossAxisownerSize,
            ownerWidth) -
        paddingAndBorderAxisCross;

		    // STEP 7: CROSS-AXIS ALIGNMENT
    // We can skip child alignment if we're just measuring the container.
    if (performLayout) {
		for i := startOfLineIndex; i < endOfLineIndex; i++ {
		  child := node.GetChild(i);
		  if (child.Style.Display == DisplayNone) {
			continue;
		  }
		  if (child.Style.PositionType == PositionTypeAbsolute) {
			// If the child is absolutely positioned and has a
			// top/left/bottom/right set, override all the previously computed
			// positions to set it correctly.
			isChildLeadingPosDefined := child.isLeadingPositionDefined(crossAxis);
			if (isChildLeadingPosDefined) {
			  child.SetLayoutPosition(
				  child.getLeadingPosition(crossAxis, availableInnerCrossDim).GetValue() +
					  node.getLeadingBorder(crossAxis) +
					  child.getLeadingMargin(crossAxis, availableInnerWidth).GetValue(),
				  pos[crossAxis]);
			}
			// If leading position is not defined or calculations result in Nan,
			// default to border + margin
			if (!isChildLeadingPosDefined ||
				FloatIsUndefined(child.Layout.Position[pos[crossAxis]])) {
			  child.SetLayoutPosition(
				  node.getLeadingBorder(crossAxis) +
					  child.getLeadingMargin(crossAxis, availableInnerWidth).GetValue(),
				  pos[crossAxis]);
			}
		  } else {
			leadingCrossDim := leadingPaddingAndBorderCross;
  
			// For a relative children, we're either using alignItems (owner) or
			// alignSelf (child) in order to determine the position in the cross
			// axis
			alignItem := nodeAlignItem(node, child);
  
			// If the child uses align stretch, we need to lay it out one more
			// time, this time forcing the cross-axis size to be the computed
			// cross size for the current line.
			if (alignItem == AlignStretch &&
				child.marginLeadingValue(crossAxis).Unit != UnitAuto &&
				child.marginTrailingValue(crossAxis).Unit != UnitAuto) {
			  // If the child defines a definite size for its cross axis, there's
			  // no need to stretch.
			  if (!NodeIsStyleDimDefined(
					  child, crossAxis, availableInnerCrossDim)) {
			childMainSize := child.Layout.measuredDimensions[dim[mainAxis]];
				childStyle := child.Style;
				childCrossSize := collectedFlexItemsValues.crossDim
				if !childStyle.AspectRatio.IsUndefined {
					var mid float32
					if isMainAxisRow {
						mid = childMainSize / childStyle.AspectRatio.GetValue()
					} else {
						mid = childMainSize * childStyle.AspectRatio.GetValue()
					}	
									childCrossSize = child.getMarginForAxis(crossAxis, availableInnerWidth).GetValue() + mid

				}
  
				childMainSize +=
					child.getMarginForAxis(mainAxis, availableInnerWidth).GetValue();
  
				childMainMeasureMode := MeasureModeExactly;
				childCrossMeasureMode := MeasureModeExactly;
				ConstrainMaxSizeForMode(
					child,
					mainAxis,
					availableInnerMainDim,
					availableInnerWidth,
					&childMainMeasureMode,
					&childMainSize);
				ConstrainMaxSizeForMode(
					child,
					crossAxis,
					availableInnerCrossDim,
					availableInnerWidth,
					&childCrossMeasureMode,
					&childCrossSize);
				var childWidth, childHeight float32
				if (isMainAxisRow) {
				  childWidth = childMainSize
				  childHeight = childCrossSize
				} else {
									  childWidth = childCrossSize
									  				  childHeight = childMainSize
				}						
  
				alignContent := node.getStyle().alignContent();
				crossAxisDoesNotGrow = alignContent != AlignStretch && isNodeFlexWrap;
				var childWidthMeasureMode, childHeightMeasureMode MeasureMode
				if FloatIsUndefined(childWidth) || (!isMainAxisRow && crossAxisDoesNotGrow) {
				  childWidthMeasureMode = MeasureModeUndefined
				}else{
										  childWidthMeasureMode = MeasureModeExactly
				}	
				if FloatIsUndefined(childHeight) || (isMainAxisRow && crossAxisDoesNotGrow) {
									  childHeightMeasureMode = MeasureModeUndefined
				}else{
									  childHeightMeasureMode = MeasureModeExactly
				}
  
				layoutNodeInternal(
					child,
					childWidth,
					childHeight,
					direction,
					childWidthMeasureMode,
					childHeightMeasureMode,
					availableInnerWidth,
					availableInnerHeight,
					true,
					LayoutPassReasonStretch,
					config,
					layoutMarkerData,
					layoutContext,
					depth,
					generationCount);
			  }
			} else {
			  remainingCrossDim := containerCrossAxis -
				  NodeDimWithMargin(child, crossAxis, availableInnerWidth);
  
			  if (child.marginLeadingValue(crossAxis).Unit == UnitAuto &&
				  child.marginTrailingValue(crossAxis).Unit == UnitAuto) {
				leadingCrossDim += FloatMax(0.0, remainingCrossDim / 2);
			  } else if (
				  child.marginTrailingValue(crossAxis).Unit == UnitAuto) {
				// No-Op
			  } else if (
				  child.marginLeadingValue(crossAxis).Unit == UnitAuto) {
				leadingCrossDim += YGFloatMax(0.0, remainingCrossDim);
			  } else if (alignItem == AlignFlexStart) {
				// No-Op
			  } else if (alignItem == AlignCenter) {
				leadingCrossDim += remainingCrossDim / 2;
			  } else {
				leadingCrossDim += remainingCrossDim;
			  }
			}
			// And we apply the position
			child.SetLayoutPosition(
				child.Layout.Position[pos[crossAxis]] + totalLineCrossDim +
					leadingCrossDim,
				pos[crossAxis]);
		  }
		}
	  }
	  var appliedCrossGap float32 = 0.0
	  if lineCount != 0 {
		appliedCrossGap = crossAxisGap
	  }
	  totalLineCrossDim += collectedFlexItemsValues.crossDim + appliedCrossGap;
	  maxLineMainDim = FloatMax(maxLineMainDim, collectedFlexItemsValues.mainDim);
	}
  // STEP 8: MULTI-LINE CONTENT ALIGNMENT
  // currentLead stores the size of the cross dim
  if (performLayout && (isNodeFlexWrap || IsBaselineLayout(node))) {
    var crossDimLead float32 = 0;
    currentLead := leadingPaddingAndBorderCross;
    if (!FloatIsUndefined(availableInnerCrossDim)) {
      remainingAlignContentDim :=
          availableInnerCrossDim - totalLineCrossDim;
      switch (node.getStyle().alignContent()) {
        case AlignFlexEnd:
          currentLead += remainingAlignContentDim;
          break;
        case AlignCenter:
          currentLead += remainingAlignContentDim / 2;
          break;
        case AlignStretch:
          if (availableInnerCrossDim > totalLineCrossDim) {
            crossDimLead = remainingAlignContentDim / lineCount;
          }
          break;
        case AlignSpaceAround:
          if (availableInnerCrossDim > totalLineCrossDim) {
            currentLead += remainingAlignContentDim / (2 * lineCount);
            if (lineCount > 1) {
              crossDimLead = remainingAlignContentDim / lineCount;
            }
          } else {
            currentLead += remainingAlignContentDim / 2;
          }
          break;
        case AlignSpaceBetween:
          if (availableInnerCrossDim > totalLineCrossDim && lineCount > 1) {
            crossDimLead = remainingAlignContentDim / (lineCount - 1);
          }
          break;
        case AlignAuto:
        case AlignFlexStart:
        case AlignBaseline:
          break;
      }
    }
    endIndex := uint32(0);
    for i := 0; i < lineCount; i++ {
      startIndex := endIndex;
      var uint32 ii;

      // compute the line's height and find the endIndex
     var lineHeight float32 = 0;
      var maxAscentForCurrentLine, maxDescentForCurrentLine float32;
      for ii = startIndex; ii < childCount; ii++ {
        child := node.GetChild(ii);
        if (child.getStyle().display() == DisplayNone) {
          continue;
        }
        if (child.getStyle().positionType() != PositionTypeAbsolute) {
          if (child.getLineIndex() != i) {
            break;
          }
          if (NodeIsLayoutDimDefined(child, crossAxis)) {
            lineHeight = FloatMax(
                lineHeight,
                child.getLayout().measuredDimensions[dim[crossAxis]] +
                    child.getMarginForAxis(crossAxis, availableInnerWidth).GetValue());
          }
          if (NodeAlignItem(node, child) == AlignBaseline) {
            ascent := Baseline(child, layoutContext) +
                child.getLeadingMargin( FlexDirectionColumn, availableInnerWidth).GetValue()
            descent :=
                child.getLayout().measuredDimensions[DimensionHeight] +
                child.getMarginForAxis(
                        FlexDirectionColumn, availableInnerWidth).GetValue() -
                ascent;
            maxAscentForCurrentLine =
                FloatMax(maxAscentForCurrentLine, ascent);
            maxDescentForCurrentLine =
                FloatMax(maxDescentForCurrentLine, descent);
            lineHeight = FloatMax(
                lineHeight, maxAscentForCurrentLine + maxDescentForCurrentLine);
          }
        }
      }
      endIndex = ii;
      lineHeight += crossDimLead;
	  if i != 0 {
		currentLead += crossAxisGap;
	  }else{
		currentLead += 0
	  }

      if (performLayout) {
        for ii := startIndex; ii < endIndex; ii++ {
          child := node.getChild(ii);
          if child.getStyle().display() == DisplayNone {
            continue;
          }
          if (child.getStyle().positionType() != PositionTypeAbsolute) {
            switch NodeAlignItem(node, child) {
              case AlignFlexStart: {
                child.SetLayoutPosition(
                    currentLead +
                        child.getLeadingMargin(crossAxis, availableInnerWidth).GetValue(),
                    pos[crossAxis]);
                break;
              }
              case AlignFlexEnd: {
                child.SetLayoutPosition(
                    currentLead + lineHeight -
                        child.getTrailingMargin(crossAxis, availableInnerWidth).GetValue() -
                        child.getLayout().measuredDimensions[dim[crossAxis]],
                    pos[crossAxis]);
                break;
              }
              case AlignCenter: {
                childHeight :=
                    child.getLayout().measuredDimensions[dim[crossAxis]];

                child.SetLayoutPosition(
                    currentLead + (lineHeight - childHeight) / 2,
                    pos[crossAxis]);
                break;
              }
              case AlignStretch: 
                child.SetLayoutPosition(
                    currentLead +
                        child.getLeadingMargin(crossAxis, availableInnerWidth).GetValue(),
                    pos[crossAxis]);

                // Remeasure child with the line height as it as been only
                // measured with the owners height yet.
                if !NodeIsStyleDimDefined(
                        child, crossAxis, availableInnerCrossDim) {
							var childWidth,childHeight float32
                 if isMainAxisRow {
						childWidth = child.getLayout().measuredDimensions[DimensionWidth] +
							child.getMarginForAxis(mainAxis, availableInnerWidth).GetValue()
						childHeight = lineHeight
				 }else{
						childWidth = lineHeight
						childHeight = child.getLayout().measuredDimensions[DimensionHeight] +
							child.getMarginForAxis(crossAxis, availableInnerWidth).GetValue()
				 }

                  if !(FloatsEqual(
                            childWidth,
                            child.getLayout().measuredDimensions[DimensionWidth]) &&
                        FloatsEqual(
                            childHeight,
                            child.getLayout().measuredDimensions[DimensionHeight])) {
                    LayoutNodeInternal(
                        child,
                        childWidth,
                        childHeight,
                        direction,
                        YGMeasureModeExactly,
                        YGMeasureModeExactly,
                        availableInnerWidth,
                        availableInnerHeight,
                        true,
                        LayoutPassReasonMultilineStretch,
                        config,
                        layoutMarkerData,
                        layoutContext,
                        depth,
                        generationCount);
                  }
                }
                break;
              
              case AlignBaseline: 
                child.SetLayoutPosition(
                    currentLead + maxAscentForCurrentLine -
                        Baseline(child, layoutContext) +
                        child.getLeadingPosition(
                                YGFlexDirectionColumn, availableInnerCrossDim).GetValue(),
                    EdgeTop);

                break;
              case AlignAuto:
              case AlignSpaceBetween:
              case AlignSpaceAround:
                break;
            }
          }
        }
      }
      currentLead += lineHeight;
    }
  }

   // STEP 9: COMPUTING FINAL DIMENSIONS

   node.SetLayoutMeasuredDimension(
	NodeBoundAxis(
		node,
		FlexDirectionRow,
		availableWidth - marginAxisRow,
		ownerWidth,
		ownerWidth),
	DimensionWidth);

node.SetLayoutMeasuredDimension(
	NodeBoundAxis(
		node,
		FlexDirectionColumn,
		availableHeight - marginAxisColumn,
		ownerHeight,
		ownerWidth),
	DimensionHeight);

// If the user didn't specify a width or height for the node, set the
// dimensions based on the children.
if (measureModeMainDim == MeasureModeUndefined || (node.getStyle().overflow() != OverflowScroll && measureModeMainDim == MeasureModeAtMost)) {
  // Clamp the size to the min/max size, if specified, and make sure it
  // doesn't go below the padding and border amount.
  node.SetLayoutMeasuredDimension(
	  NodeBoundAxis(
		  node, mainAxis, maxLineMainDim, mainAxisownerSize, ownerWidth),
	  dim[mainAxis]);

} else if (measureModeMainDim == MeasureModeAtMost &&
	node.getStyle().overflow() == OverflowScroll) {
  node.SetLayoutMeasuredDimension(
	  FloatMax(
		  FloatMin(availableInnerMainDim + paddingAndBorderAxisMain,
			  NodeBoundAxisWithinMinAndMax(node,
				  mainAxis,
				  NewFloatOptional(maxLineMainDim),
				  mainAxisownerSize).GetValue()),
		  paddingAndBorderAxisMain),
	  dim[mainAxis]);
}

if (measureModeCrossDim == MeasureModeUndefined || (node.getStyle().overflow() != OverflowScroll &&  measureModeCrossDim == MeasureModeAtMost)) {
  // Clamp the size to the min/max size, if specified, and make sure it
  // doesn't go below the padding and border amount.
  node.SetLayoutMeasuredDimension(
	  NodeBoundAxis(
		  node,
		  crossAxis,
		  totalLineCrossDim + paddingAndBorderAxisCross,
		  crossAxisownerSize,
		  ownerWidth),
	  dim[crossAxis]);

} else if (measureModeCrossDim == MeasureModeAtMost && node.getStyle().overflow() == OverflowScroll) {
  node.SetLayoutMeasuredDimension(
	  FloatMax(
		  FloatMin(
			  availableInnerCrossDim + paddingAndBorderAxisCross,
			  NodeBoundAxisWithinMinAndMax(
				  node,
				  crossAxis,
				  NewFloatOptional(totalLineCrossDim + paddingAndBorderAxisCross),
				  crossAxisownerSize).GetValue()),
		  paddingAndBorderAxisCross),
	  dim[crossAxis]);
}

// As we only wrapped in normal direction yet, we need to reverse the
// positions on wrap-reverse.
if (performLayout && node.getStyle().flexWrap() == WrapWrapReverse) {
  for i := 0; i < childCount; i++ {
	child := NodeGetChild(node, i);
	if (child.getStyle().positionType() != PositionTypeAbsolute) {
	  child.SetLayoutPosition(
		  node.getLayout().measuredDimensions[dim[crossAxis]] -
			  child.getLayout().position[pos[crossAxis]] -
			  child.getLayout().measuredDimensions[dim[crossAxis]],
		  pos[crossAxis]);
	}
  }
}
if (performLayout) {
    // STEP 10: SIZING AND POSITIONING ABSOLUTE CHILDREN
    for _, child := range node.getChildren() {
      if (child.getStyle().display() == DisplayNone ||
          child.getStyle().positionType() != PositionTypeAbsolute) {
        continue;
      }
      absolutePercentageAgainstPaddingEdge :=
          node.getConfig().experimentalFeatureEnabled(
              ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge);

      NodeAbsoluteLayoutChild(
          node,
          child,
          absolutePercentageAgainstPaddingEdge
              ? node->getLayout().measuredDimensions[YGDimensionWidth]
              : availableInnerWidth,
          isMainAxisRow ? measureModeMainDim : measureModeCrossDim,
          absolutePercentageAgainstPaddingEdge
              ? node->getLayout().measuredDimensions[YGDimensionHeight]
              : availableInnerHeight,
          direction,
          config,
          layoutMarkerData,
          layoutContext,
          depth,
          generationCount);
    }

    // STEP 11: SETTING TRAILING POSITIONS FOR CHILDREN
    const bool needsMainTrailingPos = mainAxis == YGFlexDirectionRowReverse ||
        mainAxis == YGFlexDirectionColumnReverse;
    const bool needsCrossTrailingPos = crossAxis == YGFlexDirectionRowReverse ||
        crossAxis == YGFlexDirectionColumnReverse;

    // Set trailing position if necessary.
    if (needsMainTrailingPos || needsCrossTrailingPos) {
      for (uint32_t i = 0; i < childCount; i++) {
        const YGNodeRef child = node->getChild(i);
        if (child->getStyle().display() == YGDisplayNone) {
          continue;
        }
        if (needsMainTrailingPos) {
          YGNodeSetChildTrailingPosition(node, child, mainAxis);
        }

        if (needsCrossTrailingPos) {
          YGNodeSetChildTrailingPosition(node, child, crossAxis);
        }
      }
    }
  }
		lineCount++
		startOfLineIndex = endOfLineIndex
	}
}

//NodeIsLayoutDimDefined
func NodeIsLayoutDimDefined(node *Node, axis FlexDirection) bool {
	value := node.getLayout().measuredDimensions[dim[axis]]
	return !FloatIsUndefined(value) && value >= 0.0
}
//NodeDimWithMargin
func NodeDimWithMargin(node *Node, axis FlexDirection, widthSize float32) float32 {
	return node.Layout.measuredDimensions[dim[axis]] + node.getLeadingMargin(axis, widthSize).Add(
	node.getTrailingMargin(axis, widthSize)).GetValue()
}

//ConstrainMaxSizeForMode
func ConstrainMaxSizeForMode(node *Node, axis FlexDirection, ownerAxisSize, ownerWidth float32, mode *MeasureMode, size *float32) {
	maxSize :=
	ResolveValue(node.Style.MaxDimensions[dim[axis]], ownerAxisSize).Add(FloatOptional(node.getMarginForAxis(axis, ownerWidth)))
switch (*mode) {
  case MeasureModeExactly:
  case MeasureModeAtMost:
	*size = (maxSize.isUndefined() || *size < maxSize.GetValue())
		? *size
		: maxSize.unwrap();
	break;
  case MeasureModeUndefined:
	if (!maxSize.isUndefined()) {
	  *mode = MeasureModeAtMost;
	  *size = maxSize.unwrap();
	}
	break;
}
}
//NodeBoundAxis
func NodeBoundAxis(node *Node, axis FlexDirection, value, axisSize, widthSize float32) float32 {
	return FloatMax(
		NodeBoundAxisWithinMinAndMax(
			node, axis, NewFloatOptional(value), axisSize).GetValue(),
		NodePaddingAndBorderForAxis(node, axis, widthSize));
}

//JustifyMainAxis
func JustifyMainAxis(node *Node, collectedFlexItemsValues *CollectFlexItemsRowValues, startOfLineIndex uint32, mainAxis, crossAxis FlexDirection, measureModeMainDim, measureModeCrossDim MeasureMode, mainAxisownerSize, ownerWidth, availableInnerMainDim, availableInnerCrossDim, availableInnerWidth float32, performLayout bool, layoutContext any) {
	style := node.Style;
	leadingPaddingAndBorderMain :=
		node.getLeadingPaddingAndBorder(mainAxis, ownerWidth).GetValue();
	trailingPaddingAndBorderMain :=
		node.getTrailingPaddingAndBorder(mainAxis, ownerWidth).GetValue();
	gap := node.getGapForAxis(mainAxis, ownerWidth).GetValue();
	// If we are using "at most" rules in the main axis, make sure that
	// remainingFreeSpace is 0 when min main dimension is not given
	if (measureModeMainDim == MeasureModeAtMost &&
		collectedFlexItemsValues.remainingFreeSpace > 0) {
	  if (!style.MinDimensions[dim[mainAxis]].IsUndefined() &&
		  !ResolveValue(style.MinDimensions[dim[mainAxis]], mainAxisownerSize).IsUndefined()) {
		// This condition makes sure that if the size of main dimension(after
		// considering child nodes main dim, leading and trailing padding etc)
		// falls below min dimension, then the remainingFreeSpace is reassigned
		// considering the min dimension
  
		// `minAvailableMainDim` denotes minimum available space in which child
		// can be laid out, it will exclude space consumed by padding and border.
		minAvailableMainDim :=
			ResolveValue(
				style.MinDimensions[dim[mainAxis]], mainAxisownerSize).GetValue() -
			leadingPaddingAndBorderMain - trailingPaddingAndBorderMain;
		occupiedSpaceByChildNodes :=
			availableInnerMainDim - collectedFlexItemsValues.remainingFreeSpace;
		collectedFlexItemsValues.remainingFreeSpace =
			FloatMax(0, minAvailableMainDim - occupiedSpaceByChildNodes);
	  } else {
		collectedFlexItemsValues.remainingFreeSpace = 0;
	  }
	}
  
	var numberOfAutoMarginsOnCurrentLine int = 0;
	for i := startOfLineIndex;i < collectedFlexItemsValues.endOfLineIndex;i++ {
	  child := node.Children[i]
	  if (child.Style.PositionType != PositionTypeAbsolute) {
		if (child.marginLeadingValue(mainAxis).Unit == UnitAuto) {
		  numberOfAutoMarginsOnCurrentLine++;
		}
		if (child.marginTrailingValue(mainAxis).Unit == UnitAuto) {
		  numberOfAutoMarginsOnCurrentLine++;
		}
	  }
	}
  
	// In order to position the elements in the main axis, we have two controls.
	// The space between the beginning and the first element and the space between
	// each two elements.
	var leadingMainDim float32 = 0;
	betweenMainDim := gap;
	justifyContent := node.Style.JustifyContent;
  
	if (numberOfAutoMarginsOnCurrentLine == 0) {
	  switch (justifyContent) {
		case JustifyCenter:
		  leadingMainDim = collectedFlexItemsValues.remainingFreeSpace / 2;
		  break;
		case JustifyFlexEnd:
		  leadingMainDim = collectedFlexItemsValues.remainingFreeSpace;
		  break;
		case JustifySpaceBetween:
		  if (collectedFlexItemsValues.itemsOnLine > 1) {
			betweenMainDim +=
				FloatMax(collectedFlexItemsValues.remainingFreeSpace, 0) /
				(collectedFlexItemsValues.itemsOnLine - 1);
		  }
		  break;
		case JustifySpaceEvenly:
		  // Space is distributed evenly across all elements
		  leadingMainDim = collectedFlexItemsValues.remainingFreeSpace /
			  (collectedFlexItemsValues.itemsOnLine + 1);
		  betweenMainDim += leadingMainDim;
		  break;
		case JustifySpaceAround:
		  // Space on the edges is half of the space between elements
		  leadingMainDim = 0.5 * collectedFlexItemsValues.remainingFreeSpace /
			  collectedFlexItemsValues.itemsOnLine;
		  betweenMainDim += leadingMainDim * 2;
		  break;
		case JustifyFlexStart:
		  break;
	  }
	}
  
	collectedFlexItemsValues.mainDim =
		leadingPaddingAndBorderMain + leadingMainDim;
	collectedFlexItemsValues.crossDim = 0;
  
	var maxAscentForCurrentLine,maxDescentForCurrentLine float32;
	isNodeBaselineLayout := IsBaselineLayout(node);
	for (uint32_t i = startOfLineIndex;
		 i < collectedFlexItemsValues.endOfLineIndex;
		 i++) {
	  const YGNodeRef child = node.getChild(i);
	  const YGStyle& childStyle = child->getStyle();
	  const YGLayout childLayout = child->getLayout();
	  const bool isLastChild = i == collectedFlexItemsValues.endOfLineIndex - 1;
	  // remove the gap if it is the last element of the line
	  if (isLastChild) {
		betweenMainDim -= gap;
	  }
	  if (childStyle.display() == YGDisplayNone) {
		continue;
	  }
	  if (childStyle.positionType() == YGPositionTypeAbsolute &&
		  child->isLeadingPositionDefined(mainAxis)) {
		if (performLayout) {
		  // In case the child is position absolute and has left/top being
		  // defined, we override the position to whatever the user said (and
		  // margin/border).
		  child->setLayoutPosition(
			  child->getLeadingPosition(mainAxis, availableInnerMainDim)
					  .GetValue() +
				  node.getLeadingBorder(mainAxis) +
				  child->getLeadingMargin(mainAxis, availableInnerWidth).GetValue(),
			  pos[mainAxis]);
		}
	  } else {
		// Now that we placed the element, we need to update the variables.
		// We need to do that only for relative elements. Absolute elements do not
		// take part in that phase.
		if (childStyle.positionType() != YGPositionTypeAbsolute) {
		  if (child->marginLeadingValue(mainAxis).unit == YGUnitAuto) {
			collectedFlexItemsValues.mainDim +=
				collectedFlexItemsValues.remainingFreeSpace /
				numberOfAutoMarginsOnCurrentLine;
		  }
  
		  if (performLayout) {
			child->setLayoutPosition(
				childLayout.position[pos[mainAxis]] +
					collectedFlexItemsValues.mainDim,
				pos[mainAxis]);
		  }
  
		  if (child->marginTrailingValue(mainAxis).unit == YGUnitAuto) {
			collectedFlexItemsValues.mainDim +=
				collectedFlexItemsValues.remainingFreeSpace /
				numberOfAutoMarginsOnCurrentLine;
		  }
		  bool canSkipFlex =
			  !performLayout && measureModeCrossDim == YGMeasureModeExactly;
		  if (canSkipFlex) {
			// If we skipped the flex step, then we can't rely on the measuredDims
			// because they weren't computed. This means we can't call
			// YGNodeDimWithMargin.
			collectedFlexItemsValues.mainDim += betweenMainDim +
				child->getMarginForAxis(mainAxis, availableInnerWidth).GetValue() +
				childLayout.computedFlexBasis.GetValue();
			collectedFlexItemsValues.crossDim = availableInnerCrossDim;
		  } else {
			// The main dimension is the sum of all the elements dimension plus
			// the spacing.
			collectedFlexItemsValues.mainDim += betweenMainDim +
				YGNodeDimWithMargin(child, mainAxis, availableInnerWidth);
  
			if (isNodeBaselineLayout) {
			  // If the child is baseline aligned then the cross dimension is
			  // calculated by adding maxAscent and maxDescent from the baseline.
			  const float ascent = YGBaseline(child, layoutContext) +
				  child
					  ->getLeadingMargin(
						  YGFlexDirectionColumn, availableInnerWidth)
					  .GetValue();
			  const float descent =
				  child->getLayout().measuredDimensions[YGDimensionHeight] +
				  child
					  ->getMarginForAxis(
						  YGFlexDirectionColumn, availableInnerWidth)
					  .GetValue() -
				  ascent;
  
			  maxAscentForCurrentLine =
				  YGFloatMax(maxAscentForCurrentLine, ascent);
			  maxDescentForCurrentLine =
				  YGFloatMax(maxDescentForCurrentLine, descent);
			} else {
			  // The cross dimension is the max of the elements dimension since
			  // there can only be one element in that cross dimension in the case
			  // when the items are not baseline aligned
			  collectedFlexItemsValues.crossDim = YGFloatMax(
				  collectedFlexItemsValues.crossDim,
				  YGNodeDimWithMargin(child, crossAxis, availableInnerWidth));
			}
		  }
		} else if (performLayout) {
		  child->setLayoutPosition(
			  childLayout.position[pos[mainAxis]] +
				  node.getLeadingBorder(mainAxis) + leadingMainDim,
			  pos[mainAxis]);
		}
	  }
	}
	collectedFlexItemsValues.mainDim += trailingPaddingAndBorderMain;
  
	if (isNodeBaselineLayout) {
	  collectedFlexItemsValues.crossDim =
		  maxAscentForCurrentLine + maxDescentForCurrentLine;
	}
}

//IsBaselineLayout
func IsBaselineLayout(node *Node) bool {
	if (FlexDirectionIsColumn(node.Style.FlexDirection)) {
		return false;
	  }
	  if (node.Style.AlignItems == AlignBaseline) {
		return true;
	  }
	  childCount := node.getChildCount()
	  for i := 0; i < childCount; i++ {
		child := NodeGetChild(node, i)
		if (child.Style.PositionType != PositionTypeAbsolute &&
			child.Style.AlignSelf == AlignBaseline) {
		  return true;
		}
	  }
	
	  return false;
	}

func NodeGetChild(node *Node,  index int) *Node {
		if index < len(node.Children) {
		  return node.GetChild(index)
		}
		return nil;
}
// ResolveFlexibleLength
func ResolveFlexibleLength(node *Node, collectedFlexItemsValues *CollectFlexItemsRowValues, mainAxis, crossAxis FlexDirection, mainAxisownerSize, availableInnerMainDim, availableInnerCrossDim, availableInnerWidth, availableInnerHeight float32, mainAxisOverflows bool, measureModeCrossDim MeasureMode, performLayout bool, config *Config, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32) {
	originalFreeSpace := collectedFlexItemsValues.remainingFreeSpace
	// First pass: detect the flex items whose min/max constraints trigger
	DistributeFreeSpaceFirstPass(
		collectedFlexItemsValues,
		mainAxis,
		mainAxisownerSize,
		availableInnerMainDim,
		availableInnerWidth)

	// Second pass: resolve the sizes of the flexible items
	distributedFreeSpace := DistributeFreeSpaceSecondPass(
		collectedFlexItemsValues,
		node,
		mainAxis,
		crossAxis,
		mainAxisownerSize,
		availableInnerMainDim,
		availableInnerCrossDim,
		availableInnerWidth,
		availableInnerHeight,
		mainAxisOverflows,
		measureModeCrossDim,
		performLayout,
		config,
		layoutMarkerData,
		layoutContext,
		depth,
		generationCount)

	collectedFlexItemsValues.remainingFreeSpace =
		originalFreeSpace - distributedFreeSpace
}

// DistributeFreeSpaceFirstPass
func DistributeFreeSpaceFirstPass(collectedFlexItemsValues *CollectFlexItemsRowValues, mainAxis FlexDirection, mainAxisownerSize, availableInnerMainDim, availableInnerWidth float32) {
	var flexShrinkScaledFactor, flexGrowFactor, baseMainSize, boundMainSize, deltaFreeSpace float32
	for _, currentRelativeChild := range collectedFlexItemsValues.relativeChildren {
		childFlexBasis := NodeBoundAxisWithinMinAndMax(
			currentRelativeChild,
			mainAxis,
			currentRelativeChild.Layout.computedFlexBasis,
			mainAxisownerSize).GetValue()
		if collectedFlexItemsValues.remainingFreeSpace < 0 {
			flexShrinkScaledFactor =
				-currentRelativeChild.resolveFlexShrink() * childFlexBasis

			// Is this child able to shrink?
			if !FloatIsUndefined(flexShrinkScaledFactor) &&
				flexShrinkScaledFactor != 0 {
				baseMainSize = childFlexBasis +
					collectedFlexItemsValues.remainingFreeSpace/
						collectedFlexItemsValues.totalFlexShrinkScaledFactors*
						flexShrinkScaledFactor
				boundMainSize = NodeBoundAxis(
					currentRelativeChild,
					mainAxis,
					baseMainSize,
					availableInnerMainDim,
					availableInnerWidth)
				if !FloatIsUndefined(baseMainSize) &&
					!FloatIsUndefined(boundMainSize) &&
					baseMainSize != boundMainSize {
					// By excluding this item's size and flex factor from remaining, this
					// item's min/max constraints should also trigger in the second pass
					// resulting in the item's size calculation being identical in the
					// first and second passes.
					deltaFreeSpace += boundMainSize - childFlexBasis
					collectedFlexItemsValues.totalFlexShrinkScaledFactors -=
						(-currentRelativeChild.resolveFlexShrink() *
							currentRelativeChild.Layout.computedFlexBasis.GetValue())
				}
			}
		} else if !FloatIsUndefined(collectedFlexItemsValues.remainingFreeSpace) &&
			collectedFlexItemsValues.remainingFreeSpace > 0 {
			flexGrowFactor = currentRelativeChild.resolveFlexGrow()

			// Is this child able to grow?
			if !FloatIsUndefined(flexGrowFactor) && flexGrowFactor != 0 {
				baseMainSize = childFlexBasis +
					collectedFlexItemsValues.remainingFreeSpace/
						collectedFlexItemsValues.totalFlexGrowFactors*flexGrowFactor
				boundMainSize = NodeBoundAxis(
					currentRelativeChild,
					mainAxis,
					baseMainSize,
					availableInnerMainDim,
					availableInnerWidth)

				if !FloatIsUndefined(baseMainSize) &&
					!FloatIsUndefined(boundMainSize) &&
					baseMainSize != boundMainSize {
					// By excluding this item's size and flex factor from remaining, this
					// item's min/max constraints should also trigger in the second pass
					// resulting in the item's size calculation being identical in the
					// first and second passes.
					deltaFreeSpace += boundMainSize - childFlexBasis
					collectedFlexItemsValues.totalFlexGrowFactors -= flexGrowFactor
				}
			}
		}
	}
	collectedFlexItemsValues.remainingFreeSpace -= deltaFreeSpace
}

// DistributeFreeSpaceSecondPass
func DistributeFreeSpaceSecondPass(collectedFlexItemsValues *CollectFlexItemsRowValues, node *Node, mainAxis, crossAxis FlexDirection, mainAxisownerSize, availableInnerMainDim, availableInnerCrossDim, availableInnerWidth, availableInnerHeight float32, mainAxisOverflows bool, measureModeCrossDim MeasureMode, performLayout bool, config *Config, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32) float32 {
	var childFlexBasis, flexShrinkScaledFactor, flexGrowFactor, deltaFreeSpace float32
	isMainAxisRow := flexDirectionIsRow(mainAxis)
	isNodeFlexWrap := node.Style.FlexWrap != WrapNoWrap
	for _, currentRelativeChild := range collectedFlexItemsValues.relativeChildren {
		childFlexBasis = NodeBoundAxisWithinMinAndMax(
			currentRelativeChild,
			mainAxis,
			currentRelativeChild.Layout.computedFlexBasis,
			mainAxisownerSize).GetValue()
		updatedMainSize := childFlexBasis
		if !FloatIsUndefined(collectedFlexItemsValues.remainingFreeSpace) &&
			collectedFlexItemsValues.remainingFreeSpace < 0 {
			flexShrinkScaledFactor =
				-currentRelativeChild.resolveFlexShrink() * childFlexBasis
			// Is this child able to shrink?
			if flexShrinkScaledFactor != 0 {
				var childSize float32

				if !FloatIsUndefined(
					collectedFlexItemsValues.totalFlexShrinkScaledFactors) &&
					collectedFlexItemsValues.totalFlexShrinkScaledFactors == 0 {
					childSize = childFlexBasis + flexShrinkScaledFactor
				} else {
					childSize = childFlexBasis +
						(collectedFlexItemsValues.remainingFreeSpace/
							collectedFlexItemsValues.totalFlexShrinkScaledFactors)*
							flexShrinkScaledFactor
				}

				updatedMainSize = NodeBoundAxis(
					currentRelativeChild,
					mainAxis,
					childSize,
					availableInnerMainDim,
					availableInnerWidth)
			}
		} else if !FloatIsUndefined(collectedFlexItemsValues.remainingFreeSpace) &&
			collectedFlexItemsValues.remainingFreeSpace > 0 {
			flexGrowFactor = currentRelativeChild.resolveFlexGrow()

			// Is this child able to grow?
			if !FloatIsUndefined(flexGrowFactor) && flexGrowFactor != 0 {
				updatedMainSize = NodeBoundAxis(
					currentRelativeChild,
					mainAxis,
					childFlexBasis+
						collectedFlexItemsValues.remainingFreeSpace/
							collectedFlexItemsValues.totalFlexGrowFactors*
							flexGrowFactor,
					availableInnerMainDim,
					availableInnerWidth)
			}
		}
		deltaFreeSpace += updatedMainSize - childFlexBasis

		marginMain :=
			currentRelativeChild.getMarginForAxis(mainAxis, availableInnerWidth).GetValue()
		marginCross :=
			currentRelativeChild.getMarginForAxis(crossAxis, availableInnerWidth).GetValue()

		var childCrossSize float32
		childMainSize := updatedMainSize + marginMain
		var childCrossMeasureMode MeasureMode
		childMainMeasureMode := MeasureModeExactly

		childStyle := currentRelativeChild.Style
		if !childStyle.AspectRatio.IsUndefined {
			if isMainAxisRow {
				childCrossSize = (childMainSize - marginMain) / childStyle.AspectRatio.Value
			} else {
				childCrossSize = (childMainSize - marginMain) * childStyle.AspectRatio.Value
			}
			childCrossMeasureMode = MeasureModeExactly

			childCrossSize += marginCross
		} else if !FloatIsUndefined(availableInnerCrossDim) &&
			!NodeIsStyleDimDefined(
				currentRelativeChild, crossAxis, availableInnerCrossDim) &&
			measureModeCrossDim == MeasureModeExactly &&
			!(isNodeFlexWrap && mainAxisOverflows) &&
			NodeAlignItem(node, currentRelativeChild) == AlignStretch &&
			currentRelativeChild.marginLeadingValue(crossAxis).Unit !=
				UnitAuto &&
			currentRelativeChild.marginTrailingValue(crossAxis).Unit !=
				UnitAuto {
			childCrossSize = availableInnerCrossDim
			childCrossMeasureMode = MeasureModeExactly
		} else if !NodeIsStyleDimDefined(
			currentRelativeChild, crossAxis, availableInnerCrossDim) {
			childCrossSize = availableInnerCrossDim
			if FloatIsUndefined(childCrossSize) {
				childCrossMeasureMode = MeasureModeUndefined
			} else {
				childCrossMeasureMode = MeasureModeAtMost
			}
		} else {
			childCrossSize =
				ResolveValue(
					*currentRelativeChild.getResolvedDimension(int(dim[crossAxis])),
					availableInnerCrossDim).GetValue() +
					marginCross
			isLoosePercentageMeasurement :=
				currentRelativeChild.getResolvedDimension(int(dim[crossAxis])).Unit ==
					UnitPercent &&
					measureModeCrossDim != MeasureModeExactly
			if FloatIsUndefined(childCrossSize) || isLoosePercentageMeasurement {
				childCrossMeasureMode = MeasureModeUndefined
			} else {
				childCrossMeasureMode = MeasureModeExactly
			}
		}

		ConstrainMaxSizeForMode(
			currentRelativeChild,
			mainAxis,
			availableInnerMainDim,
			availableInnerWidth,
			&childMainMeasureMode,
			&childMainSize)
		ConstrainMaxSizeForMode(
			currentRelativeChild,
			crossAxis,
			availableInnerCrossDim,
			availableInnerWidth,
			&childCrossMeasureMode,
			&childCrossSize)

		requiresStretchLayout :=
			!NodeIsStyleDimDefined(
				currentRelativeChild, crossAxis, availableInnerCrossDim) &&
				NodeAlignItem(node, currentRelativeChild) == AlignStretch &&
				currentRelativeChild.marginLeadingValue(crossAxis).Unit !=
					UnitAuto &&
				currentRelativeChild.marginTrailingValue(crossAxis).Unit != UnitAuto

		var childWidth, childHeight float32
		var childWidthMeasureMode, childHeightMeasureMode MeasureMode
		if isMainAxisRow {
			childWidth = childMainSize
			childHeight = childCrossSize
			childWidthMeasureMode = childMainMeasureMode
			childHeightMeasureMode = childCrossMeasureMode
		} else {
			childWidth = childCrossSize
			childHeight = childMainSize
			childWidthMeasureMode = childCrossMeasureMode
			childHeightMeasureMode = childMainMeasureMode
		}

		isLayoutPass := performLayout && !requiresStretchLayout
		// Recursively call the layout algorithm for this child with the updated
		// main size.
		reason := LayoutPassReasonFlexMeasure
		if isLayoutPass {
			reason = LayoutPassReasonFlexLayout
		}
		layoutNodeInternal(
			currentRelativeChild,
			childWidth,
			childHeight,
			node.Layout.Direction,
			childWidthMeasureMode,
			childHeightMeasureMode,
			availableInnerWidth,
			availableInnerHeight,
			isLayoutPass,
			reason,
			config,
			layoutMarkerData,
			layoutContext,
			depth,
			generationCount)
		node.SetLayoutHadOverflow(
			node.Layout.HadOverflow ||
				currentRelativeChild.Layout.HadOverflow)
	}
	return deltaFreeSpace
}

// CalculateCollectFlexItemsRowValues
func CalculateCollectFlexItemsRowValues(node *Node, ownerDirection Direction, mainAxisownerSize, availableInnerWidth, availableInnerMainDim float32, startOfLineIndex, lineCount uint32) CollectFlexItemsRowValues {
	var flexAlgoRowMeasurement CollectFlexItemsRowValues
	var sizeConsumedOnCurrentLineIncludingMinConstraint float32 = 0
	mainAxis := ResolveFlexDirection(node.Style.FlexDirection, node.resolveDirection(ownerDirection))
	isNodeFlexWrap := node.Style.FlexWrap != WrapNoWrap
	gap := node.getGapForAxis(mainAxis, availableInnerWidth).GetValue()
	// Add items to the current line until it's full or we run out of items.
	endOfLineIndex := startOfLineIndex
	for endOfLineIndex < uint32(len(node.Children)) {
		child := node.Children[endOfLineIndex]
		if child.Style.Display == DisplayNone || child.Style.PositionType == PositionTypeAbsolute {
			continue
		}
		isFirstElementInLine := (endOfLineIndex - startOfLineIndex) == 0
		child.setLineIndex(lineCount)
		childMarginMainAxis := child.getMarginForAxis(mainAxis, availableInnerWidth).GetValue()
		childLeadingGapMainAxis := float32(0)
		if !isFirstElementInLine {
			childLeadingGapMainAxis = gap
		}
		flexBasisWithMinAndMaxConstraints := NodeBoundAxisWithinMinAndMax(child, mainAxis, child.Layout.computedFlexBasis, mainAxisownerSize).GetValue()
		// If this is a multi-line flow and this item pushes us over the available
		// size, we've hit the end of the current line. Break out of the loop and
		// lay out the current line.
		if sizeConsumedOnCurrentLineIncludingMinConstraint+
			flexBasisWithMinAndMaxConstraints+childMarginMainAxis+
			childLeadingGapMainAxis >
			availableInnerMainDim &&
			isNodeFlexWrap && flexAlgoRowMeasurement.itemsOnLine > 0 {
			break
		}

		sizeConsumedOnCurrentLineIncludingMinConstraint +=
			flexBasisWithMinAndMaxConstraints + childMarginMainAxis +
				childLeadingGapMainAxis
		flexAlgoRowMeasurement.sizeConsumedOnCurrentLine +=
			flexBasisWithMinAndMaxConstraints + childMarginMainAxis +
				childLeadingGapMainAxis
		flexAlgoRowMeasurement.itemsOnLine++

		if child.isNodeFlexible() {
			flexAlgoRowMeasurement.totalFlexGrowFactors += child.resolveFlexGrow()

			// Unlike the grow factor, the shrink factor is scaled relative to the
			// child dimension.
			flexAlgoRowMeasurement.totalFlexShrinkScaledFactors +=
				-child.resolveFlexShrink() *
					child.Layout.computedFlexBasis.GetValue()
		}

		flexAlgoRowMeasurement.relativeChildren = append(
			flexAlgoRowMeasurement.relativeChildren, child)
		endOfLineIndex++
	}
	// The total flex factor needs to be floored to 1.
	if flexAlgoRowMeasurement.totalFlexGrowFactors > 0 &&
		flexAlgoRowMeasurement.totalFlexGrowFactors < 1 {
		flexAlgoRowMeasurement.totalFlexGrowFactors = 1
	}

	// The total flex shrink factor needs to be floored to 1.
	if flexAlgoRowMeasurement.totalFlexShrinkScaledFactors > 0 &&
		flexAlgoRowMeasurement.totalFlexShrinkScaledFactors < 1 {
		flexAlgoRowMeasurement.totalFlexShrinkScaledFactors = 1
	}
	flexAlgoRowMeasurement.endOfLineIndex = endOfLineIndex
	return flexAlgoRowMeasurement
}

// NodeComputeFlexBasisForChildren
func NodeComputeFlexBasisForChildren(node *Node, availableInnerWidth, availableInnerHeight float32, widthMeasureMode, heightMeasureMode MeasureMode, direction Direction, mainAxis FlexDirection, config *Config, performLayout bool, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32) float32 {
	var totalOuterFlexBasis float32
	var singleFlexChild *Node
	var measureModeMainDim MeasureMode
	children := node.Children
	if FlexDirectionIsRow(mainAxis) {
		measureModeMainDim = widthMeasureMode
	} else {
		measureModeMainDim = heightMeasureMode
	}
	// If there is only one child with flexGrow + flexShrink it means we can set
	// the computedFlexBasis to 0 instead of measuring and shrinking / flexing the
	// child to exactly match the remaining space
	if measureModeMainDim == MeasureModeExactly {
		for _, child := range children {
			if child.isNodeFlexible() {
				if singleFlexChild != nil || child.resolveFlexGrow() == 0 || child.resolveFlexShrink() == 0 {
					// There is already a flexible child, or this flexible child doesn't
					// have flexGrow and flexShrink, abort
					singleFlexChild = nil
					break
				}
				singleFlexChild = child
			}
		}
	}
	for _, child := range children {
		child.resolveDimension()
		if child.Style.Display == DisplayNone {
			ZeroOutLayoutRecursively(node, layoutContext)
			node.hasNewLayout = true
			node.SetDirty(true)
			continue
		}
		if performLayout {
			// Set the initial position (relative to the owner).
			childDirection := node.resolveDirection(direction)
			var mainDim, crossDim float32
			if FlexDirectionIsRow(mainAxis) {
				mainDim = availableInnerWidth
				crossDim = availableInnerHeight
			} else {
				mainDim = availableInnerHeight
				crossDim = availableInnerWidth
			}
			child.setPosition(childDirection, mainDim, crossDim, availableInnerWidth)
		}
		if child.Style.PositionType == PositionTypeAbsolute {
			continue
		}
		if child == singleFlexChild {
			child.SetLayoutComputedFlexBasisGeneration(generationCount)
			child.SetLayoutComputedFlexBasis(NewFloatOptional(0))
		} else {
			NodeComputeFlexBasisForChild(node,
				child,
				availableInnerWidth,
				widthMeasureMode,
				availableInnerHeight,
				availableInnerWidth,
				availableInnerHeight,
				heightMeasureMode,
				direction,
				config,
				layoutMarkerData,
				layoutContext,
				depth,
				generationCount)
		}
		totalOuterFlexBasis = child.Layout.computedFlexBasis.Add(child.getMarginForAxis(mainAxis, availableInnerWidth)).GetValue()
	}
	return totalOuterFlexBasis
}

// NodeComputeFlexBasisForChild
func NodeComputeFlexBasisForChild(node, child *Node, width float32, widthMode MeasureMode, height, ownerWidth, ownerHeight float32, heightMode MeasureMode, direction Direction, config *Config, layoutMarkerData *LayoutData, layoutContext any, depth uint32, generationCount uint32) {
	mainAxis := ResolveFlexDirection(node.Style.FlexDirection, direction)
	isMainAxisRow := flexDirectionIsRow(mainAxis)
	var mainAxisSize, mainAxisownerSize float32
	if isMainAxisRow {
		mainAxisSize = width
		mainAxisownerSize = ownerWidth
	} else {
		mainAxisSize = height
		mainAxisownerSize = ownerHeight
	}
	var childWidthMeasureMode, childHeightMeasureMode MeasureMode
	resolvedFlexBasis := ResolveValue(child.resolveFlexBasisPtr(), mainAxisownerSize)
	isRowStyleDimDefined := NodeIsStyleDimDefined(child, FlexDirectionRow, ownerWidth)
	isColumnStyleDimDefined := NodeIsStyleDimDefined(child, FlexDirectionColumn, ownerHeight)

	if !resolvedFlexBasis.IsUndefined && !FloatIsUndefined(mainAxisSize) {
		if child.Layout.computedFlexBasis.IsUndefined || (child.Config.experimentalFeatureEnabled(ExperimentalFeatureWebFlexBasis) && child.Layout.computedFlexBasisGeneration != generationCount) {
			paddingAndBorder := NewFloatOptional(NodePaddingAndBorderForAxis(child, mainAxis, ownerWidth))
			node.SetLayoutComputedFlexBasis(FloatOptionalMax(resolvedFlexBasis, paddingAndBorder))
		}
	} else if isMainAxisRow && isRowStyleDimDefined {
		// The width is definite, so use that as the flex basis.
		paddingAndBorder := NewFloatOptional(
			NodePaddingAndBorderForAxis(child, FlexDirectionRow, ownerWidth))

		child.SetLayoutComputedFlexBasis(FloatOptionalMax(
			ResolveValue(
				*child.getResolvedDimensions()[DimensionWidth], ownerWidth),
			paddingAndBorder))
	} else if !isMainAxisRow && isColumnStyleDimDefined {
		// The height is definite, so use that as the flex basis.
		paddingAndBorder := NewFloatOptional(NodePaddingAndBorderForAxis(child, FlexDirectionColumn, ownerWidth))
		child.SetLayoutComputedFlexBasis(FloatOptionalMax(ResolveValue(*child.getResolvedDimensions()[DimensionHeight], ownerHeight), paddingAndBorder))
	} else {
		// Compute the flex basis and hypothetical main size (i.e. the clamped flex
		// basis).
		childWidth := Undefined
		childHeight := Undefined
		childWidthMeasureMode = MeasureModeUndefined
		childHeightMeasureMode = MeasureModeUndefined
		marginRow := child.getMarginForAxis(FlexDirectionRow, ownerWidth).GetValue()
		marginColumn := child.getMarginForAxis(FlexDirectionColumn, ownerWidth).GetValue()
		if isRowStyleDimDefined {
			childWidth = ResolveValue(*child.getResolvedDimensions()[DimensionWidth], ownerWidth).GetValue() + marginRow
			childWidthMeasureMode = MeasureModeExactly
		}
		if isColumnStyleDimDefined {
			childHeight = ResolveValue(*child.getResolvedDimensions()[DimensionHeight], ownerHeight).GetValue() + marginColumn
			childHeightMeasureMode = MeasureModeExactly

		}
		// The W3C spec doesn't say anything about the 'overflow' property, but all
		// major browsers appear to implement the following logic.
		if (!isMainAxisRow && node.Style.Overflow == OverflowScroll) ||
			node.Style.Overflow != OverflowScroll {
			if FloatIsUndefined(childWidth) && !FloatIsUndefined(width) {
				childWidth = width
				childWidthMeasureMode = MeasureModeAtMost
			}
			if FloatIsUndefined(childHeight) && !FloatIsUndefined(height) {
				childHeight = height
				childHeightMeasureMode = MeasureModeAtMost

			}
		}

		childStyle := child.Style
		if !childStyle.AspectRatio.IsUndefined {
			if !isMainAxisRow && childWidthMeasureMode == MeasureModeExactly {
				childHeight = marginColumn + (childWidth-marginRow)/childStyle.AspectRatio.Value
				childHeightMeasureMode = MeasureModeExactly
			} else if isMainAxisRow && childHeightMeasureMode == MeasureModeExactly {
				childWidth = marginRow + (childHeight-marginColumn)*childStyle.AspectRatio.Value
				childWidthMeasureMode = MeasureModeExactly
			}

		}

		// If child has no defined size in the cross axis and is set to stretch, set
		// the cross axis to be measured exactly with the available inner width
		hasExactWidth := !FloatIsUndefined(width) && widthMode == MeasureModeExactly
		childWidthStretch := NodeAlignItem(node, child) == AlignStretch && childWidthMeasureMode != MeasureModeExactly
		if !isMainAxisRow && !isRowStyleDimDefined && hasExactWidth &&
			childWidthStretch {
			childWidth = width
			childWidthMeasureMode = MeasureModeExactly
			if !childStyle.AspectRatio.IsUndefined {
				childHeight =
					(childWidth - marginRow) / childStyle.AspectRatio.Value
				childHeightMeasureMode = MeasureModeExactly
			}
		}

		hasExactHeight := !FloatIsUndefined(height) && heightMode == MeasureModeExactly
		childHeightStretch := NodeAlignItem(node, child) == AlignStretch && childHeightMeasureMode != MeasureModeExactly
		if isMainAxisRow && !isColumnStyleDimDefined && hasExactHeight &&
			childHeightStretch {
			childHeight = height
			childHeightMeasureMode = MeasureModeExactly
			if !childStyle.AspectRatio.IsUndefined {
				childWidth = (childHeight - marginColumn) * childStyle.AspectRatio.Value
				childWidthMeasureMode = MeasureModeExactly
			}
		}

		ConstrainMaxSizeForMode(
			child,
			FlexDirectionRow,
			ownerWidth,
			ownerWidth,
			&childWidthMeasureMode,
			&childWidth)
		ConstrainMaxSizeForMode(
			child,
			FlexDirectionColumn,
			ownerHeight,
			ownerWidth,
			&childHeightMeasureMode,
			&childHeight)

		// Measure the child
		layoutNodeInternal(
			child,
			childWidth,
			childHeight,
			direction,
			childWidthMeasureMode,
			childHeightMeasureMode,
			ownerWidth,
			ownerHeight,
			false,
			LayoutPassReasonMeasureChild,
			config,
			layoutMarkerData,
			layoutContext,
			depth,
			generationCount)
		child.SetLayoutComputedFlexBasis(NewFloatOptional(FloatMax(
			child.Layout.measuredDimensions[dim[mainAxis]],
			NodePaddingAndBorderForAxis(child, mainAxis, ownerWidth))))
	}
	child.SetLayoutComputedFlexBasisGeneration(generationCount)
}

// ConstrainMaxSizeForMode
func ConstrainMaxSizeForMode(node *Node, axis FlexDirection, ownerAxisSize, ownerWidth float32, mode *MeasureMode, size *float32) {
	maxSize := ResolveValue(node.Style.MaxDimensions[axis], ownerAxisSize).Add(node.getMarginForAxis(axis, ownerWidth))
	switch *mode {
	case MeasureModeExactly:
	case MeasureModeAtMost:
		if maxSize.IsUndefined || *size < maxSize.Value {
			*size = *size
		} else {
			*size = maxSize.Value
		}
		break
	case MeasureModeUndefined:
		if !maxSize.IsUndefined {
			*mode = MeasureModeAtMost
			*size = maxSize.Value
		}
		break
	}
}

// NodeAlignItem
func NodeAlignItem(node, child *Node) Align {
	align := child.Style.AlignSelf
	if child.Style.AlignSelf == AlignAuto {
		align = node.Style.AlignItems
	}
	if align == AlignBaseline && flexDirectionIsColumn(node.Style.FlexDirection) {
		return AlignFlexStart
	}
	return align
}

// ZeroOutLayoutRecursively
func ZeroOutLayoutRecursively(node *Node, layoutContext any) {
	node.Layout = &Layout{}
	node.SetLayoutDimension(0, 0)
	node.SetLayoutDimension(0, 1)
	node.hasNewLayout = true
	node.iterChildrenAfterCloningIfNeeded(ZeroOutLayoutRecursively, layoutContext)
}

// NodeCalculateAvailableInnerDim
func NodeCalculateAvailableInnerDim(node *Node, dimension Dimension, availableDim, paddingAndBorder, ownerDim float32) float32 {
	var minInnerDim, maxInnerDim float32
	availableInnerDim := availableDim - paddingAndBorder
	// The value of availableInnerDim includes only the padding and border sizes. We
	// want to exclude the size of the container itself.
	if !FloatIsUndefined(availableInnerDim) {
		minDimensionOptional := ResolveValue(node.Style.MinDimensions[dimension], ownerDim)
		if minDimensionOptional.IsUndefined {
			minInnerDim = 0
		} else {
			minInnerDim = minDimensionOptional.Value - paddingAndBorder
		}
		maxDimensionOptional := ResolveValue(node.Style.MaxDimensions[dimension], ownerDim)
		if maxDimensionOptional.IsUndefined {
			maxInnerDim = math.MaxFloat32
		} else {
			maxInnerDim = maxDimensionOptional.Value - paddingAndBorder
		}
		// We never want to go smaller than padding
		availableInnerDim = max(min(availableInnerDim, maxInnerDim), minInnerDim)
	}
	return availableInnerDim
}

// NodeFixedSizeSetMeasuredDimensions
func NodeFixedSizeSetMeasuredDimensions(node *Node, availableWidth, availableHeight float32, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32) bool {
	if (!FloatIsUndefined(availableWidth) &&
		widthMeasureMode == MeasureModeAtMost && availableWidth <= 0.0) ||
		(!FloatIsUndefined(availableHeight) &&
			heightMeasureMode == MeasureModeAtMost && availableHeight <= 0.0) ||
		(widthMeasureMode == MeasureModeExactly && heightMeasureMode == MeasureModeAtMost) {
		var valueWidth, valueHeight float32
		if FloatIsUndefined(availableWidth) || (widthMeasureMode == MeasureModeAtMost && availableWidth < 0.0) {
			valueWidth = 0
		} else {
			valueWidth = availableWidth
		}
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionRow, valueWidth, ownerWidth, ownerWidth), int(DimensionWidth))
		if FloatIsUndefined(availableHeight) || (heightMeasureMode == MeasureModeAtMost && availableHeight < 0.0) {
			valueHeight = 0
		} else {
			valueHeight = availableHeight
		}
		node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionColumn, valueHeight, ownerHeight, ownerWidth), int(DimensionHeight))
		return true
	}
	return false
}

// NodeEmptyContainerSetMeasuredDimensions
func NodeEmptyContainerSetMeasuredDimensions(node *Node, availableWidth, availableHeight float32, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32) {
	padding := &node.Layout.Padding
	border := &node.Layout.Border
	width := availableWidth
	if widthMeasureMode == MeasureModeUndefined || widthMeasureMode == MeasureModeAtMost {
		width = padding[EdgeLeft] + padding[EdgeRight] + border[EdgeLeft] + border[EdgeRight]
	}
	node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionRow, width, ownerWidth, ownerWidth), int(DimensionWidth))
	height := availableHeight
	if heightMeasureMode == MeasureModeUndefined || heightMeasureMode == MeasureModeAtMost {
		height = padding[EdgeTop] + padding[EdgeBottom] + border[EdgeTop] + border[EdgeBottom]
	}
	node.SetLayoutMeasuredDimension(NodeBoundAxis(node, FlexDirectionColumn, height, ownerHeight, ownerWidth), int(DimensionHeight))
}

// NodeWithMeasureFuncSetMeasuredDimensions
func NodeWithMeasureFuncSetMeasuredDimensions(node *Node, availableWidth, availableHeight float32, widthMeasureMode, heightMeasureMode MeasureMode, ownerWidth, ownerHeight float32, layoutMarkerData *LayoutData, layoutContext any, reason LayoutPassReason) {
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
