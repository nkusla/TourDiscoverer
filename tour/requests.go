package main

type CreateTourRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Tags        string `json:"tags"`
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
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Difficulty     string  `json:"difficulty"`
	Tags           string  `json:"tags"`
	Status         string  `json:"status"`
	Price          float64 `json:"price"`
	AuthorUsername string  `json:"author_username"`
	Message        string  `json:"message"`
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
