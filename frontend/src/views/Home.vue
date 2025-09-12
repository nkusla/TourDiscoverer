<template>
  <div class="home">
    <div class="hero-section bg-primary text-white py-5">
      <div class="container">
        <div class="row">
          <div class="col-lg-8 mx-auto text-center">
            <h1 class="display-4 mb-4" v-if="!isAuthenticated">Discover Amazing Tours</h1>
            <h1 class="display-4 mb-4" v-else>Welcome back, {{ userStore.username }}!</h1>
            <p class="lead mb-4" v-if="!isAuthenticated">
              Explore incredible destinations and create unforgettable memories with our curated tours.
            </p>
            <p class="lead mb-4" v-else>
              Ready to explore new adventures? Browse tours or create your own unique experience.
            </p>
            <div class="d-grid d-md-flex justify-content-md-center gap-2">
              <router-link to="/tours" class="btn btn-light btn-lg">
                Browse Tours
              </router-link>
              <router-link v-if="isAuthenticated" to="/tour/create" class="btn btn-outline-light btn-lg">
                Create Tour
              </router-link>
              <router-link v-else to="/login" class="btn btn-outline-light btn-lg">
                Join Now
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="container py-5">
      <div class="row">
        <div class="col-lg-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="fas fa-map-marked-alt fa-3x text-primary mb-3"></i>
              <h5 class="card-title">Interactive Maps</h5>
              <p class="card-text">
                View tours on interactive maps with detailed key points and routes.
              </p>
            </div>
          </div>
        </div>
        
        <div class="col-lg-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="fas fa-plus-circle fa-3x text-success mb-3"></i>
              <h5 class="card-title">Create Tours</h5>
              <p class="card-text">
                Design your own tours by adding key points directly on the map.
              </p>
            </div>
          </div>
        </div>
        
        <div class="col-lg-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="fas fa-users fa-3x text-info mb-3"></i>
              <h5 class="card-title">Community</h5>
              <p class="card-text">
                Share your experiences and discover tours created by others.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useUserStore } from '../stores/user'

export default {
  name: 'Home',
  setup() {
    const userStore = useUserStore()
    const isAuthenticated = computed(() => userStore.isAuthenticated)
    
    return {
      userStore,
      isAuthenticated
    }
  }
}
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, #007bff, #0056b3);
}

.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-5px);
}
</style>
