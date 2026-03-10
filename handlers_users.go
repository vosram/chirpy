package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (conf *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	type jsonResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	var reqBody params
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	user, err := conf.db.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create user in database", err)
		return
	}
	respBody := jsonResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, respBody)
}
