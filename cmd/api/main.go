package main

import (
	"flag"
	"log"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	db   struct {
		dsn string
	}
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Ping no error")
	r := mux.NewRouter()

	r.HandleFunc("/games", Games).Methods("GET")

	log.Printf("Server on port :8080")
	http.ListenAndServe(":8080", r)
}

func openDB(cfg config) (*sql.DB, error) {
	log.Printf("Opening server...")
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}