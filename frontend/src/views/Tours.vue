<template>
  <div class="tours">
    <div class="container py-4">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h2>Tours</h2>
        <router-link to="/tour/create" class="btn btn-primary">
          <i class="fas fa-plus"></i> Create New Tour
        </router-link>
      </div>
      
      <!-- Filters -->
      <div class="row mb-4">
        <div class="col-md-4">
          <input 
            v-model="searchQuery" 
            type="text" 
            class="form-control" 
            placeholder="Search tours..."
          />
        </div>
        <div class="col-md-3">
          <select v-model="difficultyFilter" class="form-select">
            <option value="">All Difficulties</option>
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
            <option value="extreme">Extreme</option>
          </select>
        </div>
        <div class="col-md-3">
          <select v-model="statusFilter" class="form-select">
            <option value="">All Statuses</option>
            <option value="draft">Draft</option>
            <option value="published">Published</option>
          </select>
        </div>
        <div class="col-md-2">
          <button class="btn btn-outline-secondary w-100" @click="clearFilters">
            Clear
          </button>
        </div>
      </div>
      
      <!-- Loading state -->
      <div v-if="loading" class="text-center py-5">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      
      <!-- Tours grid -->
      <div v-else-if="filteredTours.length > 0" class="row">
        <div 
          v-for="tour in filteredTours" 
          :key="tour.id" 
          class="col-lg-4 col-md-6 mb-4"
        >
          <div class="card h-100">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-start mb-2">
                <h5 class="card-title">{{ tour.name }}</h5>
                <span 
                  class="badge"
                  :class="getStatusBadgeClass(tour.status)"
                >
                  {{ tour.status }}
                </span>
              </div>
              
              <p class="card-text text-muted">
                {{ truncateText(tour.description, 100) }}
              </p>
              
              <div class="mb-3">
                <small class="text-muted">
                  <strong>Difficulty:</strong> 
                  <span class="text-capitalize">{{ tour.difficulty }}</span>
                </small>
                <br>
                <small class="text-muted">
                  <strong>Distance:</strong> {{ tour.distance?.toFixed(1) || 0 }} km
                </small>
                <br>
                <small class="text-muted">
                  <strong>Key Points:</strong> {{ tour.key_points?.length || 0 }}
                </small>
                <br>
                <small class="text-muted">
                  <strong>Price:</strong> ${{ tour.price || 0 }}
                </small>
              </div>
              
              <div class="mb-2">
                <small class="text-muted">
                  <strong>Tags:</strong> {{ tour.tags }}
                </small>
              </div>
              
              <div class="mb-3">
                <small class="text-muted">
                  By {{ tour.author_username }} • 
                  {{ formatDate(tour.created_at) }}
                </small>
              </div>
            </div>
            
            <div class="card-footer bg-transparent">
              <div class="btn-group w-100" role="group">
                <button 
                  class="btn btn-outline-primary"
                  @click="viewTour(tour)"
                >
                  View
                </button>
                <router-link 
                  :to="`/tour/edit/${tour.id}`"
                  class="btn btn-outline-secondary"
                  v-if="canEdit(tour)"
                >
                  Edit
                </router-link>
                <button 
                  class="btn btn-outline-danger"
                  @click="deleteTour(tour)"
                  v-if="canDelete(tour)"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty state -->
      <div v-else class="text-center py-5">
        <i class="fas fa-map fa-3x text-muted mb-3"></i>
        <h4 class="text-muted">No tours found</h4>
        <p class="text-muted">
          {{ searchQuery || difficultyFilter || statusFilter 
             ? 'Try adjusting your filters' 
             : 'Create your first tour to get started' }}
        </p>
        <router-link to="/tour/create" class="btn btn-primary" v-if="!searchQuery && !difficultyFilter && !statusFilter">
          Create Tour
        </router-link>
      </div>
    </div>
    
    <!-- View Tour Modal -->
    <div 
      class="modal fade" 
      id="viewTourModal" 
      tabindex="-1"
      ref="viewModal"
    >
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ selectedTour?.name }}</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <div class="row" v-if="selectedTour">
              <div class="col-md-6">
                <h6>Tour Details</h6>
                <p><strong>Description:</strong> {{ selectedTour.description }}</p>
                <p><strong>Difficulty:</strong> {{ selectedTour.difficulty }}</p>
                <p><strong>Distance:</strong> {{ selectedTour.distance?.toFixed(2) }} km</p>
                <p><strong>Price:</strong> ${{ selectedTour.price }}</p>
                <p><strong>Tags:</strong> {{ selectedTour.tags }}</p>
                <p><strong>Status:</strong> {{ selectedTour.status }}</p>
                
                <h6 class="mt-4">Key Points</h6>
                <div class="list-group">
                  <div 
                    v-for="(point, index) in selectedTour.key_points" 
                    :key="point.id"
                    class="list-group-item"
                  >
                    <strong>{{ index + 1 }}. {{ point.name }}</strong>
                    <p class="mb-1">{{ point.description }}</p>
                    <small>{{ point.latitude.toFixed(4) }}, {{ point.longitude.toFixed(4) }}</small>
                  </div>
                </div>
              </div>
              
              <div class="col-md-6">
                <h6>Tour Route</h6>
                <LeafletMap
                  :key-points="selectedTour.key_points || []"
                  :editable="false"
                  map-height="400px"
                />
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
              Close
            </button>
            <router-link 
              :to="`/tour/edit/${selectedTour?.id}`"
              class="btn btn-primary"
              v-if="selectedTour && canEdit(selectedTour)"
              data-bs-dismiss="modal"
            >
              Edit Tour
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useTourStore } from '../stores/tour'
import { useUserStore } from '../stores/user'
import LeafletMap from '../components/Map/LeafletMap.vue'
import { Modal } from 'bootstrap'

