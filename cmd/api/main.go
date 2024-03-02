package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/games", Games).Methods("GET")

	fmt.Println("Server on port :8080")
	http.ListenAndServe(":8080", r)
}
