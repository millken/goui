package yoga

const (
	experimentalFeatureCount = 1
)

// Config describes a configuration
type Config struct {
	ExperimentalFeatures                          [experimentalFeatureCount + 1]bool
	UseWebDefaults                                bool
	UseLegacyStretchBehaviour                     bool
	ShouldDiffLayoutWithoutLegacyStretchBehaviour bool
	PointScaleFactor                              float32
	Logger                                        Logger
	Context                                       interface{}
	PrintTree                                     bool
}

var (
	defaultConfig = Config{
		ExperimentalFeatures: [experimentalFeatureCount + 1]bool{
			false,
			false,
		},
		UseWebDefaults:   false,
		PointScaleFactor: 1,
		Logger:           DefaultLog,
		Context:          nil,
	}
)

// NewConfig creates new config
func NewConfig() *Config {
	config := &Config{}
	*config = defaultConfig
	return config
}

func (c *Config) experimentalFeatureEnabled(feature ExperimentalFeature) bool {
	return c.ExperimentalFeatures[feature]
}
