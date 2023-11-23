package yoga

const (
	// This value was chosen based on empirical data:
	// 98% of analyzed layouts require less than 8 entries.
	YG_MAX_CACHED_RESULT_COUNT = 8
)

type YGLayout struct {
	Position                    [4]float32
	Dimensions                  [2]float32
	Margin                      [4]float32
	Border                      [4]float32
	Padding                     [4]float32
	ComputedFlexBasisGeneration uint32
	ComputedFlexBasis           float32
	GenerationCount             uint32
	LastOwnerDirection          YGDirection
	NextCachedMeasurementsIndex uint32
	CachedMeasurements          [YG_MAX_CACHED_RESULT_COUNT]YGCachedMeasurement
	MeasuredDimensions          [2]float32
	CachedLayout                YGCachedMeasurement
	Flags                       uint8
}

type YGCachedMeasurement struct {
	AvailableWidth    float32
	AvailableHeight   float32
	WidthMeasureMode  YGMeasureMode
	HeightMeasureMode YGMeasureMode
	ComputedWidth     float32
	ComputedHeight    float32
}

func (l *YGLayout) Direction() YGDirection {
	return YGDirection((l.Flags >> 0) & 3)
}

func (l *YGLayout) SetDirection(direction YGDirection) {
	l.Flags = (l.Flags &^ (3 << 0)) | (uint8(direction) << 0)
}

func (l *YGLayout) HadOverflow() bool {
	return (l.Flags>>2)&1 != 0
}

func (l *YGLayout) SetHadOverflow(hadOverflow bool) {
	if hadOverflow {
		l.Flags |= (1 << 2)
	} else {
		l.Flags &= ^(1 << 2)
	}
}

func (l *YGLayout) Equals(layout YGLayout) bool {
	// Implement your own equality comparison logic here
	return false
}
