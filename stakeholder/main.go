package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	database := InitDatabase()

	repository := &StakeholderRepository{database: database}
	service := &StakeholderService{repository: repository}
	validator := validator.New()
	handler := &StakeholderHandler{service: service, validator: validator}

	// Public routes
	r.HandleFunc("/stakeholder", handler.CreateStakeholder).Methods(http.MethodPost)
	
	// Protected routes (require JWT validation from API gateway)
	r.HandleFunc("/profile", handler.GetProfile).Methods(http.MethodGet)
	r.HandleFunc("/profile", handler.UpdateProfile).Methods(http.MethodPut)

	// Internal routes
	r.HandleFunc("/internal/ping", handler.Ping).Methods(http.MethodGet)
	r.HandleFunc("/internal/user", handler.CreateStakeholderFromAuth).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Stakeholder service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
