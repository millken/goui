package yoga

import "unsafe"

type YGSize struct {
	width  float32
	height float32
}
type YGCloneNodeFunc func(oldNode, owner *YGNode, childIndex int) *YGNode
type YGMeasureFunc func(node *YGNode, width float32, widthMode YGMeasureMode, height float32, heightMode YGMeasureMode) YGSize

type YGBaselineFunc func(node *YGNode, width, height float32) float32

type YGDirtiedFunc func(node *YGNode)

type YGPrintFunc func(node *YGNode)

type YGNodeCleanupFunc func(node *YGNode)

type YGLogger func(config *YGConfig, node *YGNode, level YGLogLevel, format string, args ...interface{}) int
type LogWithContextFn func(config *YGConfig, node *YGNode, level YGLogLevel, context unsafe.Pointer, format string, args ...interface{}) int
type MeasureWithContextFn func(node *YGNode, width float32, widthMode YGMeasureMode, height float32, heightMode YGMeasureMode, context unsafe.Pointer) YGSize

type BaselineWithContextFn func(node *YGNode, width, height float32, context unsafe.Pointer) float32

type PrintWithContextFn func(node *YGNode, context unsafe.Pointer)
type CloneWithContextFn func(node, owner *YGNode, childIndex int, cloneContext unsafe.Pointer) *YGNode
type YGLogWithContextFn func(config *YGConfig, node *YGNode, level YGLogLevel, context unsafe.Pointer, format string, args ...interface{}) int
type YGCloneNodeWithContextFn func(node *YGNode, owner *YGNode, childIndex int, cloneContext unsafe.Pointer) *YGNode
