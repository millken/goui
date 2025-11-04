package renderer

import "github.com/millken/go/thorvg"

type Renderer interface {
	Cleanup() // 清理资源

	// 渲染控制
	// BeginFrame 返回本帧可用于绘制的 Canvas，由具体渲染后端创建/管理
	BeginFrame(width, height int) (*thorvg.Canvas, error)
	EndFrame() // 结束渲染帧
	Present()  // 呈现到屏幕

}
