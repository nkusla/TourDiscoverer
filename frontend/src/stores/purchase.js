import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const usePurchaseStore = defineStore('purchase', () => {
  const purchasedTokens = ref([])
  const loading = ref(false)
  
  const fetchPurchasedTours = async () => {
    loading.value = true
    try {
      const response = await api.get('/api/purchases/tokens')
      purchasedTokens.value = response.data.tokens || []
      return purchasedTokens.value
    } catch (error) {
      console.error('Failed to fetch purchased tours:', error)
      throw new Error(error.response?.data?.error || 'Failed to fetch purchased tours')
    } finally {
      loading.value = false
    }
  }
  
  const getTokenDetails = async (token) => {
    try {
      const response = await api.get(`/api/purchases/tokens/${token}`)
      return response.data.token
    } catch (error) {
      console.error('Failed to get token details:', error)
      throw new Error(error.response?.data?.error || 'Failed to get token details')
    }
  }
  
  const validateAccess = async (tourId) => {
    try {
      const response = await api.get(`/api/purchases/validate/${tourId}`)
      return response.data.token
    } catch (error) {
      // If error is 404, user hasn't purchased this tour
      if (error.response?.status === 404) {
        return null
      }
      console.error('Failed to validate access:', error)
      throw new Error(error.response?.data?.error || 'Failed to validate access')
    }
  }
  
  const hasPurchased = (tourId) => {
    return purchasedTokens.value.some(token => token.tour_id === tourId && token.status === 'active')
  }
  
  const getPurchaseToken = (tourId) => {
    return purchasedTokens.value.find(token => token.tour_id === tourId && token.status === 'active')
  }
  
  return {
    purchasedTokens,
    loading,
    fetchPurchasedTours,
    getTokenDetails,
    validateAccess,
    hasPurchased,
    getPurchaseToken
  }
})