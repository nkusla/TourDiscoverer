<template>
  <div class="position-simulator">
    <div class="simulator-header">
      <h2>Position Simulator</h2>
      <p class="subtitle">Click on the map to set your current position</p>
    </div>

    <!-- Current Position Display -->
    <div v-if="currentPosition" class="current-position">
      <h3>Current Position</h3>
      <div class="position-info">
        <span><strong>Latitude:</strong> {{ currentPosition.latitude.toFixed(6) }}</span>
        <span><strong>Longitude:</strong> {{ currentPosition.longitude.toFixed(6) }}</span>
        <span><strong>Last Updated:</strong> {{ formatDate(currentPosition.updated_at) }}</span>
      </div>
    </div>

    <div v-else class="no-position">
      <p>No position set. Click on the map to set your current location.</p>
    </div>

    <!-- Map Container -->
    <div class="map-container">
      <LeafletMap
        :key-points="positionMarkers"
        :editable="true"
        :map-height="'500px'"
        :initial-center="mapCenter"
        :initial-zoom="13"
        @map-clicked="onMapClick"
        ref="mapComponent"
      />
    </div>

    <!-- Action Buttons -->
    <div class="actions">
      <button 
        @click="() => refreshPosition(true)" 
        class="btn btn-secondary"
        :disabled="isLoading"
      >
        <span v-if="isLoading">Loading...</span>
        <span v-else>Refresh Position</span>
      </button>
      
      <button 
        @click="clearPosition" 
        class="btn btn-danger"
        v-if="currentPosition"
        :disabled="isLoading"
      >
        Clear Position
      </button>
    </div>

    <!-- Status Messages -->
    <div v-if="statusMessage" :class="['status-message', statusType]">
      {{ statusMessage }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import LeafletMap from './Map/LeafletMap.vue'
import api from '../services/api'

export default {
  name: 'PositionSimulator',
  components: {
    LeafletMap
  },
  setup() {
    const currentPosition = ref(null)
    const mapCenter = ref([44.7866, 20.4489]) // Belgrade, Serbia
    const mapComponent = ref(null)
    const isLoading = ref(false)
    const statusMessage = ref('')
    const statusType = ref('success') // 'success', 'error', 'info'

    // Create marker for current position
    const positionMarkers = computed(() => {
      if (!currentPosition.value) {
        return []
      }
      
      return [{
        latitude: currentPosition.value.latitude,
        longitude: currentPosition.value.longitude,
        name: '📍 Your Position',
        description: `Last updated: ${formatDate(currentPosition.value.updated_at)}<br>Lat: ${currentPosition.value.latitude.toFixed(6)}<br>Lng: ${currentPosition.value.longitude.toFixed(6)}`,
        order: 1
      }]
    })

    const showStatus = (message, type = 'info') => {
      statusMessage.value = message
      statusType.value = type
      setTimeout(() => {
        statusMessage.value = ''
      }, 3000)
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    const updatePosition = async (latitude, longitude) => {
      isLoading.value = true
      try {
        const response = await api.post('/api/stakeholder/position', {
          latitude,
          longitude
        })
        
        currentPosition.value = response.data
        showStatus('Position updated successfully!', 'success')
        
        // Update map center to new position
        mapCenter.value = [latitude, longitude]
        
      } catch (error) {
        console.error('Error updating position:', error)
        showStatus('Failed to update position. Please try again.', 'error')
      } finally {
        isLoading.value = false
      }
    }

    const onMapClick = (event) => {
      const { latitude, longitude } = event
      updatePosition(latitude, longitude)
    }

    const refreshPosition = async (showSuccessMessage = false) => {
      isLoading.value = true
      try {
        const response = await api.get('/api/stakeholder/position')
        
        if (response.data === null) {
          currentPosition.value = null
          if (showSuccessMessage) {
            showStatus('No position found. Click on the map to set your location.', 'info')
          }
          return
        }
        
        currentPosition.value = response.data
        
        // Update map center to current position
        if (currentPosition.value) {
          mapCenter.value = [currentPosition.value.latitude, currentPosition.value.longitude]
        }
        
        if (showSuccessMessage) {
          showStatus('Position refreshed successfully!', 'success')
        }
      } catch (error) {
        // Only real errors (network, 500, etc.) will reach here
        console.error('Error fetching position:', error)
        showStatus('Failed to fetch position. Please try again.', 'error')
      } finally {
        isLoading.value = false
      }
    }

    const clearPosition = async () => {
      isLoading.value = true
      try {
        await api.delete('/api/stakeholder/position')
        currentPosition.value = null
        mapCenter.value = [44.7866, 20.4489] // Reset to Belgrade
        showStatus('Position cleared successfully!', 'success')
      } catch (error) {
        console.error('Error clearing position:', error)
        showStatus('Failed to clear position. Please try again.', 'error')
      } finally {
        isLoading.value = false
      }
    }

    onMounted(() => {
      refreshPosition(false) // Don't show messages on initial load
    })

    return {
      currentPosition,
      mapCenter,
      mapComponent,
      isLoading,
      statusMessage,
      statusType,
      positionMarkers,
      formatDate,
      onMapClick,
      refreshPosition,
      clearPosition
    }
  }
}
</script>

<style scoped>
.position-simulator {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.simulator-header {
  text-align: center;
  margin-bottom: 30px;
}

.simulator-header h2 {
  color: #333;
  margin-bottom: 10px;
}

.subtitle {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.current-position {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
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

.no-position {
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  text-align: center;
}

.no-position p {
  margin: 0;
  color: #856404;
}

.map-container {
  margin-bottom: 20px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.actions {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-bottom: 20px;
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

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #5a6268;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #c82333;
}

.status-message {
  padding: 12px 20px;
  border-radius: 6px;
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

@media (max-width: 768px) {
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