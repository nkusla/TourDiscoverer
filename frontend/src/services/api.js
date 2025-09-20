import axios from 'axios'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: process.env.NODE_ENV === 'production'
    ? 'http://localhost:8080'
    : 'http://localhost:8080', // Your API Gateway URL
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Handle common error cases
    if (error.response?.status === 401) {
      // Unauthorized - clear token and redirect to login
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    } else if (error.response?.status === 403) {
      // Forbidden
      console.error('Access forbidden')
    } else if (error.response?.status >= 500) {
      // Server error
      console.error('Server error:', error.response.data)
    }

    return Promise.reject(error)
  }
)

export default api
