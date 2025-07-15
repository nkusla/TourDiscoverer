package main

import (
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURL := os.Getenv("BLOG_DB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://blog-db:27017"
	}
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog_db").Collection("blogs")
	repo := &BlogRepository{collection: collection}
	service := &BlogService{repository: repo}
	handler := &BlogHandler{service: service}

	http.HandleFunc("/blog", handler.CreateBlog)
	http.HandleFunc("/blogs", handler.GetAllBlogs)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	log.Printf("Blog service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
