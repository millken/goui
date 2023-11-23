package goyoga

type YGMeasureFunc func(
	node *YGNode,
	width float32,
	widthMode YGMeasureMode,
	height float32,
	heightMode YGMeasureMode,
) YGSize

type YGBaselineFunc func(node *YGNode, width, height float32) float32

type YGPrintFunc func(node *YGNode)

type YGDirtiedFunc func(node *YGNode)

type YGNode struct {
	hasNewLayout_        bool
	isReferenceBaseline_ bool
	isDirty_             bool
	nodeType_            YGNodeType
	context_             interface{}
	measureFunc_         YGMeasureFunc
	baselineFunc_        YGBaselineFunc
	printFunc_           YGPrintFunc
	dirtiedFunc_         YGDirtiedFunc
	style_               YGStyle
	layout_              LayoutResults
	lineIndex_           uint32
	owner_               *YGNode
	children_            []*YGNode
	config_              *YGConfig
	resolvedDimensions_  [2]YGValue
}

func NewNode(config *YGConfig) *YGNode {
	node := &YGNode{}
	return node
}

// getHasNewLayout
func (node *YGNode) getHasNewLayout() bool {
	return node.hasNewLayout_
}

// getNodeType
func (node *YGNode) getNodeType() YGNodeType {
	return node.nodeType_
}

// hasMeasureFunc
func (node *YGNode) hasMeasureFunc() bool {
	return node.measureFunc_ != nil
}

func (node *YGNode) measure(
	width float32,
	widthMode YGMeasureMode,
	height float32,
	heightMode YGMeasureMode,
) YGSize {
	return node.measureFunc_(node, width, widthMode, height, heightMode)
}

// hasBaselineFunc
func (node *YGNode) hasBaselineFunc() bool {
	return node.baselineFunc_ != nil
}

// baseline
func (node *YGNode) baseline(width, height float32) float32 {
	return node.baselineFunc_(node, width, height)
}

// hasErrata
func (node *YGNode) hasErrata(errata YGErrata) bool {
	return node.config_.hasErrata(errata)
}

// getDirtiedFunc
func (node *YGNode) getDirtiedFunc() YGDirtiedFunc {
	return node.dirtiedFunc_
}

// getStyle
func (node *YGNode) getStyle() *YGStyle {
	return &node.style_
}

// getLayout
func (node *YGNode) getLayout() *LayoutResults {
	return &node.layout_
}

// getLineIndex
func (node *YGNode) getLineIndex() uint32 {
	return node.lineIndex_
}

// isReferenceBaseline
func (node *YGNode) isReferenceBaseline() bool {
	return node.isReferenceBaseline_
}

// getOwner
func (node *YGNode) getOwner() *YGNode {
	return node.owner_
}

// getChildren
func (node *YGNode) getChildren() []*YGNode {
	return node.children_
}

// getChild
func (node *YGNode) getChild(index uint32) *YGNode {
	return node.children_[index]
}

// getChildrenCount
func (node *YGNode) getChildrenCount() uint32 {
	return uint32(len(node.children_))
}

// getConfig
func (node *YGNode) getConfig() *YGConfig {
	return node.config_
}

// isDirty
func (node *YGNode) isDirty() bool {
	return node.isDirty_
}

// getResolvedDimensions
func (node *YGNode) getResolvedDimensions() [2]YGValue {
	return node.resolvedDimensions_
}

// getResolvedDimension
func (node *YGNode) getResolvedDimension(dimension YGDimension) YGValue {
	return node.resolvedDimensions_[dimension]
}

// computeEdgeValueForColumn
func (node *YGNode) computeEdgeValueForColumn(
	edges *[EdgeCount]CompactValue,
	edge YGEdge,
) CompactValue {
	if edges[edge].isDefined() {
		return edges[edge]
	} else if edges[YGEdgeVertical].isDefined() {
		return edges[YGEdgeVertical]
	} else {
		return edges[YGEdgeAll]
	}
}

// computeEdgeValueForRow
func (node *YGNode) computeEdgeValueForRow(
	edges *[EdgeCount]CompactValue,
	rowEdge YGEdge,
	edge YGEdge,
) CompactValue {
	if edges[rowEdge].isDefined() {
		return edges[rowEdge]
	} else if edges[edge].isDefined() {
		return edges[edge]
	} else if edges[YGEdgeHorizontal].isDefined() {
		return edges[YGEdgeHorizontal]
	} else {
		return edges[YGEdgeAll]
	}
}

