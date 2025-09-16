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
  emits: ["keypoint-added", "keypoint-removed"],
  setup(props, { emit }) {
    const mapId = ref(`map-${Math.random().toString(36).substr(2, 9)}`)
    const map = ref(null)
    const markers = ref([])
    const polyline = ref(null)

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
          const newKeyPoint = {
            id: Date.now(),
            name: `Point ${props.keyPoints.length + 1}`,
            description: "",
            latitude: e.latlng.lat,
            longitude: e.latlng.lng,
            order: props.keyPoints.length
          }
          emit("keypoint-added", newKeyPoint)
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
        const marker = L.marker([point.latitude, point.longitude])
          .bindPopup(createPopupContent(point, index))
          .addTo(map.value)

        markers.value.push(marker)
      })

      // Draw polyline if >1 point
      if (props.keyPoints.length > 1) {
        const sorted = [...props.keyPoints].sort((a, b) => (a.order || 0) - (b.order || 0))
        const coords = sorted.map((p) => [p.latitude, p.longitude])

        polyline.value = L.polyline(coords, {
          color: "#007bff",
          weight: 4,
          opacity: 0.7
        }).addTo(map.value)
      }

      // Auto-fit bounds
      if (props.keyPoints.length > 0) {
        const group = new L.featureGroup(markers.value)
        map.value.fitBounds(group.getBounds().pad(0.1))
      }
    }

    const createPopupContent = (point, index) => {
      const container = document.createElement("div")
      container.innerHTML = `
        <div>
          <h6>${point.name}</h6>
          <p>${point.description || ""}</p>
          <small>Order: ${point.order ?? index + 1}</small>
        </div>
      `

      if (props.editable) {
        const btn = document.createElement("button")
        btn.className = "btn btn-sm btn-danger mt-2"
        btn.innerText = "Remove"
        btn.onclick = () => emit("keypoint-removed", index)
        container.appendChild(btn)
      }

      return container
    }

    // Watch for keypoint updates
    watch(() => props.keyPoints, updateMarkers, { deep: true })

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
