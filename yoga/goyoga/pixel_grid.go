package goyoga

import "math"

func roundValueToPixelGrid(value float64, pointScaleFactor float64, forceCeil bool, forceFloor bool) float32 {
	scaledValue := value * pointScaleFactor
	fractial := math.Mod(scaledValue, 1.0)
	if fractial < 0.00001 {
		fractial++
	}
	if inexactEqual(fractial, 0) {
		scaledValue = scaledValue - fractial
	} else if inexactEqual(fractial, 1) {
		scaledValue = scaledValue - fractial + 1
	} else if forceCeil {
		scaledValue = scaledValue - fractial + 1
	} else if forceFloor {
		scaledValue = scaledValue - fractial
	} else {
		scaledValue = scaledValue - fractial +
			If(!math.IsNaN(fractial) &&
				(fractial > 0.5 || inexactEqual(fractial, 0.5)), 1.0, 0.0)
	}
	return If(math.IsNaN(scaledValue) || math.IsNaN(pointScaleFactor), YGUndefined, float32(scaledValue/pointScaleFactor))
}
