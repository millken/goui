package font

import (
	"fmt"
	"strconv"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// FontManager 高性能字体管理器（使用sync.Map和复合key）
type FontManager struct {
	measurers   sync.Map // map[string]*TextMeasurer，使用复合key
	parsedFonts sync.Map // map[string]*opentype.Font，缓存解析后的字体对象
}

func NewFontManager() *FontManager {
	fm := &FontManager{}
	if err := fm.LoadFont(OpenSansSans40TTF); err != nil {
		panic(err)
	}
	return fm
}

// GetMeasurer 获取或创建测量器（使用sync.Map的高性能版本）
func (fm *FontManager) GetMeasurer(fontFamily string, fontSize float64) *TextMeasurer {
	var buf [8]byte
	b := strconv.AppendFloat(buf[:0], fontSize, 'f', 1, 64)
	key := fontFamily + "@" + string(b)
	// 快速路径：直接从sync.Map读取（无锁操作）
	if measurer, ok := fm.measurers.Load(key); ok {
		return measurer.(*TextMeasurer)
	}

	// 慢速路径：需要创建新的测量器
	measurer := fm.createMeasurer(fontFamily, fontSize)
	if measurer == nil {
		return nil
	}

	// 原子性存储到sync.Map
	actual, loaded := fm.measurers.LoadOrStore(key, measurer)
	if loaded {
		// 如果已经有其他goroutine存储了，使用已存在的
		return actual.(*TextMeasurer)
	}

	return measurer
}

// createMeasurer 创建新的测量器
func (fm *FontManager) createMeasurer(fontFamily string, fontSize float64) *TextMeasurer {
	// 尝试加载指定字体,回退默认字体
	face, err := fm.loadFontFace(fontFamily, fontSize)
	if err != nil {
		panic(err)
	}
	return NewTextMeasurer(face)
}

func (fm *FontManager) LoadFont(source any) error {
	fontData, err := getSource(source)
	if err != nil {
		return err
	}

	fontFamily, err := GetFontFamily(fontData)
	if err != nil {
		return fmt.Errorf("get font family: %w", err)
	}

	// 解析字体
	parsedFont, err := opentype.Parse(fontData)
	if err != nil {
		return fmt.Errorf("parse font data: %w", err)
	}

	// 存储解析后的字体对象到缓存
	fm.parsedFonts.Store(fontFamily, parsedFont)

	return nil
}

func (fm *FontManager) LoadedFamilies() []string {
	families := make([]string, 0)
	fm.parsedFonts.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			families = append(families, k)
		}
		return true
	})
	return families
}

// 加载字体（简化版）
func (fm *FontManager) loadFontFace(fontFamily string, fontSize float64) (font.Face, error) {
	// 尝试从缓存获取解析后的字体对象
	for _, fontFamily := range []string{fontFamily, defaultFontFamily} {
		parsedFontAny, ok := fm.parsedFonts.Load(fontFamily)
		if ok {
			parsedFont := parsedFontAny.(*opentype.Font)
			face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
				Size:    fontSize,
				DPI:     96, // 使用标准的96 DPI
				Hinting: font.HintingFull,
			})
			if err != nil {
				return nil, fmt.Errorf("create font face for %s: %w", fontFamily, err)
			}
			return face, nil
		}
	}
	// 现在返回一个错误，触发默认字体加载
	return nil, fmt.Errorf("font not found: %s", fontFamily)
}

