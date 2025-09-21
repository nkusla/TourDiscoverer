<template>
  <div class="tour-editor">
    <div class="container-fluid">
      <div class="row">
        <!-- Sidebar for tour details -->
        <div class="col-md-4 col-lg-3 bg-light border-end" style="height: calc(100vh - 56px); overflow-y: auto;">
          <div class="p-3">
            <h4>{{ isEditMode ? 'Edit Tour' : 'Create New Tour' }}</h4>

            <form @submit.prevent="saveTour">
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
                      :disabled="!isEditMode"
                    />
                    <small class="text-muted" v-if="!isEditMode">Auto-set to 0 for draft</small>
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
                    :key="point.id || index"
                    class="card mb-2"
                    :class="{ 'border-primary': selectedKeyPointIndex === index }"
                  >
                    <div class="card-body p-2">
                      <div class="d-flex justify-content-between align-items-start">
                        <div class="flex-grow-1">
                          <small class="text-muted">Point {{ index + 1 }}</small>
                          <div class="fw-bold">{{ point.name }}</div>
                          <small class="text-muted">
                            {{ point.latitude.toFixed(4) }}, {{ point.longitude.toFixed(4) }}
                          </small>
                          <div v-if="point.description" class="mt-1">
                            <small class="text-secondary">{{ point.description }}</small>
                          </div>
                        </div>
                        <div class="btn-group-vertical">
                          <button
                            type="button"
                            class="btn btn-sm btn-outline-primary"
                            @click="editKeyPoint(index)"
                            title="Edit key point"
                          >
                            ✏️
                          </button>
                          <button
                            type="button"
                            class="btn btn-sm btn-outline-danger"
                            @click="removeKeyPoint(index)"
                            title="Delete key point"
                          >
                            🗑️
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Distance and Status info -->
              <div class="mb-3" v-if="totalDistance > 0">
                <small class="text-muted">
                  <strong>Total Distance:</strong> {{ totalDistance.toFixed(2) }} km
                </small>
              </div>

              <div class="mb-3" v-if="isEditMode">
                <small class="text-muted">
                  <strong>Status:</strong> {{ tourData.status }}
                </small>
              </div>

              <!-- Action buttons -->
              <div class="d-grid gap-2">
                <button
                  type="submit"
                  class="btn btn-primary"
                  :disabled="!canSave || isLoading"
                >
                  <span v-if="isLoading" class="spinner-border spinner-border-sm me-2"></span>
                  {{ isEditMode ? 'Update Tour' : 'Create Tour' }}
                </button>

                <button
                  type="button"
                  class="btn btn-outline-secondary"
                  @click="goBack"
                >
                  Cancel
                </button>

                <button
                  type="button"
                  class="btn btn-outline-warning"
                  @click="clearAll"
                  v-if="!isEditMode"
                >
                  Clear All
                </button>
              </div>

              <!-- Success/Error Messages -->
              <div v-if="successMessage" class="alert alert-success mt-3">
                {{ successMessage }}
              </div>

              <div v-if="errorMessage" class="alert alert-danger mt-3">
                {{ errorMessage }}
              </div>
            </form>
          </div>
        </div>

        <!-- Map area -->
        <div class="col-md-8 col-lg-9 p-0">
          <LeafletMap
            :key-points="keyPoints"
            :editable="true"
            :selected-point-index="selectedKeyPointIndex"
            map-height="calc(100vh - 56px)"
            @keypoint-added="addKeyPoint"
            @keypoint-moved="moveKeyPoint"
            @keypoint-clicked="selectKeyPoint"
            @map-clicked="onMapClick"
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
            <form>
              <div class="mb-3">
                <label class="form-label">Name *</label>
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
                      step="any"
                      readonly
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
                      step="any"
                      readonly
                    />
                  </div>
                </div>
              </div>

              <div class="alert alert-info">
                <small>💡 To change the location, click on the map while this dialog is open.</small>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LeafletMap from '../components/Map/LeafletMap.vue'
import { useTourStore } from '../stores/tour'
import { Modal } from 'bootstrap'

