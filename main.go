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

	// File server wrapped with middleware
	//	fileServer := http.FileServer(http.Dir(filepathRoot))
	//	mux.Handle("/app/", cfg.middlewareMetricsInc(
	//		http.StripPrefix("/app/", fileServer),
	//	))
	mux.Handle("/app/", cfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// /metrics endpoint
	// mux.HandleFunc("GET /api/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	// /reset endpoint
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

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
