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
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/tour/create">Create Tour</router-link>
          </li>
        </ul>
        
        <ul class="navbar-nav">
          <li class="nav-item" v-if="!isAuthenticated">
            <router-link class="nav-link" to="/login">Login</router-link>
          </li>
          <li class="nav-item dropdown" v-else>
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="fas fa-user me-1"></i>{{ userStore.username }}
              <span class="badge bg-light text-dark ms-1">{{ userStore.user?.role }}</span>
            </a>
            <ul class="dropdown-menu">
              <li><h6 class="dropdown-header">{{ userStore.user?.email }}</h6></li>
              <li><hr class="dropdown-divider"></li>
              <li><router-link class="dropdown-item" to="/profile">
                <i class="fas fa-user me-2"></i>Profile
              </router-link></li>
              <li><a class="dropdown-item" href="#" @click="logout">
                <i class="fas fa-sign-out-alt me-2"></i>Logout
              </a></li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

export default {
  name: 'Navbar',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const isAuthenticated = computed(() => userStore.isAuthenticated)
    
    const logout = () => {
      userStore.logout()
      router.push('/login')
    }
    
    return {
      userStore,
      isAuthenticated,
      logout
    }
  }
}
</script>
