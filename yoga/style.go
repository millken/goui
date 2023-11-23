package yoga

var (
	DefaultEdgeValuesUnit = [EdgeCount]Value{
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
		ValueUndefined,
	}
	DefaultDimensionValuesUnit = [DimensionCount]Value{
		ValueUndefined,
		ValueUndefined,
	}
	DefaultDimensionValuesAutoUnit = [2]Value{
		ValueAuto,
		ValueAuto,
	}
	DefaultStyle = Style{
		Flex:           undefinedFloatOptional,
		FlexGrow:       undefinedFloatOptional,
		FlexShrink:     undefinedFloatOptional,
		FlexBasis:      ValueAuto,
		JustifyContent: JustifyFlexStart,
		AlignItems:     AlignStretch,
		AlignContent:   AlignFlexStart,
		Direction:      DirectionInherit,
		FlexDirection:  FlexDirectionColumn,
		PositionType:   PositionTypeRelative,
		FlexWrap:       WrapNoWrap,
		Overflow:       OverflowVisible,
		Display:        DisplayFlex,
		Dimensions:     DefaultDimensionValuesAutoUnit,
		MinDimensions:  DefaultDimensionValuesUnit,
		MaxDimensions:  DefaultDimensionValuesUnit,
		Position:       DefaultEdgeValuesUnit,
		Margin:         DefaultEdgeValuesUnit,
		Padding:        DefaultEdgeValuesUnit,
		Border:         DefaultEdgeValuesUnit,
		AspectRatio:    undefinedFloatOptional,
	}
)

type Style struct {
	Direction      Direction
	FlexDirection  FlexDirection
	JustifyContent Justify
	AlignContent   Align
	AlignItems     Align
	AlignSelf      Align
	PositionType   PositionType
	FlexWrap       Wrap
	Overflow       Overflow
	Display        Display
	Flex           FloatOptional
	FlexGrow       FloatOptional
	FlexShrink     FloatOptional
	FlexBasis      Value
	Margin         [EdgeCount]Value
	Position       [EdgeCount]Value
	Padding        [EdgeCount]Value
	Border         [EdgeCount]Value
	Gap            [GutterCount]Value
	Dimensions     [DimensionCount]Value
	MinDimensions  [DimensionCount]Value
	MaxDimensions  [DimensionCount]Value
	AspectRatio    FloatOptional
}

func (s *Style) alignContent() Align {
	return s.AlignContent
}

// display
func (s *Style) display() Display {
	return s.Display
}

// positionType
func (s *Style) positionType() PositionType {
	return s.PositionType
}

// overflow
func (s *Style) overflow() Overflow {
	return s.Overflow
}

// flexWrap
func (s *Style) flexWrap() Wrap {
	return s.FlexWrap
}
