package main

import (
	"errors"
	"fmt"
	"strings"
)

type TourService struct {
	repository *TourRepository
}

func (service *TourService) CreateTour(request *CreateTourRequest, authorUsername string) (*Tour, error) {
	if !isValidDifficulty(request.Difficulty) {
		return nil, errors.New("invalid difficulty level")
	}

	tour := &Tour{
		Name:             request.Name,
		Description:      request.Description,
		Difficulty:       request.Difficulty,
		Tags:             strings.TrimSpace(request.Tags),
		Status:           TourStatusDraft,
		Price:            0,
		AuthorUsername:   authorUsername,
		TransportDetails: []Transport{},
		Distance:         0,
		KeyPoints:        []KeyPoint{},
	}

	err := service.repository.CreateTour(tour)
	if err != nil {
		return nil, err
	}

	return tour, nil
}

func (service *TourService) GetToursByAuthor(authorUsername string) ([]Tour, error) {
	tours, err := service.repository.GetToursByAuthor(authorUsername)
	if err != nil {
		return nil, err
	}
	return tours, nil
}

func (service *TourService) GetTourByID(id uint) (*Tour, error) {
	return service.repository.GetTourByID(id)
}

func (service *TourService) CreateKeyPoint(request *CreateKeyPointRequest, tourID uint, authorUsername string) (*KeyPoint, error) {
	// Check if tour exists and belongs to the author
	tour, err := service.repository.GetTourByID(tourID)
	if err != nil {
		return nil, err
	}

	if tour.AuthorUsername != authorUsername {
		return nil, fmt.Errorf("you can only add key points to your own tours")
	}

	keyPoint := &KeyPoint{
		TourID:      tourID,
		Name:        request.Name,
		Description: request.Description,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		ImageURL:    request.ImageURL,
		Order:       request.Order,
	}

	err = service.repository.CreateKeyPoint(keyPoint)
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
