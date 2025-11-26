// Package renderer defines the rendering interfaces and OpenGL/ThorVG based implementations
// used by GOUI to draw components onto a canvas.
package renderer

import "github.com/millken/thorvg"

type Renderer interface {
	Cleanup() // 清理资源

	// 渲染控制
	// BeginFrame 返回本帧可用于绘制的 Canvas，由具体渲染后端创建/管理
	BeginFrame(width, height int) (*thorvg.Canvas, error)
	EndFrame() error // 结束渲染帧

	// ResetScene 当目标大小或上下文重建时重置渲染场景
	ResetScene(width, height int) error
}
