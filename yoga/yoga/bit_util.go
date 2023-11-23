package yoga

import (
	"math/bits"
)

// std::bitset with one bit for each option defined in YG_ENUM_SEQ_DECL
type EnumBitset struct {
	bits uint
}

func NewEnumBitset() EnumBitset {
	return EnumBitset{bits: 0}
}

func (bitset *EnumBitset) Set(index int) {
	bitset.bits |= 1 << index
}

func (bitset *EnumBitset) Reset(index int) {
	bitset.bits &= ^(1 << index)
}

func (bitset *EnumBitset) Test(index int) bool {
	return (bitset.bits & (1 << index)) != 0
}

func (bitset *EnumBitset) Count() int {
	return bits.OnesCount(bitset.bits)
}

func log2ceilFn(n int) int {
	if n < 1 {
		return 0
	}
	return 1 + log2ceilFn(n/2)
}

func mask(bitWidth, index int) int {
	return ((1 << bitWidth) - 1) << index
}

// The number of bits necessary to represent enums defined with YG_ENUM_SEQ_DECL
func bitWidthFn(Enum interface{}) int {
	count := enumsCount(Enum)
	if count <= 0 {
		panic("Enums must have at least one entry")
	}
	return log2ceilFn(count - 1)
}

func getEnumData(flags int, index int, Enum interface{}) interface{} {
	bitWidth := bitWidthFn(Enum)
	return (flags & mask(bitWidth, index)) >> index
}

func setEnumData(flags *uint32, index int, newValue interface{}, Enum interface{}) {
	bitWidth := bitWidthFn(Enum)
	*flags = (*flags & ^uint32(mask(bitWidth, index))) |
		(uint32(newValue.(int)) << index & uint32(mask(bitWidth, index)))
}

func setEnumData8(flags *uint8, index int, newValue interface{}, Enum interface{}) {
	bitWidth := bitWidthFn(Enum)
	*flags = (*flags & ^uint8(mask(bitWidth, index))) |
		(uint8(newValue.(int)) << index & uint8(mask(bitWidth, index)))
}

func getBooleanData(flags int, index int) bool {
	return (flags>>index)&1 != 0
}

func setBooleanData(flags *uint8, index int, value bool) {
	if value {
		*flags |= 1 << index
	} else {
		*flags &= ^(1 << index)
	}
}

func enumsCount(Enum interface{}) int {
	// implementation details
	return 0
}
