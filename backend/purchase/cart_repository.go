package main

import (
	"gorm.io/gorm"
)

type CartRepository struct {
	database *Database
}

func NewCartRepository(db *Database) *CartRepository {
	return &CartRepository{database: db}
}

func (r *CartRepository) GetCartByUserID(userID string) (*ShoppingCart, error) {
	var cart ShoppingCart
	result := r.database.db.Preload("Items").Where("user_id = ?", userID).First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrCartNotFound
		}
		return nil, result.Error
	}
	return &cart, nil
}

func (r *CartRepository) GetCartByID(cartID uint) (*ShoppingCart, error) {
	var cart ShoppingCart
	result := r.database.db.Preload("Items").Where("id = ?", cartID).First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrCartNotFound
		}
		return nil, result.Error
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(userID string) (*ShoppingCart, error) {
	cart := &ShoppingCart{
		UserID: userID,
		Total:  0,
	}
	
	result := r.database.db.Create(cart)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return cart, nil
}

func (r *CartRepository) GetOrCreateCart(userID string) (*ShoppingCart, error) {
	cart, err := r.GetCartByUserID(userID)
	if err == ErrCartNotFound {
		return r.CreateCart(userID)
	}
	return cart, err
}

func (r *CartRepository) AddItemToCart(cartID uint, item *OrderItem) error {
	item.CartID = cartID
	result := r.database.db.Create(item)
	return result.Error
}

func (r *CartRepository) RemoveItemFromCart(cartID uint, tourID uint) error {
	result := r.database.db.Where("cart_id = ? AND tour_id = ?", cartID, tourID).Delete(&OrderItem{})
	if result.RowsAffected == 0 {
		return ErrItemNotFound
	}
	return result.Error
}

func (r *CartRepository) UpdateCartTotal(cartID uint, total float64) error {
	result := r.database.db.Model(&ShoppingCart{}).Where("id = ?", cartID).Update("total", total)
	return result.Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	// Delete all items
	err := r.database.db.Where("cart_id = ?", cartID).Delete(&OrderItem{}).Error
	if err != nil {
		return err
	}
	
	// Reset total to 0
	return r.UpdateCartTotal(cartID, 0)
}