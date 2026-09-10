package main

import (
	"log"
	"net/http"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for deleteChirp handler")
		return
	}

	// Extract access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	// Validate JWT → gives authenticated user ID
	userID, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired access token")
		return
	}

	// Extract chirp ID from path
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}

	ctx := r.Context()

	// Fetch chirp from DB
	chirp, err := cfg.Queries.GetChirp(ctx, chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	// Check ownership
	if chirp.UserID.UUID != userID {
		respondWithError(w, http.StatusForbidden, "you are not the author of this chirp")
		return
	}

	// Delete chirp
	err = cfg.Queries.DeleteChirp(ctx, chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp")
		return
	}

	// 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

/*
 func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
 	log.Println("METHOD:", r.Method)

 	if r.Method != http.MethodDelete {
 		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for deleteChirp handler")
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

 	// Extract chirpID from URL path
 	chirpIDStr := r.URL.Path[len("/api/chirps/"):]
 	chirpID, err := uuid.Parse(chirpIDStr)
 	if err != nil {
 		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
 		return
 	}

 	ctx := r.Context()

 	// Check if the chirp exists and belongs to the authenticated user
 	chirp, err := cfg.Queries.GetChirpByID(ctx, chirpID)
 	if err != nil {
 		if err == sql.ErrNoRows {
 			respondWithError(w, http.StatusNotFound, "chirp not found")
 		} else {
 			log.Printf("Error fetching chirp: %s", err)
 			respondWithError(w, http.StatusInternalServerError, "could not fetch chirp")
 		}
 		return
 	}

 	if chirp.UserID != userID {
 		respondWithError(w, http.StatusForbidden, "you are not authorized to delete this chirp")
 		return
 	}

 	// Delete the chirp
 	err = cfg.Queries.DeleteChirp(ctx, chirpID)
 	if err != nil {
 		log.Printf("Error deleting chirp: %s", err)
 		respondWithError(w, http.StatusInternalServerError, "could not delete chirp")
 		return
 	}

 	// Respond with success
 	respondWithJSON(w, http.StatusOK, map[string]string{
 		"message": "chirp deleted successfully",
 	})
 }
*/
