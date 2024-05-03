package main 

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods("GET")
	
	r.HandleFunc("/publishers", app.getPublishers).Methods("GET")
	r.HandleFunc("/publishers/{id:[0-9]+}", app.getPublisher).Methods("GET")

	// r.HandleFunc("/games", app.getGames).Methods("GET")
	r.HandleFunc("/games/{id:[0-9]+}", app.getGame).Methods("GET")
	r.HandleFunc("/games", app.postGame).Methods("POST")
	r.HandleFunc("/games/{id:[0-9]+}", app.updateGame).Methods("PATCH")
	r.HandleFunc("/games/{id:[0-9]+}", app.deleteGame).Methods("DELETE")

	return r
}