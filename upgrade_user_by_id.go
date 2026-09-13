package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerUpgradeUser CALLED")
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for upgradeUser handler")
		return
	}

	// --- API KEY VALIDATION USING auth.GetAPIKey ---
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing api key")
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	// --- END API KEY VALIDATION ---

	var params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// If event is NOT "user.upgraded", respond 204 immediately
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Data.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()

	_, err = cfg.Queries.UpgradeUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerUpgradeUser CALLED")
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for upgradeUser handler")
		return
	}

	// --- API KEY VALIDATION ---
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "missing authorization header")
		return
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
		return
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if apiKey == "" {
		respondWithError(w, http.StatusUnauthorized, "empty api key")
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	// --- END API KEY VALIDATION ---

	var params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// If event is NOT "user.upgraded", respond 204 immediately
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Data.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()

	_, err = cfg.Queries.UpgradeUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
*/
/* OLD BELOW
func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerUpgradeUser CALLED")
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for upgradeUser handler")
		return
	}

	// --- API KEY VALIDATION USING auth.GetAPIKey ---
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing api key")
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	// --- END API KEY VALIDATION ---

	var params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Polka sends form-encoded webhook bodies
		if err := r.ParseForm(); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid form body")
			return
		}

		params.Event = r.FormValue("event")
		params.Data.UserID = r.FormValue("user_id")

	} else {
		// JSON webhook bodies
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}
	}

	// OLD CODE BELOW
	//if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
	//	respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
	//	return
	//}

	// If event is NOT "user.upgraded", respond 204 immediately
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Data.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()

	_, err = cfg.Queries.UpgradeUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
*/
/*package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for upgradeUser handler")
		return
	}

	var params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// If event is NOT "user.upgraded", respond 204 immediately
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Data.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := r.Context()

	// Attempt upgrade — ignore returned user, only check error
	_, err = cfg.Queries.UpgradeUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}

	// Success → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

/* OLD VERSION Below
func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for upgradeUser handler")
		return
	}

	var params struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding upgrade user request: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	ctx := r.Context()
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		log.Printf("Error parsing user ID: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := cfg.Queries.UpgradeUser(ctx, userID)
	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not upgrade user")
		return
	}

	response := struct {
		ID          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          user.ID.String(),
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
*/
