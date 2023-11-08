package yoga

// Logger defines logging function
type Logger func(config *Config, node *Node, level LogLevel, format string, args ...interface{}) int

func Absf(v float32) float32 {
	if v < 0 {
		return -v
	} else {
		return v
	}
}