export default {
  name: 'Tours',
  components: {
    LeafletMap
  },
  setup() {
    const tourStore = useTourStore()
    const userStore = useUserStore()
    
    const loading = ref(false)
    const searchQuery = ref('')
    const difficultyFilter = ref('')
    const statusFilter = ref('')
    const selectedTour = ref(null)
    const viewModal = ref(null)
    const viewModalInstance = ref(null)
    
    const filteredTours = computed(() => {
      let tours = tourStore.tours
      
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        tours = tours.filter(tour => 
          tour.name.toLowerCase().includes(query) ||
          tour.description.toLowerCase().includes(query) ||
          tour.tags.toLowerCase().includes(query)
        )
      }
      
      if (difficultyFilter.value) {
        tours = tours.filter(tour => tour.difficulty === difficultyFilter.value)
      }
      
      if (statusFilter.value) {
        tours = tours.filter(tour => tour.status === statusFilter.value)
      }
      
      return tours
    })
    
    onMounted(async () => {
      // Initialize Bootstrap modal
      if (viewModal.value) {
        viewModalInstance.value = new Modal(viewModal.value)
      }
      
      // Load tours
      await loadTours()
    })
    
    const loadTours = async () => {
      loading.value = true
      try {
        await tourStore.fetchTours()
      } catch (error) {
        console.error('Failed to load tours:', error)
      } finally {
        loading.value = false
      }
    }
    
    const viewTour = (tour) => {
      selectedTour.value = tour
      viewModalInstance.value?.show()
    }
    
    const deleteTour = async (tour) => {
      if (confirm(`Are you sure you want to delete "${tour.name}"?`)) {
        try {
          await tourStore.deleteTour(tour.id)
          // Reload tours
          await loadTours()
        } catch (error) {
          console.error('Failed to delete tour:', error)
        }
      }
    }
    
    const canEdit = (tour) => {
      return userStore.isAuthenticated && 
             (userStore.username === tour.author_username || userStore.isAdmin)
    }
    
    const canDelete = (tour) => {
      return canEdit(tour)
    }
    
    const clearFilters = () => {
      searchQuery.value = ''
      difficultyFilter.value = ''
      statusFilter.value = ''
    }
    
    const getStatusBadgeClass = (status) => {
      switch (status) {
        case 'published':
          return 'bg-success'
        case 'draft':
          return 'bg-secondary'
        default:
          return 'bg-secondary'
      }
    }
    
    const truncateText = (text, maxLength) => {
      if (!text) return ''
      if (text.length <= maxLength) return text
      return text.substring(0, maxLength) + '...'
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return ''
      return new Date(dateString).toLocaleDateString()
    }
    
    return {
      loading,
      searchQuery,
      difficultyFilter,
      statusFilter,
      selectedTour,
      viewModal,
      filteredTours,
      viewTour,
      deleteTour,
      canEdit,
      canDelete,
      clearFilters,
      getStatusBadgeClass,
      truncateText,
      formatDate
    }
  }
}
</script>

<style scoped>
.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-2px);
}

.modal-xl .modal-body {
  max-height: 70vh;
  overflow-y: auto;
}

.list-group-item p {
  font-size: 0.9em;
  color: #6c757d;
}
</style>
