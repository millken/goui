package goyoga

import (
	"math"
	"testing"
)

func TestUint32ToFloat(t *testing.T) {

	a := math.Float32frombits(uint32(233))
	if a != float32(233) {
		t.Fatalf("except 233, got %v", a)
	}
}
