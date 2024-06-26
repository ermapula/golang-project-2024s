package main

import (
	"errors"
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
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getGames(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string
		Genres      []string
		PublisherId int
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.PublisherId = app.readInt(qs, "publisher_id", -1, v)
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "price", "release_date", "-id", "-title", "-price", "-release_date"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	games, metadata, err := app.models.Games.GetAll(input.Title, input.Genres, input.PublisherId, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"games": games, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getPublisherGames(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Title       string
		Genres      []string
		PublisherId int
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.PublisherId = id
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "price", "release_date", "-id", "-title", "-price", "-release_date"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	games, metadata, err := app.models.Games.GetAll(input.Title, input.Genres, input.PublisherId, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"games": games, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) postGame(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string    `json:"title"`
		Genres      []string  `json:"genres"`
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
		Genres:      input.Genres,
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

	err = app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Games.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateGame(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string    `json:"title"`
		Genres      []string   `json:"genres"`
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
	if input.Genres != nil {
		game.Genres = input.Genres
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

	v := validator.New()
	if model.ValidateGame(v, game); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Games.Update(game)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
