package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidatedChirp struct {
	Body string `json:"body"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type ValidResponse struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var chirp ValidatedChirp
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}
	if len(chirp.Body) > 140 {
		responseWithError(w, 400, "Chirp is too long")
		return
	}
	respData := ValidResponse{
		Valid: true,
	}
	respondWithJSON(w, 200, respData)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	customError := ErrorResponse{Error: msg}
	jsonData, err := json.Marshal(customError)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}
