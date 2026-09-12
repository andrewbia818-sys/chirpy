package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/andrewbia818-sys/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for login handler")
		return
	}

	// expires_in_seconds removed
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding login request: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	ctx := r.Context()
	user, err := cfg.Queries.GetUserByEmail(ctx, params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	matches, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !matches {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// ACCESS TOKEN — always 1 hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.Secret, time.Hour)
	if err != nil {
		log.Printf("Error creating access token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not create access token")
		return
	}

	// REFRESH TOKEN — using MakeRefreshToken() exactly
	refreshToken := auth.MakeRefreshToken()

	expiresAt := time.Now().UTC().Add(60 * 24 * time.Hour) // 60 days

	_, err = cfg.Queries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Printf("Error storing refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not store refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"id":            user.ID,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"email":         user.Email,
		"token":         accessToken,
		"refresh_token": refreshToken,
		"is_chirpy_red": user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for refresh handler")
		return
	}

	// Extract refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	ctx := r.Context()

	// Look up refresh token in DB
	refreshTokenRecord, err := cfg.Queries.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Check revoked or expired
	if refreshTokenRecord.RevokedAt.Valid ||
		time.Now().UTC().After(refreshTokenRecord.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token is revoked or expired")
		return
	}

	// Extract user ID
	userID := refreshTokenRecord.UserID.UUID

	// Create new access token (1 hour)
	newAccessToken, err := auth.MakeJWT(userID, cfg.Secret, time.Hour)
	if err != nil {
		log.Printf("Error creating new access token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not create new access token")
		return
	}

	// Respond with new access token
	respondWithJSON(w, http.StatusOK, map[string]any{
		"token": newAccessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for revoke handler")
		return
	}

	// Extract refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	ctx := r.Context()

	// Look up refresh token in DB
	rt, err := cfg.Queries.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// If already revoked, treat as unauthorized
	if rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token already revoked")
		return
	}

	_, err = cfg.Queries.RevokeRefreshToken(ctx, database.RevokeRefreshTokenParams{
		Token: refreshToken,
	})

	// Revoke the token
	_, err = cfg.Queries.RevokeRefreshToken(ctx, database.RevokeRefreshTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not revoke refresh token")
		return
	}

	// 204 No Content — successful, no body
	w.WriteHeader(http.StatusNoContent)
}
