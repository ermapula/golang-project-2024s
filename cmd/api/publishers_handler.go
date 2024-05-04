package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) getPublishers(w http.ResponseWriter, r *http.Request) {

	publishers, err := app.models.Publishers.GetAll()
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"publishers": publishers}, nil)
}

func (app *application) getPublisher(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	publisher, err := app.models.Publishers.Get(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"publisher": publisher}, nil)
}

func (app *application) convert(w http.ResponseWriter, r *http.Request) {
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
	s := "INSERT INTO games (title, genres, release_date, price, publisher_id) VALUES ("
	s += fmt.Sprintf(`'%v', '{`, input.Title)
	for i, a := range input.Genres {
		if i != len(input.Genres) - 1 {
			s += fmt.Sprintf(`"%v", `, a)
		} else {
			s += fmt.Sprintf(`"%v"}', `, a)
		}
	}
	s += fmt.Sprintf(`'%v', %v, %v);`, input.ReleaseDate.Format("2006-01-02"), input.Price, input.PublisherId)

	fmt.Fprintf(w, s)
}