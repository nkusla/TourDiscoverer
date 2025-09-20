<template>
  <div class="create-blog-page">
    <div class="container mt-4">
      <!-- Header -->
      <div class="d-flex align-items-center mb-4">
        <router-link to="/blogs" class="btn btn-outline-secondary me-3">
          ← Back to Blogs
        </router-link>
        <h2>Create New Blog Post</h2>
      </div>

      <!-- Check authentication -->
      <div v-if="!userStore.isAuthenticated" class="alert alert-warning">
        <h5>Authentication Required</h5>
        <p>You must be logged in to create a blog post.</p>
        <router-link to="/login" class="btn btn-primary">
          Go to Login
        </router-link>
      </div>

      <!-- Create blog form -->
      <div v-else>
        <CreateBlogForm @blog-created="handleBlogCreated" />
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import CreateBlogForm from '../components/Blog/CreateBlogForm.vue'

export default {
  name: 'CreateBlog',
  components: {
    CreateBlogForm
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()

    const handleBlogCreated = (newBlog) => {
      console.log('New blog created:', newBlog)
      
      // Navigate to blogs list
      router.push('/blogs')
    }

    return {
      userStore,
      handleBlogCreated
    }
  }
}
</script>

<style scoped>
.create-blog-page {
  min-height: 80vh;
}
</style>