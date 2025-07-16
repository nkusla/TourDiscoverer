package seeder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"seed/data"
)

type AuthSeeder struct {
	*BaseSeeder
}

func NewAuthSeeder() *AuthSeeder {
	return &AuthSeeder{
		BaseSeeder: NewBaseSeeder("AuthService", "AUTH_SERVICE_URL"),
	}
}

func (s *AuthSeeder) Seed() {
	for _, user := range data.Users {
		if err := s.registerUser(user); err != nil {
			log.Printf("Error registering user %s: %v\n", user["username"], err)
		} else {
			log.Printf("User %s registered successfully.\n", user["username"])
		}
	}
}

func (s *AuthSeeder) registerUser(user map[string]interface{}) error {
	jsonData, _ := json.Marshal(user)

	resp, err := s.client.Post(
		s.serviceURL+"/register",
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
