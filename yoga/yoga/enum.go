package yoga

type YGAlign int

const (
	YGAlignAuto YGAlign = iota
	YGAlignFlexStart
	YGAlignCenter
	YGAlignFlexEnd
	YGAlignStretch
	YGAlignBaseline
	YGAlignSpaceBetween
	YGAlignSpaceAround
)

func (a YGAlign) String() string {
	switch a {
	case YGAlignAuto:
		return "auto"
	case YGAlignFlexStart:
		return "flex-start"
	case YGAlignCenter:
		return "center"
	case YGAlignFlexEnd:
		return "flex-end"
	case YGAlignStretch:
		return "stretch"
	case YGAlignBaseline:
		return "baseline"
	case YGAlignSpaceBetween:
		return "space-between"
	case YGAlignSpaceAround:
		return "space-around"
	}
	return ""
}

type YGDimension int

const (
	YGDimensionWidth YGDimension = iota
	YGDimensionHeight
)

func (d YGDimension) String() string {
	switch d {
	case YGDimensionWidth:
		return "width"
	case YGDimensionHeight:
		return "height"
	}
	return ""
}

type YGDirection int

const (
	YGDirectionInherit YGDirection = iota
	YGDirectionLTR
	YGDirectionRTL
)

func (d YGDirection) String() string {
	switch d {
	case YGDirectionInherit:
		return "inherit"
	case YGDirectionLTR:
		return "ltr"
	case YGDirectionRTL:
		return "rtl"
	}
	return ""
}

type YGDisplay int

const (
	YGDisplayFlex YGDisplay = iota
	YGDisplayNone
)

func (d YGDisplay) String() string {
	switch d {
	case YGDisplayFlex:
		return "flex"
	case YGDisplayNone:
		return "none"
	}
	return ""
}

type YGEdge int

const (
	YGEdgeLeft YGEdge = iota
	YGEdgeTop
	YGEdgeRight
	YGEdgeBottom
	YGEdgeStart
	YGEdgeEnd
	YGEdgeHorizontal
	YGEdgeVertical
	YGEdgeAll
)

func (e YGEdge) String() string {
	switch e {
	case YGEdgeLeft:
		return "left"
	case YGEdgeTop:
		return "top"
	case YGEdgeRight:
		return "right"
	case YGEdgeBottom:
		return "bottom"
	case YGEdgeStart:
		return "start"
	case YGEdgeEnd:
		return "end"
	case YGEdgeHorizontal:
		return "horizontal"
	case YGEdgeVertical:
		return "vertical"
	case YGEdgeAll:
		return "all"
	}
	return ""
}

type YGErrata int

const (
	YGErrataNone YGErrata = iota
	YGErrataStretchFlexBasis
	YGErrataAll     = 2147483647
	YGErrataClassic = 2147483646
)

func (e YGErrata) String() string {
	switch e {
	case YGErrataNone:
		return "none"
	case YGErrataStretchFlexBasis:
		return "stretch-flex-basis"
	case YGErrataAll:
		return "all"
	case YGErrataClassic:
		return "classic"
	}
	return ""
}

type YGExperimentalFeature int

const (
	YGExperimentalFeatureWebFlexBasis YGExperimentalFeature = iota
	YGExperimentalFeatureAbsolutePercentageAgainstPaddingEdge
	YGExperimentalFeatureFixJNILocalRefOverflows
)

func (e YGExperimentalFeature) String() string {
	switch e {
	case YGExperimentalFeatureWebFlexBasis:
		return "web-flex-basis"
	case YGExperimentalFeatureAbsolutePercentageAgainstPaddingEdge:
		return "absolute-percentage-against-padding-edge"
	case YGExperimentalFeatureFixJNILocalRefOverflows:
		return "fix-jni-local-ref-overflows"
	}
	return ""
}

type YGFlexDirection int

const (
	YGFlexDirectionColumn YGFlexDirection = iota
	YGFlexDirectionColumnReverse
	YGFlexDirectionRow
	YGFlexDirectionRowReverse
)

func (f YGFlexDirection) String() string {
	switch f {
	case YGFlexDirectionColumn:
		return "column"
	case YGFlexDirectionColumnReverse:
		return "column-reverse"
	case YGFlexDirectionRow:
		return "row"
	case YGFlexDirectionRowReverse:
		return "row-reverse"
	}
	return ""
}

type YGGutter int

