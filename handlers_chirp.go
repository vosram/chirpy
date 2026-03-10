package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type jsonResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't parse json body", err)
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	// censor chirp then insert to db

	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanBody := censorString(params.Body, badwords)
	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanBody,
		UserID: params.UserID,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create chirp in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, jsonResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})
}
