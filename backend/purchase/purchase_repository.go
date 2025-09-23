package main

import (
	"gorm.io/gorm"
)

type PurchaseRepository struct {
	database *Database
}

func NewPurchaseRepository(db *Database) *PurchaseRepository {
	return &PurchaseRepository{database: db}
}

func (r *PurchaseRepository) CreatePurchaseToken(token *TourPurchaseToken) error {
	result := r.database.db.Create(token)
	return result.Error
}

func (r *PurchaseRepository) GetTokenByID(tokenStr string) (*TourPurchaseToken, error) {
	var token TourPurchaseToken
	result := r.database.db.Where("token = ?", tokenStr).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrTokenNotFound
		}
		return nil, result.Error
	}
	return &token, nil
}

func (r *PurchaseRepository) GetTokensByUserID(userID string) ([]TourPurchaseToken, error) {
	var tokens []TourPurchaseToken
	result := r.database.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
}

func (r *PurchaseRepository) GetTokenByUserAndTour(userID string, tourID uint) (*TourPurchaseToken, error) {
	var token TourPurchaseToken
	result := r.database.db.Where("user_id = ? AND tour_id = ? AND status = 'active'", userID, tourID).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrTokenNotFound
		}
		return nil, result.Error
	}
	return &token, nil
}

func (r *PurchaseRepository) UpdateTokenStatus(tokenStr string, status string) error {
	result := r.database.db.Model(&TourPurchaseToken{}).Where("token = ?", tokenStr).Update("status", status)
	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}
	return result.Error
}