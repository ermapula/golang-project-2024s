package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/ermapula/golang-project/pkg/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "Server port")
	flag.StringVar(&cfg.db.dsn, "DB-DSN", "postgres://ermek:adminpass@localhost/golang-project?sslmode=disable", "Postgres DSN")
	flag.Parse()
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func openDB(cfg config) (*sql.DB, error) {
	log.Printf("Opening server...")
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) run() {
	r := mux.NewRouter()

	r.HandleFunc("/publishers", app.getPublishers).Methods("GET")
	r.HandleFunc("/publishers/{id:[0-9]+}", app.getPublisher).Methods("GET")

	r.HandleFunc("/games/{gameId:[0-9]+}", app.getGame).Methods("GET")
	r.HandleFunc("/games", app.postGame).Methods("POST")
	// r.HandleFunc("/games/{gameId:[0-9]+}", app.updateGame).Methods("PUT")
	// r.HandleFunc("/games/{gameId:[0-9]+}", app.deleteGame).Methods("DELETE")

	log.Printf("Server on port :%s\n", app.config.port)
	http.ListenAndServe(app.config.port, r)
}
