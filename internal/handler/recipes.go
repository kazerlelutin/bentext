package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bentext/internal/recipe"
)

func Recipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	recipesDir := recipesDir()
	log.Printf("recipes: dossier utilisé = %s", recipesDir)
	entries, err := readRecipeFiles(recipesDir)
	if err != nil {
		log.Printf("recipes: lecture du dossier: %v", err)
		http.Error(w, "Impossible de lire les recettes", http.StatusInternalServerError)
		return
	}

	var result []*recipe.Recipe
	for _, entry := range entries {
		content, err := readFileContent(filepath.Join(recipesDir, entry.Name()))
		if err != nil {
			log.Printf("recipes: lecture %s: %v", entry.Name(), err)
			continue
		}
		slug := slugFromFilename(entry.Name())
		lang := langFromFilename(entry.Name())
		rec := recipe.Parse(content, slug, lang)
		if rec != nil {
			if img := findRecipeImage(r, publicDir(), entry.Name()); img != nil {
				rec.Image = img
			}
			result = append(result, rec)
		}
	}

	if langFilter := strings.TrimSpace(r.URL.Query().Get("lang")); langFilter != "" {
		filtered := result[:0]
		for _, rec := range result {
			if rec.Lang == langFilter {
				filtered = append(filtered, rec)
			}
		}
		result = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("recipes: encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

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

func readRecipeFiles(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e)
		}
	}
	return files, nil
}

func readFileContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func slugFromFilename(name string) string {
	stem := strings.TrimSuffix(name, ".bentext")
	ext := filepath.Ext(stem)
	if ext != "" {
		return strings.TrimSuffix(stem, ext)
	}
	return stem
}

func langFromFilename(name string) string {
	stem := strings.TrimSuffix(name, ".bentext")
	ext := strings.TrimPrefix(filepath.Ext(stem), ".")
	if ext == "" {
		return "fr"
	}
	return ext
}
