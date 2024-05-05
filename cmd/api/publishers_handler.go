package main

import (
	"net/http"
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