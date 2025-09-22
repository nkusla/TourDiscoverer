<template>
  <div class="tour-execution">
    <div v-if="!activeTourExecution" class="tour-selection">
      <div class="section-header">
        <h2>Select a Tour to Execute</h2>
        <p class="subtitle">Choose from available published and archived tours</p>
      </div>

      <div v-if="loading" class="loading">
        <p>Loading tours...</p>
      </div>

      <div v-else-if="executableTours.length === 0" class="no-tours">
        <p>No tours available for execution.</p>
      </div>

      <div v-else class="tours-grid">
        <div 
          v-for="tour in executableTours" 
          :key="tour.id" 
          class="tour-card"
          @click="selectTour(tour)"
        >
          <div class="tour-info">
            <h3>{{ tour.name }}</h3>
            <p class="description">{{ tour.description }}</p>
            <div class="tour-details">
              <span class="difficulty" :class="tour.difficulty.toLowerCase()">
                {{ tour.difficulty }}
              </span>
              <span class="distance">{{ tour.distance.toFixed(1) }} km</span>
              <span class="key-points">{{ tour.key_points?.length || 0 }} key points</span>
            </div>
            <div class="tags">
              <span v-for="tag in tour.tags.split(',')" :key="tag" class="tag">
                {{ tag.trim() }}
              </span>
            </div>
          </div>
          <button class="start-btn">Start Tour</button>
        </div>
      </div>
    </div>

    <div v-else class="active-tour">
      <div class="tour-header">
        <h2>{{ currentTour?.name || 'Active Tour' }}</h2>
        <div class="tour-status">
          <span class="status-badge active">Active</span>
          <span class="tour-time">Started: {{ formatTime(activeTourExecution.start_time) }}</span>
        </div>
      </div>

      <!-- Current Position Display -->
      <div v-if="currentPosition" class="current-position">
        <h3>Your Current Position</h3>
        <div class="position-info">
          <span><strong>Latitude:</strong> {{ currentPosition.latitude.toFixed(6) }}</span>
          <span><strong>Longitude:</strong> {{ currentPosition.longitude.toFixed(6) }}</span>
          <span><strong>Last Updated:</strong> {{ formatTime(currentPosition.updated_at) }}</span>
        </div>
      </div>

      <!-- Key Point Progress -->
      <div v-if="currentTour" class="key-points-progress">
        <h3>Key Points Progress</h3>
        <div class="key-points-list">
          <div 
            v-for="(keyPoint, index) in currentTour.key_points" 
            :key="keyPoint.id"
            class="key-point-item"
            :class="{ completed: isKeyPointCompleted(keyPoint.id) }"
          >
            <div class="key-point-number">{{ index + 1 }}</div>
            <div class="key-point-info">
              <h4>{{ keyPoint.name }}</h4>
              <p>{{ keyPoint.description }}</p>
              <div class="coordinates">
                Lat: {{ keyPoint.latitude.toFixed(6) }}, Lng: {{ keyPoint.longitude.toFixed(6) }}
              </div>
            </div>
            <div class="key-point-status">
              <span v-if="isKeyPointCompleted(keyPoint.id)" class="completed-badge">✓ Completed</span>
              <span v-else class="pending-badge">Pending</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Status Messages -->
      <div v-if="statusMessage" :class="['status-message', statusType]">
        {{ statusMessage }}
      </div>

      <!-- Action Buttons -->
      <div class="actions">
        <button @click="checkPosition" class="btn btn-primary" :disabled="checking">
          <span v-if="checking">Checking...</span>
          <span v-else>Check Position Now</span>
        </button>
        
        <button @click="completeTour" class="btn btn-success">
          Complete Tour
        </button>
        
        <button @click="abandonTour" class="btn btn-danger">
          Abandon Tour
        </button>
      </div>

      <!-- Auto-check indicator -->
      <div class="auto-check-indicator">
        <p>
          <span class="indicator-dot" :class="{ active: autoCheckActive }"></span>
          Auto-checking position every 10 seconds
        </p>
        <p><small>Next check in: {{ nextCheckIn }}s</small></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { tourExecutionService } from '../services/tourExecution'

