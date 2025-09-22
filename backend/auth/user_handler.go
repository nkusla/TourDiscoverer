package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *UserService
}

var validate = validator.New()

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	authLogger.Info("User registration attempt")

	var registerReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		authLogger.Error("Failed to decode registration request", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(registerReq); err != nil {
		authLogger.Warn("Registration validation failed: " + err.Error())
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.RegisterUser(registerReq)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyExists) || errors.Is(err, ErrEmailAlreadyExists) {
			authLogger.Warn("Registration failed: " + err.Error())
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, ErrInvalidRole) {
			authLogger.Warn("Registration failed: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			authLogger.Error("Registration failed", err)
			http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Record successful registration
	// recordRegistration() // UKLONJENO - nema metrics
	authLogger.InfoWithFields("User registered successfully", map[string]interface{}{
		"username": registerReq.Username,
		"role":     registerReq.Role,
	})

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	authLogger.Info("User login attempt")

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		authLogger.Error("Failed to decode login request", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(loginReq); err != nil {
		authLogger.Warn("Login validation failed: " + err.Error())
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.AuthenticateUser(loginReq.Username, loginReq.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			// recordLoginFailure() // UKLONJENO - nema metrics
			authLogger.Warn("Login failed: invalid credentials for user " + loginReq.Username)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if errors.Is(err, ErrUserBanned) {
			// recordLoginFailure() // UKLONJENO - nema metrics
			authLogger.Warn("Login failed: user is banned " + loginReq.Username)
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			// recordLoginFailure() // UKLONJENO - nema metrics
			authLogger.Error("Login failed", err)
			http.Error(w, "error authenticating user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Record successful login
	// recordLoginSuccess() // UKLONJENO - nema metrics
	authLogger.InfoWithFields("User logged in successfully", map[string]interface{}{
		"username": loginReq.Username,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JWTResponse{Token: tokenString})
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userRole := r.Header.Get("x-user-role")
	if userRole != "admin" {
		http.Error(w, "unauthorized: Only admins can view all users", http.StatusForbidden)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "error retrieving users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userRole := r.Header.Get("x-user-role")
	if userRole != "admin" {
		http.Error(w, "unauthorized: Only admins can block users", http.StatusForbidden)
		return
	}

	var blockReq BlockUserRequest
	err := json.NewDecoder(r.Body).Decode(&blockReq)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(blockReq); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.BlockUser(blockReq.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "error toggling user block status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Auth Service"})
}
