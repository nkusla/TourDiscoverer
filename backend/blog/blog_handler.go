package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type BlogHandler struct {
	service *BlogService
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	// Uzmi username iz header-a koji postavlja API Gateway
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "username not found in headers", http.StatusBadRequest)
		return
	}

	var blog Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Postavi author na osnovu JWT token-a
	blog.Author = username

	if err := h.service.CreateBlog(&blog); err != nil {
		http.Error(w, "failed to create blog", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	// Uzmi username iz header-a ako postoji (opciono za javni pristup)
	username := r.Header.Get("x-username")
	
	blogs, err := h.service.GetAllBlogs(username)
	if err != nil {
		http.Error(w, "failed to fetch blogs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Uzmi username iz header-a koji postavlja API Gateway
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "username not found in headers", http.StatusBadRequest)
		return
	}

	var likeReq LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&likeReq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if likeReq.BlogID == "" {
		http.Error(w, "blog_id is required", http.StatusBadRequest)
		return
	}

	isLiked, err := h.service.ToggleLike(likeReq.BlogID, username)
	if err != nil {
		http.Error(w, "failed to toggle like", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"liked": isLiked,
		"message": func() string {
			if isLiked {
				return "Blog liked successfully"
			}
			return "Blog like removed successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BlogHandler) GetLikeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Uzmi username iz header-a koji postavlja API Gateway
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "username not found in headers", http.StatusBadRequest)
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	blogID := query.Get("blog_id")

	if blogID == "" {
		http.Error(w, "blog_id query parameter is required", http.StatusBadRequest)
		return
	}

	isLiked, likeCount, err := h.service.GetLikeStatus(blogID, username)
	if err != nil {
		http.Error(w, "failed to get like status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"liked":      isLiked,
		"like_count": likeCount,
		"username":   username, // dodato za debugging
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
