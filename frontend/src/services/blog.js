import api from './api'

// Blog service for API calls
const blogService = {
  // Get all blogs (requires authentication)
  async getBlogs() {
    try {
      // Proverava da li je korisnik ulogovan - blogovi se mogu videti samo ako si ulogovan
      const token = localStorage.getItem('token'); 
      if (!token) {
        throw new Error('Authentication required to view blogs');
      }
      
      const endpoint = '/api/blogs/personalized';
      console.log('Using endpoint:', endpoint);
      
      const response = await api.get(endpoint)
      console.log('Blog response from backend:', response.data)
      return response.data
    } catch (error) {
      console.error('Error fetching blogs:', error)
      throw error
    }
  },

  // Get blog by ID
  async getBlogById(id) {
    try {
      const response = await api.get(`/api/blogs/${id}`)
      return response.data
    } catch (error) {
      console.error('Error fetching blog:', error)
      throw error
    }
  },

  // Create a new blog
  async createBlog(blogData) {
    try {
      const response = await api.post('/api/blogs', blogData)
      return response.data
    } catch (error) {
      console.error('Error creating blog:', error)
      throw error
    }
  },

  // Like a blog
  async likeBlog(blogId) {
    try {
      const response = await api.post('/api/blogs/like', { blog_id: blogId })
      return response.data
    } catch (error) {
      console.error('Error liking blog:', error)
      throw error
    }
  },

  // Unlike a blog (same endpoint, it toggles)
  async unlikeBlog(blogId) {
    try {
      const response = await api.post('/api/blogs/like', { blog_id: blogId })
      return response.data
    } catch (error) {
      console.error('Error unliking blog:', error)
      throw error
    }
  },

  // Get comments for a blog
  async getComments(blogId) {
    try {
      const response = await api.get(`/api/blogs/comments?blog_id=${blogId}`)
      return response.data
    } catch (error) {
      console.error('Error fetching comments:', error)
      throw error
    }
  },

  // Add comment to a blog
  async addComment(blogId, commentData) {
    try {
      const commentWithBlogId = {
        blog_id: blogId,
        text: commentData.content  // Convert 'content' to 'text' for backend
      }
      const response = await api.post('/api/blogs/comment', commentWithBlogId)
      return response.data
    } catch (error) {
      console.error('Error adding comment:', error)
      throw error
    }
  }
}

export default blogService