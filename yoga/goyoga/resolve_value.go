package goyoga

func resolveValue(value YGValue, ownerSize float32) FloatOptional {
	switch value.unit {
	case YGUnitPoint:
		return NewFloatOptional(value.value)
	case YGUnitPercent:
		return NewFloatOptional(value.value * ownerSize * 0.01)
	default:
		return undefinedFloatOptional
	}
}

func resolveCompactValue(value CompactValue, ownerSize float32) FloatOptional {
	return resolveValue(value.YGValue(), ownerSize)
}
