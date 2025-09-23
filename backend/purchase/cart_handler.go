package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CartHandler struct {
	service *CartService
}

func NewCartHandler(service *CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	cart, err := h.service.GetCart(userID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CartResponse{
		Cart:    cart,
		Message: "Cart retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	var request AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.AddToCart(userID, request.TourID)
	if err != nil {
		switch err {
		case ErrTourAlreadyInCart:
			h.sendErrorResponse(w, "Tour is already in cart", http.StatusConflict)
		case ErrTourNotPublished:
			h.sendErrorResponse(w, "Tour is not published", http.StatusBadRequest)
		case ErrTourArchived:
			h.sendErrorResponse(w, "Tour is archived and cannot be purchased", http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Failed to add tour to cart: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tour added to cart successfully"})
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveFromCart(userID, uint(tourID))
	if err != nil {
		switch err {
		case ErrCartNotFound:
			h.sendErrorResponse(w, "Cart not found", http.StatusNotFound)
		case ErrItemNotFound:
			h.sendErrorResponse(w, "Tour not found in cart", http.StatusNotFound)
		default:
			h.sendErrorResponse(w, "Failed to remove tour from cart: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tour removed from cart successfully"})
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	err := h.service.ClearCart(userID)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			h.sendErrorResponse(w, "Cart not found", http.StatusNotFound)
		default:
			h.sendErrorResponse(w, "Failed to clear cart: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared successfully"})
}

func (h *CartHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}