// isFlexStartPositionDefined
func (node *YGNode) isFlexStartPositionDefined(axis YGFlexDirection) bool {
	// startEdge := flexStartEdge(axis)
	// leadingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().postion(), YGEdgeStart, startEdge), node.computeEdgeValueForColumn(node.style_.postion(), startEdge))
	return false
}

// isInlineStartPositionDefined
func (node *YGNode) isInlineStartPositionDefined(axis YGFlexDirection, direction YGDirection) bool {
	return false
}

// isFlexEndPositionDefined
func (node *YGNode) isFlexEndPositionDefined(axis YGFlexDirection) bool {
	return false
}

// isInlineEndPositionDefined
func (node *YGNode) isInlineEndPositionDefined(axis YGFlexDirection, direction YGDirection) bool {
	return false
}

// getFlexStartPosition
func (node *YGNode) getFlexStartPosition(axis YGFlexDirection, axisSize float32) float32 {
	return 0
}

// getInlineStartPosition
func (node *YGNode) getInlineStartPosition(axis YGFlexDirection, direction YGDirection, axisSize float32) float32 {
	return 0
}

// getFlexEndPosition
func (node *YGNode) getFlexEndPosition(axis YGFlexDirection, axisSize float32) float32 {
	return 0
}

// getInlineEndPosition
func (node *YGNode) getInlineEndPosition(axis YGFlexDirection, direction YGDirection, axisSize float32) float32 {
	return 0
}

// getFlexStartMargin
func (node *YGNode) getFlexStartMargin(axis YGFlexDirection, widthSize float32) float32 {
	return 0
}

