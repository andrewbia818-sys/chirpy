package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filepathRoot = "./app/"

	cfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// /metrics endpoint
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	// /reset endpoint
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	// readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// validateChirp endpoint.
	mux.HandleFunc("POST /api/validate_chirp/{rest...}", handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s\n", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
