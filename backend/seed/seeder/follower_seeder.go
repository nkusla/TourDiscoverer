package seeder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"seed/data"
)

type FollowerSeeder struct {
	*BaseSeeder
}

func NewFollowerSeeder() *FollowerSeeder {
	return &FollowerSeeder{
		BaseSeeder: NewBaseSeeder("FollowerService", "FOLLOWER_SERVICE_URL"),
	}
}

func (s *FollowerSeeder) Seed() {
	for _, follower := range data.Followers {
		if err := s.createFollower(follower); err != nil {
			log.Printf("Error creating follower %s: %v\n", follower["username"], err)
		} else {
			log.Printf("Follower %s created successfully.\n", follower["username"])
		}
	}
}

func (s *FollowerSeeder) createFollower(follower map[string]interface{}) error {
	jsonData, _ := json.Marshal(follower)

	resp, err := s.client.Post(
		s.serviceURL+"/follow",
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
