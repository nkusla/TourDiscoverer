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

	tour := &Tour{
		Name:           request.Name,
		Description:    request.Description,
		Difficulty:     request.Difficulty,
		Tags:           strings.TrimSpace(request.Tags),
		Status:         TourStatusDraft,
		Price:          0,
		AuthorUsername: authorUsername,
		KeyPoints:      []KeyPoint{},
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

func isValidDifficulty(difficulty string) bool {
	validDifficulties := []string{DifficultyEasy, DifficultyMedium, DifficultyHard}
	for _, valid := range validDifficulties {
		if difficulty == valid {
			return true
		}
	}
	return false
}
