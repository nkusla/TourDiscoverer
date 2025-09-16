<template>
  <div class="create-tour">
    <div class="container-fluid">
      <div class="row">
        <!-- Sidebar for tour details -->
        <div class="col-md-4 col-lg-3 bg-light border-end" style="height: calc(100vh - 56px); overflow-y: auto;">
          <div class="p-3">
            <h4>Create New Tour</h4>
            
            <form @submit.prevent="createTour">
              <!-- Basic Tour Info -->
              <div class="mb-3">
                <label class="form-label">Tour Name *</label>
                <input 
                  v-model="tourData.name" 
                  type="text" 
                  class="form-control"
                  required
                  placeholder="Enter tour name"
                />
              </div>
              
              <div class="mb-3">
                <label class="form-label">Description *</label>
                <textarea 
                  v-model="tourData.description" 
                  class="form-control"
                  rows="3"
                  required
                  placeholder="Describe your tour"
                ></textarea>
              </div>
              
              <div class="row">
                <div class="col-6">
                  <div class="mb-3">
                    <label class="form-label">Difficulty *</label>
                    <select v-model="tourData.difficulty" class="form-select" required>
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
                      v-model.number="tourData.price" 
                      type="number" 
                      class="form-control"
                      min="0"
                      step="0.01"
                      :value="0"
                      readonly
                    />
                    <small class="text-muted">Auto-set to 0 for draft</small>
                  </div>
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Tags *</label>
                <input 
                  v-model="tourData.tags" 
                  type="text" 
                  class="form-control"
                  required
                  placeholder="e.g., nature, adventure, culture"
                />
              </div>
              
              <!-- Key Points List -->
              <div class="mb-3">
                <h6>Key Points ({{ keyPoints.length }})</h6>
                <div class="alert alert-info" v-if="keyPoints.length === 0">
                  <small>Click on the map to add key points for your tour</small>
                </div>
                
                <div class="key-points-list" style="max-height: 250px; overflow-y: auto;">
                  <div 
                    v-for="(point, index) in keyPoints" 
                    :key="point.id"
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
                      
                      <!-- Edit point details -->
                      <div class="mt-2">
                        <input 
                          v-model="point.name"
                          type="text"
                          class="form-control form-control-sm mb-1"
                          placeholder="Point name"
                        />
                        <textarea 
                          v-model="point.description"
                          class="form-control form-control-sm"
                          rows="2"
                          placeholder="Point description"
                        ></textarea>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Distance info -->
              <div class="mb-3" v-if="totalDistance > 0">
                <small class="text-muted">
                  Total Distance: {{ totalDistance.toFixed(2) }} km
                </small>
              </div>
              
              <!-- Action buttons -->
              <div class="d-grid gap-2">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="!canCreateTour"
                >
                  Create Tour (Draft)
                </button>
                <button 
                  type="button" 
                  class="btn btn-outline-secondary" 
                  @click="clearAll"
                >
                  Clear All
                </button>
                <router-link to="/tours" class="btn btn-outline-secondary">
                  Cancel
                </router-link>
              </div>
            </form>
            
            <!-- Creation status -->
            <div v-if="isCreating" class="alert alert-info mt-3">
              Creating tour...
            </div>
            
            <div v-if="createError" class="alert alert-danger mt-3">
              Error: {{ createError }}
            </div>
            
            <div v-if="createSuccess" class="alert alert-success mt-3">
              Tour created successfully!
            </div>
          </div>
        </div>
        
        <!-- Map area -->
        <div class="col-md-8 col-lg-9 p-0">
          <div class="map-header bg-primary text-white p-2">
            <small>Click on the map to add key points for your tour</small>
          </div>
          <LeafletMap
            :key-points="keyPoints"
            :editable="true"
            map-height="calc(100vh - 56px - 40px)"
            @keypoint-added="addKeyPoint"
            @keypoint-removed="removeKeyPoint"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import LeafletMap from '../components/Map/LeafletMap.vue'
import { useTourStore } from '../stores/tour'
import { useUserStore } from '../stores/user'

export default {
  name: 'CreateTour',
  components: {
    LeafletMap
  },
  setup() {
    const router = useRouter()
    const tourStore = useTourStore()
    const userStore = useUserStore()
    
    const tourData = ref({
      name: '',
      description: '',
      difficulty: '',
      price: 0,
      tags: ''
    })
    
    const keyPoints = ref([])
    const isCreating = ref(false)
    const createError = ref('')
    const createSuccess = ref(false)
    
    const canCreateTour = computed(() => {
      return tourData.value.name && 
             tourData.value.description && 
             tourData.value.difficulty && 
             tourData.value.tags &&
             keyPoints.value.length >= 2
    })
    
    const totalDistance = computed(() => {
      if (keyPoints.value.length < 2) return 0
      
      let distance = 0
      for (let i = 1; i < keyPoints.value.length; i++) {
        const prev = keyPoints.value[i - 1]
        const curr = keyPoints.value[i]
        distance += calculateDistance(prev.latitude, prev.longitude, curr.latitude, curr.longitude)
      }
      return distance
    })
    
    const calculateDistance = (lat1, lon1, lat2, lon2) => {
      const R = 6371 // Radius of the Earth in kilometers
      const dLat = (lat2 - lat1) * Math.PI / 180
      const dLon = (lon2 - lon1) * Math.PI / 180
      const a = 
        Math.sin(dLat/2) * Math.sin(dLat/2) +
        Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * 
        Math.sin(dLon/2) * Math.sin(dLon/2)
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
      const d = R * c // Distance in kilometers
      return d
    }
    
    const addKeyPoint = (newPoint) => {
      keyPoints.value.push(newPoint)
    }
    
    const removeKeyPoint = (index) => {
      keyPoints.value.splice(index, 1)
      // Update order for remaining points
      keyPoints.value.forEach((point, idx) => {
        point.order = idx
      })
    }
    
    const clearAll = () => {
      tourData.value = {
        name: '',
        description: '',
        difficulty: '',
        price: 0,
        tags: ''
      }
      keyPoints.value = []
      createError.value = ''
      createSuccess.value = false
    }
    
    const createTour = async () => {
      if (!canCreateTour.value) return
      
      isCreating.value = true
      createError.value = ''
      createSuccess.value = false
      
      try {
        const tourPayload = {
          ...tourData.value,
          key_points: keyPoints.value.map((point, index) => ({
            ...point,
            order: index
          })),
          distance: totalDistance.value,
          status: 'draft',
          price: 0 // Always 0 for draft tours
        }
        
        console.log('Creating tour with payload:', tourPayload)
        
        const result = await tourStore.createTour(tourPayload)
        
        if (result) {
          createSuccess.value = true
          setTimeout(() => {
            router.push('/tours')
          }, 2000)
        } else {
          createError.value = 'Failed to create tour'
        }
      } catch (error) {
        console.error('Error creating tour:', error)
        createError.value = error.message || 'Failed to create tour'
      } finally {
        isCreating.value = false
      }
    }
    
    return {
      tourData,
      keyPoints,
      isCreating,
      createError,
      createSuccess,
      canCreateTour,
      totalDistance,
      addKeyPoint,
      removeKeyPoint,
      clearAll,
      createTour
    }
  }
}
</script>

<style scoped>
.create-tour {
  height: calc(100vh - 56px);
}

.map-header {
  font-size: 14px;
  text-align: center;
}

.key-points-list {
  border: 1px solid #dee2e6;
  border-radius: 4px;
  padding: 8px;
}
</style>