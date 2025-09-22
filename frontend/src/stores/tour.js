import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useTourStore = defineStore('tour', () => {
  const tours = ref([])
  const currentTour = ref(null)
  const loading = ref(false)

  const fetchTours = async (params = {}) => {
    loading.value = true
    try {
      const response = await api.get('/api/tours/all', { params })
      // Backend returns { tours: [...], count: number }
      tours.value = response.data.tours || []
      return response.data.tours || []
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch tours')
    } finally {
      loading.value = false
    }
  }

  const fetchMyTours = async (params = {}) => {
    loading.value = true
    try {
      const response = await api.get('/api/tours/my', { params })
      // Backend returns { tours: [...], count: number }
      tours.value = response.data.tours || []
      return response.data.tours || []
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch tours')
    } finally {
      loading.value = false
    }
  }

  const getTour = async (id) => {
    try {
      const response = await api.get(`/api/tours/${id}`)
      currentTour.value = response.data
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch tour')
    }
  }

  const createTour = async (tourData) => {
    try {
      const response = await api.post('/api/tours', tourData)
      tours.value.unshift(response.data)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to create tour')
    }
  }

  const updateTour = async (id, tourData) => {
    try {
      const response = await api.put(`/api/tours/${id}`, tourData)

      // Update tour in the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value[index] = response.data
      }

      currentTour.value = response.data
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to update tour')
    }
  }

  const deleteTour = async (id) => {
    try {
      await api.delete(`/api/tours/${id}`)

      // Remove tour from the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value.splice(index, 1)
      }

      return true
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to delete tour')
    }
  }

  const publishTour = async (id) => {
    try {
      const response = await api.put(`/api/tours/${id}/publish`)

      // Update tour status in the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value[index].status = 'published'
      }

      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to publish tour')
    }
  }

  const archiveTour = async (id) => {
    try {
      const response = await api.put(`/api/tours/${id}/archive`)

      // Update tour status in the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value[index].status = 'archived'
      }

      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to archive tour')
    }
  }

  const unarchiveTour = async (id) => {
    try {
      const response = await api.put(`/api/tours/${id}/unarchive`)

      // Update tour status in the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value[index].status = 'published'
      }

      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to unarchive tour')
    }
  }

  // Key point management
  const addKeyPoint = async (tourId, keyPointData) => {
    try {
      const response = await api.post(`/api/tours/${tourId}/keypoint`, keyPointData)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to add key point')
    }
  }

  const updateKeyPoint = async (tourId, keyPointId, keyPointData) => {
    try {
      const response = await api.put(`/api/tours/${tourId}/keypoint/${keyPointId}`, keyPointData)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to update key point')
    }
  }

  const deleteKeyPoint = async (tourId, keyPointId) => {
    try {
      await api.delete(`/api/tours/${tourId}/keypoint/${keyPointId}`)
      return true
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to delete key point')
    }
  }

  const searchTours = async (query, filters = {}) => {
    try {
      const params = { search: query, ...filters }
      const response = await api.get('/api/tours/search', { params })
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Search failed')
    }
  }

  const getPopularTours = async (limit = 10) => {
    try {
      const response = await api.get('/api/tours/popular', { params: { limit } })
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch popular tours')
    }
  }

  const getNearbyTours = async (latitude, longitude, radius = 50) => {
    try {
      const response = await api.get('/api/tours/nearby', {
        params: { latitude, longitude, radius }
      })
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch nearby tours')
    }
  }

  return {
    tours,
    currentTour,
    loading,
    fetchTours,
    fetchMyTours,
    getTour,
    createTour,
    updateTour,
    deleteTour,
    publishTour,
    archiveTour,
    unarchiveTour,
    addKeyPoint,
    updateKeyPoint,
    deleteKeyPoint,
    searchTours,
    getPopularTours,
    getNearbyTours
  }
})
