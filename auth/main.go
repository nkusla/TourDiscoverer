package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	database := InitDB()
	SeedDB(database)

	repository := &UserRepository{database: database}
	handler := &AuthHandler{repository: repository}

	r.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	r.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Auth service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
