package main

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func (r *UserRepository) Create(user *User) error {
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
