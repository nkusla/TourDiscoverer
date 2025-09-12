<template>
  <div class="map-container" :style="{ height: mapHeight }">
    <l-map
      ref="map"
      :zoom="zoom"
      :center="center"
      @click="onMapClick"
      :options="mapOptions"
    >
      <l-tile-layer
        :url="tileLayerUrl"
        :attribution="attribution"
      />
      
      <!-- Tour route polyline -->
      <l-polyline
        v-if="keyPoints.length > 1"
        :lat-lngs="routeCoordinates"
        :color="routeColor"
        :weight="4"
        :opacity="0.7"
      />
      
      <!-- Key points markers -->
      <l-marker
        v-for="(point, index) in keyPoints"
        :key="point.id || index"
        :lat-lng="[point.latitude, point.longitude]"
        :draggable="editable"
        @dragend="onMarkerDragEnd(point, $event)"
        @click="onMarkerClick(point, index)"
      >
        <l-icon
          :icon-size="[25, 41]"
          :icon-anchor="[12, 41]"
          :popup-anchor="[1, -34]"
          :shadow-size="[41, 41]"
          :icon-url="getMarkerIcon(index)"
        />
        
        <l-popup>
          <div class="popup-content">
            <h6>{{ point.name }}</h6>
            <p v-if="point.description">{{ point.description }}</p>
            <small>Order: {{ point.order || index + 1 }}</small>
            <div class="mt-2" v-if="editable">
              <button 
                class="btn btn-sm btn-primary me-1" 
                @click="editKeyPoint(point, index)"
              >
                Edit
              </button>
              <button 
                class="btn btn-sm btn-danger" 
                @click="removeKeyPoint(index)"
              >
                Remove
              </button>
            </div>
          </div>
        </l-popup>
      </l-marker>
    </l-map>
    
    <!-- Map controls -->
    <div class="map-controls" v-if="editable">
      <div class="card">
        <div class="card-body p-2">
          <small class="text-muted">
            Click on map to add key points
          </small>
          <div class="mt-1">
            <button 
              class="btn btn-sm btn-secondary me-1"
              @click="clearRoute"
              :disabled="keyPoints.length === 0"
            >
              Clear All
            </button>
            <button 
              class="btn btn-sm btn-success"
              @click="saveRoute"
              :disabled="keyPoints.length === 0"
            >
              Save Route
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { LMap, LTileLayer, LMarker, LPolyline, LPopup, LIcon } from '@vue-leaflet/vue-leaflet'

export default {
  name: 'TourMap',
  components: {
    LMap,
    LTileLayer,
    LMarker,
    LPolyline,
    LPopup,
    LIcon
  },
  props: {
    keyPoints: {
      type: Array,
      default: () => []
    },
    editable: {
      type: Boolean,
      default: false
    },
    mapHeight: {
      type: String,
      default: '500px'
    },
    initialCenter: {
      type: Array,
      default: () => [44.7866, 20.4489] // Belgrade, Serbia
    },
    initialZoom: {
      type: Number,
      default: 10
    }
  },
  emits: ['keypoint-added', 'keypoint-updated', 'keypoint-removed', 'route-saved', 'route-cleared'],
  setup(props, { emit }) {
    const map = ref(null)
    const zoom = ref(props.initialZoom)
    const center = ref([...props.initialCenter])
    
    const tileLayerUrl = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
    const attribution = '© OpenStreetMap contributors'
    const routeColor = '#007bff'
    
    const mapOptions = {
      zoomControl: true,
      attributionControl: true
    }
    
    // Compute route coordinates from key points
    const routeCoordinates = computed(() => {
      return props.keyPoints
        .sort((a, b) => (a.order || 0) - (b.order || 0))
        .map(point => [point.latitude, point.longitude])
    })
    
    // Auto-center map when key points change
    watch(() => props.keyPoints, (newPoints) => {
      if (newPoints.length > 0) {
        // Calculate bounds of all key points
        const lats = newPoints.map(p => p.latitude)
        const lngs = newPoints.map(p => p.longitude)
        
        const minLat = Math.min(...lats)
        const maxLat = Math.max(...lats)
        const minLng = Math.min(...lngs)
        const maxLng = Math.max(...lngs)
        
        // Set center to middle of bounds
        center.value = [
          (minLat + maxLat) / 2,
          (minLng + maxLng) / 2
        ]
        
        // Adjust zoom based on spread
        const latSpread = maxLat - minLat
        const lngSpread = maxLng - minLng
        const maxSpread = Math.max(latSpread, lngSpread)
        
        if (maxSpread > 1) {
          zoom.value = 8
        } else if (maxSpread > 0.1) {
          zoom.value = 10
        } else {
          zoom.value = 13
        }
      }
    }, { immediate: true })
    
    const onMapClick = (event) => {
      if (!props.editable) return
      
      const { lat, lng } = event.latlng
      const newKeyPoint = {
        id: Date.now(), // Temporary ID
        name: `Point ${props.keyPoints.length + 1}`,
        description: '',
        latitude: lat,
        longitude: lng,
        order: props.keyPoints.length
      }
      
      emit('keypoint-added', newKeyPoint)
    }
    
    const onMarkerDragEnd = (point, event) => {
      const { lat, lng } = event.target.getLatLng()
      const updatedPoint = {
        ...point,
        latitude: lat,
        longitude: lng
      }
      
      emit('keypoint-updated', updatedPoint)
    }
    
    const onMarkerClick = (point, index) => {
      // Handle marker click if needed
    }
    
    const editKeyPoint = (point, index) => {
      // Emit event to open edit dialog
      emit('keypoint-edit', { point, index })
    }
    
    const removeKeyPoint = (index) => {
      emit('keypoint-removed', index)
    }
    
    const clearRoute = () => {
      emit('route-cleared')
    }
    
    const saveRoute = () => {
      emit('route-saved')
    }
    
    const getMarkerIcon = (index) => {
      // Return different icons based on order or use numbered markers
      return `https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-blue.png`
    }
    
    return {
      map,
      zoom,
      center,
      tileLayerUrl,
      attribution,
      mapOptions,
      routeColor,
      routeCoordinates,
      onMapClick,
      onMarkerDragEnd,
      onMarkerClick,
      editKeyPoint,
      removeKeyPoint,
      clearRoute,
      saveRoute,
      getMarkerIcon
    }
  }
}
</script>

<style scoped>
.map-container {
  position: relative;
  width: 100%;
}

.map-controls {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 1000;
  max-width: 200px;
}

.popup-content h6 {
  margin: 0 0 5px 0;
}

.popup-content p {
  margin: 0 0 5px 0;
  font-size: 0.9em;
}

:deep(.leaflet-container) {
  height: 100%;
  width: 100%;
}
</style>
