package xmath

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Numeric is a constraint that permits any integer or float type.
type Numeric interface {
	constraints.Float | constraints.Integer
}

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs[T Numeric](x T) T {
	switch any(x).(type) {
	case float32:
		return T(math.Float32frombits(math.Float32bits(float32(x)) &^ (1 << 31)))
	case float64:
		return T(math.Abs(float64(x)))
	default:
		if x < 0 {
			return -x
		}
		return x
	}
}

// Floor returns the greatest integer value less than or equal to x.
func Floor[T Numeric](x T) T {
	switch any(x).(type) {
	case float32, float64:
		return T(math.Floor(float64(x)))
	default:
		return x
	}
}

// Ceil returns the smallest integer value greater than or equal to x.
func Ceil[T Numeric](x T) T {
	switch any(x).(type) {
	case float32, float64:
		return T(math.Ceil(float64(x)))
	default:
		return x
	}
}

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin[T constraints.Float](x T) T {
	return T(math.Sin(float64(x)))
}

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//
//	Cos(±0) = 1
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos[T constraints.Float](x T) T {
	return T(math.Cos(float64(x)))
}

// IsNaN reports whether f is a "not-a-number" value.
func IsNaN[T constraints.Float](f T) bool {
	// Only NaNs satisfy f != f.
	return f != f
}

// EqualWithin returns true if a and b are within the given tolerance of each other.
func EqualWithin[T Numeric](a, b, tolerance T) bool {
	return Abs(a-b) <= tolerance
}
