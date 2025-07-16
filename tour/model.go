package main

import (
	"time"

	"gorm.io/gorm"
)

// Tour represents a tour created by an author
type Tour struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null" validate:"required"`
	Description string         `json:"description"`
	Difficulty  string         `json:"difficulty" validate:"required"` // težina
	Tags        string         `json:"tags"`                           // tagovi kao comma-separated string
	Status      string         `json:"status" gorm:"default:'draft'"`  // draft, published, etc.
	Price       float64        `json:"price" gorm:"default:0"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	KeyPoints   []KeyPoint     `json:"key_points" gorm:"foreignKey:TourID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// KeyPoint represents a key location point in a tour (ključna tačka)
type KeyPoint struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TourID      uint           `json:"tour_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null" validate:"required"`
	Description string         `json:"description"`
	Latitude    float64        `json:"latitude" gorm:"not null" validate:"required"` // geografska širina
	Longitude   float64        `json:"longitude" gorm:"not null" validate:"required"` // geografska dužina
	ImageURL    string         `json:"image_url"`                                     // slika
	Order       int            `json:"order" gorm:"default:0"`                       // redosled u turi
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Tour status constants
const (
	TourStatusDraft     = "draft"
	TourStatusPublished = "published"
	TourStatusArchived  = "archived"
)

// Difficulty constants
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)
