package yoga

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

func ResolveValue(value Value, ownerSize float32) FloatOptional {
	switch value.Unit {
	case UnitPoint:
		return NewFloatOptional(value.Value)
	case UnitPercent:
		return NewFloatOptional(value.Value * ownerSize * 0.01)
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

// FloatOptionalMax
func FloatOptionalMax(a FloatOptional, b FloatOptional) FloatOptional {
	if a.IsUndefined {
		return b
	}
	if b.IsUndefined {
		return a
	}
	if a.Value > b.Value {
		return a
	}
	return b
}

// FloatMax
func FloatMax(a float32, b float32) float32 {
	if !FloatIsUndefined(a) && !FloatIsUndefined(b) {
		return max(a, b)
	}
	if FloatIsUndefined(a) {
		return b
	}
	return a
}

// FloatsEqual
func FloatsEqual(a float32, b float32) bool {
	if FloatIsUndefined(a) && FloatIsUndefined(b) {
		return true
	}
	return a == b
}

// Any 若expr成立，则返回a；否则返回b。
func If[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
