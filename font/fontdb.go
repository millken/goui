package font

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf16"

	mmap "github.com/edsrzf/mmap-go"
)

// ID is a unique per-database face ID.
type ID uint64

var idCounter uint64 = 1

func nextID() ID {
	return ID(atomic.AddUint64(&idCounter, 1))
}

// Language is a simplified placeholder (e.g. "en", "en-US").
type Language string

// Stretch simplified enumeration for font-stretch.
// We map values to 0..8 internally and provide ToNumber() -> 1..9 to match ttf_parser::Width::to_number.
type Stretch int

const (
	StretchUltraCondensed Stretch = iota
	StretchExtraCondensed
	StretchCondensed
	StretchSemiCondensed
	StretchNormal
	StretchSemiExpanded
	StretchExpanded
	StretchExtraExpanded
	StretchUltraExpanded
)

func (s Stretch) ToNumber() int {
	// map to 1..9 (narrower => smaller)
	return int(s) + 1
}

// Weight values
type Weight uint16

const (
	WeightThin       Weight = 100
	WeightExtraLight        = 200
	WeightLight             = 300
	WeightNormal            = 400
	WeightMedium            = 500
	WeightSemiBold          = 600
	WeightBold              = 700
	WeightExtraBold         = 800
	WeightBlack             = 900
)

// Style holds Normal/Italic/Oblique
type Style int

const (
	StyleNormal Style = iota
	StyleItalic
	StyleOblique
)

// Family: either Name (specific) or Generic (serif, sans-serif, ...).
type Family struct {
	Name    string
	Generic string
}

// Source for font data.
type Source struct {
	FilePath string // if on-disk
	Data     []byte // raw bytes (may be an mmap-backed slice)
	// Note: if Data is mmap-backed, it should not be unmapped while the database is using it.
}

// FaceInfo describes a single face.
type FaceInfo struct {
	ID       ID
	Source   Source
	Index    uint32
	Families []struct {
		Name     string
		Language Language
	}
	PostScriptName string
	Style          Style
	Weight         Weight
	Stretch        Stretch
	Monospaced     bool
}

// Database holds faces and generic family mapping.
type Database struct {
	mu sync.RWMutex

	faces map[ID]*FaceInfo

	familySerif     string
	familySansSerif string
	familyCursive   string
	familyFantasy   string
	familyMonospace string

	// concurrency controls
	loadWorkers int
	// whether to attempt mmap optimization for file reads
	useMmap bool
}

// New returns an initialized Database.
func New() *Database {
	fantasy := "Impact"
	if runtime.GOOS == "darwin" {
		fantasy = "Papyrus"
	}
	// default workers: number of CPUs * 2
	w := runtime.NumCPU() * 2
	return &Database{
		faces:           make(map[ID]*FaceInfo),
		familySerif:     "Times New Roman",
		familySansSerif: "Arial",
		familyCursive:   "Comic Sans MS",
		familyFantasy:   fantasy,
		familyMonospace: "Courier New",
		loadWorkers:     w,
		useMmap:         true, // opt-in by default
	}
}

// Configure concurrency and memmap behavior.
func (db *Database) SetLoadWorkers(n int) {
	if n <= 0 {
		return
	}
	db.loadWorkers = n
}
func (db *Database) EnableMMap(enable bool) { db.useMmap = enable }

// Setters for generic families.
func (db *Database) SetSerifFamily(name string)     {
	if db == nil { return }
	db.familySerif = name
}
func (db *Database) SetSansSerifFamily(name string) {
	if db == nil { return }
	db.familySansSerif = name
}
func (db *Database) SetCursiveFamily(name string)   {
	if db == nil { return }
	db.familyCursive = name
}
func (db *Database) SetFantasyFamily(name string)   {
	if db == nil { return }
	db.familyFantasy = name
}
func (db *Database) SetMonospaceFamily(name string) {
	if db == nil { return }
	db.familyMonospace = name
}
func (db *Database) FamilyName(f Family) string {
	if db == nil {
		return ""
	}
	if f.Name != "" {
		return f.Name
	}
	switch strings.ToLower(f.Generic) {
	case "serif":
		return db.familySerif
	case "sans-serif", "sans serif":
		return db.familySansSerif
	case "cursive":
		return db.familyCursive
	case "fantasy":
		return db.familyFantasy
	case "monospace":
		return db.familyMonospace
	default:
		return f.Generic
	}
}

// LoadSystemFonts scans known system directories and loads fonts.
// On Linux this will try to use fontconfig (fc-list + fc-match) to find fonts and
// to resolve generic family aliases. If fontconfig is not available, it falls back
// to scanning well-known system font directories.
func (db *Database) LoadSystemFonts() {
	if db == nil {
		return
	}
	if runtime.GOOS == "linux" {
		if db.loadUsingFontconfig() {
			return
		}
		// fallback to dir scanning
	}

	dirs := systemFontDirs()
	seen := map[string]struct{}{}
	for _, d := range dirs {
		_ = db.LoadFontsDirConcurrent(d, seen)
	}
}

