package handler

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bentext/internal/recipe"
)

var imageExts = []string{".jpg", ".jpeg", ".png", ".gif"}

// recipeBaseName returns the recipe file name without extension and without language suffix.
// e.g. "onigiri-kimchi-mozza.fr.bentext" -> "onigiri-kimchi-mozza"
func recipeBaseName(recipeFilename string) string {
	name := strings.TrimSuffix(recipeFilename, ".bentext")
	for _, lang := range []string{".fr", ".en", ".es", ".de", ".it", ".ko", ".zh"} {
		name = strings.TrimSuffix(name, lang)
	}
	return name
}

func publicDir() string {
	cwd := filepath.Join(".", "public")
	if _, err := os.Stat(cwd); err == nil {
		return cwd
	}
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Join(filepath.Dir(exe), "public")
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	return cwd
}

// baseURL returns the full origin (scheme + host) from the request.
func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if s := r.Header.Get("X-Forwarded-Proto"); s == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

// findRecipeImage looks for an image in publicDir with the same base name as the recipe.
// Returns nil if not found or dimensions cannot be read. URL is the full URL (origin + path).
func findRecipeImage(r *http.Request, publicDir, recipeFilename string) *recipe.RecipeImage {
	base := recipeBaseName(recipeFilename)
	for _, ext := range imageExts {
		fpath := filepath.Join(publicDir, base+ext)
		if _, err := os.Stat(fpath); err != nil {
			continue
		}
		width, height, err := imageDimensions(fpath)
		if err != nil {
			continue
		}
		return &recipe.RecipeImage{
			URL:    baseURL(r) + "/public/" + base + ext,
			Width:  width,
			Height: height,
		}
	}
	return nil
}

func imageDimensions(path string) (width, height int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// PublicFileServer returns an http.Handler that serves files from the public directory.
// Use with: mux.Handle("/public/", http.StripPrefix("/public/", PublicFileServer()))
func PublicFileServer() http.Handler {
	return http.FileServer(http.Dir(publicDir()))
}
