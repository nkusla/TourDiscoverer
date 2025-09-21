package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type PurchaseHandler struct {
	service *PurchaseService
}

func NewPurchaseHandler(service *PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	tokens, err := h.service.Checkout(userID)
	if err != nil {
		switch err {
		case ErrCartNotFound:
			h.sendErrorResponse(w, "Cart not found", http.StatusNotFound)
		case ErrEmptyCart:
			h.sendErrorResponse(w, "Cart is empty", http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Checkout failed: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := CheckoutResponse{
		Tokens:  tokens,
		Message: "Checkout completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PurchaseHandler) GetUserTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	tokens, err := h.service.GetUserTokens(userID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get purchase tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := TokensResponse{
		Tokens:  tokens,
		Message: "Purchase tokens retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PurchaseHandler) GetTokenDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenStr := vars["token"]

	token, err := h.service.GetTokenDetails(tokenStr)
	if err != nil {
		switch err {
		case ErrTokenNotFound:
			h.sendErrorResponse(w, "Token not found", http.StatusNotFound)
		case ErrTokenExpired:
			h.sendErrorResponse(w, "Token has expired", http.StatusGone)
		case ErrTokenInvalid:
			h.sendErrorResponse(w, "Token is invalid", http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Failed to get token details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := TokenResponse{
		Token:   token,
		Message: "Token details retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PurchaseHandler) ValidateAccess(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-username")
	if userID == "" {
		h.sendErrorResponse(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	
	// Parse tour ID from string to uint
	var tourID uint
	if err := json.Unmarshal([]byte(tourIDStr), &tourID); err != nil {
		h.sendErrorResponse(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	token, err := h.service.ValidateAccess(userID, tourID)
	if err != nil {
		switch err {
		case ErrTokenNotFound:
			h.sendErrorResponse(w, "No valid purchase found for this tour", http.StatusNotFound)
		case ErrTokenExpired:
			h.sendErrorResponse(w, "Purchase token has expired", http.StatusGone)
		case ErrTokenInvalid:
			h.sendErrorResponse(w, "Purchase token is invalid", http.StatusBadRequest)
		default:
			h.sendErrorResponse(w, "Access validation failed: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := TokenResponse{
		Token:   token,
		Message: "Access validated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PurchaseHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Purchase service is running"})
}

func (h *PurchaseHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}