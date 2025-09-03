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

	repository := &ReviewRepository{database: database}
	service := &ReviewService{repository: repository}
	handler := NewReviewHandler(service)

	// Review routes
	r.HandleFunc("/", handler.CreateReview).Methods(http.MethodPost)
	r.HandleFunc("/{id}", handler.GetReview).Methods(http.MethodGet)
	r.HandleFunc("/{id}", handler.UpdateReview).Methods(http.MethodPut)
	r.HandleFunc("/{id}", handler.DeleteReview).Methods(http.MethodDelete)
	r.HandleFunc("/my", handler.GetMyReviews).Methods(http.MethodGet)
	
	// Tour-specific review routes
	r.HandleFunc("/tour/{tour_id}", handler.GetTourReviews).Methods(http.MethodGet)
	r.HandleFunc("/tour/{tour_id}/rating", handler.GetTourRating).Methods(http.MethodGet)

	// Health check
	r.HandleFunc("/internal/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Review service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
