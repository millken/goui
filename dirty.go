package goui

import "github.com/millken/go/thorvg"

// CheckDirty 检查组件树中是否有任何脏组件
func CheckDirty(component Component) bool {
	if !component.IsVisible() {
		return false
	}

	if component.IsDirty() {
		return true
	}

	// 递归检查子组件
	for _, child := range component.GetChildren() {
		if CheckDirty(child) {
			return true
		}
	}

	return false
}

// ClearDirtyRecursive 递归清除组件树中所有组件的脏标记
func ClearDirtyRecursive(component Component) {
	if !component.IsVisible() {
		return
	}

	component.ClearDirty()

	// 递归清除子组件的脏标记
	for _, child := range component.GetChildren() {
		ClearDirtyRecursive(child)
	}
}

// PaintRecursive 递归绘制组件树，但只绘制脏组件及其子树
func PaintRecursive(component Component, canvas *thorvg.Canvas) {
	if !component.IsVisible() {
		return
	}

	// 如果当前组件是脏的，或者其父组件是脏的（重布局情况下），则需要重绘
	shouldPaint := component.IsDirty()

	if shouldPaint {
		component.Paint(canvas)
	}

	// 递归处理子组件
	for _, child := range component.GetChildren() {
		PaintRecursive(child, canvas)
	}
}