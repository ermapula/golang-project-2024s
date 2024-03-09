package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	// "github.com/ermapula/golang-project/pkg/model"
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

func (app *application) getPublishers(w http.ResponseWriter, r *http.Request) {

	publishers, err := app.models.Publishers.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, publishers)
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

func (app *application) getGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, game)
}

func (app *application) postGame(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string    `json:"title"`
		Genre       string    `json:"genre"`
		ReleaseDate time.Time `json:"releaseDate"`
		Price       float64   `json:"price"`
		PublisherId int       `json:"publisherId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}
	game := &model.Game{
		Title:       input.Title,
		Genre:       input.Genre,
		Price:       input.Price,
		ReleaseDate: input.ReleaseDate,
		PublisherId: input.PublisherId,
	}

	err = app.models.Games.Post(game)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, game)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dest)
	if err != nil {
		return err
	}

	return nil
}
