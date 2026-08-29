package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Middleware: increments fileserverHits on every request
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Non apiCOnfig version of handlerReadiness
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Readiness handler version using *apiConfig struct
//func (cfg *apiConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("OK"))
//}

// /metrics handler: prints "Hits: x"
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hits := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

// /reset handler: resets counter to zero
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const port = "8080"
	const filepathRoot = "./app/"

	cfg := &apiConfig{}

	mux := http.NewServeMux()

	// File server wrapped with middleware
	//	fileServer := http.FileServer(http.Dir(filepathRoot))
	//	mux.Handle("/app/", cfg.middlewareMetricsInc(
	//		http.StripPrefix("/app/", fileServer),
	//	))
	mux.Handle("/app/", cfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// /metrics endpoint
	mux.HandleFunc("GET /api/metrics", cfg.handlerMetrics)

	// /reset endpoint
	mux.HandleFunc("POST /api/reset", cfg.handlerReset)

	// readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Root endpoint
	//mux.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//	w.WriteHeader(http.StatusOK)
	//	w.Write([]byte("OK"))
	//})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s\n", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
