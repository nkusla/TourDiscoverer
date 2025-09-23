package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type TourService struct {
	repository *TourRepository
}

func (service *TourService) CreateTour(request *CreateTourRequest, authorUsername string) (*Tour, error) {
	if !isValidDifficulty(request.Difficulty) {
		return nil, errors.New("invalid difficulty level")
	}

	// Convert key points from request to model
	var keyPoints []KeyPoint
	for i, kpRequest := range request.KeyPoints {
		keyPoint := KeyPoint{
			Name:        kpRequest.Name,
			Description: kpRequest.Description,
			Latitude:    kpRequest.Latitude,
			Longitude:   kpRequest.Longitude,
			ImageURL:    kpRequest.ImageURL,
			Order:       i, // Use index as order
		}
		keyPoints = append(keyPoints, keyPoint)
	}

	tour := &Tour{
		Name:             request.Name,
		Description:      request.Description,
		Difficulty:       request.Difficulty,
		Tags:             strings.TrimSpace(request.Tags),
		Status:           TourStatusDraft,
		Price:            0, // Always 0 for draft
		AuthorUsername:   authorUsername,
		TransportDetails: request.TransportDetails,
		Distance:         request.Distance,
		KeyPoints:        keyPoints, // GORM will handle the association
	}

	// Create tour with all associations in one transaction
	err := service.repository.CreateTour(tour)
	if err != nil {
		return nil, err
	}

	return tour, nil
}

