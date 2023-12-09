//go:build darwin && !ios
// +build darwin,!ios

package app

import "gitee.com/millken/goui/driver/darwin"

type nativeApplication struct {
	app darwin.NSApplication
}
