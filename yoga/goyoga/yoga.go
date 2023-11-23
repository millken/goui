package goyoga

import "math"

const (
	uvnan = 0x7FC00001
)

var (
	NaN                 = math.Float32frombits(uvnan)
	YGUndefined float32 = NaN
)

// IsNaN reports whether f is an IEEE 754 “not-a-number” value.
func IsNaN(f float32) (is bool) {
	return f != f
}

func If[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
