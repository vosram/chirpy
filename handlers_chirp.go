package main

import (
	"encoding/json"
	"errors"
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

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	type resChrip struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"Updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't get all chirps", err)
		return
	}

	finalChirps := make([]resChrip, len(dbChirps))
	for i, dbChirp := range dbChirps {
		finalChirps[i] = resChrip{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, finalChirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	type resChrip struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"Updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirpIDstr := r.PathValue("chirpID")
	if len(chirpIDstr) == 0 {
		err := errors.New("chirpID not available from request")
		responseWithError(w, http.StatusInternalServerError, "Couldn't get chirpID from request", err)
		return
	}

	chirpUUID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "chirp id was not a valid id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Couldn't find that chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, resChrip{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
