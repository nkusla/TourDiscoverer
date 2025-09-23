package main

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseService struct {
	purchaseRepository *PurchaseRepository
	cartRepository     *CartRepository
}

func NewPurchaseService(purchaseRepo *PurchaseRepository, cartRepo *CartRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepository: purchaseRepo,
		cartRepository:     cartRepo,
	}
}

func (s *PurchaseService) Checkout(userID string) ([]TourPurchaseToken, error) {
	// Get user's cart
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}

	var tokens []TourPurchaseToken

	// Create purchase token for each item
	for _, item := range cart.Items {
		token := TourPurchaseToken{
			UserID:    userID,
			TourID:    item.TourID,
			TourName:  item.TourName,
			Token:     s.generateToken(),
			Status:    "active",
			ExpiresAt: time.Now().Add(365 * 24 * time.Hour), // 1 year expiration
		}

		err := s.purchaseRepository.CreatePurchaseToken(&token)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	// Clear the cart after successful checkout
	err = s.cartRepository.ClearCart(cart.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *PurchaseService) GetUserTokens(userID string) ([]TourPurchaseToken, error) {
	return s.purchaseRepository.GetTokensByUserID(userID)
}

func (s *PurchaseService) GetTokenDetails(tokenStr string) (*TourPurchaseToken, error) {
	token, err := s.purchaseRepository.GetTokenByID(tokenStr)
	if err != nil {
		return nil, err
	}

	if token.IsExpired() {
		// Update token status to expired
		s.purchaseRepository.UpdateTokenStatus(tokenStr, "expired")
		return nil, ErrTokenExpired
	}

	if !token.IsValid() {
		return nil, ErrTokenInvalid
	}

	return token, nil
}

func (s *PurchaseService) ValidateAccess(userID string, tourID uint) (*TourPurchaseToken, error) {
	token, err := s.purchaseRepository.GetTokenByUserAndTour(userID, tourID)
	if err != nil {
		return nil, err
	}

	if token.IsExpired() {
		s.purchaseRepository.UpdateTokenStatus(token.Token, "expired")
		return nil, ErrTokenExpired
	}

	if !token.IsValid() {
		return nil, ErrTokenInvalid
	}

	return token, nil
}

func (s *PurchaseService) generateToken() string {
	return uuid.New().String()
}