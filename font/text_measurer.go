package font

import (
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 简单的文本测量器
type TextMeasurer struct {
	face       font.Face
	glyphCache *GlyphCache
	mutex      sync.RWMutex
}

// 测量结果
type TextMetrics struct {
	Width  float64
	Height float64
}

// 创建测量器
func NewTextMeasurer(face font.Face) *TextMeasurer {
	measurer := &TextMeasurer{
		face:       face,
		glyphCache: NewGlyphCache(face),
	}

	// 立即预加载常用字符
	measurer.PreloadCommonGlyphs()

	return measurer
}

// 预加载常用字符
func (tm *TextMeasurer) PreloadCommonGlyphs() {
	tm.glyphCache.PreloadCommon()
}

// 测量文本宽度和高度（支持多行）
func (tm *TextMeasurer) MeasureText(text string) TextMetrics {
	if text == "" {
		return TextMetrics{Width: 0, Height: 0}
	}

	// 按换行符分割文本
	lines := strings.Split(text, "\n")

	// 计算高度
	metrics := tm.face.Metrics()
	lineHeight := fixedToFloat(metrics.Height)
	totalHeight := lineHeight * float64(len(lines))

	// 计算每行的宽度，取最大值
	var maxWidth fixed.Int26_6

	for _, line := range lines {
		if line == "" {
			// 空行宽度为0，但仍有高度
			continue
		}

		runes := []rune(line)

		// 批量获取字形宽度
		advances := tm.glyphCache.GetAdvances(runes)

		var lineWidth fixed.Int26_6
		var prevRune rune = -1

		for i, r := range runes {
			var advance fixed.Int26_6

			if advances[i] >= 0 {
				advance = advances[i]
			} else {
				// 动态加载
				advance = tm.glyphCache.GetAdvance(r)
			}

			// 处理字距调整（对于拉丁字母等）
			if prevRune >= 0 {
				kern := tm.face.Kern(prevRune, r)
				lineWidth += kern
			}

			lineWidth += advance
			prevRune = r
		}

		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return TextMetrics{
		Width:  fixedToFloat(maxWidth),
		Height: totalHeight,
	}
}

// 工具函数：fixed 转 float
func fixedToFloat(value fixed.Int26_6) float64 {
	return float64(value) / 64.0
}
