package main

import "errors"

var (
	ErrCartNotFound       = errors.New("shopping cart not found")
	ErrItemNotFound       = errors.New("item not found in cart")
	ErrTourAlreadyInCart  = errors.New("tour already in shopping cart")
	ErrTourNotPublished   = errors.New("tour is not published")
	ErrTourArchived       = errors.New("tour is archived and cannot be purchased")
	ErrEmptyCart          = errors.New("shopping cart is empty")
	ErrTokenNotFound      = errors.New("purchase token not found")
	ErrTokenExpired       = errors.New("purchase token has expired")
	ErrTokenInvalid       = errors.New("purchase token is invalid")
	ErrUnauthorized       = errors.New("unauthorized access")
)