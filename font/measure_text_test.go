package font

import (
	"testing"
)

func TestGlobalMeasureText(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected struct {
			minWidth  float64
			maxWidth  float64
			minHeight float64
			maxHeight float64
		}
	}{
		{
			name: "default size single line",
			text: "Hello World",
			expected: struct {
				minWidth  float64
				maxWidth  float64
				minHeight float64
				maxHeight float64
			}{
				minWidth:  50,
				maxWidth:  150,
				minHeight: 15,
				maxHeight: 25,
			},
		},
		{
			name: "with custom size",
			text: "Hello World",
		},
		{
			name: "empty text",
			text: "",
			expected: struct {
				minWidth  float64
				maxWidth  float64
				minHeight float64
				maxHeight float64
			}{
				minWidth:  0,
				maxWidth:  0,
				minHeight: 0,
				maxHeight: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var metric TextMetrics

			if tc.name == "with custom size" {
				metric = MeasureText(tc.text, WithFontSize(20))
			} else {
				metric = MeasureText(tc.text)
			}

			t.Logf("Text: %q, Width: %.2f, Height: %.2f", tc.text, metric.Width, metric.Height)

			// For empty text test
			if tc.name == "empty text" {
				if metric.Width != 0 || metric.Height != 0 {
					t.Errorf("Expected zero metrics for empty text, got Width: %.2f, Height: %.2f", metric.Width, metric.Height)
				}
				return
			}

			// For non-empty text, check reasonable bounds
			if tc.expected.minWidth > 0 && metric.Width < tc.expected.minWidth {
				t.Errorf("Expected width >= %.2f, got %.2f", tc.expected.minWidth, metric.Width)
			}
			if tc.expected.maxWidth > 0 && metric.Width > tc.expected.maxWidth {
				t.Errorf("Expected width <= %.2f, got %.2f", tc.expected.maxWidth, metric.Width)
			}
			if tc.expected.minHeight > 0 && metric.Height < tc.expected.minHeight {
				t.Errorf("Expected height >= %.2f, got %.2f", tc.expected.minHeight, metric.Height)
			}
			if tc.expected.maxHeight > 0 && metric.Height > tc.expected.maxHeight {
				t.Errorf("Expected height <= %.2f, got %.2f", tc.expected.maxHeight, metric.Height)
			}

			// Ensure positive values for non-empty text
			if metric.Width <= 0 || metric.Height <= 0 {
				t.Errorf("Expected positive metrics for non-empty text, got Width: %.2f, Height: %.2f", metric.Width, metric.Height)
			}
		})
	}
}

func TestMultiLineTextMeasurement(t *testing.T) {
	fm := NewFontManager()
	measurer := fm.GetMeasurer("NonExistFont", 16)
	if measurer == nil {
		t.Fatal("Expected default measurer, got nil")
	}

	testCases := []struct {
		name     string
		text     string
		expected struct {
			minWidth    float64
			maxWidth    float64
			exactHeight float64
		}
	}{
		{
			name: "single line",
			text: "Hello",
			expected: struct {
				minWidth    float64
				maxWidth    float64
				exactHeight float64
			}{
				minWidth:    40, // approximate
				maxWidth:    80,
				exactHeight: 25, // 96 DPI下单行高度
			},
		},
		{
			name: "two lines, short first",
			text: "Hi\nWorld",
			expected: struct {
				minWidth    float64
				maxWidth    float64
				exactHeight float64
			}{
				minWidth:    40,
				maxWidth:    80,
				exactHeight: 50, // 2 lines at 96 DPI
			},
		},
		{
			name: "two lines, long first",
			text: "Hello World\nHi",
			expected: struct {
				minWidth    float64
				maxWidth    float64
				exactHeight float64
			}{
				minWidth:    80,
				maxWidth:    160,
				exactHeight: 50, // 2 lines at 96 DPI
			},
		},
		{
			name: "empty lines",
			text: "\n\n",
			expected: struct {
				minWidth    float64
				maxWidth    float64
				exactHeight float64
			}{
				minWidth:    0,
				maxWidth:    0,
				exactHeight: 75, // 3 lines at 96 DPI
			},
		},
		{
			name: "mixed content",
			text: "Hello,\n 世界",
			expected: struct {
				minWidth    float64
				maxWidth    float64
				exactHeight float64
			}{
				minWidth:    40,
				maxWidth:    100,
				exactHeight: 50, // 2 lines at 96 DPI
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metric := measurer.MeasureText(tc.text)
			t.Logf("Text: %q, Width: %.2f, Height: %.2f", tc.text, metric.Width, metric.Height)

			// Check height
			if metric.Height != tc.expected.exactHeight {
				t.Errorf("Expected height %.2f, got %.2f", tc.expected.exactHeight, metric.Height)
			}

			// Check width range
			if metric.Width < tc.expected.minWidth || metric.Width > tc.expected.maxWidth {
				t.Errorf("Expected width between %.2f and %.2f, got %.2f",
					tc.expected.minWidth, tc.expected.maxWidth, metric.Width)
			}
		})
	}
}
