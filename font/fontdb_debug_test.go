package font

import (
	"fmt"
	"strings"
	"testing"
)

// TestDebugFontQuery debugging the font query issues
func TestDebugFontQuery(t *testing.T) {
	db := New()
	db.LoadSystemFonts()

	if len(db.faces) == 0 {
		t.Skip("No fonts loaded")
		return
	}

	// Let's search for fonts that might match common names
	searchNames := []string{
		"Noto Sans",
		"Noto Serif",
		"Arial",
		"Times",
		"Courier",
		"DejaVu",
		"Liberation",
	}

	fmt.Printf("Searching for fonts with common family names...\n")
	for _, searchName := range searchNames {
		matches := []ID{}
		for id, face := range db.faces {
			for _, family := range face.Families {
				if strings.Contains(strings.ToLower(family.Name), strings.ToLower(searchName)) {
					matches = append(matches, id)
					break // found match in this face
				}
			}
		}
		if len(matches) > 0 {
			fmt.Printf("Found %d fonts matching '%s'\n", len(matches), searchName)
			// Show first few matches
			for i, id := range matches {
				if i >= 3 { // limit output
					break
				}
				face := db.faces[id]
				fmt.Printf("  ID %d: %s (Families: %v)\n", id, face.PostScriptName, face.Families)
			}
		}
	}

	// Test query with actual family names from loaded fonts
	fmt.Printf("\nTesting queries with actual loaded family names...\n")
	count := 0
	for _, face := range db.faces {
		if count >= 5 { // test first 5 faces
			break
		}
		if len(face.Families) > 0 {
			familyName := face.Families[0].Name
			query := Query{Families: []Family{{Name: familyName}}}
			result := db.Query(query)

			if result == 0 {
				fmt.Printf("WARNING: Query for exact family name '%s' returned 0 (ID should be %d)\n",
					familyName, face.ID)
			} else {
				fmt.Printf("SUCCESS: Query for '%s' returned ID %d\n", familyName, result)
			}
		}
		count++
	}

	// Test generic queries with actual family names
	fmt.Printf("\nTesting generic queries...\n")
	genericQueries := []Query{
		{Families: []Family{{Name: "Noto Sans"}}},
		{Families: []Family{{Name: "Noto Serif"}}},
		{Families: []Family{{Name: "Noto Sans Mono"}}},
		{Families: []Family{{Name: "DejaVu Sans"}}},
	}

	for i, query := range genericQueries {
		result := db.Query(query)
		fmt.Printf("Generic query %d result: %d\n", i+1, result)

		// Let's also check if these families exist by manual search
		searchName := query.Families[0].Name
		found := false
		for _, face := range db.faces {
			for _, family := range face.Families {
				if family.Name == searchName {
					found = true
					fmt.Printf("  Found '%s' manually in face ID %d\n", searchName, face.ID)
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Printf("  '%s' not found in manual search\n", searchName)
		}
	}
}

// TestDatabaseInternals tests internal database structure
func TestDatabaseInternals(t *testing.T) {
	db := New()
	db.LoadSystemFonts()

	if len(db.faces) == 0 {
		t.Skip("No fonts loaded")
		return
	}

	// Check if any faces have proper family structure
	facesWithFamilies := 0
	facesWithPostScript := 0

	for _, face := range db.faces {
		if len(face.Families) > 0 {
			facesWithFamilies++
		}
		if face.PostScriptName != "" {
			facesWithPostScript++
		}
	}

	t.Logf("Database statistics:")
	t.Logf("Total faces: %d", len(db.faces))
	t.Logf("Faces with families: %d", facesWithFamilies)
	t.Logf("Faces with PostScript names: %d", facesWithPostScript)

	// Show sample family structures
	t.Logf("\nSample family structures:")
	count := 0
	for _, face := range db.faces {
		if count >= 3 {
			break
		}
		if len(face.Families) > 0 {
			t.Logf("Face ID %d:", face.ID)
			for i, family := range face.Families {
				t.Logf("  Family %d: Name='%s', Language='%s'", i+1, family.Name, family.Language)
			}
		}
		count++
	}
}