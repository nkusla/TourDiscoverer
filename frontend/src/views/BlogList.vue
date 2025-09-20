<template>
  <div class="blog-list-page">
    <div class="container mt-4">
      <!-- Header -->
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h2>Blog Posts</h2>
        <router-link 
          v-if="userStore.isAuthenticated" 
          to="/blog/create" 
          class="btn btn-primary"
        >
          Create New Blog
        </router-link>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="text-center">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="mt-2">Loading blogs...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="alert alert-danger">
        <h5>Error loading blogs</h5>
        <p>{{ error }}</p>
        <button class="btn btn-outline-danger" @click="loadBlogs">
          Try Again
        </button>
      </div>

      <!-- Empty state -->
      <div v-else-if="blogs.length === 0" class="text-center py-5">
        <div class="mb-4">
          <i class="bi bi-journal-text display-1 text-muted"></i>
        </div>
        <h4>No blogs yet</h4>
        <p class="text-muted">Be the first to share your thoughts!</p>
        <router-link 
          v-if="userStore.isAuthenticated" 
          to="/blog/create" 
          class="btn btn-primary"
        >
          Create First Blog
        </router-link>
      </div>

      <!-- Blogs list -->
      <div v-else class="blogs-container">
        <BlogCard 
          v-for="blog in blogs" 
          :key="blog.id" 
          :blog="blog"
          @blog-updated="handleBlogUpdate"
        />
      </div>

      <!-- Load more button -->
      <div v-if="hasMore && !loading" class="text-center mt-4">
        <button class="btn btn-outline-primary" @click="loadMoreBlogs">
          Load More Blogs
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import BlogCard from '../components/Blog/BlogCard.vue'
import blogService from '../services/blog'

export default {
  name: 'BlogList',
  components: {
    BlogCard
  },
  setup() {
    const userStore = useUserStore()
    const blogs = ref([])
    const loading = ref(false)
    const error = ref('')
    const hasMore = ref(false)

    const loadBlogs = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const response = await blogService.getBlogs()
        blogs.value = response || []
        
        // For now, we don't have pagination, so set hasMore to false
        hasMore.value = false
        
      } catch (err) {
        console.error('Error loading blogs:', err)
        error.value = err.response?.data?.message || 'Failed to load blogs'
        blogs.value = []
      } finally {
        loading.value = false
      }
    }

    const loadMoreBlogs = async () => {
      // This would be implemented if pagination is added to the backend
      console.log('Load more blogs - pagination not implemented yet')
    }

    const handleBlogUpdate = () => {
      // Refresh the blogs list to get updated like counts
      loadBlogs()
    }

    onMounted(() => {
      loadBlogs()
    })

    return {
      userStore,
      blogs,
      loading,
      error,
      hasMore,
      loadBlogs,
      loadMoreBlogs,
      handleBlogUpdate
    }
  }
}
</script>

<style scoped>
.blog-list-page {
  min-height: 80vh;
}

.blogs-container {
  max-width: 800px;
  margin: 0 auto;
}

.spinner-border {
  width: 3rem;
  height: 3rem;
}

.display-1 {
  font-size: 4rem;
}

.bi-journal-text::before {
  content: "📝";
}
</style>