package main

import (
	"encoding/json"
	"github.com/ermapula/golang-project/pkg/model"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Games(w http.ResponseWriter, r *http.Request) {
	games := model.GetGames()
	respondWithJSON(w, http.StatusOK, games)
}