// loadUsingFontconfig tries to use the system fontconfig utilities (fc-list / fc-match).
// Returns true on success (fonts loaded using fontconfig), false on failure (caller should fallback).
func (db *Database) loadUsingFontconfig() bool {
	if db == nil {
		return false
	}
	_, err := exec.LookPath("fc-list")
	if err != nil {
		return false
	}
	// Get list of font files using fc-list
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "fc-list", "-f", "%{file}\n")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(out), "\n")
	// Concurrently load files listed by fc-list (dedupe)
	seen := map[string]struct{}{}
	sem := make(chan struct{}, db.loadWorkers)
	var wg sync.WaitGroup
	for _, line := range lines {
		path := strings.TrimSpace(line)
		if path == "" {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		wg.Add(1)
		sem <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()
			// Try to load; ignore errors
			_, _ = db.LoadFontFileWithOptions(p, true)
		}(path)
	}
	wg.Wait()

	// Resolve generic families using fc-match (if present)
	generics := map[string]*string{
		"serif":      &db.familySerif,
		"sans-serif": &db.familySansSerif,
		"cursive":    &db.familyCursive,
		"fantasy":    &db.familyFantasy,
		"monospace":  &db.familyMonospace,
	}
	for gen := range generics {
		// fc-match -s <pattern> prints list of matches; we'll request the family name format.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		cmd := exec.CommandContext(ctx, "fc-match", "-s", gen)
		out, err := cmd.Output()
		cancel()
		if err != nil || len(out) == 0 {
			continue
		}
		// fc-match -s typically outputs "file: \"Family\" \"style\"" or file paths.
		// We should first try to extract quoted family names, then fall back to other parsing.
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// First try: Extract quoted family name from "file: \"Family\" \"style\" format
			if strings.Contains(line, "\"") {
				quoteParts := strings.Split(line, "\"")
				if len(quoteParts) >= 2 {
					// The family name is between the first pair of quotes
					familyName := strings.TrimSpace(quoteParts[1])
					if familyName != "" {
						*(generics[gen]) = familyName
						break
					}
				}
			}

			// Second try: If line looks like a path, try to load that file to read its family name.
			if strings.Contains(line, "/") {
				// sometimes fc-match outputs "DejaVu Serif:style=Book" etc.
				filePath := line
				// if colon-separated with ":", the last element might be file path
				if parts := strings.Split(line, ":"); len(parts) > 1 {
					last := parts[len(parts)-1]
					if strings.Contains(last, "/") {
						filePath = last
					}
				}
				if _, err := os.Stat(filePath); err == nil {
					// try to read first face family
					if id, err := db.LoadFontFileWithOptions(filePath, true); err == nil && id != 0 {
						db.mu.RLock()
						fi := db.faces[id]
						db.mu.RUnlock()
						if fi != nil && len(fi.Families) > 0 {
							*(generics[gen]) = fi.Families[0].Name
						}
					}
					break
				}
			}
			// Third try: As fallback try to treat line as "Family:style" and take text before ':'.
			if idx := strings.Index(line, ":"); idx > 0 {
				name := strings.TrimSpace(line[:idx])
				if name != "" {
					*(generics[gen]) = name
					break
				}
			}
			// Final fallback: use entire line
			*(generics[gen]) = line
			break
		}
	}

	return true
}

// GetFontPathsByFamily returns file paths for fonts matching the given family name exactly.
func (db *Database) GetFontPathsByFamily(familyName string) []string {
	if db == nil {
		return nil
	}

	var paths []string
	seen := make(map[string]bool)

	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, face := range db.faces {
		for _, family := range face.Families {
			// Case-insensitive exact match
			if strings.EqualFold(family.Name, familyName) {
				path := face.Source.FilePath
				if path != "" && !seen[path] {
					paths = append(paths, path)
					seen[path] = true
				}
				break // Found match in this face, move to next face
			}
		}
	}

	return paths
}

// GetFontPathsByFamilyPartial returns file paths for fonts partially matching the family name.
func (db *Database) GetFontPathsByFamilyPartial(searchTerm string) []string {
	if db == nil {
		return nil
	}

	var paths []string
	seen := make(map[string]bool)

	db.mu.RLock()
	defer db.mu.RUnlock()

	searchLower := strings.ToLower(searchTerm)

	for _, face := range db.faces {
		for _, family := range face.Families {
			// Case-insensitive partial match
			if strings.Contains(strings.ToLower(family.Name), searchLower) {
				path := face.Source.FilePath
				if path != "" && !seen[path] {
					paths = append(paths, path)
					seen[path] = true
				}
				break // Found match in this face, move to next face
			}
		}
	}

	return paths
}

// GetFontPathsByPostScript returns file paths for fonts matching the PostScript name exactly.
func (db *Database) GetFontPathsByPostScript(postScriptName string) []string {
	if db == nil {
		return nil
	}

	var paths []string
	seen := make(map[string]bool)

	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, face := range db.faces {
		// Case-insensitive match
		if strings.EqualFold(face.PostScriptName, postScriptName) {
			path := face.Source.FilePath
			if path != "" && !seen[path] {
				paths = append(paths, path)
				seen[path] = true
			}
		}
	}

	return paths
}

