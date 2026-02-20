// Program check-recipes reports which recipe bases exist for each language (fr, en, zh, ko, ja).
// Run from repo root: go run ./cmd/check-recipes
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var langSuffixes = []string{".fr", ".en", ".es", ".de", ".it", ".ko", ".zh", ".ja"}

const ext = ".bentext"

func recipesDir() string {
	cwd := filepath.Join(".", "recipes")
	if _, err := os.Stat(cwd); err == nil {
		return cwd
	}
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Join(filepath.Dir(exe), "recipes")
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	return cwd
}

func baseName(filename string) string {
	name := strings.TrimSuffix(filename, ext)
	for _, suf := range langSuffixes {
		name = strings.TrimSuffix(name, suf)
	}
	return name
}

func langFromFilename(filename string) string {
	stem := strings.TrimSuffix(filename, ext)
	extPart := filepath.Ext(stem)
	return strings.TrimPrefix(extPart, ".")
}

func main() {
	dir := recipesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lecture %s: %v\n", dir, err)
		os.Exit(1)
	}

	// base -> set of langs
	matrix := make(map[string]map[string]bool)
	wantLangs := []string{"fr", "en", "zh", "ko", "ja"}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ext) {
			continue
		}
		base := baseName(e.Name())
		lang := langFromFilename(e.Name())
		if lang == "" {
			lang = "fr"
		}
		if matrix[base] == nil {
			matrix[base] = make(map[string]bool)
		}
		matrix[base][lang] = true
	}

	bases := make([]string, 0, len(matrix))
	for b := range matrix {
		bases = append(bases, b)
	}
	sort.Strings(bases)

	fmt.Println("Base name\tfr\ten\tzh\tko\tja\tmanquantes")
	fmt.Println(strings.Repeat("-", 70))
	for _, base := range bases {
		langs := matrix[base]
		row := base
		var missing []string
		for _, l := range wantLangs {
			if langs[l] {
				row += "\t✓"
			} else {
				row += "\t—"
				missing = append(missing, l)
			}
		}
		row += "\t" + strings.Join(missing, ", ")
		fmt.Println(row)
	}
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%d recettes (bases), langues vérifiées: %s\n", len(bases), strings.Join(wantLangs, ", "))
}
