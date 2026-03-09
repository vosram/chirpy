package main

import (
	"encoding/json"
	"net/http"
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
