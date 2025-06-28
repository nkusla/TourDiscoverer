package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repository *UserRepository
}

var validate = validator.New()

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

	if err := validate.Struct(registerReq); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	if registerReq.Role != RoleGuide && registerReq.Role != RoleTourist {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	user := User{
		Username: registerReq.Username,
		Password: string(hashedPassword),
		Email:    registerReq.Email,
		Role:     registerReq.Role,
	}

	err = h.repository.Create(&user)
	if err != nil {
		if err == ErrUserAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

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

	if err := validate.Struct(loginReq); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repository.FindByUsername(loginReq.Username)
	if err != nil {
		http.Error(w, "Username or password is incorrect", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Username or password is incorrect", http.StatusUnauthorized)
		return
	}

	tokenString, err := CreateJWT(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JWTResponse{Token: tokenString})
}

func (h *AuthHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.repository.FindAll()
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PingResponse{Message: "pong", Service: "Auth Service"})
}