export default {
  name: 'TourExecution',
  setup() {
    const executableTours = ref([])
    const activeTourExecution = ref(null)
    const currentTour = ref(null)
    const currentPosition = ref(null)
    const loading = ref(false)
    const checking = ref(false)
    const statusMessage = ref('')
    const statusType = ref('info')
    const autoCheckActive = ref(false)
    const nextCheckIn = ref(10)
    
    let autoCheckInterval = null
    let countdownInterval = null

    const completedKeyPoints = computed(() => {
      return activeTourExecution.value?.key_point_completions || []
    })

    const isKeyPointCompleted = (keyPointId) => {
      return completedKeyPoints.value.some(completion => completion.key_point_id === keyPointId)
    }

    const showStatus = (message, type = 'info') => {
      statusMessage.value = message
      statusType.value = type
      setTimeout(() => {
        statusMessage.value = ''
      }, 5000)
    }

    const formatTime = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    const loadExecutableTours = async () => {
      loading.value = true
      try {
        const response = await tourExecutionService.getExecutableTours()
        executableTours.value = response.data.tours || []
      } catch (error) {
        console.error('Error loading executable tours:', error)
        showStatus('Failed to load tours', 'error')
      } finally {
        loading.value = false
      }
    }

    const checkActiveTourExecution = async () => {
      try {
        const response = await tourExecutionService.getActiveTourExecution()
        if (response.data) {
          activeTourExecution.value = response.data
          await loadCurrentTour()
          startAutoCheck()
        }
      } catch (error) {
        // No active tour execution - this is normal
        console.log('No active tour execution')
      }
    }

    const loadCurrentTour = async () => {
      if (!activeTourExecution.value) return
      
      try {
        const response = await tourExecutionService.getTourById(activeTourExecution.value.tour_id)
        currentTour.value = response.data
      } catch (error) {
        console.error('Error loading current tour:', error)
        showStatus('Failed to load tour details', 'error')
      }
    }

    const selectTour = async (tour) => {
      try {
        // First get current position from Position Simulator
        const positionResponse = await tourExecutionService.getTouristPosition()
        if (!positionResponse.data) {
          showStatus('Please set your position first using the Position Simulator', 'error')
          return
        }

        const position = positionResponse.data
        
        // Start tour execution
        loading.value = true
        const response = await tourExecutionService.startTourExecution(
          tour.id, 
          position.latitude, 
          position.longitude
        )
        
        activeTourExecution.value = response.data
        currentTour.value = tour
        currentPosition.value = position
        
        showStatus(`Tour "${tour.name}" started successfully!`, 'success')
        startAutoCheck()
        
      } catch (error) {
        console.error('Error starting tour:', error)
        showStatus('Failed to start tour: ' + (error.response?.data?.message || 'Unknown error'), 'error')
      } finally {
        loading.value = false
      }
    }

    const checkPosition = async () => {
      if (!activeTourExecution.value) return

      checking.value = true
      try {
        // Get current position
        const positionResponse = await tourExecutionService.getTouristPosition()
        if (!positionResponse.data) {
          showStatus('Position not found. Please update your position.', 'error')
          return
        }

        currentPosition.value = positionResponse.data
        
        // Check proximity to key points
        const proximityResponse = await tourExecutionService.checkProximity(
          activeTourExecution.value.id,
          currentPosition.value.latitude,
          currentPosition.value.longitude
        )

        if (proximityResponse.data.key_point_reached) {
          showStatus(proximityResponse.data.message, 'success')
          // Reload active execution to get updated key point completions
          await checkActiveTourExecution()
        } else {
          showStatus('No key points nearby', 'info')
        }
        
      } catch (error) {
        console.error('Error checking position:', error)
        showStatus('Failed to check position: ' + (error.response?.data?.message || 'Unknown error'), 'error')
      } finally {
        checking.value = false
      }
    }

    const completeTour = async () => {
      if (!confirm('Are you sure you want to complete this tour?')) return

      try {
        await tourExecutionService.endTourExecution(activeTourExecution.value.id, 'completed')
        showStatus('Tour completed successfully!', 'success')
        activeTourExecution.value = null
        currentTour.value = null
        currentPosition.value = null
        stopAutoCheck()
        await loadExecutableTours() // Reload tours
      } catch (error) {
        console.error('Error completing tour:', error)
        showStatus('Failed to complete tour: ' + (error.response?.data?.message || 'Unknown error'), 'error')
      }
    }

    const abandonTour = async () => {
      if (!confirm('Are you sure you want to abandon this tour? Your progress will be lost.')) return

      try {
        await tourExecutionService.endTourExecution(activeTourExecution.value.id, 'abandoned')
        showStatus('Tour abandoned', 'info')
        activeTourExecution.value = null
        currentTour.value = null
        currentPosition.value = null
        stopAutoCheck()
        await loadExecutableTours() // Reload tours
      } catch (error) {
        console.error('Error abandoning tour:', error)
        showStatus('Failed to abandon tour: ' + (error.response?.data?.message || 'Unknown error'), 'error')
      }
    }

    const startAutoCheck = () => {
      stopAutoCheck() // Clear any existing intervals
      autoCheckActive.value = true
      nextCheckIn.value = 10
      
      // Check position every 10 seconds
      autoCheckInterval = setInterval(() => {
        checkPosition()
        nextCheckIn.value = 10 // Reset countdown
      }, 10000)
      
      // Update countdown every second
      countdownInterval = setInterval(() => {
        nextCheckIn.value = nextCheckIn.value - 1
        if (nextCheckIn.value <= 0) {
          nextCheckIn.value = 10
        }
      }, 1000)
    }

    const stopAutoCheck = () => {
      if (autoCheckInterval) {
        clearInterval(autoCheckInterval)
        autoCheckInterval = null
      }
      if (countdownInterval) {
        clearInterval(countdownInterval)
        countdownInterval = null
      }
      autoCheckActive.value = false
    }

    onMounted(async () => {
      await checkActiveTourExecution()
      if (!activeTourExecution.value) {
        await loadExecutableTours()
      }
    })

    onUnmounted(() => {
      stopAutoCheck()
    })

    return {
      executableTours,
      activeTourExecution,
      currentTour,
      currentPosition,
      loading,
      checking,
      statusMessage,
      statusType,
      autoCheckActive,
      nextCheckIn,
      completedKeyPoints,
      isKeyPointCompleted,
      formatTime,
      selectTour,
      checkPosition,
      completeTour,
      abandonTour
    }
  }
}
</script>

