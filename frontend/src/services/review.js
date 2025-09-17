import api from './api'

// Simple Review Service
export const reviewService = {
  // Create a new review
  createReview: async (reviewData) => {
    const response = await api.post('/api/reviews', reviewData)
    return response.data
  },

  // Get reviews for a specific tour
  getTourReviews: async (tourId, page = 1, pageSize = 10) => {
    const response = await api.get(`/api/reviews/tour/${tourId}`, {
      params: { page, page_size: pageSize }
    })
    return response.data
  },

  // Get tour average rating
  getTourRating: async (tourId) => {
    const response = await api.get(`/api/reviews/tour/${tourId}/rating`)
    return response.data
  },

  // Delete a review
  deleteReview: async (reviewId) => {
    await api.delete(`/api/reviews/${reviewId}`)
  }
}

export default reviewService