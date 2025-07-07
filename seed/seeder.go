package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type SeederService struct {
	authURL     string
	followerURL string
	client      *http.Client
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewSeederService() *SeederService {
	return &SeederService{
		authURL:     getEnvOrDefault("AUTH_SERVICE_URL", ""),
		followerURL: getEnvOrDefault("FOLLOWER_SERVICE_URL", ""),
		client:      &http.Client{},
	}
}

func (s *SeederService) waitForServices() error {
	services := map[string]string{
		"auth":     s.authURL + "/internal/ping",
		"follower": s.followerURL + "/internal/ping",
	}

	for name, url := range services {
		for i := 0; i < 15; i++ { // Wait up to 10 seconds
			resp, err := s.client.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				log.Printf("%s service is ready", name)
				break
			}
			if resp != nil {
				resp.Body.Close()
			}

			if i == 29 {
				return fmt.Errorf("%s service not ready after 15 seconds", name)
			}
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func (s *SeederService) SeedAll() {
	if err := s.waitForServices(); err != nil {
		log.Printf("Error waiting for services: %v", err)
		return
	}

	s.seedUsers()
	s.seedFollowers()
}

func (s *SeederService) seedUsers() {
	for _, user := range users {
		if err := s.registerUser(user); err != nil {
			log.Printf("Error registering user %s: %v\n", user["username"], err)
		} else {
			log.Printf("User %s registered successfully.\n", user["username"])
		}
	}
}

func (s *SeederService) registerUser(user map[string]interface{}) error {
	jsonData, _ := json.Marshal(user)

	resp, err := s.client.Post(
		s.authURL+"/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register user %s, status code: %d", user["username"], resp.StatusCode)
	}

	return nil
}

func (s *SeederService) seedFollowers() {
	for _, follower := range followers {
		if err := s.createFollower(follower); err != nil {
			log.Printf("Error creating follower %s: %v\n", follower["username"], err)
		} else {
			log.Printf("Follower %s created successfully.\n", follower["username"])
		}
	}
}

func (s *SeederService) createFollower(follower map[string]interface{}) error {
	jsonData, _ := json.Marshal(follower)

	resp, err := s.client.Post(
		s.followerURL+"/follow",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create follower %s, status code: %d", follower["username"], resp.StatusCode)
	}

	return nil
}
