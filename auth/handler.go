package main

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if registerReq.Username == "" || registerReq.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if registerReq.Role == "" {
		registerReq.Role = "tourist"
	}

	// Here you would typically save the user to a database.

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Here you would typically validate the user credentials against a database.
	// For simplicity, we assume the credentials are valid.

	tokenString, err := CreateJWT(loginReq.Username, "tourist")
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JWTResponse{Token: tokenString})
}

func (h *AuthHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
