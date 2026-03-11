package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/chirpy/internal/auth"
	"github.com/vosram/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedP, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hashedP,
	})
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

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type jsonResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody params
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't decode request body", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	passMatch, err := auth.CheckPasswordHash(reqBody.Password, dbUser.HashedPassword)
	if err != nil || !passMatch {
		responseWithError(w, http.StatusInternalServerError, "Error checking password match", err)
		return
	}

	respondWithJSON(w, http.StatusOK, jsonResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
