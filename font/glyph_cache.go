package font

import (
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type GlyphCache struct {
	face  font.Face
	cache map[rune]fixed.Int26_6 // rune -> advance
	mutex sync.RWMutex
}

func NewGlyphCache(face font.Face) *GlyphCache {
	return &GlyphCache{
		face:  face,
		cache: make(map[rune]fixed.Int26_6),
	}
}

// 预加载常用字符（类似 Skia 的做法）
func (gc *GlyphCache) PreloadCommon() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	// ASCII 基础字符集
	for r := rune(0x20); r <= 0x7E; r++ {
		advance, ok := gc.face.GlyphAdvance(r)
		if ok {
			gc.cache[r] = advance
		}
	}

	// 常用标点
	commonPunct := []rune{'，', '。', '！', '？', '；', '：', '“', '”', '‘', '’', '（', '）', '【', '】', '《', '》'}
	for _, r := range commonPunct {
		advance, ok := gc.face.GlyphAdvance(r)
		if ok {
			gc.cache[r] = advance
		}
	}

	// 常用汉字前 100 个（覆盖 40% 使用频率）
	commonChinese := []rune{
		'的', '一', '是', '在', '不', '了', '有', '和', '人', '这',
		'中', '大', '为', '上', '个', '国', '我', '以', '要', '他',
		'时', '来', '用', '们', '生', '到', '作', '地', '于', '出',
		'就', '分', '对', '成', '会', '可', '主', '发', '年', '动',
		'同', '工', '也', '能', '下', '过', '子', '说', '产', '种',
		'面', '而', '方', '后', '多', '定', '行', '学', '法', '所',
		'民', '得', '经', '十', '三', '之', '进', '着', '等', '部',
		'度', '家', '电', '力', '里', '如', '水', '化', '高', '自',
		'二', '理', '起', '小', '物', '现', '实', '加', '量', '都',
		'两', '体', '制', '机', '当', '使', '点', '从', '业', '本',
	}

	for _, r := range commonChinese {
		advance, ok := gc.face.GlyphAdvance(r)
		if ok {
			gc.cache[r] = advance
		}
	}
}

// 获取字形宽度（带缓存）
func (gc *GlyphCache) GetAdvance(r rune) fixed.Int26_6 {
	gc.mutex.RLock()
	advance, exists := gc.cache[r]
	gc.mutex.RUnlock()

	if exists {
		return advance
	}

	// 动态加载并缓存
	advance, ok := gc.face.GlyphAdvance(r)
	if !ok {
		// 使用默认宽度
		advance = fixed.I(500) // 默认 8px 左右
	}

	gc.mutex.Lock()
	gc.cache[r] = advance
	gc.mutex.Unlock()

	return advance
}

// 批量获取（优化性能）
func (gc *GlyphCache) GetAdvances(runes []rune) []fixed.Int26_6 {
	advances := make([]fixed.Int26_6, len(runes))

	gc.mutex.RLock()
	defer gc.mutex.RUnlock()

	for i, r := range runes {
		if advance, exists := gc.cache[r]; exists {
			advances[i] = advance
		} else {
			// 标记为需要动态加载
			advances[i] = -1
		}
	}

	return advances
}
