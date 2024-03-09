package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ermapula/golang-project/pkg/model"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func (app *application) getPublisher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid publisher ID")
		return
	}

	publisher, err := app.models.Publishers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return 
	}

	app.respondWithJSON(w, http.StatusOK, publisher)
}