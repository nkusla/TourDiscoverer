package main

import (
	"time"
)

type ReviewService struct {
	repository *ReviewRepository
}

func (s *ReviewService) CreateReview(req CreateReviewRequest, username string) (*Review, error) {
	// Parse visit date
	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return nil, ErrInvalidVisitDate
	}

	// Check if review already exists for this tour by this user
	exists, err := s.repository.CheckExistingReview(req.TourID, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateReview
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidRating
	}

	// Create review
	review := &Review{
		TourID:          req.TourID,
		TouristUsername: username,
		Rating:          req.Rating,
		Comment:         req.Comment,
		VisitDate:       visitDate,
		ReviewDate:      time.Now(),
	}

	err = s.repository.CreateReview(review)
	if err != nil {
		return nil, err
	}

	// Add images if provided
	if len(req.Images) > 0 {
		images := make([]ReviewImage, len(req.Images))
		for i, imageURL := range req.Images {
			images[i] = ReviewImage{
				ReviewID: review.ID,
				ImageURL: imageURL,
			}
		}
		err = s.repository.CreateReviewImages(images)
		if err != nil {
			return nil, err
		}
		review.Images = images
	}

	return review, nil
}

func (s *ReviewService) GetReviewByID(id uint) (*Review, error) {
	return s.repository.GetReviewByID(id)
}

func (s *ReviewService) GetReviewsByTourID(tourID uint, page, pageSize int) (*ReviewListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	reviews, totalCount, err := s.repository.GetReviewsByTourID(tourID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Calculate average rating
	avgRating, err := s.repository.GetAverageRating(tourID)
	if err != nil {
		avgRating = 0
	}

	// Convert to response format
	reviewResponses := make([]ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = review.ToResponse()
	}

	return &ReviewListResponse{
		Reviews:       reviewResponses,
		TotalCount:    totalCount,
		AverageRating: avgRating,
	}, nil
}

func (s *ReviewService) GetReviewsByUsername(username string, page, pageSize int) (*ReviewListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	reviews, totalCount, err := s.repository.GetReviewsByUsername(username, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	reviewResponses := make([]ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = review.ToResponse()
	}

	return &ReviewListResponse{
		Reviews:    reviewResponses,
		TotalCount: totalCount,
	}, nil
}

func (s *ReviewService) UpdateReview(id uint, req UpdateReviewRequest, username string) (*Review, error) {
	review, err := s.repository.GetReviewByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns this review
	if review.TouristUsername != username {
		return nil, ErrUnauthorizedReview
	}

	// Update fields if provided
	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, ErrInvalidRating
		}
		review.Rating = *req.Rating
	}

	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	// Update images if provided
	if req.Images != nil {
		// Delete existing images
		err = s.repository.DeleteReviewImages(review.ID)
		if err != nil {
			return nil, err
		}

		// Create new images
		if len(req.Images) > 0 {
			images := make([]ReviewImage, len(req.Images))
			for i, imageURL := range req.Images {
				images[i] = ReviewImage{
					ReviewID: review.ID,
					ImageURL: imageURL,
				}
			}
			err = s.repository.CreateReviewImages(images)
			if err != nil {
				return nil, err
			}
			review.Images = images
		} else {
			review.Images = []ReviewImage{}
		}
	}

	err = s.repository.UpdateReview(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) DeleteReview(id uint, username string) error {
	review, err := s.repository.GetReviewByID(id)
	if err != nil {
		return err
	}

	// Check if user owns this review
	if review.TouristUsername != username {
		return ErrUnauthorizedReview
	}

	return s.repository.DeleteReview(id)
}

func (s *ReviewService) GetTourAverageRating(tourID uint) (float64, error) {
	return s.repository.GetAverageRating(tourID)
}
