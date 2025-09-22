package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var authLogger *AppLogger

func main() {
	authLogger = NewAppLogger("auth-service")
	defer authLogger.Close()

	authLogger.Info("Starting auth service initialization")

	r := mux.NewRouter().StrictSlash(true)
	database := InitDatabase()
	SeedAdmins(database)

	repository := &UserRepository{database: database}
	service := &UserService{repository: repository}
	handler := &UserHandler{service: service}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authLogger.HTTPMiddleware(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})(w, r)
		})
	})

	// Dodaj Prometheus middleware
	//r.Use(prometheusMiddleware)

	r.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	r.HandleFunc("/user", handler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/block", handler.BlockUser).Methods(http.MethodPost)
	r.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/internal/ping", handler.Ping).Methods(http.MethodGet)

	// Prometheus metrics uklonjen - sporo kompajliranje
	// r.Handle("/metrics", metricsHandler()).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authLogger.InfoWithFields("Auth service starting", map[string]interface{}{
		"port": port,
	})

	if err := http.ListenAndServe(":"+port, r); err != nil {
		authLogger.Error("Failed to start auth service", err)
	}
}
