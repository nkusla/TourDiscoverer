package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbPort := os.Getenv("BLOG_DB_PORT")
	dbHost := os.Getenv("BLOG_DB_HOST")
	dbUsername := os.Getenv("BLOG_DB_USER")
	dbPassword := os.Getenv("BLOG_DB_PASSWORD")
	dbName := os.Getenv("BLOG_DB_NAME")
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", dbUsername, dbPassword, dbHost, dbPort, dbName)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog_db").Collection("blogs")
	repo := &BlogRepository{collection: collection}
	service := &BlogService{repository: repo}
	handler := &BlogHandler{service: service}

	// Create Gorilla Mux router
	r := mux.NewRouter().StrictSlash(true)

	// Blog routes
	r.HandleFunc("/", handler.CreateBlog).Methods(http.MethodPost)
	r.HandleFunc("/", handler.GetAllBlogs).Methods(http.MethodGet)
	r.HandleFunc("/like", handler.ToggleLike).Methods(http.MethodPost)
	r.HandleFunc("/like-status", handler.GetLikeStatus).Methods(http.MethodGet)

	// Comment routes
	commentCollection := client.Database("blog_db").Collection("comments")
	commentRepo := &CommentRepository{collection: commentCollection}
	commentService := &CommentService{repository: commentRepo}
	commentHandler := &CommentHandler{service: commentService}

	r.HandleFunc("/comment", commentHandler.CreateComment).Methods(http.MethodPost)
	r.HandleFunc("/comments", commentHandler.GetComments).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	log.Printf("Blog service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
