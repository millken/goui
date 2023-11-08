package yoga

type FloatOptional struct {
	IsUndefined bool
	Value       float32
}

var (
	undefinedFloatOptional = FloatOptional{IsUndefined: true, Value: Undefined}
)

func NewFloatOptional(value float32) FloatOptional {
	if value == Undefined {
		return undefinedFloatOptional
	}
	return FloatOptional{IsUndefined: false, Value: value}
}

func (f FloatOptional) IsDefined() bool {
	return !f.IsUndefined
}

func (f FloatOptional) GetValue() float32 {
	return f.Value
}

func (f FloatOptional) Equals(other FloatOptional) bool {
	return f.IsUndefined == other.IsUndefined && f.Value == other.Value
}

func (f FloatOptional) EqualsFloat(other float32) bool {
	return f.IsDefined() && f.Value == other
}

func (f FloatOptional) SetValue(value float32) {
	if value == Undefined {
		f.IsUndefined = true
		f.Value = Undefined
		return
	}
	f.IsUndefined = false
	f.Value = value
}

// Add
func (f FloatOptional) Add(other FloatOptional) FloatOptional {
	if f.IsUndefined || other.IsUndefined {
		return undefinedFloatOptional
	}
	return NewFloatOptional(f.Value + other.Value)
}

// Mult
func (f FloatOptional) Mult(other FloatOptional) FloatOptional {
	if f.IsUndefined || other.IsUndefined {
		return undefinedFloatOptional
	}
	return NewFloatOptional(f.Value * other.Value)
}
