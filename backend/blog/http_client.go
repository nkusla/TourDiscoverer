package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

type User struct {
	Username string `json:"username"`
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPClient) GetFollowingUsers(username string) ([]User, error) {
	followerServiceURL := os.Getenv("FOLLOWER_SERVICE_URL")
	if followerServiceURL == "" {
		// Fallback za development
		followerServiceURL = "http://localhost:8082"
	}

	url := fmt.Sprintf("%s/user/%s/following", followerServiceURL, username)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get following users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Ako korisnik ne postoji ili ne prati nikoga, vrati praznu listu
		return []User{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("follower service returned status %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return users, nil
}