<style scoped>
.tour-execution {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.section-header {
  text-align: center;
  margin-bottom: 30px;
}

.section-header h2 {
  color: #333;
  margin-bottom: 10px;
}

.subtitle {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.loading, .no-tours {
  text-align: center;
  padding: 40px;
  color: #666;
}

.tours-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.tour-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.tour-card:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0,123,255,0.1);
}

.tour-info h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.description {
  color: #666;
  margin-bottom: 15px;
  line-height: 1.4;
}

.tour-details {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.difficulty, .distance, .key-points {
  font-size: 14px;
  padding: 4px 8px;
  border-radius: 4px;
  background: #f8f9fa;
}

.difficulty.easy { background: #d4edda; color: #155724; }
.difficulty.medium { background: #fff3cd; color: #856404; }
.difficulty.hard { background: #f8d7da; color: #721c24; }

.tags {
  margin-bottom: 15px;
}

.tag {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  margin-right: 5px;
}

.start-btn {
  background: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  margin-top: 10px;
}

.start-btn:hover {
  background: #0056b3;
}

.active-tour {
  background: white;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.tour-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.tour-status {
  text-align: right;
}

.status-badge {
  background: #28a745;
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 14px;
  display: block;
  margin-bottom: 5px;
}

.tour-time {
  font-size: 14px;
  color: #666;
}

.current-position {
  background: #f8f9fa;
  border-radius: 6px;
  padding: 20px;
  margin-bottom: 30px;
}

.current-position h3 {
  margin: 0 0 15px 0;
  color: #495057;
}

.position-info {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.position-info span {
  background: white;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #dee2e6;
  font-size: 14px;
}

.key-points-progress {
  margin-bottom: 30px;
}

.key-points-progress h3 {
  margin-bottom: 20px;
  color: #333;
}

.key-points-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.key-point-item {
  display: flex;
  align-items: center;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  border-left: 4px solid #dee2e6;
}

.key-point-item.completed {
  background: #d4edda;
  border-left-color: #28a745;
}

.key-point-number {
  background: #007bff;
  color: white;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 15px;
  flex-shrink: 0;
}

.key-point-item.completed .key-point-number {
  background: #28a745;
}

.key-point-info {
  flex-grow: 1;
}

.key-point-info h4 {
  margin: 0 0 5px 0;
  color: #333;
}

.key-point-info p {
  margin: 0 0 5px 0;
  color: #666;
  font-size: 14px;
}

.coordinates {
  font-size: 12px;
  color: #999;
}

.key-point-status {
  margin-left: 15px;
}

.completed-badge {
  background: #28a745;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.pending-badge {
  background: #6c757d;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-message {
  padding: 12px 20px;
  border-radius: 6px;
  margin-bottom: 20px;
  text-align: center;
  font-weight: 500;
}

.status-message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status-message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.status-message.info {
  background: #d1ecf1;
  color: #0c5460;
  border: 1px solid #bee5eb;
}

.actions {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #0056b3;
}

.btn-success {
  background: #28a745;
  color: white;
}

.btn-success:hover {
  background: #218838;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover {
  background: #c82333;
}

.auto-check-indicator {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #dee2e6;
}

.auto-check-indicator p {
  margin: 0 0 5px 0;
  color: #666;
}

.indicator-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #6c757d;
  margin-right: 8px;
}

.indicator-dot.active {
  background: #28a745;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

@media (max-width: 768px) {
  .tour-header {
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }
  
  .position-info {
    flex-direction: column;
    gap: 10px;
  }
  
  .actions {
    flex-direction: column;
    align-items: center;
  }
  
  .btn {
    width: 200px;
  }
}
</style>