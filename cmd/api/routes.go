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

	r.HandleFunc("/games", app.requirePermission("games:read", app.getGames)).Methods("GET")
	r.HandleFunc("/games/{id:[0-9]+}", app.requirePermission("games:read", app.getGame)).Methods("GET")
	r.HandleFunc("/games", app.requirePermission("games:write", app.requireActivatedUser(app.postGame))).Methods("POST")
	r.HandleFunc("/games/{id:[0-9]+}", app.requirePermission("games:write", app.requireActivatedUser(app.updateGame))).Methods("PATCH")
	r.HandleFunc("/games/{id:[0-9]+}", app.requirePermission("games:write", app.requireActivatedUser(app.deleteGame))).Methods("DELETE")

	r.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	r.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

	r.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	return app.recoverPanic(app.authenticate(r))
}