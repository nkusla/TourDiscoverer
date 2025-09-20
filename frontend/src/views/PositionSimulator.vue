<template>
  <div class="position-simulator-page">
    <PositionSimulator />
  </div>
</template>

<script>
import PositionSimulator from '../components/PositionSimulator.vue'

export default {
  name: 'PositionSimulatorView',
  components: {
    PositionSimulator
  }
}
</script>

<style scoped>
.position-simulator-page {
  min-height: calc(100vh - 70px); /* Adjust based on navbar height */
  background: #f8f9fa;
  padding: 20px 0;
}
</style>