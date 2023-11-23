package yoga

var (
	YGValueZero      = YGValue{0, YGUnitPoint}
	YGValueUndefined = YGValue{YGUndefined, YGUnitUndefined}
	YGValueAuto      = YGValue{YGUndefined, YGUnitAuto}
)

type YGValue struct {
	value float32
	unit  YGUnit
}

func (v YGValue) IsUndefined() bool {
	return v.unit == YGUnitUndefined
}

func (v *YGValue) equal(other *YGValue) bool {
	if v.unit != other.unit {
		return false
	}
	switch v.unit {
	case YGUnitUndefined:
	case YGUnitAuto:
		return true
	case YGUnitPoint:
	case YGUnitPercent:
		return v.value == other.value
	}
	return false
}

func (v *YGValue) notEqual(other *YGValue) bool {
	return !v.equal(other)
}
