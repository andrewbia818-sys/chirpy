package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/andrewbia818-sys/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPut {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for updateUser handler")
		return
	}

	// Extract access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	// Validate access token → gives us the authenticated user ID
	userID, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired access token")
		return
	}

	// Parse request body
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	ctx := r.Context()

	// Update user in DB
	updatedUser, err := cfg.Queries.UpdateUser(ctx, database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	// Respond with updated user (no password)
	respondWithJSON(w, http.StatusOK, map[string]any{
		"id":            updatedUser.ID,
		"created_at":    updatedUser.CreatedAt,
		"updated_at":    updatedUser.UpdatedAt,
		"email":         updatedUser.Email,
		"is_chirpy_red": updatedUser.IsChirpyRed,
	})
}
