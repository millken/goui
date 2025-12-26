package font

import (
	"fmt"
	"runtime"
	"testing"
)

// TestLoadSystemFonts tests the system font loading functionality
func TestLoadSystemFonts(t *testing.T) {
	db := New()

	// Load system fonts - this should not panic
	db.LoadSystemFonts()

	// Check that some fonts were loaded (at least on most systems)
	// On some CI systems, system fonts might not be available
	if len(db.faces) == 0 {
		t.Log("No system fonts loaded - this might be expected on some systems")
		return
	}

	t.Logf("Loaded %d font faces", len(db.faces))

	// Print some information about loaded fonts for debugging
	count := 0
	for id, face := range db.faces {
		if count >= 5 { // Limit output for readability
			break
		}
		t.Logf("Font ID %d: %s (Families: %v)", id, face.PostScriptName, face.Families)
		count++
	}
}

// TestDefaultFonts tests default font family selection
func TestDefaultFonts(t *testing.T) {
	db := New()

	// Load system fonts first
	db.LoadSystemFonts()

	// Test generic family names
	testCases := []struct {
		generic string
		name    string
	}{
		{"serif", db.familySerif},
		{"sans-serif", db.familySansSerif},
		{"sans serif", db.familySansSerif},
		{"cursive", db.familyCursive},
		{"fantasy", db.familyFantasy},
		{"monospace", db.familyMonospace},
	}

	for _, tc := range testCases {
		family := Family{Generic: tc.generic}
		result := db.FamilyName(family)

		if result == "" {
			t.Errorf("FamilyName for %s returned empty string", tc.generic)
		}

		t.Logf("Generic '%s' resolved to: %s", tc.generic, result)
	}
}

// TestFontQuery tests CSS-like font queries
func TestFontQuery(t *testing.T) {
	db := New()

	// Load system fonts
	db.LoadSystemFonts()

	if len(db.faces) == 0 {
		t.Skip("Skipping query tests - no fonts loaded")
		return
	}

	// Test queries for common font families (using fonts that are actually available)
	testQueries := []struct {
		name    string
		query   Query
		wantID  ID // 0 means we don't care if it matches, just that it doesn't crash
	}{
		{
			name: "Noto Sans query",
			query: Query{
				Families: []Family{{Name: "Noto Sans"}},
			},
		},
		{
			name: "Noto Serif query",
			query: Query{
				Families: []Family{{Name: "Noto Serif"}},
			},
		},
		{
			name: "DejaVu Sans query",
			query: Query{
				Families: []Family{{Name: "DejaVu Sans"}},
			},
		},
		{
			name: "Sans-serif generic",
			query: Query{
				Families: []Family{{Generic: "sans-serif"}},
			},
		},
		{
			name: "Serif generic",
			query: Query{
				Families: []Family{{Generic: "serif"}},
			},
		},
		{
			name: "Monospace generic",
			query: Query{
				Families: []Family{{Generic: "monospace"}},
			},
		},
		{
			name: "Bold weight query",
			query: Query{
				Families: []Family{{Name: "Noto Sans"}},
				Weight:   WeightBold,
			},
		},
		{
			name: "Italic style query",
			query: Query{
				Families: []Family{{Name: "Noto Serif"}},
				Style:    StyleItalic,
			},
		},
		{
			name: "Multiple families fallback",
			query: Query{
				Families: []Family{
					{Name: "NonExistentFont"},
					{Name: "AnotherNonExistentFont"},
					{Name: "Noto Sans"},
				},
			},
		},
	}

	for _, tc := range testQueries {
		t.Run(tc.name, func(t *testing.T) {
			result := db.Query(tc.query)
			t.Logf("Query '%s' result: %d", tc.name, result)

			// If we got a result, verify it exists in the database
			if result != 0 {
				if _, exists := db.faces[result]; !exists {
					t.Errorf("Query returned ID %d which doesn't exist in database", result)
				}
			}
		})
	}
}

// TestOSpecificDefaults tests OS-specific default font settings
func TestOSpecificDefaults(t *testing.T) {
	db := New()

	// Check OS-specific defaults
	switch runtime.GOOS {
	case "darwin":
		if db.familyFantasy != "Papyrus" {
			t.Errorf("On macOS, fantasy font should be 'Papyrus', got '%s'", db.familyFantasy)
		}
	case "windows", "linux":
		if db.familyFantasy != "Impact" {
			t.Errorf("On Windows/Linux, fantasy font should be 'Impact', got '%s'", db.familyFantasy)
		}
	}

	// Check common defaults that should be the same across platforms
	expectedDefaults := map[string]string{
		"familySerif":     "Times New Roman",
		"familySansSerif": "Arial",
		"familyCursive":   "Comic Sans MS",
		"familyMonospace": "Courier New",
	}

	actualDefaults := map[string]string{
		"familySerif":     db.familySerif,
		"familySansSerif": db.familySansSerif,
		"familyCursive":   db.familyCursive,
		"familyMonospace": db.familyMonospace,
	}

	for field, expected := range expectedDefaults {
		if actual := actualDefaults[field]; actual != expected {
			t.Logf("%s: expected '%s', got '%s' (this may be system-dependent)", field, expected, actual)
		}
	}
}

