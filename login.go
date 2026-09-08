package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for login handler")
		return
	}

	// Accept optional expires_in_seconds
	var params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding login request: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate required fields
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// Look up user
	ctx := r.Context()
	user, err := cfg.Queries.GetUserByEmail(ctx, params.Email)
	if err != nil {
		log.Printf("Error fetching user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Verify password using CheckPasswordHash
	matches, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if !matches {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Determine expiration time
	const maxExpiration = time.Hour
	expiration := maxExpiration

	if params.ExpiresInSeconds != nil {
		sec := time.Duration(*params.ExpiresInSeconds) * time.Second
		if sec < maxExpiration {
			expiration = sec
		}
	}

	// Create JWT
	token, err := auth.MakeJWT(user.ID, cfg.Secret, expiration)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not create token")
		return
	}

	// Respond with required shape
	respondWithJSON(w, http.StatusOK, map[string]any{
		"id":         user.ID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"email":      user.Email,
		"token":      token,
	})
}
