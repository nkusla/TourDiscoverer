package main

import (
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

func (repo *TourRepository) GetTourByID(id uint) (*Tour, error) {
	var tour Tour
	result := repo.database.Preload("KeyPoints").Where("id = ?", id).First(&tour)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tour, nil
}

func (repo *TourRepository) UpdateTour(tour *Tour) error {
	result := repo.database.Save(tour)
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
