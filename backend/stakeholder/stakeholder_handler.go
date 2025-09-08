package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

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
	var profilePictureData string

	// Check if it's multipart form data (file upload) or JSON
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle file upload
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		// Get form values
		req.FirstName = r.FormValue("first_name")
		req.LastName = r.FormValue("last_name")
		req.Biography = r.FormValue("biography")
		req.Motto = r.FormValue("motto")

		// Handle file upload if present
		file, header, err := r.FormFile("profile_picture")
		if err == nil {
			defer file.Close()
			
			// Read file content
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading uploaded file", http.StatusInternalServerError)
				return
			}
			
			// Convert to base64 with data URL format
			mimeType := header.Header.Get("Content-Type")
			if mimeType == "" {
				// Try to detect MIME type from file extension
				ext := strings.ToLower(filepath.Ext(header.Filename))
				switch ext {
				case ".jpg", ".jpeg":
					mimeType = "image/jpeg"
				case ".png":
					mimeType = "image/png"
				case ".gif":
					mimeType = "image/gif"
				case ".webp":
					mimeType = "image/webp"
				default:
					mimeType = "image/jpeg" // default
				}
			}
			
			profilePictureData = fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(fileBytes))
		} else {
			// No file uploaded, keep existing profile picture
			profilePictureData = r.FormValue("existing_profile_picture")
		}
		
		req.ProfilePicture = profilePictureData
	} else {
		// Handle JSON request (backward compatibility)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
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
