package main

import (
	"errors"
	"net/http"
)

var (
	ErrReviewNotFound       = errors.New("review not found")
	ErrUnauthorizedReview   = errors.New("unauthorized to access this review")
	ErrInvalidRating        = errors.New("rating must be between 1 and 5")
	ErrInvalidVisitDate     = errors.New("invalid visit date format")
	ErrDuplicateReview      = errors.New("review already exists for this tour by this user")
	ErrTourNotFound         = errors.New("tour not found")
)

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(message string, statusCode int) APIError {
	return APIError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func GetErrorStatusCode(err error) int {
	switch err {
	case ErrReviewNotFound, ErrTourNotFound:
		return http.StatusNotFound
	case ErrUnauthorizedReview:
		return http.StatusForbidden
	case ErrInvalidRating, ErrInvalidVisitDate:
		return http.StatusBadRequest
	case ErrDuplicateReview:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
