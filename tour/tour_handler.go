package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type TourHandler struct {
	//service *TourService
}

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

var validate = validator.New()

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *TourHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Tour Service"})
}
