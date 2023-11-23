package goyoga

import "unsafe"

func uint32ToFloat(u uint32) float32 {
	return *(*float32)(unsafe.Pointer(&u))
}

func floatToUint32(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}
