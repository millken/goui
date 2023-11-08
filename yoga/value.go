package yoga

var (
	ValueUndefined = Value{Undefined, UnitUndefined}
	ValueAuto      = Value{Undefined, UnitAuto}
	ValueZero      = Value{0, UnitPoint}
)

type Value struct {
	Value float32
	Unit  Unit
}

func (v Value) IsUndefined() bool {
	return v.Unit == UnitUndefined
}
