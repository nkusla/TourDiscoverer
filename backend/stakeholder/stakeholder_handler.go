package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type StakeholderHandler struct {
	service   *StakeholderService
	validator *validator.Validate
}

func (h *StakeholderHandler) CreateStakeholder(w http.ResponseWriter, r *http.Request) {
	var req CreateStakeholderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stakeholder, err := h.service.CreateStakeholder(
		req.Username,
		req.FirstName,
		req.LastName,
		req.ProfilePicture,
		req.Biography,
		req.Motto,
	)
	if err != nil {
		if err.Error() == "stakeholder already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract username from headers (set by API gateway after JWT validation)
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "Username not found in request", http.StatusUnauthorized)
		return
	}

	stakeholder, err := h.service.GetStakeholderProfile(username)
	if err != nil {
		if err.Error() == "profile not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Extract username from headers (set by API gateway after JWT validation)
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "Username not found in request", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	stakeholder, err := h.service.UpdateStakeholderProfile(
		username,
		req.FirstName,
		req.LastName,
		req.ProfilePicture,
		req.Biography,
		req.Motto,
	)
	if err != nil {
		if err.Error() == "profile not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stakeholder service is running"))
}

func (h *StakeholderHandler) CreateStakeholderFromAuth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Create stakeholder with just username - other fields can be updated later via profile
	stakeholder, err := h.service.CreateStakeholder(
		req.Username,
		"", // empty first_name - user can update via profile
		"", // empty last_name - user can update via profile
		"", // empty profile_picture - user can update via profile
		"", // empty biography - user can update via profile
		"", // empty motto - user can update via profile
	)
	if err != nil {
		if err.Error() == "stakeholder already exists" {
			// If stakeholder already exists, that's okay for this internal endpoint
			w.WriteHeader(http.StatusCreated)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stakeholder)
}
