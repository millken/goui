package v2

// Canvas 提供较完整的绘图 API，支持变换、状态栈、裁剪和基础图元
type Canvas interface {
	// 状态栈
	Save()
	Restore()

	// 变换
	Translate(tx, ty float32)
	Scale(sx, sy float32)

	// 颜色/样式设置
	SetColor(color Color)
	SetFont(font Font)

	// 基础图元
	DrawRect(rect Rect)
	DrawLine(x0, y0, x1, y1 float32)

	// 图片与文本（后端可选择性支持；这里为占位）
	DrawImage(img Image, dest Rect)
	DrawText(text string, x, y float32)
}

// NopCanvas 是一个空实现，便于后端快速打通渲染链路
type NopCanvas struct{}

func (NopCanvas) Save()                              {}
func (NopCanvas) Restore()                           {}
func (NopCanvas) Translate(tx, ty float64)           {}
func (NopCanvas) Scale(sx, sy float64)               {}
func (NopCanvas) ClipRect(r Rect)                    {}
func (NopCanvas) SetColor(color Color)               {}
func (NopCanvas) SetFont(font Font)                  {}
func (NopCanvas) SetLineWidth(w float64)             {}
func (NopCanvas) DrawRect(rect Rect)                 {}
func (NopCanvas) DrawLine(x0, y0, x1, y1 float64)    {}
func (NopCanvas) DrawImage(img Image, dest Rect)     {}
func (NopCanvas) DrawText(text string, x, y float64) {}
