<template>
  <div class="blog-card card mb-4">
    <div class="card-header">
      <div class="d-flex justify-content-between align-items-center">
        <h5 class="card-title mb-0">{{ blog.title }}</h5>
        <small class="text-muted">{{ formatDate(blog.created_at) }}</small>
      </div>
      <small class="text-muted">by {{ blog.author }}</small>
    </div>
    
    <div class="card-body">
      <p class="card-text">{{ blog.description }}</p>
      
      <!-- Images if any -->
      <div v-if="blog.images && blog.images.length > 0" class="blog-images mb-3">
        <div class="row">
          <div v-for="(image, index) in blog.images" :key="index" class="col-md-4 mb-2">
            <img :src="image" :alt="`Blog image ${index + 1}`" class="img-fluid rounded">
          </div>
        </div>
      </div>

      <!-- Like section -->
      <div class="blog-actions d-flex align-items-center mb-3">
        <button 
          class="btn btn-sm me-2"
          :class="isLiked ? 'btn-danger' : 'btn-outline-danger'"
          @click="toggleLike"
          :disabled="loading"
        >
          <i class="bi bi-heart-fill" v-if="isLiked"></i>
          <i class="bi bi-heart" v-else></i>
          {{ likeCount }}
        </button>
        <span class="text-muted">{{ likeCount }} likes</span>
      </div>

      <!-- Comments section -->
      <div class="comments-section">
        <h6>Comments ({{ comments?.length || 0 }})</h6>
        
        <!-- Add comment form -->
        <div v-if="canComment" class="add-comment mb-3">
          <div class="input-group">
            <input 
              v-model="newComment" 
              type="text" 
              class="form-control" 
              placeholder="Add a comment..."
              @keypress.enter="addComment"
            >
            <button 
              class="btn btn-primary" 
              @click="addComment"
              :disabled="!newComment.trim() || loading"
            >
              Post
            </button>
          </div>
        </div>

        <!-- Comments list -->
        <div v-if="comments && comments.length > 0" class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment mb-2 p-2 border rounded">
            <div class="d-flex justify-content-between">
              <strong>{{ comment.author }}</strong>
              <small class="text-muted">{{ formatDate(comment.created_at) }}</small>
            </div>
            <p class="mb-0">{{ comment.text }}</p>
          </div>
        </div>
        
        <div v-else class="text-muted">
          No comments yet. Be the first to comment!
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useUserStore } from '../../stores/user'
import blogService from '../../services/blog'

export default {
  name: 'BlogCard',
  props: {
    blog: {
      type: Object,
      required: true
    }
  },
  emits: ['blog-updated'],
  setup(props, { emit }) {
    const userStore = useUserStore()
    const comments = ref([])
    const newComment = ref('')
    const loading = ref(false)
    
    // Create reactive copies for like status and count
    const isLiked = ref(props.blog.is_liked_by_user || false)
    const likeCount = ref(props.blog.like_count || 0)
    
    console.log('Blog prop received:', props.blog)
    console.log('Initial isLiked:', props.blog.is_liked_by_user)
    
    // Watch for prop changes in case blog is updated externally
    watch(() => props.blog.is_liked_by_user, (newValue) => {
      console.log('is_liked_by_user changed to:', newValue)
      isLiked.value = newValue || false
    })
    
    watch(() => props.blog.like_count, (newValue) => {
      likeCount.value = newValue || 0
    })

    const canComment = computed(() => {
      return userStore.isAuthenticated
    })

    const formatDate = (timestamp) => {
      if (!timestamp) return ''
      const date = new Date(timestamp * 1000) // Convert from Unix timestamp
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
    }

    const toggleLike = async () => {
      if (!userStore.isAuthenticated) {
        alert('Please login to like posts')
        return
      }

      loading.value = true
      try {
        let response
        const wasLiked = isLiked.value
        
        if (wasLiked) {
          response = await blogService.unlikeBlog(props.blog.id)
        } else {
          response = await blogService.likeBlog(props.blog.id)
        }
        
        // Update local reactive state to avoid page refresh
        isLiked.value = !wasLiked
        likeCount.value = wasLiked ? 
          Math.max(0, likeCount.value - 1) : 
          likeCount.value + 1
        
      } catch (error) {
        console.error('Error toggling like:', error)
        alert('Failed to update like status')
      } finally {
        loading.value = false
      }
    }

    const addComment = async () => {
      if (!newComment.value.trim()) return

      loading.value = true
      try {
        const commentData = {
          content: newComment.value.trim()
        }
        
        await blogService.addComment(props.blog.id, commentData)
        newComment.value = ''
        await loadComments() // Refresh comments
      } catch (error) {
        console.error('Error adding comment:', error)
        alert('Failed to add comment')
      } finally {
        loading.value = false
      }
    }

    const loadComments = async () => {
      try {
        const result = await blogService.getComments(props.blog.id)
        comments.value = result || []
      } catch (error) {
        console.error('Error loading comments:', error)
        comments.value = []
      }
    }

    onMounted(() => {
      loadComments()
    })

    return {
      comments,
      newComment,
      loading,
      canComment,
      isLiked,
      likeCount,
      formatDate,
      toggleLike,
      addComment
    }
  }
}
</script>

<style scoped>
.blog-card {
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.blog-images img {
  max-height: 200px;
  object-fit: cover;
}

.comment {
  background-color: #f8f9fa;
}

.btn:disabled {
  opacity: 0.6;
}
</style>