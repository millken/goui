package yoga

import (
	"fmt"
	"os"
)

// DefaultLog is default logging function
func DefaultLog(config *Config, node *Node, level LogLevel, format string,
	args ...interface{}) int {
	switch level {
	case LogLevelError, LogLevelFatal:
		n, _ := fmt.Fprintf(os.Stderr, format, args...)
		return n
	case LogLevelWarn, LogLevelInfo, LogLevelDebug, LogLevelVerbose:
		fallthrough
	default:
		n, _ := fmt.Printf(format, args...)
		return n
	}
}
