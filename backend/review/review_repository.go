package main

import (
	"gorm.io/gorm"
)

type ReviewRepository struct {
	database *gorm.DB
}

func (r *ReviewRepository) CreateReview(review *Review) error {
	return r.database.Create(review).Error
}

func (r *ReviewRepository) GetReviewByID(id uint) (*Review, error) {
	var review Review
	err := r.database.Preload("Images").First(&review, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetReviewsByTourID(tourID uint, limit, offset int) ([]Review, int64, error) {
	var reviews []Review
	var totalCount int64

	// Get total count
	err := r.database.Model(&Review{}).Where("tour_id = ?", tourID).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Get reviews with pagination
	err = r.database.Preload("Images").
		Where("tour_id = ?", tourID).
		Order("review_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, totalCount, err
}

func (r *ReviewRepository) GetReviewsByUsername(username string, limit, offset int) ([]Review, int64, error) {
	var reviews []Review
	var totalCount int64

	// Get total count
	err := r.database.Model(&Review{}).Where("tourist_username = ?", username).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Get reviews with pagination
	err = r.database.Preload("Images").
		Where("tourist_username = ?", username).
		Order("review_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, totalCount, err
}

func (r *ReviewRepository) UpdateReview(review *Review) error {
	return r.database.Save(review).Error
}

func (r *ReviewRepository) DeleteReview(id uint) error {
	result := r.database.Delete(&Review{}, id)
	if result.RowsAffected == 0 {
		return ErrReviewNotFound
	}
	return result.Error
}

func (r *ReviewRepository) CheckExistingReview(tourID uint, username string) (bool, error) {
	var count int64
	err := r.database.Model(&Review{}).
		Where("tour_id = ? AND tourist_username = ?", tourID, username).
		Count(&count).Error
	return count > 0, err
}

func (r *ReviewRepository) GetAverageRating(tourID uint) (float64, error) {
	var avg struct {
		Average float64
	}
	
	err := r.database.Model(&Review{}).
		Select("AVG(rating) as average").
		Where("tour_id = ?", tourID).
		Scan(&avg).Error
	
	return avg.Average, err
}

func (r *ReviewRepository) DeleteReviewImages(reviewID uint) error {
	return r.database.Where("review_id = ?", reviewID).Delete(&ReviewImage{}).Error
}

func (r *ReviewRepository) CreateReviewImages(images []ReviewImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.database.Create(&images).Error
}
