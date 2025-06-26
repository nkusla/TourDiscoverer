package main

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"-" gorm:"not null"` // "-" excludes from JSON
	Email    string `json:"email" gorm:"not null;unique"`
	Role     string `json:"role" gorm:"not null;default:'tourist'"`
	IsBanned bool   `json:"is_banned" gorm:"default:false"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
