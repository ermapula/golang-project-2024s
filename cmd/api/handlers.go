package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ermapula/golang-project/pkg/model"
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
	id, err := app.readIDParam(r)

	if err != nil {
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
	id, err := app.readIDParam(r)

	if err != nil {
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

func (app *application) deleteGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}

	err = app.models.Games.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
func (app *application) updateGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid game id")
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "404 Not Found")
		return 
	}
	
	var input struct {
		Title       *string    `json:"title"`
		Genre       *string    `json:"genre"`
		ReleaseDate *time.Time `json:"releaseDate"`
		Price       *float64   `json:"price"`
		PublisherId *int       `json:"publisherId"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return 
	}

	if input.Title != nil {
		game.Title = *input.Title
	}
	if input.Genre != nil {
		game.Genre = *input.Genre
	}
	if input.Price != nil {
		game.Price = *input.Price
	}
	if input.ReleaseDate != nil {
		game.ReleaseDate = *input.ReleaseDate
	}
	if input.PublisherId != nil {
		game.PublisherId = *input.PublisherId
	}

	err = app.models.Games.Update(game)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, game)
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
