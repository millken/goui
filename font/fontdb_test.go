package font

import (
	"os"
	"testing"
)

// TestNewDatabase tests the New function
func TestNewDatabase(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("New() returned nil")
	}

	if db.faces == nil {
		t.Error("New() did not initialize faces map")
	}

	if db.loadWorkers <= 0 {
		t.Error("New() did not set loadWorkers correctly")
	}

	if !db.useMmap {
		t.Error("New() did not enable mmap by default")
	}
}

// TestConfigure tests configuration functions
func TestConfigure(t *testing.T) {
	db := New()

	// Test SetLoadWorkers
	db.SetLoadWorkers(5)
	if db.loadWorkers != 5 {
		t.Errorf("SetLoadWorkers(5) failed, got %d", db.loadWorkers)
	}

	db.SetLoadWorkers(-1) // should not change
	if db.loadWorkers != 5 {
		t.Error("SetLoadWorkers(-1) should not change the value")
	}

	// Test EnableMMap
	db.EnableMMap(false)
	if db.useMmap {
		t.Error("EnableMMap(false) failed")
	}

	db.EnableMMap(true)
	if !db.useMmap {
		t.Error("EnableMMap(true) failed")
	}
}

// TestFamilySetters tests family setter functions
func TestFamilySetters(t *testing.T) {
	db := New()

	db.SetSerifFamily("Test Serif")
	if db.familySerif != "Test Serif" {
		t.Errorf("SetSerifFamily failed, got %s", db.familySerif)
	}

	db.SetSansSerifFamily("Test Sans")
	if db.familySansSerif != "Test Sans" {
		t.Errorf("SetSansSerifFamily failed, got %s", db.familySansSerif)
	}

	db.SetCursiveFamily("Test Cursive")
	if db.familyCursive != "Test Cursive" {
		t.Errorf("SetCursiveFamily failed, got %s", db.familyCursive)
	}

	db.SetFantasyFamily("Test Fantasy")
	if db.familyFantasy != "Test Fantasy" {
		t.Errorf("SetFantasyFamily failed, got %s", db.familyFantasy)
	}

	db.SetMonospaceFamily("Test Mono")
	if db.familyMonospace != "Test Mono" {
		t.Errorf("SetMonospaceFamily failed, got %s", db.familyMonospace)
	}
}

// TestFamilySettersNil tests nil database handling
func TestFamilySettersNil(t *testing.T) {
	var db *Database = nil

	// These should not panic
	db.SetSerifFamily("Test Serif")
	db.SetSansSerifFamily("Test Sans")
	db.SetCursiveFamily("Test Cursive")
	db.SetFantasyFamily("Test Fantasy")
	db.SetMonospaceFamily("Test Mono")
}

// TestFamilyName tests the FamilyName function
func TestFamilyName(t *testing.T) {
	db := New()

	// Test with specific family name
	family := Family{Name: "Arial"}
	result := db.FamilyName(family)
	if result != "Arial" {
		t.Errorf("FamilyName with specific name failed, got %s", result)
	}

	// Test with generic family names
	family = Family{Generic: "serif"}
	result = db.FamilyName(family)
	if result != db.familySerif {
		t.Errorf("FamilyName for serif failed, got %s", result)
	}

	family = Family{Generic: "sans-serif"}
	result = db.FamilyName(family)
	if result != db.familySansSerif {
		t.Errorf("FamilyName for sans-serif failed, got %s", result)
	}

	family = Family{Generic: "SANS SERIF"}
	result = db.FamilyName(family)
	if result != db.familySansSerif {
		t.Errorf("FamilyName for SANS SERIF (case insensitive) failed, got %s", result)
	}

	family = Family{Generic: "cursive"}
	result = db.FamilyName(family)
	if result != db.familyCursive {
		t.Errorf("FamilyName for cursive failed, got %s", result)
	}

	family = Family{Generic: "fantasy"}
	result = db.FamilyName(family)
	if result != db.familyFantasy {
		t.Errorf("FamilyName for fantasy failed, got %s", result)
	}

	family = Family{Generic: "monospace"}
	result = db.FamilyName(family)
	if result != db.familyMonospace {
		t.Errorf("FamilyName for monospace failed, got %s", result)
	}

	// Test unknown generic family
	family = Family{Generic: "unknown"}
	result = db.FamilyName(family)
	if result != "unknown" {
		t.Errorf("FamilyName for unknown generic family failed, got %s", result)
	}
}

// TestFamilyNameNil tests nil database handling
func TestFamilyNameNil(t *testing.T) {
	var db *Database = nil
	family := Family{Name: "Arial"}
	result := db.FamilyName(family)
	if result != "" {
		t.Errorf("FamilyName with nil database should return empty string, got %s", result)
	}
}

