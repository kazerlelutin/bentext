package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"bentext/internal/recipe"
)

// ConvertBentxt accepts text in bentxt format and returns the recipe in JSON.
// POST /api/convert/bentxt
// Query: lang (optional, default "fr"), slug (optional, default "")
func ConvertBentxt(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "fr"
	}
	slug := r.URL.Query().Get("slug")

	rec := recipe.Parse(string(body), slug, lang)
	if rec == nil {
		http.Error(w, "Invalid bentxt format (at least 3 sections separated by ---)", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rec); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
	}
}
