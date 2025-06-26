package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	handler := &AuthHandler{}

	r.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	r.HandleFunc("/ping", handler.Ping).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Auth service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
