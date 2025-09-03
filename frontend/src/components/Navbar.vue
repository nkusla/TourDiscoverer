<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
    <div class="container">
      <router-link class="navbar-brand" to="/">
        <strong>Tour Discoverer</strong>
      </router-link>
      
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
        <span class="navbar-toggler-icon"></span>
      </button>
      
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <router-link class="nav-link" to="/">Home</router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/tours">Tours</router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/tour/create">Create Tour</router-link>
          </li>
        </ul>
        
        <ul class="navbar-nav">
          <li class="nav-item" v-if="!isAuthenticated">
            <router-link class="nav-link" to="/login">Login</router-link>
          </li>
          <li class="nav-item dropdown" v-else>
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              {{ userStore.username }}
            </a>
            <ul class="dropdown-menu">
              <li><a class="dropdown-item" href="#" @click="logout">Logout</a></li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
import { computed } from 'vue'
import { useUserStore } from '../stores/user'

export default {
  name: 'Navbar',
  setup() {
    const userStore = useUserStore()
    
    const isAuthenticated = computed(() => userStore.isAuthenticated)
    
    const logout = () => {
      userStore.logout()
    }
    
    return {
      userStore,
      isAuthenticated,
      logout
    }
  }
}
</script>
