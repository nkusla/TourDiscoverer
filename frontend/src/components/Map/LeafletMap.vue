<template>
  <div class="leaflet-map" :id="mapId" :style="{ height: mapHeight, width: '100%' }"></div>
</template>

<script>
import { ref, onMounted, watch, nextTick } from "vue"
import * as L from "leaflet"
import "leaflet/dist/leaflet.css"

export default {
  name: "LeafletMap",
  props: {
    keyPoints: {
      type: Array,
      default: () => []
    },
    editable: {
      type: Boolean,
      default: false
    },
    selectedPointIndex: {
      type: Number,
      default: -1
    },
    mapHeight: {
      type: String,
      default: "500px"
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
  emits: ["keypoint-added", "keypoint-removed", "keypoint-moved", "keypoint-clicked", "map-clicked"],
  setup(props, { emit }) {
    const mapId = ref(`map-${Math.random().toString(36).substr(2, 9)}`)
    const map = ref(null)
    const markers = ref([])
    const polyline = ref(null)
    const isDraggingMarker = ref(false)

    // Custom marker icons
    const createMarkerIcon = (isSelected = false, isFirst = false, isLast = false) => {
      let iconUrl, iconSize, iconAnchor
      
      if (isFirst) {
        iconUrl = 'data:image/svg+xml;base64,' + btoa(`
          <svg width="25" height="41" viewBox="0 0 25 41" xmlns="http://www.w3.org/2000/svg">
            <path d="M12.5 0C5.6 0 0 5.6 0 12.5c0 12.5 12.5 28.5 12.5 28.5s12.5-16 12.5-28.5C25 5.6 19.4 0 12.5 0z" fill="#28a745"/>
            <circle cx="12.5" cy="12.5" r="6" fill="white"/>
            <text x="12.5" y="17" text-anchor="middle" font-family="Arial" font-size="10" fill="#28a745" font-weight="bold">S</text>
          </svg>
        `)
        iconSize = [25, 41]
        iconAnchor = [12, 41]
      } else if (isLast) {
        iconUrl = 'data:image/svg+xml;base64,' + btoa(`
          <svg width="25" height="41" viewBox="0 0 25 41" xmlns="http://www.w3.org/2000/svg">
            <path d="M12.5 0C5.6 0 0 5.6 0 12.5c0 12.5 12.5 28.5 12.5 28.5s12.5-16 12.5-28.5C25 5.6 19.4 0 12.5 0z" fill="#dc3545"/>
            <circle cx="12.5" cy="12.5" r="6" fill="white"/>
            <text x="12.5" y="17" text-anchor="middle" font-family="Arial" font-size="10" fill="#dc3545" font-weight="bold">E</text>
          </svg>
        `)
        iconSize = [25, 41]
        iconAnchor = [12, 41]
      } else {
        const color = isSelected ? '#007bff' : '#6c757d'
        iconUrl = 'data:image/svg+xml;base64,' + btoa(`
          <svg width="20" height="32" viewBox="0 0 20 32" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 0C4.5 0 0 4.5 0 10c0 10 10 22 10 22s10-12 10-22C20 4.5 15.5 0 10 0z" fill="${color}"/>
            <circle cx="10" cy="10" r="5" fill="white"/>
          </svg>
        `)
        iconSize = [20, 32]
        iconAnchor = [10, 32]
      }
      
      return L.icon({
        iconUrl,
        iconSize,
        iconAnchor,
        popupAnchor: [0, -32]
      })
    }

    const initializeMap = async () => {
      await nextTick()

      // Fix default marker icons
      delete L.Icon.Default.prototype._getIconUrl
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: new URL("leaflet/dist/images/marker-icon-2x.png", import.meta.url).href,
        iconUrl: new URL("leaflet/dist/images/marker-icon.png", import.meta.url).href,
        shadowUrl: new URL("leaflet/dist/images/marker-shadow.png", import.meta.url).href
      })

      const mapElement = document.getElementById(mapId.value)
      if (!mapElement) return

      // Create map
      map.value = L.map(mapElement).setView(props.initialCenter, props.initialZoom)

      // Add tile layer
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: "© OpenStreetMap contributors"
      }).addTo(map.value)

      // Add click handler for editable maps
      if (props.editable) {
        map.value.on("click", (e) => {
          if (isDraggingMarker.value) {
            isDraggingMarker.value = false
            return
          }
          
          emit("map-clicked", {
            latitude: e.latlng.lat,
            longitude: e.latlng.lng
          })
        })
      }

      updateMarkers()
    }

    const updateMarkers = () => {
      if (!map.value) return

      // Clear old markers
      markers.value.forEach((m) => map.value.removeLayer(m))
      markers.value = []

      // Clear polyline
      if (polyline.value) {
        map.value.removeLayer(polyline.value)
        polyline.value = null
      }

      // Add new markers
      props.keyPoints.forEach((point, index) => {
        const isSelected = index === props.selectedPointIndex
        const isFirst = index === 0
        const isLast = index === props.keyPoints.length - 1
        
        const marker = L.marker([point.latitude, point.longitude], {
          icon: createMarkerIcon(isSelected, isFirst, isLast),
          draggable: props.editable
        }).addTo(map.value)

        // Create popup content
        const popupContent = `
          <div>
            <strong>${point.name}</strong><br>
            ${point.description ? point.description + '<br>' : ''}
            <small>Lat: ${point.latitude.toFixed(4)}, Lng: ${point.longitude.toFixed(4)}</small>
            ${props.editable ? '<br><small>🖱️ Drag to move</small>' : ''}
          </div>
        `
        marker.bindPopup(popupContent)

        // Handle marker events
        if (props.editable) {
          marker.on('click', () => {
            emit('keypoint-clicked', index)
          })
          
          marker.on('dragstart', () => {
            isDraggingMarker.value = true
          })
          
          marker.on('dragend', (e) => {
            const newPos = e.target.getLatLng()
            emit('keypoint-moved', index, {
              latitude: newPos.lat,
              longitude: newPos.lng
            })
            
            // Update popup content with new coordinates
            const updatedPopupContent = `
              <div>
                <strong>${point.name}</strong><br>
                ${point.description ? point.description + '<br>' : ''}
                <small>Lat: ${newPos.lat.toFixed(4)}, Lng: ${newPos.lng.toFixed(4)}</small>
                <br><small>🖱️ Drag to move</small>
              </div>
            `
            marker.setPopupContent(updatedPopupContent)
            
            setTimeout(() => {
              isDraggingMarker.value = false
            }, 100)
          })
        } else {
          marker.on('click', () => {
            emit('keypoint-clicked', index)
          })
        }

        markers.value.push(marker)
      })

      // Draw polyline if >1 point
      if (props.keyPoints.length > 1) {
        const sorted = [...props.keyPoints].sort((a, b) => (a.order || 0) - (b.order || 0))
        const coords = sorted.map((p) => [p.latitude, p.longitude])

        polyline.value = L.polyline(coords, {
          color: "#007bff",
          weight: 4,
          opacity: 0.7,
          dashArray: props.editable ? '10, 5' : null
        }).addTo(map.value)
      }

      // Auto-fit bounds
      if (props.keyPoints.length > 0) {
        const group = new L.featureGroup(markers.value)
        if (polyline.value) {
          group.addLayer(polyline.value)
        }
        if (group.getBounds().isValid()) {
          map.value.fitBounds(group.getBounds(), { padding: [20, 20] })
        }
      }
    }

    // Watch for keypoint updates
    watch(() => props.keyPoints, updateMarkers, { deep: true })
    watch(() => props.selectedPointIndex, updateMarkers)

    onMounted(() => {
      initializeMap()
    })

    return { mapId }
  }
}
</script>

<style>
.leaflet-map {
  border-radius: 8px;
}
</style>
