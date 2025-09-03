package main

import "time"

type CreateReviewRequest struct {
	TourID    uint     `json:"tour_id" validate:"required"`
	Rating    int      `json:"rating" validate:"required,min=1,max=5"`
	Comment   string   `json:"comment" validate:"required"`
	VisitDate string   `json:"visit_date" validate:"required"` // Format: "2006-01-02"
	Images    []string `json:"images"`                         // Array of image URLs
}

type UpdateReviewRequest struct {
	Rating  *int     `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string  `json:"comment,omitempty"`
	Images  []string `json:"images,omitempty"` // Will replace all existing images
}

type ReviewResponse struct {
	ID              uint                `json:"id"`
	TourID          uint                `json:"tour_id"`
	TouristUsername string              `json:"tourist_username"`
	Rating          int                 `json:"rating"`
	Comment         string              `json:"comment"`
	VisitDate       time.Time           `json:"visit_date"`
	ReviewDate      time.Time           `json:"review_date"`
	Images          []ReviewImageResponse `json:"images"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type ReviewImageResponse struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
}

type ReviewListResponse struct {
	Reviews      []ReviewResponse `json:"reviews"`
	TotalCount   int64            `json:"total_count"`
	AverageRating float64         `json:"average_rating"`
}

// ToResponse converts a Review model to ReviewResponse
func (r *Review) ToResponse() ReviewResponse {
	images := make([]ReviewImageResponse, len(r.Images))
	for i, img := range r.Images {
		images[i] = ReviewImageResponse{
			ID:       img.ID,
			ImageURL: img.ImageURL,
		}
	}

	return ReviewResponse{
		ID:              r.ID,
		TourID:          r.TourID,
		TouristUsername: r.TouristUsername,
		Rating:          r.Rating,
		Comment:         r.Comment,
		VisitDate:       r.VisitDate,
		ReviewDate:      r.ReviewDate,
		Images:          images,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
