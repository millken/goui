package ui

import (
	"gitee.com/millken/goui/internal/skia"
	"gitee.com/millken/goui/internal/xmath"
)

// PointMode controls how DrawPoints() renders the points passed to it.
type PointMode byte

// Possible values for PointMode.
const (
	PointsPointMode PointMode = iota
	LinesPointMode
	PolygonPointMode
)

// ClipOp holds the type of clip operation to perform.
type ClipOp byte

// Possible values for ClipOp.
const (
	DifferenceClipOp ClipOp = iota
	IntersectClipOp
)

type Canvas struct {
	canvas  skia.Canvas
	surface *surface
}

// Save pushes the current transformation matrix and clip onto a stack and returns the current count. Multiple save
// calls should be balanced by an equal number of calls to Restore().
func (c *Canvas) Save() int {
	return skia.CanvasSave(c.canvas)
}

// Restore removes changes to the transformation matrix and clip since the last call to Save() or SaveWithOpacity().
// Does nothing if the stack is empty.
func (c *Canvas) Restore() {
	skia.CanvasRestore(c.canvas)
}

// Translate the coordinate system.
func (c *Canvas) Translate(dx, dy float32) {
	skia.CanvasTranslate(c.canvas, dx, dy)
}

// Matrix returns the current transform matrix.
func (c *Canvas) Matrix() Matrix {
	return skia.CanvasGetTotalMatrix(c.canvas)
}

// SetMatrix replaces the current matrix with the given matrix.
func (c *Canvas) SetMatrix(matrix Matrix) {
	skia.CanvasSetMatrix(c.canvas, matrix)
}

// Scale the coordinate system.
func (c *Canvas) Scale(x, y float32) {
	skia.CanvasScale(c.canvas, x, y)
}

// Rotate the coordinate system.
func (c *Canvas) Rotate(degrees float32) {
	skia.CanvasRotateRadians(c.canvas, degrees*xmath.DegreesToRadians)
}

// Skew the coordinate system. A positive value of sx skews the drawing right as y-axis values increase; a positive
// value of sy skews the drawing down as x-axis values increase.
func (c *Canvas) Skew(sx, sy float32) {
	skia.CanvasSkew(c.canvas, sx, sy)
}

// ClipRect replaces the clip with the intersection of difference of the current clip and rect.
func (c *Canvas) ClipRect(rect Rect, op ClipOp, antialias bool) {
	skia.CanavasClipRectWithOperation(c.canvas, rect, skia.ClipOp(op), antialias)
}

// ClipPath replaces the clip with the intersection of difference of the current clip and path.
func (c *Canvas) ClipPath(path *Path, op ClipOp, antialias bool) {
	skia.CanavasClipPathWithOperation(c.canvas, path.path, skia.ClipOp(op), antialias)
}

// DrawPaint fills the clip with Paint. Any MaskFilter or PathEffect in the Paint is ignored.
func (c *Canvas) DrawPaint(paint *Paint) {
	skia.CanvasDrawPaint(c.canvas, paint.paint)
}

// DrawRect draws the rectangle with Paint.
func (c *Canvas) DrawRect(rect Rect, paint *Paint) {
	skia.CanvasDrawRect(c.canvas, rect, paint.paint)
}

// DrawRoundedRect draws a rounded rectangle with Paint.
func (c *Canvas) DrawRoundedRect(rect Rect, radiusX, radiusY float32, paint *Paint) {
	skia.CanvasDrawRoundRect(c.canvas, rect, radiusX, radiusY, paint.paint)
}

// DrawArc draws an arc. startAngle and sweepAngle are in degrees. If useCenter is true, this will draw a wedge that
// includes lines from the oval center to the arc end points. If useCenter is false, then just and arc between the end
// points will be drawn.
func (c *Canvas) DrawArc(oval Rect, startAngle, sweepAngle float32, paint *Paint, useCenter bool) {
	skia.CanvasDrawArc(c.canvas, oval, startAngle, sweepAngle, useCenter, paint.paint)
}

// DrawCircle draws the circle with Paint.
func (c *Canvas) DrawCircle(cx, cy, radius float32, paint *Paint) {
	skia.CanvasDrawCircle(c.canvas, cx, cy, radius, paint.paint)
}

// DrawOval draws the oval with Paint.
func (c *Canvas) DrawOval(rect Rect, paint *Paint) {
	skia.CanvasDrawOval(c.canvas, rect, paint.paint)
}

// DrawPath draws the path with Paint.
func (c *Canvas) DrawPath(path *Path, paint *Paint) {
	skia.CanvasDrawPath(c.canvas, path.path, paint.paint)
}

// DrawPoint draws a point.
func (c *Canvas) DrawPoint(x, y float32, paint *Paint) {
	skia.CanvasDrawPoint(c.canvas, x, y, paint.paint)
}

// DrawPoints draws the points using the given mode.
func (c *Canvas) DrawPoints(pts []Point, paint *Paint, mode PointMode) {
	skia.CanvasDrawPoints(c.canvas, skia.PointMode(mode), pts, paint.paint)
}

// DrawLine draws a line.
func (c *Canvas) DrawLine(sx, sy, ex, ey float32, paint *Paint) {
	skia.CanvasDrawLine(c.canvas, sx, sy, ex, ey, paint.paint)
}

// Flush any drawing.
func (c *Canvas) Flush() {
	c.surface.flush(true)
}

// RestoreToCount restores the transformation matrix and clip to the state they where in when the specified count was
// returned from a call to Save() or SaveWithOpacity(). Does nothing if count is greater than the current state stack
// count. Restores the state to the initial values if count is <= 1.
func (c *Canvas) RestoreToCount(count int) {
	skia.CanvasRestoreToCount(c.canvas, count)
}

// DrawImageInRect draws the image into the area specified by the rect, scaling if necessary. paint may be nil.
func (c *Canvas) DrawImageInRect(img *Image, rect Rect, sampling *SamplingOptions, paint *Paint) {
	c.DrawImageRectInRect(img, Rect{Size: img.Size()}, rect, sampling, paint)
}

// DrawImageRectInRect draws a portion of the image into the area specified, scaling if necessary. srcRect should be in
// raw pixel coordinates, not logical coordinates. dstRect should be in logical coordinates. paint may be nil.
func (c *Canvas) DrawImageRectInRect(img *Image, srcRect, dstRect Rect, sampling *SamplingOptions, paint *Paint) {
	skia.CanvasDrawImageRect(c.canvas, img.ref().contextImg(c.surface), srcRect, dstRect, sampling.skSamplingOptions(),
		paint.paintOrNil())
}
