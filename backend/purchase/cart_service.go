package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CartService struct {
	cartRepository *CartRepository
}

func NewCartService(cartRepo *CartRepository) *CartService {
	return &CartService{
		cartRepository: cartRepo,
	}
}

func (s *CartService) GetCart(userID string) (*ShoppingCart, error) {
	return s.cartRepository.GetOrCreateCart(userID)
}

func (s *CartService) AddToCart(userID string, tourID uint) error {
	log.Printf("AddToCart: userID=%s, tourID=%d", userID, tourID)
	
	// Get or create cart
	cart, err := s.cartRepository.GetOrCreateCart(userID)
	if err != nil {
		log.Printf("AddToCart: failed to get/create cart: %v", err)
		return err
	}

	// Check if tour is already in cart
	if cart.HasTour(tourID) {
		log.Printf("AddToCart: tour %d already in cart", tourID)
		return ErrTourAlreadyInCart
	}

	// Fetch tour details from tour service
	tourInfo, err := s.fetchTourInfo(tourID)
	if err != nil {
		log.Printf("AddToCart: failed to fetch tour info: %v", err)
		return err
	}

	log.Printf("AddToCart: fetched tour info - ID=%d, Name=%s, Status=%s, Price=%f", 
		tourInfo.ID, tourInfo.Name, tourInfo.Status, tourInfo.Price)

	// Validate tour status
	if tourInfo.Status != "published" {
		log.Printf("AddToCart: tour %d is not published (status: %s)", tourID, tourInfo.Status)
		return ErrTourNotPublished
	}

	// Create order item
	item := &OrderItem{
		TourID:   tourID,
		TourName: tourInfo.Name,
		Price:    tourInfo.Price,
	}

	// Add item to cart
	err = s.cartRepository.AddItemToCart(cart.ID, item)
	if err != nil {
		return err
	}

	// Recalculate total
	return s.recalculateCartTotal(cart.ID)
}

func (s *CartService) RemoveFromCart(userID string, tourID uint) error {
	// Get cart
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	// Remove item
	err = s.cartRepository.RemoveItemFromCart(cart.ID, tourID)
	if err != nil {
		return err
	}

	// Recalculate total
	return s.recalculateCartTotal(cart.ID)
}

func (s *CartService) ClearCart(userID string) error {
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	return s.cartRepository.ClearCart(cart.ID)
}

func (s *CartService) recalculateCartTotal(cartID uint) error {
	cart, err := s.cartRepository.GetCartByID(cartID)
	if err != nil {
		return err
	}

	cart.CalculateTotal()
	return s.cartRepository.UpdateCartTotal(cartID, cart.Total)
}

// Helper function to fetch tour info from tour service
func (s *CartService) fetchTourInfo(tourID uint) (*TourInfo, error) {
	tourServiceURL := GetEnv("TOUR_SERVICE_URL", "http://tour-service:3006")
	url := fmt.Sprintf("%s/%d", tourServiceURL, tourID)
	
	log.Printf("fetchTourInfo: requesting URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tour info: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("fetchTourInfo: tour service response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tour service returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tour response: %v", err)
	}

	log.Printf("fetchTourInfo: response body: %s", string(body))

	var tourInfo TourInfo
	err = json.Unmarshal(body, &tourInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tour info: %v", err)
	}

	return &tourInfo, nil
}

// TourInfo represents basic tour information from tour service
type TourInfo struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}