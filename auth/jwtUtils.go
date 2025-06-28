package main

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	return value
}

func CreateJWT(username string, role string) (string, error) {
	jwtExpiration := GetEnv("JWT_EXPIRATION")

	expirationTime, err := time.ParseDuration(jwtExpiration)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtSecret := []byte(GetEnv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
