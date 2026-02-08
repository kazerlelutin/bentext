package handler

import (
	"encoding/json"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Bentext API",
		"routes": []map[string]string{
			{"method": "GET", "path": "/", "description": "This list of routes"},
			{"method": "GET", "path": "/health", "description": "Health check"},
			{"method": "GET", "path": "/api/recipes", "description": "List all recipes (from recipes/ folder)"},
			{"method": "POST", "path": "/api/convert/bentxt", "description": "Convert bentxt text (body) to JSON"},
		},
	})
}
