package main

import (
	"encoding/json"
	"log"
	"net/http"
	//"golang.org/x/tools/go/cfg"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed for createUser handler")
		return
	}

	var params User
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	ctx := r.Context()
	user, err := cfg.CreateUser(ctx, params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{
		"id":         user.ID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"email":      user.Email,
	})
	//    w.Header().Set("Content-Type", "application/json")
	//    w.WriteHeader(http.StatusCreated)
	//   if err := json.NewEncoder(w).Encode(user); err != nil {
	//        log.Printf("Error encoding response: %s", err)
	//   }

}
