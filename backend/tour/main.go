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

	repository := &TourRepository{database: database}
	service := &TourService{repository: repository}
	handler := &TourHandler{service: service}

	r.HandleFunc("/", handler.CreateTour).Methods(http.MethodPost)
	r.HandleFunc("/{id}", handler.GetTourByID).Methods(http.MethodGet)
	r.HandleFunc("/{id}", handler.UpdateTour).Methods(http.MethodPut)
	r.HandleFunc("/", handler.GetMyTours).Methods(http.MethodGet)
	r.HandleFunc("/{id}/keypoint", handler.CreateKeyPoint).Methods(http.MethodPost)
	r.HandleFunc("/{id}/publish", handler.PublishTour).Methods(http.MethodPut)
	r.HandleFunc("/{id}/archive", handler.ArchiveTour).Methods(http.MethodPut)
	r.HandleFunc("/{id}/unarchive", handler.UnarchiveTour).Methods(http.MethodPut)

	r.HandleFunc("/internal/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Tour service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
