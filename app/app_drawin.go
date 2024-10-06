//go:build darwin && !ios
// +build darwin,!ios

package app

import "github.com/millken/goui/internal/driver/darwin"

type nativeApplication struct {
	app darwin.NSApplication
}
