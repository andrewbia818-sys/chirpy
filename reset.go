package main

import (
	"log"
	"net/http"
)

// /reset handler: resets counter to zero
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed for reset handler", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Running on platform:", cfg.Platform)

	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "FORBIDDEN")
		return
	}
	cfg.fileserverHits.Store(0)

	//if cfg.Queries != nil {
	cfg.DeleteUsers(r.Context())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0, all users deleted"))
}
