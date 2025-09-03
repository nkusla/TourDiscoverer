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
      const response = await api.get('/tours', { params })
      tours.value = response.data
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch tours')
    } finally {
      loading.value = false
    }
  }
  
  const getTour = async (id) => {
    try {
      const response = await api.get(`/tours/${id}`)
      currentTour.value = response.data
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch tour')
    }
  }
  
  const createTour = async (tourData) => {
    try {
      const response = await api.post('/tours', tourData)
      tours.value.unshift(response.data)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to create tour')
    }
  }
  
  const updateTour = async (id, tourData) => {
    try {
      const response = await api.put(`/tours/${id}`, tourData)
      
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
      await api.delete(`/tours/${id}`)
      
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
      const response = await api.patch(`/tours/${id}/publish`)
      
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
  
  const unpublishTour = async (id) => {
    try {
      const response = await api.patch(`/tours/${id}/unpublish`)
      
      // Update tour status in the list
      const index = tours.value.findIndex(tour => tour.id === id)
      if (index !== -1) {
        tours.value[index].status = 'draft'
      }
      
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to unpublish tour')
    }
  }
  
  // Key point management
  const addKeyPoint = async (tourId, keyPointData) => {
    try {
      const response = await api.post(`/tours/${tourId}/keypoints`, keyPointData)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to add key point')
    }
  }
  
  const updateKeyPoint = async (tourId, keyPointId, keyPointData) => {
    try {
      const response = await api.put(`/tours/${tourId}/keypoints/${keyPointId}`, keyPointData)
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to update key point')
    }
  }
  
  const deleteKeyPoint = async (tourId, keyPointId) => {
    try {
      await api.delete(`/tours/${tourId}/keypoints/${keyPointId}`)
      return true
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to delete key point')
    }
  }
  
  const searchTours = async (query, filters = {}) => {
    try {
      const params = { search: query, ...filters }
      const response = await api.get('/tours/search', { params })
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Search failed')
    }
  }
  
  const getPopularTours = async (limit = 10) => {
    try {
      const response = await api.get('/tours/popular', { params: { limit } })
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to fetch popular tours')
    }
  }
  
  const getNearbyTours = async (latitude, longitude, radius = 50) => {
    try {
      const response = await api.get('/tours/nearby', {
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
    getTour,
    createTour,
    updateTour,
    deleteTour,
    publishTour,
    unpublishTour,
    addKeyPoint,
    updateKeyPoint,
    deleteKeyPoint,
    searchTours,
    getPopularTours,
    getNearbyTours
  }
})
