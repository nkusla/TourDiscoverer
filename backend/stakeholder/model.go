package main

import "time"

type Stakeholder struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string    `json:"username" gorm:"uniqueIndex;not null"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	Biography      string    `json:"biography"`
	Motto          string    `json:"motto"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type TouristPosition struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null;uniqueIndex"`
	Latitude  float64   `json:"latitude" gorm:"not null"`
	Longitude float64   `json:"longitude" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"`
}
