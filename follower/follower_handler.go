package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
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

func (h *FollowerHandler) DeleteFollowRelationship(w http.ResponseWriter, r *http.Request) {
	var request UnfollowUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.UnfollowUser(request.Follower, request.Followee)
	if err != nil {
		http.Error(w, "error deleting follow relationship: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		http.Error(w, "username query parameter is required", http.StatusBadRequest)
		return
	}

	followers, err := h.service.GetFollowers(username)
	if err != nil {
		http.Error(w, "error retrieving followers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowerHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		http.Error(w, "username query parameter is required", http.StatusBadRequest)
		return
	}

	following, err := h.service.GetFollowing(username)
	if err != nil {
		http.Error(w, "error retrieving following: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func (h *FollowerHandler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	follower := vars["follower"]
	followee := vars["followee"]

	if follower == "" || followee == "" {
		http.Error(w, "follower and followee parameters are required", http.StatusBadRequest)
		return
	}

	isFollowing, err := h.service.IsFollowing(follower, followee)
	if err != nil {
		http.Error(w, "error checking following status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isFollowing {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *FollowerHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Follower Service"})
}
