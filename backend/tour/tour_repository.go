package main

import (
	"time"

	"gorm.io/gorm"
)

type TourRepository struct {
	database *gorm.DB
}

func (repo *TourRepository) CreateTour(tour *Tour) error {
	result := repo.database.Create(tour)
	return result.Error
}

func (repo *TourRepository) GetToursByAuthor(authorUsername string) ([]Tour, error) {
	var tours []Tour
	result := repo.database.Preload("KeyPoints").Where("author_username = ?", authorUsername).Find(&tours)
	return tours, result.Error
}

func (repo *TourRepository) GetAllPublishedTours() ([]Tour, error) {
	var tours []Tour
	result := repo.database.Preload("KeyPoints").Where("status = ?", TourStatusPublished).Find(&tours)
	return tours, result.Error
}

func (repo *TourRepository) GetTourByID(id uint) (*Tour, error) {
	var tour Tour
	result := repo.database.Preload("KeyPoints").Where("id = ?", id).First(&tour)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tour, nil
}

func (repo *TourRepository) UpdateTour(tour *Tour) error {
	result := repo.database.Session(&gorm.Session{FullSaveAssociations: true}).Save(tour)
	return result.Error
}

func (repo *TourRepository) DeleteTour(id uint) error {
	result := repo.database.Delete(&Tour{}, id)
	return result.Error
}

func (repo *TourRepository) CreateKeyPoint(keyPoint *KeyPoint) error {
	result := repo.database.Create(keyPoint)
	return result.Error
}

func (repo *TourRepository) DeleteKeyPointsByTourID(tourID uint) error {
	// Use Unscoped() to force hard delete instead of soft delete
	result := repo.database.Unscoped().Where("tour_id = ?", tourID).Delete(&KeyPoint{})
	return result.Error
}

// TourExecution Repository Methods

func (repo *TourRepository) StartTourExecution(execution *TourExecution) error {
	result := repo.database.Create(execution)
	return result.Error
}

func (repo *TourRepository) GetActiveTourExecution(touristUsername string) (*TourExecution, error) {
	var execution TourExecution
	result := repo.database.Preload("KeyPointCompletions").
		Where("tourist_username = ? AND status = ?", touristUsername, ExecutionStatusActive).
		First(&execution)
	if result.Error != nil {
		return nil, result.Error
	}
	return &execution, nil
}

func (repo *TourRepository) GetTourExecutionByID(id uint) (*TourExecution, error) {
	var execution TourExecution
	result := repo.database.Preload("KeyPointCompletions").First(&execution, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &execution, nil
}

func (repo *TourRepository) UpdateTourExecution(execution *TourExecution) error {
	result := repo.database.Save(execution)
	return result.Error
}

func (repo *TourRepository) EndTourExecution(id uint, status string) (*TourExecution, error) {
	var execution TourExecution
	result := repo.database.First(&execution, id)
	if result.Error != nil {
		return nil, result.Error
	}

	execution.Status = status
	now := time.Now()
	execution.EndTime = &now

	result = repo.database.Save(&execution)
	if result.Error != nil {
		return nil, result.Error
	}

	return &execution, nil
}

func (repo *TourRepository) CreateKeyPointCompletion(completion *KeyPointCompletion) error {
	result := repo.database.Create(completion)
	return result.Error
}

func (repo *TourRepository) GetKeyPointCompletion(executionID, keyPointID uint) (*KeyPointCompletion, error) {
	var completion KeyPointCompletion
	result := repo.database.
		Where("tour_execution_id = ? AND key_point_id = ?", executionID, keyPointID).
		First(&completion)
	if result.Error != nil {
		return nil, result.Error
	}
	return &completion, nil
}

func (repo *TourRepository) GetExecutableToursForTourist() ([]Tour, error) {
	// For now, tourists can execute published and archived tours
	// In KT3, will add purchase check here
	var tours []Tour
	result := repo.database.Preload("KeyPoints").
		Where("status IN (?)", []string{TourStatusPublished, TourStatusArchived}).
		Find(&tours)
	return tours, result.Error
}
