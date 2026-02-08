package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"bentext/internal/handler"
)

const (
	rateLimitRequests = 100
	rateLimitWindow   = time.Minute
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/recipes", handler.Recipes)
	mux.HandleFunc("/api/recipes/", handler.Recipes)
	mux.HandleFunc("/api/convert/bentxt", handler.ConvertBentxt)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/", handler.Home)

	port := ":8080"
	log.Printf("API démarrée sur http://localhost%s", port)
	if err := http.ListenAndServe(port, cors(rateLimit(mux))); err != nil {
		log.Fatal(err)
	}
}

// rateLimit limits each IP to rateLimitRequests per rateLimitWindow (429 when exceeded).
func rateLimit(next http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string][]time.Time)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rateLimitWindow)
		n := 0
		for _, t := range clients[ip] {
			if t.After(cutoff) {
				clients[ip][n] = t
				n++
			}
		}
		clients[ip] = clients[ip][:n]
		if len(clients[ip]) >= rateLimitRequests {
			mu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		clients[ip] = append(clients[ip], now)
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
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
