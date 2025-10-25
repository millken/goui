package v2

// Canvas 提供高级绘图API（当前为调试最小实现）
type Canvas interface {
	// 最小绘制接口：绘制矩形（带颜色）
	DrawRect(rect Rect, color Color)
}

// NopCanvas 是一个空实现，便于后端快速打通渲染链路
type NopCanvas struct{}

func (NopCanvas) DrawRect(rect Rect, color Color) {}
