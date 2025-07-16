package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	database := InitDatabase()
	SeedTour(database)

	// repository := &TourRepository{database: database}
	// service := &TourService{repository: repository}
	handler := &TourHandler{}

	r.HandleFunc("/createTour", handler.CreateTour).Methods(http.MethodPost)
	r.HandleFunc("/ping", handler.Ping).Methods(http.MethodGet)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Tour service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
