<template>
  <div class="modal fade" :id="modalId" tabindex="-1" ref="modal">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Leave a Review</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitReview">
            <!-- Rating -->
            <div class="mb-3">
              <label class="form-label">Rating *</label>
              <div class="rating-input">
                <div class="form-check form-check-inline" v-for="star in 5" :key="star">
                  <input 
                    class="form-check-input" 
                    type="radio" 
                    :id="'rating' + star"
                    :value="star" 
                    v-model="rating"
                  >
                  <label class="form-check-label" :for="'rating' + star">
                    {{ star }} <i class="fas fa-star text-warning"></i>
                  </label>
                </div>
              </div>
              <div class="form-text">Selected: {{ rating || 'None' }}/5</div>
            </div>

            <!-- Visit Date -->
            <div class="mb-3">
              <label class="form-label">Visit Date *</label>
              <input
                v-model="visitDate"
                type="date"
                class="form-control"
                :max="today"
                required
              />
            </div>

            <!-- Comment -->
            <div class="mb-3">
              <label class="form-label">Comment *</label>
              <textarea
                v-model="comment"
                class="form-control"
                rows="4"
                placeholder="Share your experience..."
                required
              ></textarea>
            </div>

            <!-- Images -->
            <div class="mb-3">
              <label class="form-label">Images (Optional)</label>
              <input
                ref="fileInput"
                type="file"
                class="form-control"
                accept="image/*"
                multiple
                @change="handleFileChange"
              />
              <div class="form-text">You can select up to 5 images (JPG, PNG, GIF)</div>
              
              <!-- Image previews -->
              <div v-if="imagePreviews.length" class="mt-2">
                <div v-for="(preview, index) in imagePreviews" :key="index" class="d-inline-block me-2 mb-2">
                  <div class="position-relative">
                    <img :src="preview" alt="Preview" class="img-thumbnail" style="width: 80px; height: 80px; object-fit: cover;">
                    <button type="button" class="btn btn-sm btn-danger position-absolute top-0 end-0" style="transform: translate(50%, -50%);" @click="removeImage(index)">×</button>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary" 
            @click="submitReview"
            :disabled="!canSubmit || submitting"
          >
            <span v-if="submitting" class="spinner-border spinner-border-sm me-2"></span>
            Submit Review
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { reviewService } from '../../services/review'
import * as bootstrap from 'bootstrap'

export default {
  name: 'ReviewForm',
  props: {
    tourId: {
      type: Number,
      required: true
    },
    modalId: {
      type: String,
      default: 'reviewModal'
    }
  },
  emits: ['review-submitted'],
  setup(props, { emit }) {
    const rating = ref(0)
    const visitDate = ref('')
    const comment = ref('')
    const selectedFiles = ref([])
    const imagePreviews = ref([])
    const submitting = ref(false)
    const fileInput = ref(null)

    const today = new Date().toISOString().split('T')[0]

    const canSubmit = computed(() => {
      return rating.value > 0 && visitDate.value && comment.value.trim().length > 0
    })

    const handleFileChange = (event) => {
      const files = Array.from(event.target.files)
      if (files.length > 5) {
        alert('You can only select up to 5 images')
        return
      }

      selectedFiles.value = files
      imagePreviews.value = []

      files.forEach(file => {
        if (file.type.startsWith('image/')) {
          const reader = new FileReader()
          reader.onload = (e) => {
            imagePreviews.value.push(e.target.result)
          }
          reader.readAsDataURL(file)
        }
      })
    }

    const removeImage = (index) => {
      imagePreviews.value.splice(index, 1)
      const newFiles = Array.from(selectedFiles.value)
      newFiles.splice(index, 1)
      selectedFiles.value = newFiles
      
      // Update file input
      if (fileInput.value) {
        const dt = new DataTransfer()
        newFiles.forEach(file => dt.items.add(file))
        fileInput.value.files = dt.files
      }
    }

    const resetForm = () => {
      rating.value = 0
      visitDate.value = ''
      comment.value = ''
      selectedFiles.value = []
      imagePreviews.value = []
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }

    const submitReview = async () => {
      if (!canSubmit.value) return

      submitting.value = true
      try {
        // Convert files to base64 if any
        const imageUrls = []
        for (let file of selectedFiles.value) {
          if (file.type.startsWith('image/')) {
            const base64 = await fileToBase64(file)
            imageUrls.push(base64)
          }
        }

        const reviewData = {
          tour_id: props.tourId,
          rating: parseInt(rating.value),
          comment: comment.value.trim(),
          visit_date: visitDate.value,
          images: imageUrls
        }

        await reviewService.createReview(reviewData)
        emit('review-submitted')
        resetForm()
        
        // Close modal
        const modal = bootstrap.Modal.getInstance(document.getElementById(props.modalId))
        modal?.hide()
        
        alert('Review submitted successfully!')
      } catch (error) {
        console.error('Failed to submit review:', error)
        alert('Failed to submit review. Please try again.')
      } finally {
        submitting.value = false
      }
    }

    const fileToBase64 = (file) => {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => resolve(reader.result)
        reader.onerror = error => reject(error)
      })
    }

    return {
      rating,
      visitDate,
      comment,
      selectedFiles,
      imagePreviews,
      submitting,
      fileInput,
      today,
      canSubmit,
      handleFileChange,
      removeImage,
      submitReview
    }
  }
}
</script>

<style scoped>
.rating-input .form-check {
  margin-right: 1rem;
}

.rating-input .form-check-label {
  cursor: pointer;
  font-weight: 500;
}

.rating-input .form-check-input:checked + .form-check-label {
  color: #ffc107;
  font-weight: bold;
}

.position-relative {
  display: inline-block;
}
</style>