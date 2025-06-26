package main

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"-" gorm:"not null"` // "-" excludes from JSON
	Mail     string `json:"mail" gorm:"not null;unique"`
	Role     string `json:"role" gorm:"not null;default:'tourist'"`
	IsBanned bool   `json:"is_banned" gorm:"default:false"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
