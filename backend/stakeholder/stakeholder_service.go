package main

import (
	"errors"
	"gorm.io/gorm"
)

type StakeholderService struct {
	repository *StakeholderRepository
}

func (s *StakeholderService) CreateStakeholder(username, firstName, lastName, profilePicture, biography, motto string) (*Stakeholder, error) {
	// Check if stakeholder already exists
	existingStakeholder, err := s.repository.GetByUsername(username)
	if err == nil && existingStakeholder != nil {
		return nil, errors.New("stakeholder already exists")
	}

	stakeholder := &Stakeholder{
		Username:       username,
		FirstName:      firstName,
		LastName:       lastName,
		ProfilePicture: profilePicture,
		Biography:      biography,
		Motto:          motto,
	}

	err = s.repository.Create(stakeholder)
	if err != nil {
		return nil, err
	}

	return stakeholder, nil
}

func (s *StakeholderService) GetStakeholderProfile(username string) (*Stakeholder, error) {
	stakeholder, err := s.repository.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return stakeholder, nil
}

func (s *StakeholderService) UpdateStakeholderProfile(username, firstName, lastName, profilePicture, biography, motto string) (*Stakeholder, error) {
	stakeholder, err := s.repository.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	// Update fields only if they are provided (not empty)
	if firstName != "" {
		stakeholder.FirstName = firstName
	}
	if lastName != "" {
		stakeholder.LastName = lastName
	}
	if profilePicture != "" {
		stakeholder.ProfilePicture = profilePicture
	}
	if biography != "" {
		stakeholder.Biography = biography
	}
	if motto != "" {
		stakeholder.Motto = motto
	}

	err = s.repository.Update(stakeholder)
	if err != nil {
		return nil, err
	}

	return stakeholder, nil
}
