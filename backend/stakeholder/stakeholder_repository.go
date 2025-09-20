package main

import "gorm.io/gorm"

type StakeholderRepository struct {
	database *gorm.DB
}

func (r *StakeholderRepository) Create(stakeholder *Stakeholder) error {
	return r.database.Create(stakeholder).Error
}

func (r *StakeholderRepository) GetByUsername(username string) (*Stakeholder, error) {
	var stakeholder Stakeholder
	err := r.database.Where("username = ?", username).First(&stakeholder).Error
	if err != nil {
		return nil, err
	}
	return &stakeholder, nil
}

func (r *StakeholderRepository) Update(stakeholder *Stakeholder) error {
	return r.database.Save(stakeholder).Error
}

func (r *StakeholderRepository) Delete(username string) error {
	return r.database.Where("username = ?", username).Delete(&Stakeholder{}).Error
}

func (r *StakeholderRepository) UpdateTouristPosition(username string, latitude, longitude float64) (*TouristPosition, error) {
	var position TouristPosition

	// Try to find existing position
	err := r.database.Where("username = ?", username).First(&position).Error

	if err == gorm.ErrRecordNotFound {
		// Create new position
		position = TouristPosition{
			Username:  username,
			Latitude:  latitude,
			Longitude: longitude,
		}
		result := r.database.Create(&position)
		return &position, result.Error
	} else if err != nil {
		return nil, err
	}

	// Update existing position
	position.Latitude = latitude
	position.Longitude = longitude
	result := r.database.Save(&position)
	return &position, result.Error
}

func (r *StakeholderRepository) GetTouristPosition(username string) (*TouristPosition, error) {
	var position TouristPosition
	result := r.database.Where("username = ?", username).First(&position)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, ErrPositionNotFound
	}

	return &position, result.Error
}

func (r *StakeholderRepository) DeleteTouristPosition(username string) error {
	result := r.database.Where("username = ?", username).Delete(&TouristPosition{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrPositionNotFound
	}
	
	return nil
}
