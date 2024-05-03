package main

import (
	"net/http"
	"time"

	"github.com/ermapula/golang-project/pkg/model"
	"github.com/ermapula/golang-project/pkg/validator"
)

func (app *application) getGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
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
		app.badRequestResponse(w, r, err)
		return
	}
	game := &model.Game{
		Title:       input.Title,
		Genre:       input.Genre,
		Price:       input.Price,
		ReleaseDate: input.ReleaseDate,
		PublisherId: input.PublisherId,
	}

	v := validator.New()

	if model.ValidateGame(v, game); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
	}

	err = app.models.Games.Post(game)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
}

func (app *application) deleteGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Games.Delete(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"result": "success"}, nil)
}
func (app *application) updateGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		app.notFoundResponse(w, r)
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
		app.badRequestResponse(w, r, err)
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
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
}

