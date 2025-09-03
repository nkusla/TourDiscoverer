<template>
  <div id="app">
    <Navbar />
    <main class="container-fluid p-0">
      <router-view />
    </main>
  </div>
</template>

<script>
import Navbar from './components/Navbar.vue'

export default {
  name: 'App',
  components: {
    Navbar
  }
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

body {
  margin: 0;
  padding: 0;
}

.leaflet-container {
  height: 100%;
  width: 100%;
}
</style>
