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
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
