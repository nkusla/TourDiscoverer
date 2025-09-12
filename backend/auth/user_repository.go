package main

import (
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func (r *UserRepository) Create(user *User) error {
	err := r.database.Create(user).Error

	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, UserUsernamePrimaryKey) {
			return ErrUsernameAlreadyExists
		} else if strings.Contains(errMsg, UserEmailUniqueIndex) {
			return ErrEmailAlreadyExists
		}
	}

	return nil
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

func (r *UserRepository) FindAll() ([]*User, error) {
	var users []*User
	err := r.database.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
