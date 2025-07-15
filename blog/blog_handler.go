package main

import (
    "encoding/json"
    "net/http"
)

type BlogHandler struct {
    service *BlogService
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
    var blog Blog
    if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    if err := h.service.CreateBlog(&blog); err != nil {
        http.Error(w, "failed to create blog", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}
func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
    blogs, err := h.service.GetAllBlogs()
    if err != nil {
        http.Error(w, "failed to fetch blogs", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blogs)
}