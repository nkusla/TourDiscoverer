package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	if err := InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer CloseDatabase()

	repository := &FollowerRepository{driver: &Driver}
	service := &FollowerService{repository: repository}
	handler := &FollowerHandler{service: service}

	r.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/follow", handler.CreateFollowRelationship).Methods(http.MethodPost)
	r.HandleFunc("/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Follower service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