// GetDefaultFontPath returns the file path for a generic font family.
// Returns empty string if no exact match is found.
func (db *Database) GetDefaultFontPath(genericFamily string) string {
	if db == nil {
		return ""
	}

	// Resolve generic family to actual family name
	familyName := db.FamilyName(Family{Generic: genericFamily})
	if familyName == "" {
		return ""
	}

	// Get paths for the resolved family name
	paths := db.GetFontPathsByFamily(familyName)
	if len(paths) > 0 {
		return paths[0] // Return first match
	}

	return ""
}

// systemFontDirs returns typical system font directories per OS.
func systemFontDirs() []string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		sysroot := os.Getenv("SYSTEMROOT")
		if sysroot == "" {
			sysroot = `C:\Windows`
		}
		var dirs []string
		dirs = append(dirs, filepath.Join(sysroot, "Fonts"))
		if home != "" {
			dirs = append(dirs, filepath.Join(home, "AppData", "Local", "Microsoft", "Windows", "Fonts"))
			dirs = append(dirs, filepath.Join(home, "AppData", "Roaming", "Microsoft", "Windows", "Fonts"))
		}
		return dirs
	case "darwin":
		var dirs []string
		dirs = append(dirs, "/Library/Fonts", "/System/Library/Fonts", "/Network/Library/Fonts")
		if home != "" {
			dirs = append(dirs, filepath.Join(home, "Library", "Fonts"))
		}
		return dirs
	default:
		// Linux / Unix-like
		var dirs []string
		dirs = append(dirs, "/usr/share/fonts", "/usr/local/share/fonts")
		if home != "" {
			dirs = append(dirs, filepath.Join(home, ".fonts"), filepath.Join(home, ".local", "share", "fonts"))
		}
		return dirs
	}
}

// LoadFontsDirConcurrent recursively loads font files in dir using concurrency and a worker pool.
// 'seen' avoids duplicate canonical paths.
func (db *Database) LoadFontsDirConcurrent(root string, seen map[string]struct{}) error {
	paths := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore access errors
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".ttf", ".otf", ".ttc", ".otc":
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		// ignore walk errors to mimic Rust behavior of skipping unreadable dirs
	}

	// worker pool
	sem := make(chan struct{}, db.loadWorkers)
	var wg sync.WaitGroup

	for _, p := range paths {
		// canonicalize and dedupe
		canon, err := filepath.EvalSymlinks(p)
		if err != nil {
			canon = p
		}
		db.mu.RLock()
		_, already := seen[canon]
		db.mu.RUnlock()
		if already {
			continue
		}
		db.mu.Lock()
		seen[canon] = struct{}{}
		db.mu.Unlock()

		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".ttc" || ext == ".otc" {
				_ = db.LoadFontCollection(path)
			} else {
				_, _ = db.LoadFontFileWithOptions(path, db.useMmap)
			}
		}(p)
	}

	wg.Wait()
	return nil
}

// LoadFontFile loads a single-file font (.ttf/.otf) and returns its ID.
func (db *Database) LoadFontFile(path string) (ID, error) {
	if db == nil {
		return 0, errors.New("database is nil")
	}
	return db.LoadFontFileWithOptions(path, db.useMmap)
}

