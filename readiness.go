package main

import (
	"net/http"
)

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