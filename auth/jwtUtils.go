package main

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required!")
	}

	return secret
}

func CreateJWT(username string, role string) (string, error) {
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtSecret := []byte(GetJWTSecret())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
