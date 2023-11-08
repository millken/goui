package yoga

// This value was chosen based on empiracle data. Even the most complicated
// layouts should not require more than 16 entries to fit within the cache.
const maxCachedResultCount = 16

type Layout struct {
	Position                           [4]float32
	Dimensions                         [2]float32
	Margin                             [6]float32
	Border                             [6]float32
	Padding                            [6]float32
	Direction                          Direction
	HadOverflow                        bool
	computedFlexBasisGeneration        uint32
	computedFlexBasis                  FloatOptional
	generationCount                    uint32
	lastOwnerDirection                 Direction
	nextCachedMeasurementsIndex        uint32
	cachedMeasurements                 [maxCachedResultCount]CachedMeasurement
	measuredDimensions                 [2]float32
	cachedLayout                       CachedMeasurement
	didUseLegacyFlag                   bool
	doesLegacyStretchFlagAffectsLayout bool
}

var (
	DefaultLayout = &Layout{
		Position:                           [4]float32{0, 0, 0, 0},
		Dimensions:                         defaultDimensionValues,
		Margin:                             [6]float32{0, 0, 0, 0, 0, 0},
		Border:                             [6]float32{0, 0, 0, 0, 0, 0},
		Padding:                            [6]float32{0, 0, 0, 0, 0, 0},
		Direction:                          DirectionInherit,
		HadOverflow:                        false,
		computedFlexBasisGeneration:        0,
		computedFlexBasis:                  undefinedFloatOptional,
		generationCount:                    0,
		lastOwnerDirection:                 DirectionInherit,
		nextCachedMeasurementsIndex:        0,
		cachedMeasurements:                 [maxCachedResultCount]CachedMeasurement{},
		measuredDimensions:                 defaultDimensionValues,
		cachedLayout:                       defaultCachedMeasurement,
		didUseLegacyFlag:                   false,
		doesLegacyStretchFlagAffectsLayout: false,
	}
)
