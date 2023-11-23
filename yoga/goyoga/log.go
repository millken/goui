package goyoga

import (
	"fmt"
	"os"
)

func DefaultLogger(config *YGConfig,
	node *YGNode,
	level YGLogLevel,
	format string,
	args ...interface{}) int {
	switch level {
	case YGLogLevelError, YGLogLevelFatal:
		n, _ := fmt.Fprintf(os.Stderr, format, args...)
		return n
	case YGLogLevelWarn, YGLogLevelInfo, YGLogLevelDebug, YGLogLevelVerbose:
		fallthrough
	default:
		n, _ := fmt.Printf(format, args...)
		return n
	}
}
