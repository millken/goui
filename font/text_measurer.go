package font

import (
	"sync"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 简单的文本测量器
type TextMeasurer struct {
	face       font.Face
	glyphCache *GlyphCache
	mutex      sync.RWMutex

	// 优化：ASCII 快路径（预热一次）
	asciiAdv [128]fixed.Int26_6
	// 是否启用字距（默认关闭以提升性能）
	Kerning bool
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
		Kerning:    true,
	}

	// 立即预加载常用字符
	measurer.PreloadCommonGlyphs()
	// 预热 ASCII 宽度
	for i := 0; i < 128; i++ {
		measurer.asciiAdv[i] = measurer.glyphCache.GetAdvance(rune(i))
	}
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
	metrics := tm.face.Metrics()
	lineHeight := fixedToFloat(metrics.Height)

	var (
		maxWidth  fixed.Int26_6
		lineWidth fixed.Int26_6
		prevRune  rune = -1
		lines          = 1
	)

	for i := 0; i < len(text); {
		b := text[i]
		// ASCII 快路径
		if b < 0x80 {
			i++
			if b == '\n' {
				if lineWidth > maxWidth {
					maxWidth = lineWidth
				}
				lines++
				lineWidth = 0
				prevRune = -1
				continue
			}
			adv := tm.asciiAdv[b]
			if tm.Kerning && prevRune >= 0 {
				lineWidth += tm.face.Kern(prevRune, rune(b))
			}
			lineWidth += adv
			prevRune = rune(b)
			continue
		}

		// 非 ASCII
		r, sz := utf8.DecodeRuneInString(text[i:])
		i += sz
		if r == '\n' {
			if lineWidth > maxWidth {
				maxWidth = lineWidth
			}
			lines++
			lineWidth = 0
			prevRune = -1
			continue
		}
		adv := tm.glyphCache.GetAdvance(r)
		if tm.Kerning && prevRune >= 0 {
			lineWidth += tm.face.Kern(prevRune, r)
		}
		lineWidth += adv
		prevRune = r
	}

	if lineWidth > maxWidth {
		maxWidth = lineWidth
	}

	return TextMetrics{
		Width:  fixedToFloat(maxWidth),
		Height: lineHeight * float64(lines),
	}
}

// 工具函数：fixed 转 float
func fixedToFloat(value fixed.Int26_6) float64 {
	return float64(value) / 64.0
}
