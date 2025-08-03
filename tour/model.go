package main

import (
	"time"

	"gorm.io/gorm"
)

type Tour struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"not null" validate:"required"`
	Description      string         `json:"description" validate:"required"`
	Difficulty       string         `json:"difficulty" validate:"required"`
	Tags             string         `json:"tags" validate:"required"`
	Status           string         `json:"status" gorm:"default:'draft'"`
	Price            float64        `json:"price" gorm:"default:0"`
	TransportDetails []Transport    `json:"transport_details" gorm:"type:jsonb;serializer:json"`
	Distance         float64        `json:"distance" gorm:"default:0"`
	AuthorUsername   string         `json:"author_username" gorm:"not null"`
	KeyPoints        []KeyPoint     `json:"key_points" gorm:"foreignKey:TourID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

type Transport struct {
	Duration      uint   `json:"duration"`
	TransportType string `json:"transport_type"`
}

type KeyPoint struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TourID      uint           `json:"tour_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null" validate:"required"`
	Description string         `json:"description"`
	Latitude    float64        `json:"latitude" gorm:"not null" validate:"required"`
	Longitude   float64        `json:"longitude" gorm:"not null" validate:"required"`
	ImageURL    string         `json:"image_url"`
	Order       int            `json:"order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (t *Tour) CanBePublished() bool {
	if t.Name == "" || t.Description == "" || t.Difficulty == "" || t.Tags == "" {
		return false
	}

	if t.Status != TourStatusDraft {
		return false
	}

	if len(t.KeyPoints) < 2 {
		return false
	}

	if len(t.TransportDetails) < 1 {
		return false
	}

	return true
}

func (t *Tour) CanBeArchived() bool {
	return t.Status == TourStatusPublished
}

func (t *Tour) CanBeUnarchived() bool {
	return t.Status == TourStatusArchived
}
