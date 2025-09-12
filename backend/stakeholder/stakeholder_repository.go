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
