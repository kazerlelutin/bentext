package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"bentext/internal/ingredients"
	"bentext/internal/recipe"
)

var allowedRecipeLangs = map[string]struct{}{
	"fr": {}, "en": {}, "ja": {}, "zh": {}, "ko": {},
}

var recipeSlugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func Recipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, "/api/recipes")
	rest = strings.Trim(rest, "/")
	if rest != "" {
		serveRecipeByPath(w, r, rest)
		return
	}

	serveRecipeList(w, r)
}

func serveRecipeList(w http.ResponseWriter, r *http.Request) {
	recipesDir := recipesDir()
	log.Printf("recipes: dossier utilisé = %s", recipesDir)
	entries, err := readRecipeFiles(recipesDir)
	if err != nil {
		log.Printf("recipes: lecture du dossier: %v", err)
		http.Error(w, "Impossible de lire les recettes", http.StatusInternalServerError)
		return
	}

	includeBentext := wantBentextInJSON(r)

	var result []*recipe.Recipe
	for _, entry := range entries {
		content, err := readFileUnderRoot(recipesDir, entry.Name())
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
			enrichIngredientIcons(rec)
			if includeBentext {
				rec.Bentext = content
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

func serveRecipeByPath(w http.ResponseWriter, r *http.Request, rest string) {
	parts := strings.Split(rest, "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	lang := strings.TrimSpace(parts[0])
	slug := strings.TrimSpace(parts[1])
	if _, ok := allowedRecipeLangs[lang]; !ok || !recipeSlugRe.MatchString(slug) {
		http.NotFound(w, r)
		return
	}

	dir := recipesDir()
	fname := slug + "." + lang + ".bentext"
	content, err := readFileUnderRoot(dir, fname)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		log.Printf("recipes: lecture %s: %v", fname, err)
		http.Error(w, "Impossible de lire la recette", http.StatusInternalServerError)
		return
	}

	format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("format")))
	if format == "bentext" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte(content)); err != nil {
			log.Printf("recipes: écriture bentext: %v", err)
		}
		return
	}
	if format != "" && format != "json" {
		http.Error(w, "Paramètre format invalide (json ou bentext)", http.StatusBadRequest)
		return
	}

	rec := recipe.Parse(content, slug, lang)
	if rec == nil {
		http.Error(w, "Recette invalide", http.StatusUnprocessableEntity)
		return
	}
	if img := findRecipeImage(r, publicDir(), fname); img != nil {
		rec.Image = img
	}
	enrichIngredientIcons(rec)
	if wantBentextInJSON(r) {
		rec.Bentext = content
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rec); err != nil {
		log.Printf("recipes: encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

// wantBentextInJSON is true when the client asks to embed the raw .bentext in the JSON body.
func wantBentextInJSON(r *http.Request) bool {
	q := r.URL.Query()
	v := strings.ToLower(strings.TrimSpace(q.Get("bentext")))
	if v == "1" || v == "true" || v == "yes" {
		return true
	}
	inc := strings.ToLower(strings.TrimSpace(q.Get("include")))
	return inc == "bentext"
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

// readFileUnderRoot reads a single file name inside root (no path segments, no traversal).
func readFileUnderRoot(root, name string) (string, error) {
	if name == "" || name != filepath.Base(name) {
		return "", fmt.Errorf("invalid path")
	}
	rootAbs, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return "", err
	}
	fullAbs, err := filepath.Abs(filepath.Join(rootAbs, name))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootAbs, fullAbs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid path")
	}
	b, err := os.ReadFile(fullAbs)
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

func enrichIngredientIcons(rec *recipe.Recipe) {
	for i := range rec.Ingredients {
		x, y, ok := ingredients.Lookup(rec.Ingredients[i].Name)
		if ok {
			rec.Ingredients[i].Icon = &recipe.SpriteCoords{X: x, Y: y}
		}
		for j := range rec.Ingredients[i].Alternatives {
			x, y, ok := ingredients.Lookup(rec.Ingredients[i].Alternatives[j].Name)
			if ok {
				rec.Ingredients[i].Alternatives[j].Icon = &recipe.SpriteCoords{X: x, Y: y}
			}
		}
	}
}
