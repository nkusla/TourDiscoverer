package main

import "errors"

var (
	ErrTourNotFound       = errors.New("tour not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTourNotPublishable = errors.New("tour cannot be published")
)
