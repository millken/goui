package xmath

import "testing"

func TestFloor(t *testing.T) {
	t.Logf("%v", Floor(int(3)))
	if got := Floor(float32(3.3)); got != 3 {
		t.Failed()
	}
	if got := Floor(float64(3.3)); got != 3 {
		t.Failed()
	}
	if got := Floor(int(3)); got != 3 {
		t.Fatal()
	}
	if got := Floor(float32(4.4)); got != 3 {
		t.Failed()
	}
}
