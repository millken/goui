package yoga

import "unsafe"

type CollectFlexItemsRowValues struct {
	itemsOnLine                  uint32
	sizeConsumedOnCurrentLine    float32
	totalFlexGrowFactors         float32
	totalFlexShrinkScaledFactors float32
	endOfLineIndex               uint32
	relativeChildren             []*Node
	remainingFreeSpace           float32
	mainDim                      float32
	crossDim                     float32
}

func uint32ToFloat(u uint32) float32 {
	return *(*float32)(unsafe.Pointer(&u))
}

func floatToUint32(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}

func FlexDirectionIsRow(direction FlexDirection) bool {
	return direction == FlexDirectionRow || direction == FlexDirectionRowReverse
}

func FlexDirectionIsColumn(direction FlexDirection) bool {
	return direction == FlexDirectionColumn || direction == FlexDirectionColumnReverse
}

func ResolveFlexDirection(flexDirection FlexDirection, direction Direction) FlexDirection {
	if direction == DirectionRTL {
		if flexDirection == FlexDirectionRow {
			return FlexDirectionRowReverse
		} else if flexDirection == FlexDirectionRowReverse {
			return FlexDirectionRow
		}
	}

	return flexDirection
}

func FlexDirectionCross(flexDirection FlexDirection, direction Direction) FlexDirection {
	if FlexDirectionIsColumn(flexDirection) {
		return ResolveFlexDirection(FlexDirectionRow, direction)
	} else {
		return FlexDirectionColumn
	}
}

func ResolveValue(value YGValue, ownerSize float32) YGFloatOptional {
	switch value.unit {
	case YGUnitPoint:
		return NewYGFloatOptional(value.value)
	case YGUnitPercent:
		return NewYGFloatOptional(value.value * ownerSize * 0.01)
	default:
		return undefinedFloatOptional
	}
}

func ResolveValueMargin(value Value, ownerSize float32) FloatOptional {
	if value.Unit == UnitAuto {
		return NewFloatOptional(0)
	}
	return ResolveValue(value, ownerSize)
}

func If[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}

func FloatIsUndefined[T float32 | float64](value T) bool {
	return value != value
}

func YGFloatsEqual[T float32 | float64](a, b T) bool {
	if FloatIsUndefined(a) && FloatIsUndefined(b) {
		return true
	}
	return a == b
}
