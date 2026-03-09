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
		Valid bool `json:"valid"`
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
	respData := validResponse{
		Valid: true,
	}
	respondWithJSON(w, 200, respData)
}
