package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type FollowerHandler struct {
	service *FollowerService
}

var validate = validator.New()

func (h *FollowerHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateUser(request.Username)
	if err != nil {
		http.Error(w, "error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FollowerHandler) CreateFollowRelationship(w http.ResponseWriter, r *http.Request) {
	var request FollowUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.FollowUser(request.Follower, request.Followee)
	if err != nil {
		http.Error(w, "error creating follow relationship: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FollowerHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Follower Service"})
}
