package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/andrewbia818-sys/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// helper function ensures JSON fields with lower case
func convertChirpToResponse(c database.Chirp) chirpResponse {
	return chirpResponse{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID.UUID, // because sqlc uses uuid.NullUUID
	}
}

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

// helper function replaceBadWords  with ****, respond with a new cleaned body
func replaceBadWords(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")

	for i, w := range words {
		lw := strings.ToLower(w)
		for _, bad := range badWords {
			if lw == bad {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for createChirp handler")
		return
	}

	// Extract bearer token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	// Validate JWT
	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	// Decode request body
	var params struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Length rule
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean the chirp
	cleaned := replaceBadWords(params.Body)

	// Insert into database using sqlc
	chirp, err := cfg.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	// Respond with created chirp
	resp := convertChirpToResponse(chirp)
	respondWithJSON(w, http.StatusCreated, resp)
}