// TestFontProperties tests that loaded fonts have correct properties
func TestFontProperties(t *testing.T) {
	db := New()
	db.LoadSystemFonts()

	if len(db.faces) == 0 {
		t.Skip("Skipping properties tests - no fonts loaded")
		return
	}

	// Check a few fonts for reasonable properties
	count := 0
	for id, face := range db.faces {
		if count >= 10 { // Check first 10 fonts
			break
		}

		// Validate face properties
		if len(face.Families) == 0 {
			t.Errorf("Font ID %d has no family names", id)
		}

		if face.PostScriptName == "" {
			t.Errorf("Font ID %d has no PostScript name", id)
		}

		// Check weight is reasonable
		if int(face.Weight) < 100 || int(face.Weight) > 900 {
			t.Errorf("Font ID %d has invalid weight: %d", id, face.Weight)
		}

		// Check style is valid
		validStyles := []Style{StyleNormal, StyleItalic, StyleOblique}
		styleValid := false
		for _, s := range validStyles {
			if face.Style == s {
				styleValid = true
				break
			}
		}
		if !styleValid {
			t.Errorf("Font ID %d has invalid style: %v", id, face.Style)
		}

		// Check stretch is valid
		if face.Stretch < StretchUltraCondensed || face.Stretch > StretchUltraExpanded {
			t.Errorf("Font ID %d has invalid stretch: %v", id, face.Stretch)
		}

		t.Logf("Font ID %d: %s - Weight: %d, Style: %v, Stretch: %v, Monospaced: %t",
			id, face.PostScriptName, face.Weight, face.Style, face.Stretch, face.Monospaced)

		count++
	}
}

// TestConcurrentFontLoading tests concurrent loading of fonts
func TestConcurrentFontLoading(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Fontconfig-based concurrent loading test is Linux-specific")
		return
	}

	db := New()

	// Test with different worker counts
	originalWorkers := db.loadWorkers

	db.SetLoadWorkers(1)
	db.LoadSystemFonts()
	loadedWith1Worker := len(db.faces)

	// Reset and test with more workers
	db = New()
	db.SetLoadWorkers(4)
	db.LoadSystemFonts()
	loadedWith4Workers := len(db.faces)

	t.Logf("Loaded %d fonts with 1 worker, %d fonts with 4 workers",
		loadedWith1Worker, loadedWith4Workers)

	// Restore original
	db.SetLoadWorkers(originalWorkers)
}

// BenchmarkFontQuery benchmarks font query performance
func BenchmarkFontQuery(b *testing.B) {
	db := New()
	db.LoadSystemFonts()

	if len(db.faces) == 0 {
		b.Skip("No fonts loaded for benchmark")
		return
	}

	queries := []Query{
		{Families: []Family{{Generic: "sans-serif"}}},
		{Families: []Family{{Generic: "serif"}, {Name: "Arial"}}},
		{Families: []Family{{Generic: "monospace"}}, Weight: WeightBold},
		{Families: []Family{{Generic: "sans-serif"}}, Style: StyleItalic},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		_ = db.Query(query)
	}
}

// TestPrintSystemInfo prints system font information for debugging
func TestPrintSystemInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping system info test in short mode")
		return
	}

	db := New()
	db.LoadSystemFonts()

	fmt.Printf("=== System Font Information ===\n")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Number of fonts loaded: %d\n", len(db.faces))

	fmt.Printf("\nDefault font families:\n")
	fmt.Printf("Serif: %s\n", db.familySerif)
	fmt.Printf("Sans-serif: %s\n", db.familySansSerif)
	fmt.Printf("Cursive: %s\n", db.familyCursive)
	fmt.Printf("Fantasy: %s\n", db.familyFantasy)
	fmt.Printf("Monospace: %s\n", db.familyMonospace)

	if len(db.faces) > 0 {
		fmt.Printf("\nSample fonts (first 10):\n")
		count := 0
		for id, face := range db.faces {
			if count >= 10 {
				break
			}
			primaryFamily := ""
			if len(face.Families) > 0 {
				primaryFamily = face.Families[0].Name
			}
			fmt.Printf("ID %d: %s (%s) - Weight:%d Style:%v\n",
				id, primaryFamily, face.PostScriptName, face.Weight, face.Style)
			count++
		}
	}
}