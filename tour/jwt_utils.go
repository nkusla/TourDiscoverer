package main

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GetUsernameFromJWT(tokenString string) (string, error) {
	username, _, err := GetUserInfoFromJWT(tokenString)
	return username, err
}

func GetUserInfoFromJWT(tokenString string) (string, string, error) {

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		return "", "", errors.New("JWT_SECRET environment variable not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, claims.Role, nil
	}

	return "", "", errors.New("invalid token")
}
