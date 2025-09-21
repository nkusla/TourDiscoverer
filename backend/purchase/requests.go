package main

type AddToCartRequest struct {
	TourID uint `json:"tour_id" validate:"required"`
}

type CheckoutResponse struct {
	Tokens  []TourPurchaseToken `json:"tokens"`
	Message string              `json:"message"`
}

type CartResponse struct {
	Cart    *ShoppingCart `json:"cart"`
	Message string        `json:"message"`
}

type TokenResponse struct {
	Token   *TourPurchaseToken `json:"token"`
	Message string             `json:"message"`
}

type TokensResponse struct {
	Tokens  []TourPurchaseToken `json:"tokens"`
	Message string              `json:"message"`
}