// TestStretchToNumber tests the Stretch.ToNumber method
func TestStretchToNumber(t *testing.T) {
	tests := []struct {
		stretch  Stretch
		expected int
	}{
		{StretchUltraCondensed, 1},
		{StretchExtraCondensed, 2},
		{StretchCondensed, 3},
		{StretchSemiCondensed, 4},
		{StretchNormal, 5},
		{StretchSemiExpanded, 6},
		{StretchExpanded, 7},
		{StretchExtraExpanded, 8},
		{StretchUltraExpanded, 9},
	}

	for _, test := range tests {
		result := test.stretch.ToNumber()
		if result != test.expected {
			t.Errorf("Stretch.ToNumber() for %v failed, expected %d, got %d", test.stretch, test.expected, result)
		}
	}
}

// TestReadFileOrMMap tests the readFileOrMMap function
func TestReadFileOrMMap(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfont*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some test data
	testData := "test font data"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tmpFile.Close()

	// Test with tryMmap = true
	data, _, err := readFileOrMMap(tmpFile.Name(), true)
	if err != nil {
		t.Fatal("readFileOrMMap failed:", err)
	}
	if len(data) == 0 {
		t.Error("readFileOrMMap returned empty data")
	}
	if string(data) != testData {
		t.Errorf("readFileOrMMap returned wrong data, expected %s, got %s", testData, string(data))
	}

	// Test with tryMmap = false
	data2, isMmaped2, err := readFileOrMMap(tmpFile.Name(), false)
	if err != nil {
		t.Fatal("readFileOrMMap failed:", err)
	}
	if isMmaped2 {
		t.Error("readFileOrMMap with tryMmap=false should not return mmaped data")
	}
	if string(data2) != testData {
		t.Errorf("readFileOrMMap returned wrong data, expected %s, got %s", testData, string(data2))
	}

	// Test with empty path (should return error)
	_, _, err = readFileOrMMap("", true)
	if err == nil {
		t.Error("readFileOrMMap with empty path should return error")
	}

	// Test with non-existent file
	_, _, err = readFileOrMMap("/non/existent/file.ttf", true)
	if err == nil {
		t.Error("readFileOrMMap with non-existent file should return error")
	}
}

// TestLoadFontFileNil tests nil database handling
func TestLoadFontFileNil(t *testing.T) {
	var db *Database = nil

	// These should return error without panicking
	_, err := db.LoadFontFile("/path/to/font.ttf")
	if err == nil {
		t.Error("LoadFontFile with nil database should return error")
	}

	_, err = db.LoadFontFileWithOptions("/path/to/font.ttf", true)
	if err == nil {
		t.Error("LoadFontFileWithOptions with nil database should return error")
	}

	_, err = db.LoadFontFileWithOptions("", true)
	if err == nil {
		t.Error("LoadFontFileWithOptions with empty path should return error")
	}
}

// TestSystemFontDirs tests the systemFontDirs function
func TestSystemFontDirs(t *testing.T) {
	dirs := systemFontDirs()
	if len(dirs) == 0 {
		t.Error("systemFontDirs returned empty slice")
	}

	// Only check that the function doesn't panic and returns a reasonable list
	// The actual directories vary by operating system
}

// TestPushFaceInfo tests the PushFaceInfo function
func TestPushFaceInfo(t *testing.T) {
	db := New()

	fi := &FaceInfo{
		Source: Source{FilePath: "test.ttf"},
		Index:  0,
		Families: []struct {
			Name     string
			Language Language
		}{{Name: "TestFont", Language: "en"}},
		PostScriptName: "TestFont",
		Style:          StyleNormal,
		Weight:         WeightNormal,
		Stretch:        StretchNormal,
		Monospaced:     false,
	}

	id := db.PushFaceInfo(fi)
	if id == 0 {
		t.Error("PushFaceInfo returned ID 0")
	}

	if fi.ID != id {
		t.Error("PushFaceInfo did not set FaceInfo ID correctly")
	}

	// Check that the face is actually stored
	retrieved := db.faces[id]
	if retrieved != fi {
		t.Error("PushFaceInfo did not store FaceInfo correctly")
	}
}

// TestQuery tests basic Query functionality
func TestQuery(t *testing.T) {
	db := New()

	// Empty query should return 0
	query := Query{}
	result := db.Query(query)
	if result != 0 {
		t.Errorf("Empty query should return 0, got %d", result)
	}

	// Query with non-existent family should return 0
	query = Query{
		Families: []Family{{Name: "NonExistentFont"}},
	}
	result = db.Query(query)
	if result != 0 {
		t.Errorf("Query for non-existent font should return 0, got %d", result)
	}

	// Add a test face
	fi := &FaceInfo{
		Source: Source{FilePath: "test.ttf"},
		Index:  0,
		Families: []struct {
			Name     string
			Language Language
		}{{Name: "TestFont", Language: "en"}},
		PostScriptName: "TestFont",
		Style:          StyleNormal,
		Weight:         WeightNormal,
		Stretch:        StretchNormal,
		Monospaced:     false,
	}
	id := db.PushFaceInfo(fi)

	// Query for the added font
	query = Query{
		Families: []Family{{Name: "TestFont"}},
	}
	result = db.Query(query)
	if result != id {
		t.Errorf("Query for added font should return %d, got %d", id, result)
	}
}