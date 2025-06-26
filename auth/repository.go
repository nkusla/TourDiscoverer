package main

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

var ErrUserAlreadyExists = errors.New("User with given username already exists")

func (r *UserRepository) Create(user *User) error {
	// Check if user already exists
	var existingUser User
	err := r.database.Where("username = ?", user.Username).First(&existingUser).Error

	if err == nil {
		return ErrUserAlreadyExists
	}

	// User doesn't exist, create it
	return r.database.Create(user).Error
}

func (r *UserRepository) Update(user *User) error {
	return r.database.Save(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.database.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
