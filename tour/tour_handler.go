package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TourHandler struct {
	service *TourService
}

var validate = validator.New()

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	var request CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		h.sendErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can create tours", http.StatusForbidden)
		return
	}

	tour, err := h.service.CreateTour(&request, username)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateTourResponse{
		ID:             tour.ID,
		Name:           tour.Name,
		Description:    tour.Description,
		Difficulty:     tour.Difficulty,
		Tags:           tour.Tags,
		Status:         tour.Status,
		Price:          tour.Price,
		AuthorUsername: tour.AuthorUsername,
		Message:        "Tour created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TourHandler) UpdateTour(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	var request UpdateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		h.sendErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can update tours", http.StatusForbidden)
		return
	}

	tour, err := h.service.UpdateTour(uint(id), &request, username)
	if err != nil {
		if errors.Is(err, ErrTourNotFound) {
			h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			h.sendErrorResponse(w, "Unauthorized: You can only update your own tours", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrTourNotEditable) {
			h.sendErrorResponse(w, "Tour is not editable: "+err.Error(), http.StatusBadRequest)
			return
		}
		h.sendErrorResponse(w, "Failed to update tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateTourResponse{
		ID:             tour.ID,
		Name:           tour.Name,
		Description:    tour.Description,
		Difficulty:     tour.Difficulty,
		Tags:           tour.Tags,
		Status:         tour.Status,
		Price:          tour.Price,
		AuthorUsername: tour.AuthorUsername,
		Message:        "Tour updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TourHandler) GetMyTours(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can access their tours", http.StatusForbidden)
		return
	}

	tours, err := h.service.GetToursByAuthor(username)
	if err != nil {
		h.sendErrorResponse(w, "Failed to fetch tours: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetToursResponse{
		Tours: tours,
		Count: len(tours),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TourHandler) GetTourByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tour, err := h.service.GetTourByID(uint(id))
	if err != nil {
		h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) CreateKeyPoint(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	var request CreateKeyPointRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		h.sendErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can create key points", http.StatusForbidden)
		return
	}

	keyPoint, err := h.service.CreateKeyPoint(&request, uint(id), username)
	if err != nil {
		if errors.Is(err, ErrTourNotFound) {
			h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
		}
		if errors.Is(err, ErrUnauthorized) {
			h.sendErrorResponse(w, "Unauthorized: You can only create key points for your own tours", http.StatusForbidden)
		}
		if errors.Is(err, ErrTourNotEditable) {
			h.sendErrorResponse(w, "Tour is not editable: "+err.Error(), http.StatusBadRequest)
		} else {
			h.sendErrorResponse(w, "Failed to create key point: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := CreateKeyPointResponse{
		ID:          keyPoint.ID,
		TourID:      keyPoint.TourID,
		Name:        keyPoint.Name,
		Description: keyPoint.Description,
		Latitude:    keyPoint.Latitude,
		Longitude:   keyPoint.Longitude,
		ImageURL:    keyPoint.ImageURL,
		Order:       keyPoint.Order,
		Message:     "Key point created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TourHandler) PublishTour(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can publish tours", http.StatusForbidden)
		return
	}

	err = h.service.PublishTour(uint(id), username)
	if err != nil {
		switch {
		case errors.Is(err, ErrTourNotFound):
			h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
		case errors.Is(err, ErrUnauthorized):
			h.sendErrorResponse(w, "Unauthorized: You can only publish your own tours", http.StatusForbidden)
		case errors.Is(err, ErrTourNotPublishable):
			h.sendErrorResponse(w, "Tour cannot be published: "+err.Error(), http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Failed to publish tour: "+err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourHandler) ArchiveTour(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can archive tours", http.StatusForbidden)
		return
	}

	err = h.service.ArchiveTour(uint(id), username)
	if err != nil {
		switch {
		case errors.Is(err, ErrTourNotFound):
			h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
		case errors.Is(err, ErrUnauthorized):
			h.sendErrorResponse(w, "Unauthorized: You can only archive your own tours", http.StatusForbidden)
		case errors.Is(err, ErrTourNotArchivable):
			h.sendErrorResponse(w, "Tour cannot be archived: "+err.Error(), http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Failed to archive tour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourHandler) UnarchiveTour(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-username")
	userRole := r.Header.Get("x-user-role")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	// Check if user is a guide (author)
	if userRole != RoleGuide {
		h.sendErrorResponse(w, "Only guides can unarchive tours", http.StatusForbidden)
		return
	}

	err = h.service.UnarchiveTour(uint(id), username)
	if err != nil {
		switch {
		case errors.Is(err, ErrTourNotFound):
			h.sendErrorResponse(w, "Tour not found", http.StatusNotFound)
		case errors.Is(err, ErrUnauthorized):
			h.sendErrorResponse(w, "Unauthorized: You can only unarchive your own tours", http.StatusForbidden)
		case errors.Is(err, ErrTourNotUnarchivable):
			h.sendErrorResponse(w, "Tour cannot be unarchived: "+err.Error(), http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Failed to unarchive tour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TourHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Tour Service"})
}

func (h *TourHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
