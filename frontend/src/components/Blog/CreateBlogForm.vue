<template>
  <div class="create-blog-form">
    <div class="card">
      <div class="card-header">
        <h4>Create New Blog Post</h4>
      </div>
      <div class="card-body">
        <form @submit.prevent="submitBlog">
          <!-- Title -->
          <div class="mb-3">
            <label for="title" class="form-label">Title *</label>
            <input 
              id="title"
              v-model="form.title" 
              type="text" 
              class="form-control"
              :class="{ 'is-invalid': errors.title }"
              placeholder="Enter blog title"
              required
            >
            <div v-if="errors.title" class="invalid-feedback">
              {{ errors.title }}
            </div>
          </div>

          <!-- Description -->
          <div class="mb-3">
            <label for="description" class="form-label">Description *</label>
            <textarea 
              id="description"
              v-model="form.description" 
              class="form-control"
              :class="{ 'is-invalid': errors.description }"
              rows="6"
              placeholder="Write your blog content here... (Markdown supported)"
              required
            ></textarea>
            <div v-if="errors.description" class="invalid-feedback">
              {{ errors.description }}
            </div>
            <small class="form-text text-muted">
              You can use Markdown formatting (e.g., **bold**, *italic*, # headers)
            </small>
          </div>

          <!-- Images -->
          <div class="mb-3">
            <label for="images" class="form-label">Images (Optional)</label>
            <div class="input-group mb-2">
              <input 
                v-model="imageUrl"
                type="url" 
                class="form-control"
                placeholder="Enter image URL"
              >
              <button 
                type="button" 
                class="btn btn-outline-secondary"
                @click="addImage"
                :disabled="!imageUrl.trim()"
              >
                Add Image
              </button>
            </div>
            
            <!-- Image preview -->
            <div v-if="form.images.length > 0" class="images-preview">
              <div class="row">
                <div v-for="(image, index) in form.images" :key="index" class="col-md-3 mb-2">
                  <div class="position-relative">
                    <img :src="image" :alt="`Preview ${index + 1}`" class="img-fluid rounded">
                    <button 
                      type="button"
                      class="btn btn-sm btn-danger position-absolute top-0 end-0"
                      @click="removeImage(index)"
                    >
                      ×
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Form actions -->
          <div class="d-flex justify-content-between">
            <button 
              type="button" 
              class="btn btn-secondary"
              @click="resetForm"
            >
              Reset
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="loading || !isFormValid"
            >
              <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
              {{ loading ? 'Creating...' : 'Create Blog' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import blogService from '../../services/blog'

export default {
  name: 'CreateBlogForm',
  emits: ['blog-created'],
  setup(props, { emit }) {
    const router = useRouter()
    const userStore = useUserStore()
    
    const loading = ref(false)
    const imageUrl = ref('')
    
    const form = reactive({
      title: '',
      description: '',
      images: []
    })
    
    const errors = reactive({
      title: '',
      description: ''
    })

    const isFormValid = computed(() => {
      return form.title.trim() && form.description.trim()
    })

    const validateForm = () => {
      errors.title = ''
      errors.description = ''
      
      if (!form.title.trim()) {
        errors.title = 'Title is required'
        return false
      }
      
      if (form.title.length < 3) {
        errors.title = 'Title must be at least 3 characters long'
        return false
      }
      
      if (!form.description.trim()) {
        errors.description = 'Description is required'
        return false
      }
      
      if (form.description.length < 10) {
        errors.description = 'Description must be at least 10 characters long'
        return false
      }
      
      return true
    }

    const addImage = () => {
      if (imageUrl.value.trim() && !form.images.includes(imageUrl.value.trim())) {
        form.images.push(imageUrl.value.trim())
        imageUrl.value = ''
      }
    }

    const removeImage = (index) => {
      form.images.splice(index, 1)
    }

    const resetForm = () => {
      form.title = ''
      form.description = ''
      form.images = []
      imageUrl.value = ''
      errors.title = ''
      errors.description = ''
    }

    const submitBlog = async () => {
      if (!validateForm()) return

      loading.value = true
      
      try {
        const blogData = {
          title: form.title.trim(),
          description: form.description.trim(),
          images: form.images.length > 0 ? form.images : undefined
        }
        
        const newBlog = await blogService.createBlog(blogData)
        
        // Emit success event
        emit('blog-created', newBlog)
        
        // Reset form
        resetForm()
        
        // Show success message
        alert('Blog created successfully!')
        
        // Navigate to blogs page
        router.push('/blogs')
        
      } catch (error) {
        console.error('Error creating blog:', error)
        
        if (error.response?.status === 401) {
          alert('You must be logged in to create a blog')
        } else if (error.response?.status === 400) {
          alert('Invalid blog data. Please check your input.')
        } else {
          alert('Failed to create blog. Please try again.')
        }
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      errors,
      loading,
      imageUrl,
      isFormValid,
      addImage,
      removeImage,
      resetForm,
      submitBlog
    }
  }
}
</script>

<style scoped>
.create-blog-form {
  max-width: 800px;
  margin: 0 auto;
}

.images-preview img {
  max-height: 150px;
  object-fit: cover;
}

.position-relative .btn {
  margin: 5px;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}

textarea {
  resize: vertical;
  min-height: 120px;
}
</style>