func (service *TourService) UpdateTour(id uint, request *UpdateTourRequest, authorUsername string) (*Tour, error) {
	// Start a transaction
	tx := service.repository.database.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the existing tour
	var tour Tour
	result := tx.Preload("KeyPoints").Where("id = ?", id).First(&tour)
	if result.Error != nil {
		tx.Rollback()
		return nil, ErrTourNotFound
	}

	if tour.AuthorUsername != authorUsername {
		tx.Rollback()
		return nil, ErrUnauthorized
	}

	if tour.Status != TourStatusDraft {
		tx.Rollback()
		return nil, ErrTourNotEditable
	}

	// Update basic tour fields
	tour.Name = request.Name
	tour.Description = request.Description
	tour.Difficulty = request.Difficulty
	tour.Tags = strings.TrimSpace(request.Tags)
	tour.Price = request.Price
	tour.Distance = request.Distance

	if request.TransportDetails != nil {
		tour.TransportDetails = request.TransportDetails
	} else {
		tour.TransportDetails = []Transport{}
	}

	// Handle KeyPoints association properly
	if request.KeyPoints != nil {
		// Step 1: Delete all existing key points for this tour (hard delete)
		result := tx.Unscoped().Where("tour_id = ?", id).Delete(&KeyPoint{})
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}

		// Step 2: Clear the KeyPoints slice in the tour struct
		tour.KeyPoints = []KeyPoint{}

		// Step 3: Create new key points from the request
		newKeyPoints := make([]KeyPoint, 0, len(request.KeyPoints))
		for _, kpRequest := range request.KeyPoints {
			keyPoint := KeyPoint{
				Name:        kpRequest.Name,
				Description: kpRequest.Description,
				Latitude:    kpRequest.Latitude,
				Longitude:   kpRequest.Longitude,
				ImageURL:    kpRequest.ImageURL,
				Order:       kpRequest.Order,
				TourID:      id,
			}

			// Create each key point in the transaction
			result := tx.Create(&keyPoint)
			if result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}

			newKeyPoints = append(newKeyPoints, keyPoint)
		}

		// Step 4: Update the tour's KeyPoints slice with the new ones
		tour.KeyPoints = newKeyPoints
	}

	// Update the tour itself (without trying to save associations again)
	// Note: We don't include transport_details in Select/Updates to avoid JSONB serialization issues
	// GORM will handle it properly when we save the entire model
	result = tx.Model(&tour).Select("name", "description", "difficulty", "tags", "price", "distance").Updates(map[string]interface{}{
		"name":        tour.Name,
		"description": tour.Description,
		"difficulty":  tour.Difficulty,
		"tags":        tour.Tags,
		"price":       tour.Price,
		"distance":    tour.Distance,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Update transport_details separately using Save to ensure proper JSONB serialization
	result = tx.Model(&tour).Select("transport_details").Updates(Tour{TransportDetails: tour.TransportDetails})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Fetch the updated tour with key points (fresh from database)
	updatedTour, err := service.repository.GetTourByID(id)
	if err != nil {
		return nil, err
	}

	return updatedTour, nil
}

func (service *TourService) GetToursByAuthor(authorUsername string) ([]Tour, error) {
	tours, err := service.repository.GetToursByAuthor(authorUsername)
	if err != nil {
		return nil, err
	}
	return tours, nil
}

func (service *TourService) GetAllPublishedTours() ([]Tour, error) {
	tours, err := service.repository.GetAllPublishedTours()
	if err != nil {
		return nil, err
	}
	return tours, nil
}

func (service *TourService) GetTourByID(id uint) (*Tour, error) {
	return service.repository.GetTourByID(id)
}

func (service *TourService) CreateKeyPoint(request *CreateKeyPointRequest, tourID uint, authorUsername string) (*KeyPoint, error) {
	tour, err := service.repository.GetTourByID(tourID)
	if err != nil {
		return nil, ErrTourNotFound
	}

	if tour.AuthorUsername != authorUsername {
		return nil, ErrUnauthorized
	}

	if tour.Status != TourStatusDraft {
		return nil, ErrTourNotEditable
	}

	keyPoint := &KeyPoint{
		TourID:      tour.ID,
		Name:        request.Name,
		Description: request.Description,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		ImageURL:    request.ImageURL,
		Order:       request.Order,
	}

	tour.AddKeyPoint(keyPoint)
	err = service.repository.UpdateTour(tour)
	if err != nil {
		return nil, err
	}

	return keyPoint, nil
}

func (service *TourService) PublishTour(tourID uint, authorUsername string) error {
	tour, err := service.repository.GetTourByID(tourID)
	if err != nil {
		return ErrTourNotFound
	}

	if tour.AuthorUsername != authorUsername {
		return ErrUnauthorized
	}

	if !tour.CanBePublished() {
		return ErrTourNotPublishable
	}

	tour.Status = TourStatusPublished
	err = service.repository.UpdateTour(tour)
	if err != nil {
		return err
	}

	blogPayload := map[string]interface{}{
		"title":       tour.Name,
		"description": tour.Description,
		"author":      tour.AuthorUsername,
		"tour_id":     tour.ID,
		"tags":        tour.Tags,
	}

	importBytes, _ := json.Marshal(blogPayload)

	blogHost := os.Getenv("BLOG_SERVICE_HOST")
	blogPort := os.Getenv("BLOG_SERVICE_PORT")
	if blogHost == "" {
		blogHost = "blog-service"
	}
	if blogPort == "" {
		blogPort = "3002"
	}
	blogServiceURL := fmt.Sprintf("http://%s:%s/", blogHost, blogPort)

	req, err := http.NewRequest(http.MethodPost, blogServiceURL, strings.NewReader(string(importBytes)))
	if err != nil {
		println("SAGA: Failed to create HTTP request:", err.Error())
		tour.Status = TourStatusDraft
		_ = service.repository.UpdateTour(tour)
		return errors.New("failed to create blog post for published tour (SAGA rollback)")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-username", tour.AuthorUsername)
	req.Header.Set("x-user-role", RoleGuide)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("SAGA: Network error when calling Blog service:", err.Error())
	}
	if resp != nil {
		println("SAGA: Blog service response status:", resp.StatusCode)
		defer resp.Body.Close()
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		println("SAGA: Blog service response body:", string(bodyBytes[:n]))
	}
	if err != nil || (resp != nil && resp.StatusCode >= 300) {
		tour.Status = TourStatusDraft
		_ = service.repository.UpdateTour(tour)
		return errors.New("failed to create blog post for published tour (SAGA rollback)")
	}

	return nil
}

func (service *TourService) ArchiveTour(tourID uint, authorUsername string) error {
	tour, err := service.repository.GetTourByID(tourID)
	if err != nil {
		return ErrTourNotFound
	}

	if tour.AuthorUsername != authorUsername {
		return ErrUnauthorized
	}

	if !tour.CanBeArchived() {
		return ErrTourNotArchivable
	}

	tour.Status = TourStatusArchived

	return service.repository.UpdateTour(tour)
}

func (service *TourService) UnarchiveTour(tourID uint, authorUsername string) error {
	tour, err := service.repository.GetTourByID(tourID)
	if err != nil {
		return ErrTourNotFound
	}

	if tour.AuthorUsername != authorUsername {
		return ErrUnauthorized
	}

	if !tour.CanBeUnarchived() {
		return ErrTourNotUnarchivable
	}

	tour.Status = TourStatusPublished

	return service.repository.UpdateTour(tour)
}

func isValidDifficulty(difficulty string) bool {
	validDifficulties := []string{DifficultyEasy, DifficultyMedium, DifficultyHard}
	for _, valid := range validDifficulties {
		if difficulty == valid {
			return true
		}
	}
	return false
}

// TourExecution Service Methods

const ProximityThresholdMeters = 1000.0 // 1 km proximity threshold for simulation

func (service *TourService) StartTourExecution(request *StartTourExecutionRequest, touristUsername string) (*TourExecution, error) {
	// Check if tourist already has an active tour execution
	activeExecution, _ := service.repository.GetActiveTourExecution(touristUsername)
	if activeExecution != nil {
		return nil, errors.New("tourist already has an active tour execution")
	}

	// Verify tour exists and can be executed
	tour, err := service.repository.GetTourByID(request.TourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}

	if tour.Status != TourStatusPublished && tour.Status != TourStatusArchived {
		return nil, errors.New("tour is not available for execution")
	}

	err = service.checkTourPurchase(touristUsername, request.TourID)
	if err != nil {
		if err == ErrTourNotPurchased {
			return nil, errors.New("tour must be purchased before execution")
		}
	}

	// Create new tour execution
	execution := &TourExecution{
		TourID:          request.TourID,
		TouristUsername: touristUsername,
		Status:          ExecutionStatusActive,
		StartLatitude:   request.Latitude,
		StartLongitude:  request.Longitude,
	}

	err = service.repository.StartTourExecution(execution)
	if err != nil {
		return nil, err
	}

	return execution, nil
}

func (service *TourService) GetActiveTourExecution(touristUsername string) (*TourExecution, error) {
	return service.repository.GetActiveTourExecution(touristUsername)
}

func (service *TourService) EndTourExecution(executionID uint, status string, touristUsername string) (*TourExecution, error) {
	// Verify execution belongs to the tourist
	execution, err := service.repository.GetTourExecutionByID(executionID)
	if err != nil {
		return nil, errors.New("tour execution not found")
	}

	if execution.TouristUsername != touristUsername {
		return nil, errors.New("unauthorized access to tour execution")
	}

	if execution.Status != ExecutionStatusActive {
		return nil, errors.New("tour execution is not active")
	}

	if status != ExecutionStatusCompleted && status != ExecutionStatusAbandoned {
		return nil, errors.New("invalid end status")
	}

	return service.repository.EndTourExecution(executionID, status)
}

func (service *TourService) CheckProximity(executionID uint, latitude, longitude float64, touristUsername string) (*CheckProximityResponse, error) {
	// Verify execution belongs to the tourist
	execution, err := service.repository.GetTourExecutionByID(executionID)
	if err != nil {
		return nil, errors.New("tour execution not found")
	}

	if execution.TouristUsername != touristUsername {
		return nil, errors.New("unauthorized access to tour execution")
	}

	if execution.Status != ExecutionStatusActive {
		return nil, errors.New("tour execution is not active")
	}

	// Update last activity regardless of proximity check result
	execution.LastActivity = time.Now()
	service.repository.UpdateTourExecution(execution)

	// Get tour to access key points
	tour, err := service.repository.GetTourByID(execution.TourID)
	if err != nil {
		return &CheckProximityResponse{
			KeyPointReached: false,
			LastActivity:    execution.LastActivity,
			Message:         "Failed to load tour data",
		}, nil
	}

	// Check proximity to each key point
	fmt.Printf("DEBUG: Tourist position: lat=%f, lng=%f\n", latitude, longitude)
	fmt.Printf("DEBUG: Checking proximity for %d key points\n", len(tour.KeyPoints))

	for _, keyPoint := range tour.KeyPoints {
		// Check if this key point is already completed
		_, err := service.repository.GetKeyPointCompletion(execution.ID, keyPoint.ID)
		if err == nil {
			// Already completed, skip
			fmt.Printf("DEBUG: KeyPoint '%s' already completed, skipping\n", keyPoint.Name)
			continue
		}

		// Calculate distance using the existing haversine formula
		distance := Calculator.HaversineDistance(latitude, longitude, keyPoint.Latitude, keyPoint.Longitude)
		distanceMeters := distance * 1000 // Convert km to meters

		fmt.Printf("DEBUG: KeyPoint '%s' at lat=%f, lng=%f - Distance: %.2f meters (threshold: %.2f)\n",
			keyPoint.Name, keyPoint.Latitude, keyPoint.Longitude, distanceMeters, ProximityThresholdMeters)

		if distanceMeters <= ProximityThresholdMeters {
			// Tourist is within proximity of this key point
			completion := &KeyPointCompletion{
				TourExecutionID: execution.ID,
				KeyPointID:      keyPoint.ID,
				Latitude:        latitude,
				Longitude:       longitude,
			}

			err = service.repository.CreateKeyPointCompletion(completion)
			if err != nil {
				return &CheckProximityResponse{
					KeyPointReached: false,
					LastActivity:    execution.LastActivity,
					Message:         "Failed to record key point completion",
				}, nil
			}

			return &CheckProximityResponse{
				KeyPointReached:    true,
				KeyPoint:           &keyPoint,
				KeyPointCompletion: completion,
				LastActivity:       execution.LastActivity,
				Message:            fmt.Sprintf("Key point '%s' reached!", keyPoint.Name),
			}, nil
		}
	}

	return &CheckProximityResponse{
		KeyPointReached: false,
		LastActivity:    execution.LastActivity,
		Message:         "No key points nearby",
	}, nil
}

func (service *TourService) GetExecutableToursForTourist() ([]Tour, error) {
	return service.repository.GetExecutableToursForTourist()
}

func (service *TourService) GetPurchasedToursForTourist(userID string) ([]Tour, error) {
	// Get purchased tour IDs for this user
	purchasedTourIds, err := service.getPurchasedTourIds(userID)
	if err != nil {
		fmt.Printf("Warning: Could not get purchased tours for user %s: %v\n", userID, err)
		// Return empty list if purchase service is not available
		return []Tour{}, nil
	}

	// Get tours by IDs from repository
	return service.repository.GetPurchasedToursForTourist(purchasedTourIds)
}

func (service *TourService) getPurchasedTourIds(userID string) ([]uint, error) {

	purchaseHost := os.Getenv("PURCHASE_SERVICE_HOST")
	purchasePort := os.Getenv("PURCHASE_SERVICE_PORT")
	if purchaseHost == "" {
		purchaseHost = "purchase-service"
	}
	if purchasePort == "" {
		purchasePort = "8084"
	}

	purchaseURL := fmt.Sprintf("http://%s:%s/tokens", purchaseHost, purchasePort)

	req, err := http.NewRequest(http.MethodGet, purchaseURL, nil)
	if err != nil {
		fmt.Printf("Failed to create purchase tokens request: %v\n", err)
		return nil, fmt.Errorf("failed to get purchased tours: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-username", userID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to call purchase service: %v\n", err)
		return nil, fmt.Errorf("failed to get purchased tours: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// User has no purchases
		return []uint{}, nil
	} else if resp.StatusCode >= 300 {
		fmt.Printf("Purchase service returned error status: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("purchase service failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var tokensResponse struct {
		Tokens []struct {
			ID       uint   `json:"id"`
			TourID   uint   `json:"tour_id"`
			TourName string `json:"tour_name"`
			Status   string `json:"status"`
		} `json:"tokens"`
		Message string `json:"message"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokensResponse)
	if err != nil {
		fmt.Printf("Failed to decode purchase response: %v\n", err)
		return nil, fmt.Errorf("failed to parse purchase response: %w", err)
	}

	// Extract tour IDs from valid/active tokens
	var tourIds []uint
	for _, token := range tokensResponse.Tokens {
		if token.Status == "active" {
			tourIds = append(tourIds, token.TourID)
		}
	}

	fmt.Printf("Found %d purchased tours for user %s: %v\n", len(tourIds), userID, tourIds)
	return tourIds, nil
}

func (service *TourService) checkTourPurchase(userID string, tourID uint) error {
	purchaseHost := os.Getenv("PURCHASE_SERVICE_HOST")
	purchasePort := os.Getenv("PURCHASE_SERVICE_PORT")
	if purchaseHost == "" {
		purchaseHost = "purchase-service"
	}
	if purchasePort == "" {
		purchasePort = "8084"
	}

	purchaseURL := fmt.Sprintf("http://%s:%s/validate/%d", purchaseHost, purchasePort, tourID)

	req, err := http.NewRequest(http.MethodGet, purchaseURL, nil)
	if err != nil {
		fmt.Printf("Failed to create purchase validation request: %v\n", err)
		return fmt.Errorf("failed to validate purchase: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-username", userID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to call purchase service: %v\n", err)
		return fmt.Errorf("failed to validate purchase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {

		return ErrTourNotPurchased
	} else if resp.StatusCode >= 300 {

		fmt.Printf("Purchase service returned error status: %d\n", resp.StatusCode)
		return fmt.Errorf("purchase validation failed with status: %d", resp.StatusCode)
	}

	fmt.Printf("Tour purchase validated successfully for user %s and tour %d\n", userID, tourID)
	return nil
}
