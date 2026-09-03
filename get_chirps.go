package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/andrewbia818-sys/chirpy/internal/database"
)

func convertChirps(chirps []database.Chirp) []chirpResponse {
	result := make([]chirpResponse, 0, len(chirps))

	for _, chirp := range chirps {
		result = append(result, chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.UUID, // because sqlc uses uuid.NullUUID
		})
	}
	return result
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for getChirps handler")
		return
	}

	chirps, err := cfg.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps")
		return
	}
	// convert to chirpReponse type using convertChirps function and return as JSON
	// with 200 OK status code
	resp := convertChirps(chirps)

	respondWithJSON(w, http.StatusOK, resp)
}

// handlerGetChirp is for the GET /chirps/{chirpID} endpoint.
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for getChirp handler")
		return
	}

	// Extract chirpID from path
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "chirp ID is required")
		return
	}

	// Parse UUID
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}

	// Query database
	chirp, err := cfg.GetChirp(r.Context(), chirpID)
	if err != nil {
		// sqlc returns sql.ErrNoRows for not found
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found")
			return
		}

		log.Printf("Error retrieving chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirp")
		return
	}

	// Convert to response type
	resp := convertChirpToResponse(chirp)

	// Return JSON
	respondWithJSON(w, http.StatusOK, resp)
}
