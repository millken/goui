// SPDX-License-Identifier: GNU Lesser General Public License v2.1

//go:build !android && !ios && !js

package goui

import "github.com/go-gl/glfw/v3.3/glfw"

func glfwEnabled(enabled bool) int {
	if enabled {
		return glfw.True
	}
	return glfw.False
}
