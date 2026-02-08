package main

import (
	"log"
	"net/http"

	"bentext/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/recipes", handler.Recipes)
	mux.HandleFunc("/api/recipes/", handler.Recipes)
	mux.HandleFunc("/api/convert/bentxt", handler.ConvertBentxt)
	mux.HandleFunc("/api/hello", handler.Hello)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/", handler.Home)

	port := ":8080"
	log.Printf("API démarrée sur http://localhost%s", port)
	if err := http.ListenAndServe(port, cors(mux)); err != nil {
		log.Fatal(err)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
