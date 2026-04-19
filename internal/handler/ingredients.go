package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"bentext/internal/ingredients"
)

const spriteSize = 32

// IngredientsLookup handles GET /api/ingredients/lookup?q=...
// Returns 404 when no sprite matches the query.
func IngredientsLookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "Paramètre q requis", http.StatusBadRequest)
		return
	}
	x, y, ok := ingredients.Lookup(q)
	if !ok {
		http.Error(w, "Aucune icône pour cet ingrédient", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"matchedAlias": q,
		"sprite":       map[string]int{"x": x, "y": y, "w": spriteSize, "h": spriteSize},
		"imageUrl":     "/public/ingredients.png",
	})
}

// IngredientsSprite handles GET /api/ingredients/sprite
// Returns the sprite sheet URL, size, and byAlias map (alias -> {x, y}).
func IngredientsSprite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	byAlias := ingredients.AllByAlias()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"imageUrl":    "/public/ingredients.png",
		"spriteSize":  map[string]int{"w": spriteSize, "h": spriteSize},
		"byAlias":     byAlias,
	})
}
