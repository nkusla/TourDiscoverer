package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	database := &Database{}
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	repository := &FollowerRepository{db: database}
	service := &FollowerService{repository: repository}
	handler := &FollowerHandler{service: service}

	r.HandleFunc("/follow", handler.CreateFollowRelationship).Methods(http.MethodPost)
	r.HandleFunc("/unfollow", handler.DeleteFollowRelationship).Methods(http.MethodDelete)
	r.HandleFunc("/followers/{username}", handler.GetFollowers).Methods(http.MethodGet)
	r.HandleFunc("/following/{username}", handler.GetFollowing).Methods(http.MethodGet)

	r.HandleFunc("/internal/user", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/internal/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Follower service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
