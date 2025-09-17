<template>
  <div class="modal fade" :id="modalId" tabindex="-1" ref="modal">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            Tour Reviews
            <span v-if="reviewsData" class="text-muted">
              ({{ reviewsData.total_count }} reviews - Avg: {{ reviewsData.average_rating?.toFixed(1) }}/5)
            </span>
          </h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>
        <div class="modal-body" style="max-height: 500px; overflow-y: auto;">
          <!-- Loading -->
          <div v-if="loading" class="text-center py-4">
            <div class="spinner-border" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>

          <!-- No Reviews -->
          <div v-else-if="!reviews.length" class="text-center py-4 text-muted">
            <i class="fas fa-comments fa-2x mb-3"></i>
            <h6>No reviews yet</h6>
            <p>Be the first to review this tour!</p>
          </div>

          <!-- Reviews List -->
          <div v-else>
            <div v-for="review in reviews" :key="review.id" class="review-item mb-4 pb-3 border-bottom">
              <!-- Review Header -->
              <div class="d-flex justify-content-between align-items-start mb-2">
                <div>
                  <h6 class="mb-1">{{ review.tourist_username }}</h6>
                  <div class="rating mb-1">
                    <span v-for="star in 5" :key="star">
                      <i 
                        class="fas fa-star"
                        :class="star <= review.rating ? 'text-warning' : 'text-muted'"
                      ></i>
                    </span>
                    <span class="ms-2 text-muted">({{ review.rating }}/5)</span>
                  </div>
                  <small class="text-muted">
                    Visited: {{ formatDate(review.visit_date) }} • 
                    Reviewed: {{ formatDate(review.review_date) }}
                  </small>
                </div>
              </div>

              <!-- Review Content -->
              <p class="mb-2">{{ review.comment }}</p>
              
              <!-- Images -->
              <div v-if="review.images && review.images.length" class="review-images">
                <div class="row g-2">
                  <div v-for="image in review.images" :key="image.id" class="col-auto">
                    <img 
                      :src="image.image_url" 
                      alt="Review image"
                      class="img-thumbnail"
                      style="width: 80px; height: 80px; object-fit: cover; cursor: pointer;"
                      @click="showImageModal(image.image_url)"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Pagination -->
            <div v-if="totalPages > 1" class="text-center mt-3">
              <button 
                class="btn btn-sm btn-outline-primary me-2"
                @click="changePage(currentPage - 1)"
                :disabled="currentPage === 1"
              >
                Previous
              </button>
              <span class="mx-2">Page {{ currentPage }} of {{ totalPages }}</span>
              <button 
                class="btn btn-sm btn-outline-primary ms-2"
                @click="changePage(currentPage + 1)"
                :disabled="currentPage === totalPages"
              >
                Next
              </button>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Image Modal -->
  <div class="modal fade" id="imageViewModal" tabindex="-1">
    <div class="modal-dialog modal-lg modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Review Image</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>
        <div class="modal-body text-center">
          <img :src="selectedImage" alt="Review image" class="img-fluid" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { reviewService } from '../../services/review'

export default {
  name: 'ReviewList',
  props: {
    tourId: {
      type: Number,
      required: true
    },
    modalId: {
      type: String,
      default: 'reviewListModal'
    }
  },
  setup(props) {
    const reviews = ref([])
    const reviewsData = ref(null)
    const loading = ref(false)
    const currentPage = ref(1)
    const pageSize = ref(5)
    const selectedImage = ref('')

    const totalPages = computed(() => {
      if (!reviewsData.value?.total_count) return 1
      return Math.ceil(reviewsData.value.total_count / pageSize.value)
    })

    const loadReviews = async () => {
      if (!props.tourId) return
      
      loading.value = true
      try {
        const data = await reviewService.getTourReviews(props.tourId, currentPage.value, pageSize.value)
        reviews.value = data.reviews || []
        reviewsData.value = data
      } catch (error) {
        console.error('Failed to load reviews:', error)
        reviews.value = []
        reviewsData.value = null
      } finally {
        loading.value = false
      }
    }

    const changePage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadReviews()
      }
    }

    const showImageModal = (imageUrl) => {
      selectedImage.value = imageUrl
      const modal = new bootstrap.Modal(document.getElementById('imageViewModal'))
      modal.show()
    }

    const formatDate = (dateString) => {
      if (!dateString) return ''
      return new Date(dateString).toLocaleDateString()
    }

    // Watch for tour ID changes
    watch(() => props.tourId, () => {
      currentPage.value = 1
      loadReviews()
    }, { immediate: true })

    // Expose method to refresh reviews
    const refresh = () => {
      currentPage.value = 1
      loadReviews()
    }

    return {
      reviews,
      reviewsData,
      loading,
      currentPage,
      totalPages,
      selectedImage,
      changePage,
      showImageModal,
      formatDate,
      refresh
    }
  }
}
</script>

<style scoped>
.review-item:last-child {
  border-bottom: none !important;
}

.rating {
  font-size: 0.9em;
}

.review-images img {
  border-radius: 4px;
  transition: transform 0.2s;
}

.review-images img:hover {
  transform: scale(1.05);
}
</style>