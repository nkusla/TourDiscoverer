<template>
  <div class="tour-editor">
    <div class="container-fluid">
      <div class="row">
        <!-- Sidebar for tour details -->
        <div class="col-md-4 col-lg-3 bg-light border-end">
          <div class="p-3">
            <h4>{{ isEditing ? 'Edit Tour' : 'Create Tour' }}</h4>
            
            <form @submit.prevent="saveTour">
              <!-- Basic Tour Info -->
              <div class="mb-3">
                <label class="form-label">Tour Name</label>
                <input 
                  v-model="tour.name" 
                  type="text" 
                  class="form-control"
                  required
                  placeholder="Enter tour name"
                />
              </div>
              
              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea 
                  v-model="tour.description" 
                  class="form-control"
                  rows="3"
                  required
                  placeholder="Describe your tour"
                ></textarea>
              </div>
              
              <div class="row">
                <div class="col-6">
                  <div class="mb-3">
                    <label class="form-label">Difficulty</label>
                    <select v-model="tour.difficulty" class="form-select" required>
                      <option value="">Select...</option>
                      <option value="easy">Easy</option>
                      <option value="medium">Medium</option>
                      <option value="hard">Hard</option>
                      <option value="extreme">Extreme</option>
                    </select>
                  </div>
                </div>
                
                <div class="col-6">
                  <div class="mb-3">
                    <label class="form-label">Price ($)</label>
                    <input 
                      v-model.number="tour.price" 
                      type="number" 
                      class="form-control"
                      min="0"
                      step="0.01"
                    />
                  </div>
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Tags</label>
                <input 
                  v-model="tour.tags" 
                  type="text" 
                  class="form-control"
                  required
                  placeholder="e.g., nature, adventure, culture"
                />
              </div>
              
              <!-- Key Points List -->
              <div class="mb-3">
                <h6>Key Points ({{ tour.key_points.length }})</h6>
                <div class="key-points-list" style="max-height: 200px; overflow-y: auto;">
                  <div 
                    v-for="(point, index) in tour.key_points" 
                    :key="point.id || index"
                    class="card mb-2"
                  >
                    <div class="card-body p-2">
                      <div class="d-flex justify-content-between align-items-start">
                        <div class="flex-grow-1">
                          <small class="text-muted">Point {{ index + 1 }}</small>
                          <div class="fw-bold">{{ point.name }}</div>
                          <small class="text-muted">
                            {{ point.latitude.toFixed(4) }}, {{ point.longitude.toFixed(4) }}
                          </small>
                        </div>
                        <button 
                          type="button"
                          class="btn btn-sm btn-outline-danger"
                          @click="removeKeyPoint(index)"
                        >
                          ×
                        </button>
                      </div>
                    </div>
                  </div>
                  
                  <div v-if="tour.key_points.length === 0" class="text-muted text-center py-3">
                    <small>Click on the map to add key points</small>
                  </div>
                </div>
              </div>
              
              <!-- Distance info -->
              <div class="mb-3" v-if="tour.distance > 0">
                <small class="text-muted">
                  Total Distance: {{ tour.distance.toFixed(2) }} km
                </small>
              </div>
              
              <!-- Action buttons -->
              <div class="d-grid gap-2">
                <button type="submit" class="btn btn-primary">
                  {{ isEditing ? 'Update Tour' : 'Create Tour' }}
                </button>
                <button type="button" class="btn btn-outline-secondary" @click="goBack">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
        
        <!-- Map area -->
        <div class="col-md-8 col-lg-9 p-0">
          <TourMap
            :key-points="tour.key_points"
            :editable="true"
            map-height="calc(100vh - 56px)"
            @keypoint-added="addKeyPoint"
            @keypoint-updated="updateKeyPoint"
            @keypoint-removed="removeKeyPoint"
            @keypoint-edit="openEditDialog"
            @route-cleared="clearRoute"
          />
        </div>
      </div>
    </div>
    
    <!-- Edit Key Point Modal -->
    <div 
      class="modal fade" 
      id="editKeyPointModal" 
      tabindex="-1"
      ref="editModal"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Edit Key Point</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveKeyPointEdit">
              <div class="mb-3">
                <label class="form-label">Name</label>
                <input 
                  v-model="editingKeyPoint.name" 
                  type="text" 
                  class="form-control"
                  required
                />
              </div>
              
              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea 
                  v-model="editingKeyPoint.description" 
                  class="form-control"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="row">
                <div class="col-6">
                  <div class="mb-3">
                    <label class="form-label">Latitude</label>
                    <input 
                      v-model.number="editingKeyPoint.latitude" 
                      type="number" 
                      class="form-control"
                      step="0.000001"
                      required
                    />
                  </div>
                </div>
                
                <div class="col-6">
                  <div class="mb-3">
                    <label class="form-label">Longitude</label>
                    <input 
                      v-model.number="editingKeyPoint.longitude" 
                      type="number" 
                      class="form-control"
                      step="0.000001"
                      required
                    />
                  </div>
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Order</label>
                <input 
                  v-model.number="editingKeyPoint.order" 
                  type="number" 
                  class="form-control"
                  min="0"
                />
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
              Cancel
            </button>
            <button type="button" class="btn btn-primary" @click="saveKeyPointEdit">
              Save Changes
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TourMap from '../components/Map/TourMap.vue'
import { useTourStore } from '../stores/tour'
import { Modal } from 'bootstrap'

