package yoga

type YGFloatOptional struct {
	value float32
}

func NewYGFloatOptional(value float32) YGFloatOptional {
	return YGFloatOptional{value}
}

func (opt YGFloatOptional) unwrap() float32 {
	return opt.value
}

func (opt YGFloatOptional) isUndefined() bool {
	return IsNaN(opt.value)
}

func (lhs YGFloatOptional) equal(rhs YGFloatOptional) bool {
	return lhs.unwrap() == rhs.unwrap() ||
		(lhs.isUndefined() && rhs.isUndefined())
}

func (lhs YGFloatOptional) notEqual(rhs YGFloatOptional) bool {
	return !lhs.equal(rhs)
}

func (lhs YGFloatOptional) lessThan(rhs YGFloatOptional) bool {
	return lhs.unwrap() < rhs.unwrap()
}

func (lhs YGFloatOptional) lessThanOrEqual(rhs YGFloatOptional) bool {
	return lhs.lessThan(rhs) || lhs.equal(rhs)
}

func (lhs YGFloatOptional) greaterThan(rhs YGFloatOptional) bool {
	return !lhs.lessThanOrEqual(rhs)
}

func (lhs YGFloatOptional) greaterThanOrEqual(rhs YGFloatOptional) bool {
	return !lhs.lessThan(rhs)
}

func (lhs YGFloatOptional) add(rhs YGFloatOptional) YGFloatOptional {
	if lhs.isUndefined() || rhs.isUndefined() {
		return YGFloatOptional{NAN}
	}
	return YGFloatOptional{lhs.unwrap() + rhs.unwrap()}
}

func (lhs YGFloatOptional) subtract(rhs YGFloatOptional) YGFloatOptional {
	if lhs.isUndefined() || rhs.isUndefined() {
		return YGFloatOptional{NAN}
	}
	return YGFloatOptional{lhs.unwrap() - rhs.unwrap()}
}

func (lhs YGFloatOptional) multiply(rhs YGFloatOptional) YGFloatOptional {
	if lhs.isUndefined() || rhs.isUndefined() {
		return YGFloatOptional{NAN}
	}
	return YGFloatOptional{lhs.unwrap() * rhs.unwrap()}
}

func (opt YGFloatOptional) isZero() bool {
	return opt.equal(YGFloatOptional{0})
}

func (lhs YGFloatOptional) divide(rhs YGFloatOptional) YGFloatOptional {
	if lhs.isUndefined() || rhs.isUndefined() || rhs.isZero() {
		return YGFloatOptional{NAN}
	}
	return YGFloatOptional{lhs.unwrap() / rhs.unwrap()}
}

func (lhs YGFloatOptional) equalFloat32(rhs float32) bool {
	return lhs.equal(YGFloatOptional{rhs})
}
