package main

import (
	"time"
)

type ShoppingCart struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    string      `json:"user_id" gorm:"not null;index"`
	Items     []OrderItem `json:"items" gorm:"foreignKey:CartID"`
	Total     float64     `json:"total" gorm:"default:0"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null"`
	TourID    uint      `json:"tour_id" gorm:"not null"`
	TourName  string    `json:"tour_name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

type TourPurchaseToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	TourID    uint      `json:"tour_id" gorm:"not null"`
	TourName  string    `json:"tour_name" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Status    string    `json:"status" gorm:"default:'active'"` // active, expired, used
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Methods for ShoppingCart
func (cart *ShoppingCart) CalculateTotal() {
	total := 0.0
	for _, item := range cart.Items {
		total += item.Price
	}
	cart.Total = total
}

func (cart *ShoppingCart) HasTour(tourID uint) bool {
	for _, item := range cart.Items {
		if item.TourID == tourID {
			return true
		}
	}
	return false
}

// Methods for TourPurchaseToken
func (token *TourPurchaseToken) IsValid() bool {
	return token.Status == "active" && time.Now().Before(token.ExpiresAt)
}

func (token *TourPurchaseToken) IsExpired() bool {
	return time.Now().After(token.ExpiresAt)
}