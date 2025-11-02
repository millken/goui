package font

import (
	"sync"
	"testing"
)

func TestFontManager(t *testing.T) {
	fm := NewFontManager()
	measurer := fm.GetMeasurer("NonExistFont", 16)
	if measurer == nil {
		t.Fatal("Expected default measurer, got nil")
	}
	t.Logf("%v", fm.LoadedFamilies())
	text := "Hello,\n 世界啊啊!"
	metric := measurer.MeasureText(text)
	t.Logf("Text: %s, Width: %.2f, Height: %.2f", text, metric.Width, metric.Height)
}

// BenchmarkGetMeasurer_CacheHit 测试缓存命中的性能
func BenchmarkGetMeasurer_CacheHit(b *testing.B) {
	fm := NewFontManager()

	// 预热缓存
	fm.GetMeasurer("default", 14)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fm.GetMeasurer("default", 14) // 缓存命中
		}
	})
}

// BenchmarkGetMeasurer_Mixed 测试混合缓存命中和未命中的性能
func BenchmarkGetMeasurer_Mixed(b *testing.B) {
	fm := NewFontManager()
	fonts := []string{"font1", "font2", "font3"}
	sizes := []float64{12, 14, 16}

	// 预热一部分缓存
	fm.GetMeasurer("font1", 12)
	fm.GetMeasurer("font2", 14)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			font := fonts[i%len(fonts)]
			size := sizes[i%len(sizes)]
			fm.GetMeasurer(font, size)
			i++
		}
	})
}

// BenchmarkGetMeasurer_Concurrent 测试并发访问的性能
func BenchmarkGetMeasurer_Concurrent(b *testing.B) {
	fm := NewFontManager()

	// 预热
	fm.GetMeasurer("default", 14)

	b.ResetTimer()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/10; j++ {
				fm.GetMeasurer("default", 14).MeasureText("aaa")
			}
		}()
	}
	wg.Wait()
}

// TestGetMeasurer_ConcurrentCorrectness 测试并发访问的正确性
func TestGetMeasurer_ConcurrentCorrectness(t *testing.T) {
	fm := NewFontManager()

	var wg sync.WaitGroup
	results := make([]*TextMeasurer, 100)

	// 并发创建相同的测量器
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = fm.GetMeasurer("concurrent_test", 16)
		}(i)
	}
	wg.Wait()

	// 验证所有结果都是同一个实例
	first := results[0]
	if first == nil {
		t.Fatal("Expected non-nil measurer")
	}

	for i, result := range results {
		if result != first {
			t.Errorf("Result %d is different from first result", i)
		}
	}
}