export default {
  name: 'TourEditor',
  components: {
    TourMap
  },
  props: {
    id: String
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    const tourStore = useTourStore()
    
    const editModal = ref(null)
    const editModalInstance = ref(null)
    
    const tour = ref({
      name: '',
      description: '',
      difficulty: '',
      tags: '',
      price: 0,
      key_points: [],
      distance: 0
    })
    
    const editingKeyPoint = ref({})
    const editingIndex = ref(-1)
    
    const isEditing = computed(() => !!props.id)
    
    onMounted(async () => {
      // Initialize Bootstrap modal
      if (editModal.value) {
        editModalInstance.value = new Modal(editModal.value)
      }
      
      // Load tour data if editing
      if (isEditing.value) {
        try {
          const tourData = await tourStore.getTour(props.id)
          tour.value = { ...tourData }
        } catch (error) {
          console.error('Failed to load tour:', error)
          // Handle error (show notification, redirect, etc.)
        }
      }
    })
    
    const addKeyPoint = (keyPoint) => {
      tour.value.key_points.push(keyPoint)
      updateDistance()
    }
    
    const updateKeyPoint = (updatedPoint) => {
      const index = tour.value.key_points.findIndex(p => p.id === updatedPoint.id)
      if (index !== -1) {
        tour.value.key_points[index] = updatedPoint
        updateDistance()
      }
    }
    
    const removeKeyPoint = (index) => {
      tour.value.key_points.splice(index, 1)
      updateDistance()
    }
    
    const openEditDialog = ({ point, index }) => {
      editingKeyPoint.value = { ...point }
      editingIndex.value = index
      editModalInstance.value?.show()
    }
    
    const saveKeyPointEdit = () => {
      if (editingIndex.value !== -1) {
        tour.value.key_points[editingIndex.value] = { ...editingKeyPoint.value }
        updateDistance()
      }
      editModalInstance.value?.hide()
    }
    
    const clearRoute = () => {
      tour.value.key_points = []
      tour.value.distance = 0
    }
    
    const updateDistance = () => {
      // Simple distance calculation (you can implement more sophisticated calculation)
      if (tour.value.key_points.length < 2) {
        tour.value.distance = 0
        return
      }
      
      let totalDistance = 0
      for (let i = 0; i < tour.value.key_points.length - 1; i++) {
        const point1 = tour.value.key_points[i]
        const point2 = tour.value.key_points[i + 1]
        totalDistance += calculateDistance(
          point1.latitude, point1.longitude,
          point2.latitude, point2.longitude
        )
      }
      
      tour.value.distance = totalDistance
    }
    
    const calculateDistance = (lat1, lon1, lat2, lon2) => {
      const R = 6371 // Earth's radius in km
      const dLat = (lat2 - lat1) * Math.PI / 180
      const dLon = (lon2 - lon1) * Math.PI / 180
      const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
                Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
                Math.sin(dLon/2) * Math.sin(dLon/2)
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
      return R * c
    }
    
    const saveTour = async () => {
      try {
        if (isEditing.value) {
          // For editing, send all data including key points as before
          const tourData = {
            ...tour.value,
            key_points: tour.value.key_points.map((point, index) => ({
              ...point,
              order: point.order || index
            }))
          }
          await tourStore.updateTour(props.id, tourData)
        } else {
          // For creation, create tour first without key points
          const tourData = {
            name: tour.value.name,
            description: tour.value.description,
            difficulty: tour.value.difficulty,
            tags: tour.value.tags,
            price: tour.value.price || 0
          }
          
          // Create the tour
          const createdTour = await tourStore.createTour(tourData)
          
          // Then add key points one by one
          if (tour.value.key_points.length > 0) {
            for (let i = 0; i < tour.value.key_points.length; i++) {
              const keyPoint = tour.value.key_points[i]
              const keyPointData = {
                name: keyPoint.name,
                description: keyPoint.description || '',
                latitude: keyPoint.latitude,
                longitude: keyPoint.longitude,
                image_url: keyPoint.image_url || '',
                order: keyPoint.order || i
              }
              
              await tourStore.addKeyPoint(createdTour.id, keyPointData)
            }
          }
        }
        
        // Redirect to tours list
        router.push('/tours')
      } catch (error) {
        console.error('Failed to save tour:', error)
        // Handle error (show notification)
      }
    }
    
    const goBack = () => {
      router.back()
    }
    
    return {
      tour,
      editingKeyPoint,
      editingIndex,
      editModal,
      isEditing,
      addKeyPoint,
      updateKeyPoint,
      removeKeyPoint,
      openEditDialog,
      saveKeyPointEdit,
      clearRoute,
      saveTour,
      goBack
    }
  }
}
</script>

<style scoped>
.tour-editor {
  height: calc(100vh - 56px);
}

.key-points-list {
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
  padding: 0.5rem;
}

.key-points-list .card {
  border: 1px solid #e0e0e0;
}

.modal-backdrop {
  z-index: 1040;
}

.modal {
  z-index: 1050;
}
</style>
