package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	service   *ReviewService
	validator *validator.Validate
}

func NewReviewHandler(service *ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateReview creates a new review for a tour
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid request body", http.StatusBadRequest))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeErrorResponse(w, NewAPIError("Validation failed: "+err.Error(), http.StatusBadRequest))
		return
	}

	// Get username from context (set by auth middleware)
	username := r.Header.Get("X-Username")
	if username == "" {
		h.writeErrorResponse(w, NewAPIError("Unauthorized", http.StatusUnauthorized))
		return
	}

	review, err := h.service.CreateReview(req, username)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	h.writeSuccessResponse(w, review.ToResponse(), http.StatusCreated)
}

// GetReview retrieves a specific review by ID
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid review ID", http.StatusBadRequest))
		return
	}

	review, err := h.service.GetReviewByID(uint(id))
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	h.writeSuccessResponse(w, review.ToResponse(), http.StatusOK)
}

// GetTourReviews retrieves all reviews for a specific tour
func (h *ReviewHandler) GetTourReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourID, err := strconv.ParseUint(vars["tour_id"], 10, 32)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid tour ID", http.StatusBadRequest))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	reviewList, err := h.service.GetReviewsByTourID(uint(tourID), page, pageSize)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	h.writeSuccessResponse(w, reviewList, http.StatusOK)
}

// GetMyReviews retrieves all reviews by the authenticated user
func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Username")
	if username == "" {
		h.writeErrorResponse(w, NewAPIError("Unauthorized", http.StatusUnauthorized))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	reviewList, err := h.service.GetReviewsByUsername(username, page, pageSize)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	h.writeSuccessResponse(w, reviewList, http.StatusOK)
}

// UpdateReview updates an existing review
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid review ID", http.StatusBadRequest))
		return
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid request body", http.StatusBadRequest))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeErrorResponse(w, NewAPIError("Validation failed: "+err.Error(), http.StatusBadRequest))
		return
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		h.writeErrorResponse(w, NewAPIError("Unauthorized", http.StatusUnauthorized))
		return
	}

	review, err := h.service.UpdateReview(uint(id), req, username)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	h.writeSuccessResponse(w, review.ToResponse(), http.StatusOK)
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid review ID", http.StatusBadRequest))
		return
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		h.writeErrorResponse(w, NewAPIError("Unauthorized", http.StatusUnauthorized))
		return
	}

	err = h.service.DeleteReview(uint(id), username)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTourRating retrieves the average rating for a specific tour
func (h *ReviewHandler) GetTourRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourID, err := strconv.ParseUint(vars["tour_id"], 10, 32)
	if err != nil {
		h.writeErrorResponse(w, NewAPIError("Invalid tour ID", http.StatusBadRequest))
		return
	}

	avgRating, err := h.service.GetTourAverageRating(uint(tourID))
	if err != nil {
		h.writeErrorResponse(w, NewAPIError(err.Error(), GetErrorStatusCode(err)))
		return
	}

	response := map[string]interface{}{
		"tour_id":        tourID,
		"average_rating": avgRating,
	}

	h.writeSuccessResponse(w, response, http.StatusOK)
}

// Ping health check endpoint
func (h *ReviewHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "pong",
		"service": "Review Service",
	}
	h.writeSuccessResponse(w, response, http.StatusOK)
}

// Helper methods
func (h *ReviewHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ReviewHandler) writeErrorResponse(w http.ResponseWriter, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   err.Message,
		"status":  err.StatusCode,
	})
}
