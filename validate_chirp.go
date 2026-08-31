package main

import (
	"net/http"
	"encoding/json"
	"log"
)

// helper function
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
    response, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(code)
    w.Write(response)
    return nil
}
// helper function
func respondWithError(w http.ResponseWriter, code int, msg string) error {
    return respondWithJSON(w, code, map[string]string{"error": msg})
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
    	respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for validateChirp handler")
    	return
	}

    type chirp struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := chirp{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding chirp: %s", err)
        respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }

    // Validate chirp length (140 character rule)
    if len(params.Body) > 140 {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    // Success response { "valid": true }
    respondWithJSON(w, http.StatusOK, map[string]bool{
        "valid": true,
    })
}