// getInlineStartMargin
func (node *YGNode) getInlineStartMargin(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getFlexEndMargin
func (node *YGNode) getFlexEndMargin(axis YGFlexDirection, widthSize float32) float32 {
	return 0
}

// getInlineEndMargin
func (node *YGNode) getInlineEndMargin(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getFlexStartBorder
func (node *YGNode) getFlexStartBorder(flexDirection YGFlexDirection, direction YGDirection) float32 {
	return 0
}

// getInlineStartBorder
func (node *YGNode) getInlineStartBorder(flexDirection YGFlexDirection, direction YGDirection) float32 {
	return 0
}

// getFlexEndBorder
func (node *YGNode) getFlexEndBorder(flexDirection YGFlexDirection, direction YGDirection) float32 {
	return 0
}

// getInlineEndBorder
func (node *YGNode) getInlineEndBorder(flexDirection YGFlexDirection, direction YGDirection) float32 {
	return 0
}

// getFlexStartPadding
func (node *YGNode) getFlexStartPadding(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getInlineStartPadding
func (node *YGNode) getInlineStartPadding(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getFlexEndPadding
func (node *YGNode) getFlexEndPadding(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getInlineEndPadding
func (node *YGNode) getInlineEndPadding(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getFlexStartPaddingAndBorder
func (node *YGNode) getFlexStartPaddingAndBorder(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getInlineStartPaddingAndBorder
func (node *YGNode) getInlineStartPaddingAndBorder(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getFlexEndPaddingAndBorder
func (node *YGNode) getFlexEndPaddingAndBorder(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getInlineEndPaddingAndBorder
func (node *YGNode) getInlineEndPaddingAndBorder(axis YGFlexDirection, direction YGDirection, widthSize float32) float32 {
	return 0
}

// getMarginForAxis
func (node *YGNode) getMarginForAxis(axis YGFlexDirection, widthSize float32) float32 {
	return 0
}

// getGapForAxis
func (node *YGNode) getGapForAxis(axis YGFlexDirection) float32 {
	return 0
}

// setContext
func (node *YGNode) setContext(context interface{}) {
	node.context_ = context
}

// setPrintFunc
func (node *YGNode) setPrintFunc(printFunc YGPrintFunc) {
	node.printFunc_ = printFunc
}

// setHasNewLayout
func (node *YGNode) setHasNewLayout(hasNewLayout bool) {
	node.hasNewLayout_ = hasNewLayout
}

// setNodeType
func (node *YGNode) setNodeType(nodeType YGNodeType) {
	node.nodeType_ = nodeType
}

// TODO
func (node *YGNode) setMeasureFunc(measureFunc YGMeasureFunc) {
	node.measureFunc_ = measureFunc
}

// setBaselineFunc
func (node *YGNode) setBaselineFunc(baselineFunc YGBaselineFunc) {
	node.baselineFunc_ = baselineFunc
}

// setDirtiedFunc
func (node *YGNode) setDirtiedFunc(dirtiedFunc YGDirtiedFunc) {
	node.dirtiedFunc_ = dirtiedFunc
}

// setStyle
func (node *YGNode) setStyle(style YGStyle) {
	node.style_ = style
}

// setLayout
func (node *YGNode) setLayout(layout LayoutResults) {
	node.layout_ = layout
}

// setLineIndex
func (node *YGNode) setLineIndex(lineIndex uint32) {
	node.lineIndex_ = lineIndex
}

// setIsReferenceBaseline
func (node *YGNode) setIsReferenceBaseline(isReferenceBaseline bool) {
	node.isReferenceBaseline_ = isReferenceBaseline
}

// setOwner
func (node *YGNode) setOwner(owner *YGNode) {
	node.owner_ = owner
}

// setChildren
func (node *YGNode) setChildren(children []*YGNode) {
	node.children_ = children
}

// setConfig TODO
func (node *YGNode) setConfig(config *YGConfig) {
	node.config_ = config
}

// setDirty TODO
func (node *YGNode) setDirty(isDirty bool) {
	node.isDirty_ = isDirty
}

// setLayoutLastOwnerDirection TODO
func (node *YGNode) setLayoutLastOwnerDirection(direction YGDirection) {
}

// setLayoutComputedFlexBasis
func (node *YGNode) setLayoutComputedFlexBasis(computedFlexBasis FloatOptional) {
}

// setLayoutComputedFlexBasisGeneration
func (node *YGNode) setLayoutComputedFlexBasisGeneration(computedFlexBasisGeneration uint32) {
}

// setLayoutMeasuredDimension
func (node *YGNode) setLayoutMeasuredDimension(
	measuredDimension float32,
	dimension YGDimension,
) {
}

// setLayoutHadOverflow
func (node *YGNode) setLayoutHadOverflow(hadOverflow bool) {
}

// setLayoutDimension
func (node *YGNode) setLayoutDimension(dimensionValue float32, dimension YGDimension) {
}

// setLayoutDirection
func (node *YGNode) setLayoutDirection(direction YGDirection) {
}

// setLayoutMargin
func (node *YGNode) setLayoutMargin(margin float32, edge YGEdge) {
}

// setLayoutBorder
func (node *YGNode) setLayoutBorder(border float32, edge YGEdge) {
}

// setLayoutPadding
func (node *YGNode) setLayoutPadding(padding float32, edge YGEdge) {
}

// setLayoutPosition
func (node *YGNode) setLayoutPosition(position float32, edge YGEdge) {
}

// setPosition
func (node *YGNode) setPosition(
	direction YGDirection,
	mainSize float32,
	crossSize float32,
	ownerWidth float32) {
}

// getFlexStartMarginValue
func (node *YGNode) getFlexStartMarginValue(axis YGFlexDirection) YGValue {
	return YGValueUndefined
}

// marginTrailingValue
func (node *YGNode) getMarginTrailingValue(axis YGFlexDirection) YGValue {
	return YGValueUndefined
}

// resolveFlexBasisPtr
func (node *YGNode) resolveFlexBasisPtr() *YGValue {
	return nil
}

// resolveDimension
func (node *YGNode) resolveDimension() {
}

// resolveDirection
func (node *YGNode) resolveDirection(ownerDirection YGDirection) YGDirection {
	return YGDirectionInherit
}

// clearChildren
func (node *YGNode) clearChildren() {
}

// replaceChild
func (node *YGNode) replaceChild(oldChild, newChild *YGNode) {
}

func (node *YGNode) replaceChildIdx(child *YGNode, index uint32) {
}

// insertChild
func (node *YGNode) insertChild(child *YGNode, index uint32) {
}

func (node *YGNode) removeChild(child *YGNode) {
}

func (node *YGNode) removeChildIdx(index uint32) {
}

// cloneChildrenIfNeeded
func (node *YGNode) cloneChildrenIfNeeded() {
}

// markDirtyAndPropagate
func (node *YGNode) markDirtyAndPropagate() {
}

// resolveFlexGrow
func (node *YGNode) resolveFlexGrow() float32 {
	return 0
}

// resolveFlexShrink
func (node *YGNode) resolveFlexShrink() float32 {
	return 0
}

// isNodeFlexible
func (node *YGNode) isNodeFlexible() bool {
	return false
}

// reset
func (node *YGNode) reset() {
}