const (
	YGGutterColumn YGGutter = iota
	YGGutterRow
	YGGutterAll
)

func (g YGGutter) String() string {
	switch g {
	case YGGutterColumn:
		return "column"
	case YGGutterRow:
		return "row"
	case YGGutterAll:
		return "all"
	}
	return ""
}

type YGJustify int

const (
	YGJustifyFlexStart YGJustify = iota
	YGJustifyCenter
	YGJustifyFlexEnd
	YGJustifySpaceBetween
	YGJustifySpaceAround
	YGJustifySpaceEvenly
)

func (j YGJustify) String() string {
	switch j {
	case YGJustifyCenter:
		return "center"
	case YGJustifyFlexEnd:
		return "flex-end"
	case YGJustifySpaceBetween:
		return "space-between"
	case YGJustifySpaceAround:
		return "space-around"
	case YGJustifyFlexStart:
		return "flex-start"
	case YGJustifySpaceEvenly:
		return "space-evenly"
	}
	return ""
}

type YGLogLevel int

const (
	YGLogLevelError YGLogLevel = iota
	YGLogLevelWarn
	YGLogLevelInfo
	YGLogLevelDebug
	YGLogLevelVerbose
	YGLogLevelFatal
)

func (l YGLogLevel) String() string {
	switch l {
	case YGLogLevelError:
		return "error"
	case YGLogLevelWarn:
		return "warn"
	case YGLogLevelInfo:
		return "info"
	case YGLogLevelDebug:
		return "debug"
	case YGLogLevelVerbose:
		return "verbose"
	case YGLogLevelFatal:
		return "fatal"
	}
	return ""
}

type YGMeasureMode int

const (
	YGMeasureModeUndefined YGMeasureMode = iota
	YGMeasureModeExactly
	YGMeasureModeAtMost
)

func (m YGMeasureMode) String() string {
	switch m {
	case YGMeasureModeUndefined:
		return "undefined"
	case YGMeasureModeExactly:
		return "exactly"
	case YGMeasureModeAtMost:
		return "at-most"
	}
	return ""
}

type YGNodeType int

const (
	YGNodeTypeDefault YGNodeType = iota
	YGNodeTypeText
)

func (n YGNodeType) String() string {
	switch n {
	case YGNodeTypeDefault:
		return "default"
	case YGNodeTypeText:
		return "text"
	}
	return ""
}

type YGOverflow int

const (
	YGOverflowVisible YGOverflow = iota
	YGOverflowHidden
	YGOverflowScroll
)

func (o YGOverflow) String() string {
	switch o {
	case YGOverflowVisible:
		return "visible"
	case YGOverflowHidden:
		return "hidden"
	case YGOverflowScroll:
		return "scroll"
	}
	return ""
}

type YGPositionType int

const (
	YGPositionTypeStatic YGPositionType = iota
	YGPositionTypeRelative
	YGPositionTypeAbsolute
)

func (p YGPositionType) String() string {
	switch p {
	case YGPositionTypeStatic:
		return "static"
	case YGPositionTypeRelative:
		return "relative"
	case YGPositionTypeAbsolute:
		return "absolute"
	}
	return ""
}

type YGPrintOptions int

const (
	YGPrintOptionsLayout YGPrintOptions = 1 << iota
	YGPrintOptionsStyle
	YGPrintOptionsChildren
)

func (p YGPrintOptions) String() string {
	switch p {
	case YGPrintOptionsLayout:
		return "layout"
	case YGPrintOptionsStyle:
		return "style"
	case YGPrintOptionsChildren:
		return "children"
	}
	return ""
}

type YGUnit int

const (
	YGUnitUndefined YGUnit = iota
	YGUnitPoint
	YGUnitPercent
	YGUnitAuto
)

func (u YGUnit) String() string {
	switch u {
	case YGUnitUndefined:
		return "undefined"
	case YGUnitPoint:
		return "point"
	case YGUnitPercent:
		return "percent"
	case YGUnitAuto:
		return "auto"
	}
	return ""
}

type YGWrap int

const (
	YGWrapNoWrap YGWrap = iota
	YGWrapWrap
	YGWrapWrapReverse
)

func (w YGWrap) String() string {
	switch w {
	case YGWrapNoWrap:
		return "no-wrap"
	case YGWrapWrap:
		return "wrap"
	case YGWrapWrapReverse:
		return "wrap-reverse"
	}
	return ""
}