export default {
  name: 'TourEditor',
  components: {
    LeafletMap
  },
  props: {
    id: String // Tour ID for editing
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    const tourStore = useTourStore()

    // Reactive data
    const tourData = ref({
      name: '',
      description: '',
      difficulty: '',
      price: 0,
      tags: ''
    })

    const keyPoints = ref([])
    const selectedKeyPointIndex = ref(-1)
    const editingKeyPoint = ref({})
    const editingKeyPointIndex = ref(-1)
    const isEditingPosition = ref(false)

    const isLoading = ref(false)
    const successMessage = ref('')
    const errorMessage = ref('')

    // Modal reference
    const editModal = ref(null)
    let editModalInstance = null

    // Computed
    const isEditMode = computed(() => !!props.id)

    const totalDistance = computed(() => {
      if (keyPoints.value.length < 2) return 0

      let distance = 0
      for (let i = 0; i < keyPoints.value.length - 1; i++) {
        const p1 = keyPoints.value[i]
        const p2 = keyPoints.value[i + 1]
        distance += calculateDistance(p1.latitude, p1.longitude, p2.latitude, p2.longitude)
      }
      return distance
    })

    const canSave = computed(() => {
      return tourData.value.name &&
             tourData.value.description &&
             tourData.value.difficulty &&
             tourData.value.tags
    })

    // Methods
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

    const loadTour = async () => {
      if (!props.id) return

      try {
        isLoading.value = true
        const tour = await tourStore.getTour(props.id)

        tourData.value = {
          name: tour.name,
          description: tour.description,
          difficulty: tour.difficulty,
          price: tour.price,
          tags: tour.tags
        }

        keyPoints.value = tour.key_points || []
      } catch (error) {
        errorMessage.value = 'Failed to load tour: ' + error.message
      } finally {
        isLoading.value = false
      }
    }

    const addKeyPoint = (newPoint) => {
      const keyPoint = {
        id: Date.now(),
        name: newPoint.name || `Point ${keyPoints.value.length + 1}`,
        description: '',
        latitude: newPoint.latitude,
        longitude: newPoint.longitude,
        order: keyPoints.value.length
      }
      keyPoints.value.push(keyPoint)
      selectedKeyPointIndex.value = keyPoints.value.length - 1
    }

    const removeKeyPoint = (index) => {
      if (confirm('Are you sure you want to delete this key point?')) {
        keyPoints.value.splice(index, 1)
        // Update order for remaining points
        keyPoints.value.forEach((point, idx) => {
          point.order = idx
        })
        selectedKeyPointIndex.value = -1
      }
    }

    const editKeyPoint = (index) => {
      editingKeyPointIndex.value = index
      editingKeyPoint.value = { ...keyPoints.value[index] }
      isEditingPosition.value = true
      selectedKeyPointIndex.value = index

      if (editModalInstance) {
        editModalInstance.show()
      }
    }

    const selectKeyPoint = (index) => {
      selectedKeyPointIndex.value = index
    }

    const moveKeyPoint = (index, newPosition) => {
      if (index >= 0 && index < keyPoints.value.length) {
        keyPoints.value[index].latitude = newPosition.latitude
        keyPoints.value[index].longitude = newPosition.longitude

        // If we're editing a key point, update the editing data
        if (isEditingPosition.value && editingKeyPointIndex.value === index) {
          editingKeyPoint.value.latitude = newPosition.latitude
          editingKeyPoint.value.longitude = newPosition.longitude
        }
      }
    }

    const onMapClick = (position) => {
      if (isEditingPosition.value && editingKeyPointIndex.value >= 0) {
        // Update the position of the key point being edited
        moveKeyPoint(editingKeyPointIndex.value, position)
      } else {
        // Add new key point
        addKeyPoint(position)
      }
    }

    const saveKeyPointEdit = () => {
      if (editingKeyPointIndex.value >= 0) {
        keyPoints.value[editingKeyPointIndex.value] = { ...editingKeyPoint.value }
        editingKeyPointIndex.value = -1
        isEditingPosition.value = false

        if (editModalInstance) {
          editModalInstance.hide()
        }
      }
    }

    const saveTour = async () => {
      if (!canSave.value) return

      isLoading.value = true
      errorMessage.value = ''
      successMessage.value = ''

      try {
        const tourPayload = {
          ...tourData.value,
          key_points: keyPoints.value.map((point, index) => ({
            ...point,
            order: index
          })),
          distance: totalDistance.value,
          status: 'draft',
          price: isEditMode.value ? tourData.value.price : 0
        }

        let result
        if (isEditMode.value) {
          result = await tourStore.updateTour(props.id, tourPayload)
          successMessage.value = 'Tour updated successfully!'
        } else {
          result = await tourStore.createTour(tourPayload)
          successMessage.value = 'Tour created successfully!'
        }

        // Redirect to tours list after a delay
        setTimeout(() => {
          router.push('/tours')
        }, 2000)

      } catch (error) {
        errorMessage.value = 'Error saving tour: ' + error.message
      } finally {
        isLoading.value = false
      }
    }

    const clearAll = () => {
      if (confirm('Are you sure you want to clear all data?')) {
        tourData.value = {
          name: '',
          description: '',
          difficulty: '',
          price: 0,
          tags: ''
        }
        keyPoints.value = []
        selectedKeyPointIndex.value = -1
        errorMessage.value = ''
        successMessage.value = ''
      }
    }

    const goBack = () => {
      router.push('/tours')
    }

    // Lifecycle
    onMounted(async () => {
      // Initialize Bootstrap modal
      if (editModal.value) {
        editModalInstance = new Modal(editModal.value)
      }

      // Load tour data if editing
      if (isEditMode.value) {
        await loadTour()
      }
    })

    return {
      tourData,
      keyPoints,
      selectedKeyPointIndex,
      editingKeyPoint,
      isLoading,
      successMessage,
      errorMessage,
      editModal,
      isEditMode,
      totalDistance,
      canSave,
      addKeyPoint,
      removeKeyPoint,
      editKeyPoint,
      selectKeyPoint,
      moveKeyPoint,
      onMapClick,
      saveKeyPointEdit,
      saveTour,
      clearAll,
      goBack
    }
  }
}
</script>

<style scoped>
.tour-editor {
  height: 100vh;
}

.key-points-list .card {
  transition: all 0.2s ease;
}

.key-points-list .card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.btn-group-vertical .btn {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  line-height: 1;
}

.alert {
  border-radius: 0.375rem;
}

.modal-content {
  border-radius: 0.5rem;
}
</style>