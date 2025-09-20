package main

import (
	"encoding/json"
	"net/http"
)

type CommentHandler struct {
	service *CommentService
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Uzmi username iz header-a koji postavlja API Gateway
	username := r.Header.Get("x-username")
	if username == "" {
		http.Error(w, "username not found in headers", http.StatusBadRequest)
		return
	}

	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Postavi author na osnovu JWT token-a
	comment.Author = username

	if err := h.service.CreateComment(&comment); err != nil {
		http.Error(w, "failed to create comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	blogID := r.URL.Query().Get("blog_id")
	if blogID == "" {
		http.Error(w, "blog_id is required", http.StatusBadRequest)
		return
	}
	comments, err := h.service.GetComments(blogID)
	if err != nil {
		http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
