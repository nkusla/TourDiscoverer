import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export const useCartStore = defineStore('cart', () => {
  const cart = ref(null)
  const loading = ref(false)
  
  const cartItems = computed(() => cart.value?.items || [])
  const total = computed(() => cart.value?.total || 0)
  const itemCount = computed(() => cartItems.value.length)
  
  const fetchCart = async () => {
    loading.value = true
    try {
      const response = await api.get('/api/purchases/cart')
      cart.value = response.data.cart
    } catch (error) {
      console.error('Failed to fetch cart:', error)
      // Initialize empty cart if none exists
      cart.value = { items: [], total: 0 }
    } finally {
      loading.value = false
    }
  }
  
  const addToCart = async (tourId) => {
    try {
      await api.post('/api/purchases/cart/items', { tour_id: tourId })
      await fetchCart() // Refresh cart
      return true
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to add tour to cart'
      throw new Error(message)
    }
  }
  
  const removeFromCart = async (tourId) => {
    try {
      await api.delete(`/api/purchases/cart/items/${tourId}`)
      await fetchCart() // Refresh cart
      return true
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to remove tour from cart'
      throw new Error(message)
    }
  }
  
  const clearCart = async () => {
    try {
      await api.delete('/api/purchases/cart')
      await fetchCart() // Refresh cart
      return true
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to clear cart'
      throw new Error(message)
    }
  }
  
  const checkout = async () => {
    try {
      const response = await api.post('/api/purchases/cart/checkout')
      await fetchCart() // Should be empty after checkout
      return response.data.tokens
    } catch (error) {
      const message = error.response?.data?.error || 'Checkout failed'
      throw new Error(message)
    }
  }
  
  const isInCart = (tourId) => {
    return cartItems.value.some(item => item.tour_id === tourId)
  }
  
  return {
    cart,
    cartItems,
    total,
    itemCount,
    loading,
    fetchCart,
    addToCart,
    removeFromCart,
    clearCart,
    checkout,
    isInCart
  }
})