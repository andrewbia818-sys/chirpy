package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/andrewbia818-sys/chirpy/internal/database"
	//"golang.org/x/tools/go/cfg"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for createUser handler")
		return
	}

	// Expect: { "email": "...", "password": "..." }
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// Hash the password
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	// Create user in DB
	ctx := r.Context()
	user, err := cfg.Queries.CreateUser(ctx, database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	})

	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{
		"id":            user.ID,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"email":         user.Email,
		"is_chirpy_red": user.IsChirpyRed,
	})
}
