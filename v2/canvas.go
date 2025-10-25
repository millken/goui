package v2

// Canvas 提供高级绘图API
type Canvas interface {
	// 绘图状态
	SetColor(color Color)       // 设置绘图颜色
	SetFont(font Font)          // 设置字体
	SetLineWidth(width float64) // 设置线宽

	// 基本图形绘制
	DrawRect(rect Rect)                      // 绘制矩形
	DrawCircle(center Point, radius float64) // 绘制圆形
	DrawLine(start, end Point)               // 绘制直线

	// 文本绘制
	DrawText(text string, position Point) // 绘制文本
	MeasureText(text string) Size         // 测量文本尺寸

	// 图像绘制
	DrawImage(image Image, rect Rect) // 绘制图像

	// 变换操作
	PushTransform()           // 保存当前变换
	PopTransform()            // 恢复变换
	Translate(dx, dy float64) // 平移
	Rotate(angle float64)     // 旋转
	Scale(sx, sy float64)     // 缩放

	// 组合操作
	BeginPath()         // 开始路径
	MoveTo(point Point) // 移动路径点
	LineTo(point Point) // 添加直线段
	ClosePath()         // 闭合路径
	Fill()              // 填充路径
	Stroke()            // 描边路径
}
