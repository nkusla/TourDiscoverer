package main

import "errors"

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidCredentials = errors.New("incorrect username or password")
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
)

const (
	UserUsernamePrimaryKey = "users_pkey"
	UserEmailUniqueIndex   = "idx_user_email"
)
