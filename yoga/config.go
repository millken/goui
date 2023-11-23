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

	errata Errata
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
		errata:           ErrataNone,
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

// cloneNode
func (c *Config) cloneNode(node *Node) *Node {
	newNode := &Node{}
	*newNode = *node
	newNode.Owner = nil
	newNode.Children = nil
	newNode.NextChild = nil
	newNode.Config = c
	return newNode
}

// SetErrata
func (c *Config) SetErrata(errata Errata) {
	c.errata = errata
}

// AddErrata
func (c *Config) AddErrata(errata Errata) {
	c.errata |= errata
}

// RemoveErrata
func (c *Config) RemoveErrata(errata Errata) {
	c.errata &= ^errata
}

// GetErrata
func (c *Config) GetErrata() Errata {
	return c.errata
}

// HasErrata
func (c *Config) HasErrata(errata Errata) bool {
	return c.errata&errata != 0
}
