package font

import (
	_ "embed"
	"fmt"
	"os"
	"sync/atomic"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
)

//go:embed "OPPO_Sans_4.0.ttf"
var OpenSansSans40TTF []byte

var (
	defaultFontFamily string  = "OPPO Sans 4.0"
	defaultFontSize   float64 = 14
	defaultFontFace   atomic.Value
	globalFontManager *FontManager = NewFontManager()
)

//	func init() {
//		face, err := loadDefaultFont(14)
//		if err == nil {
//			defaultFontFace.Store(face)
//		}
//	}
func SetDefaultFontFace(face font.Face) {
	defaultFontFace.Store(face)
}

// GetDefaultFontFace 原子获取默认字体Face
func GetDefaultFontFace() font.Face {
	if face := defaultFontFace.Load(); face != nil {
		if f, ok := face.(font.Face); ok {
			return f
		}
	}
	return nil
}

type measureOptions struct {
	fontSize float64
}

type MeasureOption func(*measureOptions)

func WithFontSize(size float64) MeasureOption {
	return func(opts *measureOptions) {
		opts.fontSize = size
	}
}

// MeasureText measures the width of the given text string in pixels.
func MeasureText(text string, opts ...MeasureOption) TextMetrics {
	options := measureOptions{
		fontSize: defaultFontSize,
	}

	for _, opt := range opts {
		opt(&options)
	}

	// 使用全局字体管理器获取测量器
	measurer := globalFontManager.GetMeasurer(defaultFontFamily, options.fontSize)
	if measurer == nil {
		// 如果获取测量器失败，返回零值
		return TextMetrics{}
	}

	return measurer.MeasureText(text)
}

// GetDefaultFontData returns the default font data.
func GetDefaultFontData() []byte {
	return OpenSansSans40TTF
}

func getSource(source any) ([]byte, error) {
	switch v := source.(type) {
	case []byte:
		return v, nil
	case string:
		return os.ReadFile(v)
	}
	return nil, fmt.Errorf("unsupported font source type")
}

func GetFontFamily(source any) (string, error) {
	fontData, err := getSource(source)
	if err != nil {
		return "", err
	}
	// 解析字体
	font, err := sfnt.Parse(fontData)
	if err != nil {
		return "", err
	}

	// 获取字体名称信息
	return font.Name(nil, sfnt.NameIDFamily)
}
