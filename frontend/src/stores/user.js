import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token'))
  
  const isAuthenticated = computed(() => !!token.value)
  const username = computed(() => user.value?.username)
  const isAdmin = computed(() => user.value?.role === 'admin')
  
  const login = async (credentials) => {
    try {
      const response = await api.post('/api/auth/login', credentials)
      const { token: authToken, user: userData } = response.data
      
      token.value = authToken
      user.value = userData
      
      localStorage.setItem('token', authToken)
      localStorage.setItem('user', JSON.stringify(userData))
      
      // Set token for future API calls
      api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
      
      return userData
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Login failed')
    }
  }
  
  const register = async (userData) => {
    try {
      console.log('Attempting registration with:', userData)
      const response = await api.post('/api/auth/register', userData)
      console.log('Registration response:', response.data)
      
      // After successful registration, automatically log the user in
      if (response.data) {
        console.log('Attempting auto-login after registration')
        const loginResponse = await api.post('/api/auth/login', {
          username: userData.username,
          password: userData.password
        })
        
        console.log('Auto-login response:', loginResponse.data)
        const { token: authToken, user: userInfo } = loginResponse.data
        
        token.value = authToken
        user.value = userInfo
        
        localStorage.setItem('token', authToken)
        localStorage.setItem('user', JSON.stringify(userInfo))
        
        // Set token for future API calls
        api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
        
        return userInfo
      }
      
      return response.data
    } catch (error) {
      console.error('Registration error:', error)
      throw new Error(error.response?.data?.message || 'Registration failed')
    }
  }
  
  const logout = () => {
    user.value = null
    token.value = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    delete api.defaults.headers.common['Authorization']
  }
  
  const loadUserFromStorage = () => {
    try {
      const storedUser = localStorage.getItem('user')
      const storedToken = localStorage.getItem('token')
      
      if (storedUser && storedToken && storedUser !== 'undefined' && storedUser !== 'null') {
        const parsedUser = JSON.parse(storedUser)
        user.value = parsedUser
        token.value = storedToken
        api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      }
    } catch (error) {
      console.warn('Error loading user from storage:', error)
      // Clear corrupted data
      localStorage.removeItem('user')
      localStorage.removeItem('token')
      user.value = null
      token.value = null
      delete api.defaults.headers.common['Authorization']
    }
  }
  
  const updateProfile = async (profileData) => {
    try {
      const response = await api.put('/auth/profile', profileData)
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Profile update failed')
    }
  }
  
  // Initialize store
  loadUserFromStorage()
  
  return {
    user,
    token,
    isAuthenticated,
    username,
    isAdmin,
    login,
    register,
    logout,
    updateProfile,
    loadUserFromStorage
  }
})
