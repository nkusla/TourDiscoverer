package main

import (
	"errors"
	"strings"
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
		TransportDetails: []Transport{},
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
	result = tx.Model(&tour).Select("name", "description", "difficulty", "tags", "price", "distance", "transport_details").Updates(map[string]interface{}{
		"name":              tour.Name,
		"description":       tour.Description,
		"difficulty":        tour.Difficulty,
		"tags":              tour.Tags,
		"price":             tour.Price,
		"distance":          tour.Distance,
		"transport_details": tour.TransportDetails,
	})
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

	return service.repository.UpdateTour(tour)
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
