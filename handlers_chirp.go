package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/chirpy/internal/database"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type validatedChirp struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var chirp validatedChirp
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	if len(chirp.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	respData := validResponse{
		CleanedBody: censorString(chirp.Body, badwords),
	}
	respondWithJSON(w, 200, respData)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type jsonBody struct {
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
	var chirp jsonBody
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't parse json body", err)
		return
	}

	if len(chirp.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	// censor chirp then insert to db

	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanBody := censorString(chirp.Body, badwords)
	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanBody,
		UserID: chirp.UserID,
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create chirp in database", err)
		return
	}

	res := jsonResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, res)
}
