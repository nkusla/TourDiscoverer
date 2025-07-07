package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *UserRepository
}

func (s *UserService) RegisterUser(req RegisterRequest) error {
	if req.Role != RoleGuide && req.Role != RoleTourist {
		return ErrInvalidRole
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     req.Role,
	}

	err = s.repository.Create(user)
	if err != nil {
		return err
	}

	err = s.registerUserInFollowerService(user.Username)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return nil
}

func (s *UserService) registerUserInFollowerService(username string) error {
	followerServiceURL := os.Getenv("FOLLOWER_SERVICE_URL")
	if followerServiceURL == "" {
		return fmt.Errorf("follower service URL is not configured")
	}

	payload := map[string]string{
		"username": username,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		followerServiceURL+"/internal/user",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register user in follower service, status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if user.IsBanned {
		return "", ErrUserBanned
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := CreateJWT(user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetAllUsers() ([]*User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) BlockUser(username string) error {
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return ErrUserNotFound
	}

	user.IsBanned = !user.IsBanned
	err = s.repository.Update(user)
	if err != nil {
		return err
	}

	return nil
}
