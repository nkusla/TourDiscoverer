package main

import "time"

type CreateTourRequest struct {
	Name        string                  `json:"name" validate:"required"`
	Description string                  `json:"description"`
	Difficulty  string                  `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Tags        string                  `json:"tags"`
	KeyPoints   []CreateKeyPointRequest `json:"key_points"`
	Distance    float64                 `json:"distance"`
	Status      string                  `json:"status"`
	Price       float64                 `json:"price"`
}

type UpdateTourRequest struct {
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	Difficulty       string                  `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Tags             string                  `json:"tags"`
	TransportDetails []Transport             `json:"transport_details"`
	Price            float64                 `json:"price" validate:"omitempty,gt=0"`
	KeyPoints        []CreateKeyPointRequest `json:"key_points"`
	Distance         float64                 `json:"distance"`
	Status           string                  `json:"status"`
}

type CreateKeyPointRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	ImageURL    string  `json:"image_url"`
	Order       int     `json:"order"`
}

type CreateTourResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Difficulty     string     `json:"difficulty"`
	Tags           string     `json:"tags"`
	Status         string     `json:"status"`
	Price          float64    `json:"price"`
	Distance       float64    `json:"distance"`
	KeyPoints      []KeyPoint `json:"key_points"`
	AuthorUsername string     `json:"author_username"`
	Message        string     `json:"message"`
}

type CreateKeyPointResponse struct {
	ID          uint    `json:"id"`
	TourID      uint    `json:"tour_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"image_url"`
	Order       int     `json:"order"`
	Message     string  `json:"message"`
}

type TransportDetailsRequest struct {
	Duration      uint   `json:"duration" validate:"required"`
	TransportType string `json:"transport_type" validate:"required oneof=walking driving biking"`
}

type GetToursResponse struct {
	Tours []Tour `json:"tours"`
	Count int    `json:"count"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

// TourExecution related requests
type StartTourExecutionRequest struct {
	TourID    uint    `json:"tour_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type StartTourExecutionResponse struct {
	ID              uint      `json:"id"`
	TourID          uint      `json:"tour_id"`
	TouristUsername string    `json:"tourist_username"`
	Status          string    `json:"status"`
	StartTime       time.Time `json:"start_time"`
	StartLatitude   float64   `json:"start_latitude"`
	StartLongitude  float64   `json:"start_longitude"`
	Message         string    `json:"message"`
}

type EndTourExecutionRequest struct {
	Status string `json:"status" validate:"required,oneof=completed abandoned"`
}

type EndTourExecutionResponse struct {
	ID              uint       `json:"id"`
	TourID          uint       `json:"tour_id"`
	TouristUsername string     `json:"tourist_username"`
	Status          string     `json:"status"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Message         string     `json:"message"`
}

type CheckProximityRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type CheckProximityResponse struct {
	KeyPointReached    bool                `json:"key_point_reached"`
	KeyPoint           *KeyPoint           `json:"key_point,omitempty"`
	KeyPointCompletion *KeyPointCompletion `json:"key_point_completion,omitempty"`
	LastActivity       time.Time           `json:"last_activity"`
	Message            string              `json:"message"`
}
