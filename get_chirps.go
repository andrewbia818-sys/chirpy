package main

import (
	"log"
	"net/http"

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
	// iterate through chirps and convert to chirpReponse type
	// using func convertChirpToResponse
	resp := convertChirps(chirps)

	respondWithJSON(w, http.StatusOK, resp)
}

//respondWithJSON(w, http.StatusOK, chirps)
//}
