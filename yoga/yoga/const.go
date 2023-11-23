package yoga

var (
	YGUndefined float32 = NAN
)

const (
	// AlignCount as defined in yoga/YGEnums.h:14
	AlignCount = 8
	// DimensionCount as defined in yoga/YGEnums.h:27
	DimensionCount = 2
	// DirectionCount as defined in yoga/YGEnums.h:34
	DirectionCount = 3
	// DisplayCount as defined in yoga/YGEnums.h:42
	DisplayCount = 2
	// EdgeCount as defined in yoga/YGEnums.h:49
	EdgeCount = 9
	// ExperimentalFeatureCount as defined in yoga/YGEnums.h:63
	ExperimentalFeatureCount = 1
	// FlexDirectionCount as defined in yoga/YGEnums.h:69
	FlexDirectionCount = 4
	// JustifyCount as defined in yoga/YGEnums.h:78
	JustifyCount = 6
	// LogLevelCount as defined in yoga/YGEnums.h:89
	LogLevelCount = 6
	// MeasureModeCount as defined in yoga/YGEnums.h:100
	MeasureModeCount = 3
	// NodeTypeCount as defined in yoga/YGEnums.h:108
	NodeTypeCount = 2
	// OverflowCount as defined in yoga/YGEnums.h:115
	OverflowCount = 3
	// PositionTypeCount as defined in yoga/YGEnums.h:123
	PositionTypeCount = 2
	// PrintOptionsCount as defined in yoga/YGEnums.h:130
	PrintOptionsCount = 3
	// UnitCount as defined in yoga/YGEnums.h:138
	UnitCount = 4
	// WrapCount as defined in yoga/YGEnums.h:147
	WrapCount = 3
)

type Align int32

const (
	AlignAuto Align = iota
	AlignFlexStart
	AlignCenter
	AlignFlexEnd
	AlignStretch
	AlignBaseline
	AlignSpaceBetween
	AlignSpaceAround
)

func (a Align) String() string {
	switch a {
	case AlignCenter:
		return "center"
	case AlignFlexEnd:
		return "flex-end"
	case AlignStretch:
		return "stretch"
	case AlignBaseline:
		return "base-line"
	case AlignFlexStart:
		return "flex-start"
	case AlignAuto:
		return "auto"
	}
	return ""
}

type Dimension int32

const (
	DimensionWidth Dimension = iota
	DimensionHeight
)

type Direction int32

const (
	DirectionInherit Direction = iota
	DirectionLTR
	DirectionRTL
)

type Display int32

const (
	DisplayFlex Display = iota
	DisplayNone
)

type Edge int32

const (
	EdgeLeft Edge = iota
	EdgeTop
	EdgeRight
	EdgeBottom
	EdgeStart
	EdgeEnd
	EdgeHorizontal
	EdgeVertical
	EdgeAll
)

// ExperimentalFeature as declared in yoga/YGEnums.h:66
type ExperimentalFeature int32

// ExperimentalFeature enumeration from yoga/YGEnums.h:66
const (
	ExperimentalFeatureWebFlexBasis ExperimentalFeature = iota
)

type FlexDirection int32

const (
	FlexDirectionColumn FlexDirection = iota
	FlexDirectionColumnReverse
	FlexDirectionRow
	FlexDirectionRowReverse
)

func (f FlexDirection) String() string {
	switch f {
	case FlexDirectionColumn:
		return "column"
	case FlexDirectionColumnReverse:
		return "column-reverse"
	case FlexDirectionRow:
		return "row"
	case FlexDirectionRowReverse:
		return "row-reverse"
	}
	return ""
}

type Gutter uint8

const (
	GutterColumn Gutter = iota
	GutterRow
	GutterAll
	GutterCount = 3
)

type Justify int32

const (
	JustifyFlexStart Justify = iota
	JustifyCenter
	JustifyFlexEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

func (j Justify) String() string {
	switch j {
	case JustifyCenter:
		return "center"
	case JustifyFlexEnd:
		return "flex-end"
	case JustifySpaceBetween:
		return "space-between"
	case JustifySpaceAround:
		return "space-around"
	case JustifyFlexStart:
		return "flex-start"
	}
	return ""
}

// LogLevel as declared in yoga/YGEnums.h:97
type LogLevel int32

// LogLevel enumeration from yoga/YGEnums.h:97
const (
	LogLevelError   LogLevel = iota
	LogLevelWarn    LogLevel = 1
	LogLevelInfo    LogLevel = 2
	LogLevelDebug   LogLevel = 3
	LogLevelVerbose LogLevel = 4
	LogLevelFatal   LogLevel = 5
)

type MeasureMode int32

const (
	MeasureModeUndefined MeasureMode = iota
	MeasureModeExactly
	MeasureModeAtMost
)

type NodeType int32

const (
	NodeTypeDefault NodeType = iota
	NodeTypeText
)

type Overflow int32

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

func (o Overflow) String() string {
	switch o {
	case OverflowVisible:
		return "visible"
	case OverflowHidden:
		return "hidden"
	case OverflowScroll:
		return "scroll"
	}
	return ""
}

type PositionType int32

const (
	PositionTypeRelative PositionType = iota
	PositionTypeAbsolute
)

func (p PositionType) String() string {
	switch p {
	case PositionTypeRelative:
		return "relative"
	case PositionTypeAbsolute:
		return "absolute"
	}
	return ""
}

type PrintOptions int32

const (
	PrintOptionsLayout PrintOptions = 1 << iota
	PrintOptionsStyle
	PrintOptionsChildren
)

type Unit int32

const (
	UnitUndefined Unit = iota
	UnitPoint
	UnitPercent
	UnitAuto
)

func (u Unit) String() string {
	switch u {
	case UnitPoint:
		return "px"
	case UnitPercent:
		return "%"
	}
	return ""
}

type Wrap int32

const (
	WrapNoWrap Wrap = iota
	WrapWrap
	WrapWrapReverse
)

type Errata int

const (
	ErrataNone             Errata = 0
	ErrataStretchFlexBasis Errata = 1
	ErrataAll              Errata = 2147483647
	ErrataClassic          Errata = 2147483646
)
