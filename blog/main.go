package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	client, err := mongo.Connect(nil, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog_db").Collection("blogs")
	repo := &BlogRepository{collection: collection}
	service := &BlogService{repository: repo}
	handler := &BlogHandler{service: service}

	http.HandleFunc("/blog", handler.CreateBlog)
	http.HandleFunc("/blogs", handler.GetAllBlogs)

	commentCollection := client.Database("blog_db").Collection("comments")
	commentRepo := &CommentRepository{collection: commentCollection}
	commentService := &CommentService{repository: commentRepo}
	commentHandler := &CommentHandler{service: commentService}

	http.HandleFunc("/comment", commentHandler.CreateComment)
	http.HandleFunc("/comments", commentHandler.GetComments)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	log.Printf("Blog service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
