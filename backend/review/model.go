package main

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	TourID           uint           `json:"tour_id" gorm:"not null" validate:"required"`
	TouristUsername  string         `json:"tourist_username" gorm:"not null" validate:"required"`
	Rating           int            `json:"rating" gorm:"not null" validate:"required,min=1,max=5"`
	Comment          string         `json:"comment" validate:"required"`
	VisitDate        time.Time      `json:"visit_date" gorm:"not null" validate:"required"`
	ReviewDate       time.Time      `json:"review_date" gorm:"autoCreateTime"`
	Images           []ReviewImage  `json:"images" gorm:"foreignKey:ReviewID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

type ReviewImage struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	ReviewID uint   `json:"review_id" gorm:"not null"`
	ImageURL string `json:"image_url" gorm:"not null" validate:"required"`
}

// ValidateRating checks if the rating is within the valid range (1-5)
func (r *Review) ValidateRating() bool {
	return r.Rating >= 1 && r.Rating <= 5
}

// AddImage adds an image to the review
func (r *Review) AddImage(imageURL string) {
	image := ReviewImage{
		ReviewID: r.ID,
		ImageURL: imageURL,
	}
	r.Images = append(r.Images, image)
}
