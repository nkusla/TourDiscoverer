<template>
  <div class="home">
    <div class="hero-section bg-primary text-white py-5">
      <div class="container">
        <div class="row">
          <div class="col-lg-8 mx-auto text-center">
            <h1 class="display-4 mb-4">Discover Amazing Tours</h1>
            <p class="lead mb-4">
              Explore incredible destinations and create unforgettable memories with our curated tours.
            </p>
            <router-link to="/tours" class="btn btn-light btn-lg">
              Browse Tours
            </router-link>
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
export default {
  name: 'Home'
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
