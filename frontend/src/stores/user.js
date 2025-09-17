import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import { getUserFromToken, isTokenExpired } from '../utils/jwt'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token'))

  const isAuthenticated = computed(() => !!token.value)
  const username = computed(() => user.value?.username)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isGuide = computed(() => user.value?.role === 'guide')
  const isTourist = computed(() => user.value?.role === 'tourist')
  const canCreateTours = computed(() => user.value?.role === 'guide' || user.value?.role === 'admin')

  const login = async (credentials) => {
    try {
      const response = await api.post('/api/auth/login', credentials)
      const { token: authToken } = response.data

      // Extract user data from JWT token
      const userData = getUserFromToken(authToken)
      if (!userData) {
        throw new Error('Invalid token received')
      }

      token.value = authToken
      user.value = userData

      localStorage.setItem('token', authToken)
      // No need to store user separately since we decode from token

      // Set token for future API calls
      api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`

      return userData
    } catch (error) {
      if (error.response?.status === 401) {
        throw new Error('Invalid username or password')
      }
      else if (error.response?.status === 403) {
        throw new Error('Your account is blocked. Please contact support.')
      }

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
        const { token: authToken } = loginResponse.data

        // Extract user data from JWT token
        const userInfo = getUserFromToken(authToken)
        if (!userInfo) {
          throw new Error('Invalid token received after registration')
        }

        token.value = authToken
        user.value = userInfo

        localStorage.setItem('token', authToken)
        // No need to store user separately since we decode from token

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
    // No need to remove user since we only store token

    delete api.defaults.headers.common['Authorization']
  }

  const loadUserFromStorage = () => {
    try {
      const storedToken = localStorage.getItem('token')

      if (storedToken) {
        // Check if token is expired
        if (isTokenExpired(storedToken)) {
          console.warn('Stored token is expired, clearing storage')
          localStorage.removeItem('token')
          user.value = null
          token.value = null
          delete api.defaults.headers.common['Authorization']
          return
        }

        // Extract user data from JWT token
        const userData = getUserFromToken(storedToken)
        if (userData) {
          user.value = userData
          token.value = storedToken
          api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
        } else {
          console.warn('Could not extract user data from token')
          localStorage.removeItem('token')
          user.value = null
          token.value = null
          delete api.defaults.headers.common['Authorization']
        }
      }
    } catch (error) {
      console.warn('Error loading user from storage:', error)
      // Clear corrupted data
      localStorage.removeItem('token')
      user.value = null
      token.value = null
      delete api.defaults.headers.common['Authorization']
    }
  }

  const updateProfile = async (profileData) => {
    try {
      const response = await api.put('/auth/profile', profileData)
      // Note: Profile updates would need to be reflected in a new token
      // For now, we'll just return the response but user data comes from token
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Profile update failed')
    }
  }

  const fetchUserProfile = async () => {
    try {
      if (!user.value?.username) {
        throw new Error('No authenticated user')
      }

      // Try to fetch full profile from stakeholder service
      const response = await api.get('/api/stakeholder/profile')
      return response.data
    } catch (error) {
      console.error('Error fetching user profile:', error)

      // If profile doesn't exist (404), create one
      if (error.response?.status === 404) {
        try {
          console.log('Profile not found, creating new profile...')
          await api.post('/api/stakeholder', {
            username: user.value.username,
            first_name: '',
            last_name: '',
            profile_picture: '',
            biography: '',
            motto: ''
          })

          // Try fetching again after creation
          const response = await api.get('/api/stakeholder/profile')
          return response.data
        } catch (createError) {
          console.error('Error creating profile:', createError)
          throw new Error('Failed to create user profile')
        }
      }

      throw new Error(error.response?.data?.message || 'Failed to fetch user profile')
    }
  }

  const updateUserProfile = async (profileData, file = null) => {
    try {
      if (!user.value?.username) {
        throw new Error('No authenticated user')
      }

      let response

      if (file) {
        // Create FormData for file upload
        const formData = new FormData()
        formData.append('first_name', profileData.first_name || '')
        formData.append('last_name', profileData.last_name || '')
        formData.append('biography', profileData.biography || '')
        formData.append('motto', profileData.motto || '')
        formData.append('profile_picture', file)

        response = await api.put('/api/stakeholder/profile', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
      } else {
        // Regular JSON update (backward compatibility)
        response = await api.put('/api/stakeholder/profile', profileData)
      }

      return response.data
    } catch (error) {
      console.error('Error updating user profile:', error)
      throw new Error(error.response?.data?.message || 'Failed to update user profile')
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
    isGuide,
    isTourist,
    canCreateTours,
    login,
    register,
    logout,
    updateProfile,
    fetchUserProfile,
    updateUserProfile,
    loadUserFromStorage
  }
})
