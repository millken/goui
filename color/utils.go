package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// HSBAToRGB converts HSB components to equivalent RGB values.
// The saturation and brightness components should be floating-point values between zero and one.
// The hue component can be any floating-point number.
func HSBAToRGB(hue, saturation, brightness, alpha float32) int32 {
	var r, g, b int32

	if saturation == 0 {
		r = int32(brightness*255.0 + 0.5)
		g = r
		b = r
	} else {
		h := (hue - float32(math.Floor(float64(hue)))) * 6.0
		f := h - float32(math.Floor(float64(h)))
		p := brightness * (1.0 - saturation)
		q := brightness * (1.0 - saturation*f)
		t := brightness * (1.0 - saturation*(1.0-f))

		switch int(h) {
		case 0:
			r = int32(brightness*255.0 + 0.5)
			g = int32(t*255.0 + 0.5)
			b = int32(p*255.0 + 0.5)
		case 1:
			r = int32(q*255.0 + 0.5)
			g = int32(brightness*255.0 + 0.5)
			b = int32(p*255.0 + 0.5)
		case 2:
			r = int32(p*255.0 + 0.5)
			g = int32(brightness*255.0 + 0.5)
			b = int32(t*255.0 + 0.5)
		case 3:
			r = int32(p*255.0 + 0.5)
			g = int32(q*255.0 + 0.5)
			b = int32(brightness*255.0 + 0.5)
		case 4:
			r = int32(t*255.0 + 0.5)
			g = int32(p*255.0 + 0.5)
			b = int32(brightness*255.0 + 0.5)
		case 5:
			r = int32(brightness*255.0 + 0.5)
			g = int32(p*255.0 + 0.5)
			b = int32(q*255.0 + 0.5)
		default:
			r = 0
			g = 0
			b = 0
		}
	}

	return (int32(alpha*255.0) << 24) | (r << 16) | (g << 8) | b
}

// ParseColor 解析常见颜色字符串为 RGBA。
// 支持格式：
//   - #RRGGBB
//   - #RRGGBBAA
//   - rgb(r,g,b)
//   - rgba(r,g,b,a) 其中 a 支持 0..1 或 0..255
func ParseColor(s string) (r, g, b, a uint8, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 0, 0, false
	}
	if strings.HasPrefix(s, "#") {
		hex := strings.TrimPrefix(s, "#")
		switch len(hex) {
		case 6:
			if v, err := strconv.ParseUint(hex, 16, 32); err == nil {
				r = uint8(v >> 16)
				g = uint8((v >> 8) & 0xFF)
				b = uint8(v & 0xFF)
				a = 255
				return r, g, b, a, true
			}
		case 8:
			if v, err := strconv.ParseUint(hex, 16, 32); err == nil {
				r = uint8(v >> 24)
				g = uint8((v >> 16) & 0xFF)
				b = uint8((v >> 8) & 0xFF)
				a = uint8(v & 0xFF)
				return r, g, b, a, true
			}
		}
		return 0, 0, 0, 0, false
	}
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "rgb(") && strings.HasSuffix(lower, ")") {
		body := strings.TrimSuffix(strings.TrimPrefix(lower, "rgb("), ")")
		parts := splitCSV(body)
		if len(parts) != 3 {
			return 0, 0, 0, 0, false
		}
		r8, ok1 := parseByte(parts[0])
		g8, ok2 := parseByte(parts[1])
		b8, ok3 := parseByte(parts[2])
		if !ok1 || !ok2 || !ok3 {
			return 0, 0, 0, 0, false
		}
		return r8, g8, b8, 255, true
	}
	if strings.HasPrefix(lower, "rgba(") && strings.HasSuffix(lower, ")") {
		body := strings.TrimSuffix(strings.TrimPrefix(lower, "rgba("), ")")
		parts := splitCSV(body)
		if len(parts) != 4 {
			return 0, 0, 0, 0, false
		}
		r8, ok1 := parseByte(parts[0])
		g8, ok2 := parseByte(parts[1])
		b8, ok3 := parseByte(parts[2])
		if !ok1 || !ok2 || !ok3 {
			return 0, 0, 0, 0, false
		}
		r, g, b = r8, g8, b8
		a8, ok2 := parseAlpha(parts[3])
		if !ok2 {
			return 0, 0, 0, 0, false
		}
		return r, g, b, a8, true
	}
	return 0, 0, 0, 0, false
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func parseByte(s string) (uint8, bool) {
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 || v > 255 {
		return 0, false
	}
	return uint8(v), true
}

func parseAlpha(s string) (uint8, bool) {
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, false
		}
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		return uint8(f * 255), true
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return 0, false
	}
	if v > 255 {
		v = 255
	}
	return uint8(v), true
}

// NOTE: parseRGB currently returns same component for r,g,b. For correctness,
// we provide separate integer parsing; keeping simple util, prefer #RRGGBB in practice.
func ParseRGB(rStr, gStr, bStr string) (uint8, uint8, uint8, error) {
	r, err := strconv.Atoi(rStr)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.Atoi(gStr)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.Atoi(bStr)
	if err != nil {
		return 0, 0, 0, err
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("rgb out of range")
	}
	return uint8(r), uint8(g), uint8(b), nil
}
