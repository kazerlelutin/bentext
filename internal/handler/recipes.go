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
	for i, entry := range entries {
		content, err := readFileContent(filepath.Join(recipesDir, entry.Name()))
		if err != nil {
			log.Printf("recipes: lecture %s: %v", entry.Name(), err)
			continue
		}
		lang := langFromFilename(entry.Name())
		rec := recipe.Parse(content, i, lang)
		if rec != nil {
			result = append(result, rec)
		}
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

func langFromFilename(name string) string {
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	if ext == "bentext" {
		return "fr"
	}
	return ext
}