// LoadFontFileWithOptions allows to choose whether to attempt mmap for reading the file.
func (db *Database) LoadFontFileWithOptions(path string, tryMmap bool) (ID, error) {
	if db == nil {
		return 0, errors.New("database is nil")
	}
	if path == "" {
		return 0, errors.New("empty file path")
	}
	data, mmaped, err := readFileOrMMap(path, tryMmap)
	if err != nil {
		return 0, err
	}
	// Detect collection
	if bytes.HasPrefix(data, []byte("ttcf")) {
		// Treat as collection
		if err := db.loadFontCollectionFromBytes(path, data, mmaped); err != nil {
			return 0, err
		}
		// For collections we don't have a single ID to return; return 0
		return 0, nil
	}
	fi, err := parseFaceInfo(Source{FilePath: path, Data: data}, data, 0)
	if err != nil {
		// If mmaped and parse failed, unmap to free resources (best-effort)
		if mmaped {
			// unmapped mapping will be GC'd only if we unmap; but in this implementation we keep mapped slice
			// and don't unmap here because other faces might rely on it; leaving as-is.
		}
		return 0, err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	fi.ID = nextID()
	db.faces[fi.ID] = fi
	return fi.ID, nil
}

// LoadFontCollection loads a font collection file path (.ttc/.otc).
func (db *Database) LoadFontCollection(path string) error {
	data, mmaped, err := readFileOrMMap(path, db.useMmap)
	if err != nil {
		return err
	}
	return db.loadFontCollectionFromBytes(path, data, mmaped)
}

func (db *Database) loadFontCollectionFromBytes(path string, data []byte, mmaped bool) error {
	if !bytes.HasPrefix(data, []byte("ttcf")) {
		return errors.New("not a collection")
	}
	if len(data) < 12 {
		return errors.New("truncated collection")
	}
	// ttcf header: 0..3 'ttcf', 4..7 version, 8..11 numFonts, then numFonts uint32 offsets starting at 12
	numFonts := binary.BigEndian.Uint32(data[8:12])
	offset := 12
	if len(data) < offset+int(numFonts)*4 {
		return errors.New("truncated collection offsets")
	}
	for i := uint32(0); i < numFonts; i++ {
		offs := binary.BigEndian.Uint32(data[offset+int(i)*4 : offset+int(i)*4+4])
		if offs >= uint32(len(data)) {
			continue
		}
		// parse single face starting at offs: pass the whole data slice but with index i
		fi, err := parseFaceInfo(Source{FilePath: path, Data: data}, data[offs:], uint32(i))
		if err != nil {
			continue
		}
		db.mu.Lock()
		fi.ID = nextID()
		db.faces[fi.ID] = fi
		db.mu.Unlock()
		// keep going for other faces
	}
	// Note: we keep the mapped data alive by storing it as Source.Data in each FaceInfo above
	// (parseFaceInfo stored the same 'data' slice passed in). If mmaped==true the underlying
	// slice is mmap-backed and will remain until unmapped (not done here).
	_ = mmaped
	return nil
}

// PushFaceInfo allows user to insert FaceInfo (ID will be set).
func (db *Database) PushFaceInfo(fi *FaceInfo) ID {
	db.mu.Lock()
	defer db.mu.Unlock()
	fi.ID = nextID()
	db.faces[fi.ID] = fi
	return fi.ID
}

// WithFaceData runs a function with face raw bytes and index.
func (db *Database) WithFaceData(id ID, f func(data []byte, index uint32) error) error {
	db.mu.RLock()
	fi, ok := db.faces[id]
	db.mu.RUnlock()
	if !ok {
		return errors.New("face not found")
	}
	if len(fi.Source.Data) > 0 {
		return f(fi.Source.Data, fi.Index)
	}
	if fi.Source.FilePath != "" {
		data, _, err := readFileOrMMap(fi.Source.FilePath, false)
		if err != nil {
			return err
		}
		return f(data, fi.Index)
	}
	return errors.New("no source data")
}

// Query performs a CSS-like query and returns best matched ID (0 if none).
type Query struct {
	Families []Family
	Weight   Weight
	Stretch  Stretch
	Style    Style
}

func (db *Database) Query(q Query) ID {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Iterate families in order of preference
	for _, fam := range q.Families {
		name := db.FamilyName(fam)
		if name == "" {
			continue
		}
		// gather candidates
		cands := make([]*FaceInfo, 0)
		for _, fi := range db.faces {
			for _, f := range fi.Families {
				// case-sensitive/insensitive? original Rust compares direct strings (family.0 == name)
				if f.Name == name {
					cands = append(cands, fi)
					break
				}
			}
		}
		if len(cands) == 0 {
			continue
		}
		best := findBestMatch(cands, q)
		if best != nil {
			return best.ID
		}
	}
	return 0
}

// findBestMatch is a faithful port of the Rust find_best_match algorithm.
//
// It implements the CSS-like matching steps for font-stretch, font-style and font-weight
// as in the original Rust library. It returns the best candidate or nil if no match.
func findBestMatch(candidates []*FaceInfo, query Query) *FaceInfo {
	if len(candidates) == 0 {
		return nil
	}

	// Step 4: initial matching set = all indices
	matchingSet := make([]int, 0, len(candidates))
	for i := range candidates {
		matchingSet = append(matchingSet, i)
	}

	// Step 4a (font-stretch)
	// Check if any candidate has exact stretch
	queryStretchNum := query.Stretch.ToNumber()
	matchExactStretch := false
	for _, idx := range matchingSet {
		if candidates[idx].Stretch.ToNumber() == queryStretchNum {
			matchExactStretch = true
			break
		}
	}

	var matchingStretchNum int
	if matchExactStretch {
		matchingStretchNum = queryStretchNum
	} else {
		// If query.stretch <= Normal -> prefer narrower first, then wider.
		normalNum := StretchNormal.ToNumber()
		if queryStretchNum <= normalNum {
			// find narrower: candidate.stretch < query.stretch with minimal (query - candidate)
			bestIdx := -1
			bestDiff := int(^uint(0) >> 1) // max int
			for _, idx := range matchingSet {
				cn := candidates[idx].Stretch.ToNumber()
				if cn < queryStretchNum {
					diff := queryStretchNum - cn
					if diff < bestDiff {
						bestDiff = diff
						bestIdx = idx
					}
				}
			}
			if bestIdx != -1 {
				matchingStretchNum = candidates[bestIdx].Stretch.ToNumber()
			} else {
				// pick closest wider (min candidate.stretch - query.stretch)
				bestIdx = -1
				bestDiff = int(^uint(0) >> 1)
				for _, idx := range matchingSet {
					cn := candidates[idx].Stretch.ToNumber()
					diff := cn - queryStretchNum
					if diff < bestDiff {
						bestDiff = diff
						bestIdx = idx
					}
				}
				if bestIdx == -1 {
					return nil
				}
				matchingStretchNum = candidates[bestIdx].Stretch.ToNumber()
			}
		} else {
			// query.stretch > Normal -> prefer wider first, then narrower.
			bestIdx := -1
			bestDiff := int(^uint(0) >> 1)
			for _, idx := range matchingSet {
				cn := candidates[idx].Stretch.ToNumber()
				if cn > queryStretchNum {
					diff := cn - queryStretchNum
					if diff < bestDiff {
						bestDiff = diff
						bestIdx = idx
					}
				}
			}
			if bestIdx != -1 {
				matchingStretchNum = candidates[bestIdx].Stretch.ToNumber()
			} else {
				// pick closest narrower
				bestIdx = -1
				bestDiff = int(^uint(0) >> 1)
				for _, idx := range matchingSet {
					cn := candidates[idx].Stretch.ToNumber()
					diff := queryStretchNum - cn
					if diff < bestDiff {
						bestDiff = diff
						bestIdx = idx
					}
				}
				if bestIdx == -1 {
					return nil
				}
				matchingStretchNum = candidates[bestIdx].Stretch.ToNumber()
			}
		}
	}

	// retain only those with matching_stretch
	newSet := matchingSet[:0]
	for _, idx := range matchingSet {
		if candidates[idx].Stretch.ToNumber() == matchingStretchNum {
			newSet = append(newSet, idx)
		}
	}
	matchingSet = newSet
	if len(matchingSet) == 0 {
		return nil
	}

	// Step 4b (font-style)
	var stylePreference [3]Style
	switch query.Style {
	case StyleItalic:
		stylePreference = [3]Style{StyleItalic, StyleOblique, StyleNormal}
	case StyleOblique:
		stylePreference = [3]Style{StyleOblique, StyleItalic, StyleNormal}
	default:
		stylePreference = [3]Style{StyleNormal, StyleOblique, StyleItalic}
	}

	var matchingStyle *Style
	for _, pref := range stylePreference {
		found := false
		for _, idx := range matchingSet {
			if candidates[idx].Style == pref {
				found = true
				break
			}
		}
		if found {
			tmp := pref
			matchingStyle = &tmp
			break
		}
	}
	if matchingStyle == nil {
		// No matching style found -> no match
		return nil
	}

	// retain only those with matching_style
	newSet = matchingSet[:0]
	for _, idx := range matchingSet {
		if candidates[idx].Style == *matchingStyle {
			newSet = append(newSet, idx)
		}
	}
	matchingSet = newSet
	if len(matchingSet) == 0 {
		return nil
	}

	// Step 4c (font-weight)
	weight := int(query.Weight)
	// helper lambdas
	anyExactWeight := func(w int) bool {
		for _, idx := range matchingSet {
			if int(candidates[idx].Weight) == w {
				return true
			}
		}
		return false
	}

	var matchingWeight int
	if anyExactWeight(weight) {
		matchingWeight = weight
	} else if (400 <= weight && weight < 450) && anyExactWeight(500) {
		// Check 500 first
		matchingWeight = 500
	} else if (450 <= weight && weight <= 500) && anyExactWeight(400) {
		// Check 400 first
		matchingWeight = 400
	} else if weight <= 500 {
		// Closest weight, first checking thinner values and then fatter ones.
		bestIdx := -1
		bestDiff := int(^uint(0) >> 1)
		for _, idx := range matchingSet {
			cw := int(candidates[idx].Weight)
			if cw <= weight {
				diff := weight - cw
				if diff < bestDiff {
					bestDiff = diff
					bestIdx = idx
				}
			}
		}
		if bestIdx != -1 {
			matchingWeight = int(candidates[bestIdx].Weight)
		} else {
			// choose minimal cw - weight
			bestIdx = -1
			bestDiff = int(^uint(0) >> 1)
			for _, idx := range matchingSet {
				cw := int(candidates[idx].Weight)
				diff := cw - weight
				if diff < bestDiff {
					bestDiff = diff
					bestIdx = idx
				}
			}
			if bestIdx == -1 {
				return nil
			}
			matchingWeight = int(candidates[bestIdx].Weight)
		}
	} else {
		// weight > 500: Closest weight, first checking fatter values and then thinner ones.
		bestIdx := -1
		bestDiff := int(^uint(0) >> 1)
		for _, idx := range matchingSet {
			cw := int(candidates[idx].Weight)
			if cw >= weight {
				diff := cw - weight
				if diff < bestDiff {
					bestDiff = diff
					bestIdx = idx
				}
			}
		}
		if bestIdx != -1 {
			matchingWeight = int(candidates[bestIdx].Weight)
		} else {
			bestIdx = -1
			bestDiff = int(^uint(0) >> 1)
			for _, idx := range matchingSet {
				cw := int(candidates[idx].Weight)
				diff := weight - cw
				if diff < bestDiff {
					bestDiff = diff
					bestIdx = idx
				}
			}
			if bestIdx == -1 {
				return nil
			}
			matchingWeight = int(candidates[bestIdx].Weight)
		}
	}

	// retain only those with matching_weight
	newSet = matchingSet[:0]
	for _, idx := range matchingSet {
		if int(candidates[idx].Weight) == matchingWeight {
			newSet = append(newSet, idx)
		}
	}
	matchingSet = newSet
	if len(matchingSet) == 0 {
		return nil
	}

	// Step 4d font-size ignored.

	// Return the first element from matchingSet
	return candidates[matchingSet[0]]
}

///////////////////////////////////////////////////////////////////////////////
// I/O helpers: try mmap, fallback to read-all
///////////////////////////////////////////////////////////////////////////////

// readFileOrMMap tries to memory-map the file if tryMmap is true, otherwise reads into memory.
// Returns (data, mmapedFlag, error). If mmapedFlag is true the returned byte slice is backed by an mmap
// and must not be unmapped while used. We retain the mapping in memory by keeping the []byte.
func readFileOrMMap(path string, tryMmap bool) ([]byte, bool, error) {
	// Validate path
	if path == "" {
		return nil, false, errors.New("empty file path")
	}

	if tryMmap {
		// Try to open file and mmap
		f, err := os.Open(path)
		if err != nil {
			// If we can't open the file, fall back to reading
		} else {
			defer f.Close() // Ensure file is always closed
			mi, err2 := mmap.Map(f, mmap.RDONLY, 0)
			if err2 == nil {
				// keep mapping (do not unmap here) to provide mmap-backed []byte
				// Note: mmap-go's Map requires the file not to be closed on some platforms,
				// but Map copies internal mmap metadata; safest is to keep file opened until unmapped.
				// Here we close the file descriptor but keep the mapping slice; in practice this
				// is supported on Linux. If that causes trouble on some systems, disable mmap.
				return mi, true, nil
			}
			// else fallback to reading file
		}
	}
	// Fallback: read full file into memory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	return data, false, nil
}

///////////////////////////////////////////////////////////////////////////////
// TrueType table parsing (minimal) - extract 'name', 'OS/2' and 'post' tables //
///////////////////////////////////////////////////////////////////////////////

// parseFaceInfo parses families, postscript name, OS/2 and post from font bytes.
// data should be the font binary (for collections, a slice starting at face offset works).
func parseFaceInfo(src Source, data []byte, index uint32) (*FaceInfo, error) {
	// parse table directory
	tableDir, err := parseTableDirectory(data)
	if err != nil {
		return nil, err
	}

	// name table
	nameBytes := tableDir["name"]
	if len(nameBytes) == 0 {
		return nil, errors.New("missing name table")
	}
	families, postScript := parseNameTable(nameBytes)
	if len(families) == 0 {
		return nil, errors.New("font doesn't have a family name")
	}

	// OS/2 table (optional)
	var style Style = StyleNormal
	var weight Weight = WeightNormal
	stretch := StretchNormal
	if os2b, ok := tableDir["OS/2"]; ok && len(os2b) >= 8 {
		w, wid := parseOS2Table(os2b)
		weight = Weight(w)
		// map width class to our Stretch enum (approximate)
		switch wid {
		case 1:
			stretch = StretchUltraCondensed
		case 2:
			stretch = StretchExtraCondensed
		case 3:
			stretch = StretchCondensed
		case 4:
			stretch = StretchSemiCondensed
		case 5:
			stretch = StretchNormal
		case 6:
			stretch = StretchSemiExpanded
		case 7:
			stretch = StretchExpanded
		case 8:
			stretch = StretchExtraExpanded
		case 9:
			stretch = StretchUltraExpanded
		default:
			// clamp/fallback
			if wid < 1 {
				stretch = StretchUltraCondensed
			} else if wid > 9 {
				stretch = StretchUltraExpanded
			} else {
				stretch = StretchNormal
			}
		}
	}

	// post table (monospaced / italic)
	monospaced := false
	if postb, ok := tableDir["post"]; ok && len(postb) >= 16 {
		monospaced = parsePostTableMonospaced(postb)
		if parsePostTableItalic(postb) {
			style = StyleItalic
		}
	}

	fi := &FaceInfo{
		Source: src,
		Index:  index,
		Families: make([]struct {
			Name     string
			Language Language
		}, 0, len(families)),
		PostScriptName: postScript,
		Style:          style,
		Weight:         weight,
		Stretch:        stretch,
		Monospaced:     monospaced,
	}
	for _, f := range families {
		fi.Families = append(fi.Families, struct {
			Name     string
			Language Language
		}{Name: f.Name, Language: f.Lang})
	}
	return fi, nil
}

// parseTableDirectory parses the top-level sfnt table directory and returns raw table bytes keyed by tag.
func parseTableDirectory(data []byte) (map[string][]byte, error) {
	if len(data) < 12 {
		return nil, errors.New("font too small")
	}
	// sfnt version and numTables at offset 4
	numTables := int(binary.BigEndian.Uint16(data[4:6]))
	// table records start at offset 12
	off := 12
	if len(data) < off+numTables*16 {
		return nil, errors.New("truncated table records")
	}
	tables := make(map[string][]byte)
	for i := 0; i < numTables; i++ {
		rec := data[off+i*16 : off+i*16+16]
		tag := string(rec[0:4])
		tableOffset := int(binary.BigEndian.Uint32(rec[8:12]))
		tableLen := int(binary.BigEndian.Uint32(rec[12:16]))
		if tableOffset+tableLen <= len(data) && tableOffset >= 0 && tableLen >= 0 {
			tables[tag] = data[tableOffset : tableOffset+tableLen]
		} else {
			// store nil for missing/truncated tables
			tables[tag] = nil
		}
	}
	return tables, nil
}

type nameEntry struct {
	Name string
	Lang Language
}

// nameRec represents a name table record (platform, encoding, language, name id, length, offset)
type nameRec struct {
	PlatformID uint16
	EncodingID uint16
	LanguageID uint16
	NameID     uint16
	Length     uint16
	Offset     uint16
}

// parseNameTable extracts family names (typographic family first, otherwise family) and postscript name.
func parseNameTable(data []byte) ([]nameEntry, string) {
	if len(data) < 6 {
		return nil, ""
	}
	// name table header: format(2), count(2), stringOffset(2)
	count := int(binary.BigEndian.Uint16(data[2:4]))
	strOff := int(binary.BigEndian.Uint16(data[4:6]))
	recordsOff := 6
	if len(data) < recordsOff+count*12 {
		return nil, ""
	}
	records := make([]nameRec, 0, count)
	for i := 0; i < count; i++ {
		r := data[recordsOff+i*12 : recordsOff+i*12+12]
		rec := nameRec{
			PlatformID: binary.BigEndian.Uint16(r[0:2]),
			EncodingID: binary.BigEndian.Uint16(r[2:4]),
			LanguageID: binary.BigEndian.Uint16(r[4:6]),
			NameID:     binary.BigEndian.Uint16(r[6:8]),
			Length:     binary.BigEndian.Uint16(r[8:10]),
			Offset:     binary.BigEndian.Uint16(r[10:12]),
		}
		records = append(records, rec)
	}

	// helper for reading a record string
	readString := func(rec nameRec) (string, bool) {
		start := strOff + int(rec.Offset)
		end := start + int(rec.Length)
		if start < 0 || end > len(data) {
			return "", false
		}
		buf := data[start:end]
		// platform 0 (Unicode) or 3 (Windows) -> UTF-16BE
		if rec.PlatformID == 0 || rec.PlatformID == 3 {
			if len(buf)%2 != 0 {
				return "", false
			}
			u16s := make([]uint16, 0, len(buf)/2)
			for i := 0; i < len(buf); i += 2 {
				u16s = append(u16s, binary.BigEndian.Uint16(buf[i:i+2]))
			}
			runes := utf16.Decode(u16s)
			return string(runes), true
		}
		// platform 1 Macintosh -> MacRoman single byte, map to UTF-16 (via table)
		if rec.PlatformID == 1 {
			// map each byte via MAC_ROMAN table to rune (uint16) and construct string
			runes16 := make([]uint16, 0, len(buf))
			for _, b := range buf {
				runes16 = append(runes16, macRomanToUTF16[b])
			}
			runes := utf16.Decode(runes16)
			return string(runes), true
		}
		return "", false
	}

	// collect typographic family (nameID 16) first
	var families []nameEntry
	for _, rec := range records {
		if rec.NameID == 16 {
			if s, ok := readString(rec); ok && s != "" {
				families = append(families, nameEntry{Name: s, Lang: languageFromRecord(rec)})
			}
		}
	}
	// fallback to nameID 1 (Font Family)
	if len(families) == 0 {
		for _, rec := range records {
			if rec.NameID == 1 {
				if s, ok := readString(rec); ok && s != "" {
					families = append(families, nameEntry{Name: s, Lang: languageFromRecord(rec)})
				}
			}
		}
	}
	// Ensure English-US first if present
	if len(families) > 1 {
		for i, f := range families {
			if f.Lang == "en-US" {
				if i != 0 {
					families[0], families[i] = families[i], families[0]
				}
				break
			}
		}
	}

	// postscript name (nameID 6)
	postScript := ""
	for _, rec := range records {
		if rec.NameID == 6 {
			if s, ok := readString(rec); ok && s != "" {
				postScript = s
				break
			}
		}
	}

	return families, postScript
}

func languageFromRecord(rec interface{}) Language {
	// We simplified: prefer en-US detection for known codes, otherwise "und".
	// In a more complete implementation we'd map platformID+languageID properly.
	switch r := rec.(type) {
	case nameRec:
		// platform 3 (Windows) languageID corresponds to LANGID, Windows English US is 0x0409 (1033)
		if r.PlatformID == 3 && r.LanguageID == 0x0409 {
			return "en-US"
		}
		// Macintosh English (language id 0) often means English
		if r.PlatformID == 1 && r.LanguageID == 0 {
			return "en"
		}
		_ = r
	}
	return "und"
}

// parseOS2Table extracts usWeightClass and usWidthClass (if present).
func parseOS2Table(data []byte) (weight uint16, width uint16) {
	if len(data) < 8 {
		return uint16(WeightNormal), uint16(5) // StretchNormal = 5
	}
	// according to OS/2 spec: version u16, xAvgCharWidth s16, usWeightClass u16 at offset 4, usWidthClass u16 at offset 6
	if len(data) >= 8 {
		weight = binary.BigEndian.Uint16(data[4:6])
		width = binary.BigEndian.Uint16(data[6:8])
		return weight, width
	}
	return uint16(WeightNormal), uint16(5) // StretchNormal = 5
}

// parsePostTableMonospaced returns true if isFixedPitch (non-zero at offset 12..15)
func parsePostTableMonospaced(data []byte) bool {
	if len(data) < 16 {
		return false
	}
	return binary.BigEndian.Uint32(data[12:16]) != 0
}

// parsePostTableItalic checks italicAngle != 0 at offset 4..7
func parsePostTableItalic(data []byte) bool {
	if len(data) < 8 {
		return false
	}
	// italicAngle is fixed 16.16; zero means not italic
	return binary.BigEndian.Uint32(data[4:8]) != 0
}

///////////////////////////////////////////////////////////////////////////////
// MacRoman table mapping (to UTF-16 code points) - taken from Rust version  //
///////////////////////////////////////////////////////////////////////////////

var macRomanToUTF16 = [256]uint16{
	0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007,
	0x0008, 0x0009, 0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F,
	0x0010, 0x2318, 0x21E7, 0x2325, 0x2303, 0x0015, 0x0016, 0x0017,
	0x0018, 0x0019, 0x001A, 0x001B, 0x001C, 0x001D, 0x001E, 0x001F,
	0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
	0x0028, 0x0029, 0x002A, 0x002B, 0x002C, 0x002D, 0x002E, 0x002F,
	0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037,
	0x0038, 0x0039, 0x003A, 0x003B, 0x003C, 0x003D, 0x003E, 0x003F,
	0x0040, 0x0041, 0x0042, 0x0043, 0x0044, 0x0045, 0x0046, 0x0047,
	0x0048, 0x0049, 0x004A, 0x004B, 0x004C, 0x004D, 0x004E, 0x004F,
	0x0050, 0x0051, 0x0052, 0x0053, 0x0054, 0x0055, 0x0056, 0x0057,
	0x0058, 0x0059, 0x005A, 0x005B, 0x005C, 0x005D, 0x005E, 0x005F,
	0x0060, 0x0061, 0x0062, 0x0063, 0x0064, 0x0065, 0x0066, 0x0067,
	0x0068, 0x0069, 0x006A, 0x006B, 0x006C, 0x006D, 0x006E, 0x006F,
	0x0070, 0x0071, 0x0072, 0x0073, 0x0074, 0x0075, 0x0076, 0x0077,
	0x0078, 0x0079, 0x007A, 0x007B, 0x007C, 0x007D, 0x007E, 0x007F,
	0x00C4, 0x00C5, 0x00C7, 0x00C9, 0x00D1, 0x00D6, 0x00DC, 0x00E1,
	0x00E0, 0x00E2, 0x00E4, 0x00E3, 0x00E5, 0x00E7, 0x00E9, 0x00E8,
	0x00EA, 0x00EB, 0x00ED, 0x00EC, 0x00EE, 0x00EF, 0x00F1, 0x00F3,
	0x00F2, 0x00F4, 0x00F6, 0x00F5, 0x00FA, 0x00F9, 0x00FB, 0x00FC,
	0x2020, 0x00B0, 0x00A2, 0x00A3, 0x00A7, 0x2022, 0x00B6, 0x00DF,
	0x00AE, 0x00A9, 0x2122, 0x00B4, 0x00A8, 0x2260, 0x00C6, 0x00D8,
	0x221E, 0x00B1, 0x2264, 0x2265, 0x00A5, 0x00B5, 0x2202, 0x2211,
	0x220F, 0x03C0, 0x222B, 0x00AA, 0x00BA, 0x03A9, 0x00E6, 0x00F8,
	0x00BF, 0x00A1, 0x00AC, 0x221A, 0x0192, 0x2248, 0x2206, 0x00AB,
	0x00BB, 0x2026, 0x00A0, 0x00C0, 0x00C3, 0x00D5, 0x0152, 0x0153,
	0x2013, 0x2014, 0x201C, 0x201D, 0x2018, 0x2019, 0x00F7, 0x25CA,
	0x00FF, 0x0178, 0x2044, 0x20AC, 0x2039, 0x203A, 0xFB01, 0xFB02,
	0x2021, 0x00B7, 0x201A, 0x201E, 0x2030, 0x00C2, 0x00CA, 0x00C1,
	0x00CB, 0x00C8, 0x00CD, 0x00CE, 0x00CF, 0x00CC, 0x00D3, 0x00D4,
	0xF8FF, 0x00D2, 0x00DA, 0x00DB, 0x00D9, 0x0131, 0x02C6, 0x02DC,
	0x00AF, 0x02D8, 0x02D9, 0x02DA, 0x00B8, 0x02DD, 0x02DB, 0x02C7,
}
