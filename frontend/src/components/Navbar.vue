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
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/">Home</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/tours">Tours</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/blogs">Blogs</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/tour/create">Create Tour</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/blog/create">Create Blog</router-link>
          </li>
          <li class="nav-item" v-if="isTourist">
            <router-link class="nav-link" to="/position-simulator">Position Simulator</router-link>
          </li>
          <li class="nav-item" v-if="isTourist">
            <router-link class="nav-link" to="/tour-execution">Execute Tours</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/recommendations">
              Discover People
            </router-link>
          </li>
          <li class="nav-item" v-if="isAdmin">
            <router-link class="nav-link" to="/users">Users</router-link>
          </li>
          <li class="nav-item" v-if="isTourist">
            <router-link class="nav-link" to="/cart">
              <i class="fas fa-shopping-cart me-1"></i>
              Cart
              <span v-if="cartItemCount > 0" class="badge bg-warning text-dark ms-1">
                {{ cartItemCount }}
              </span>
            </router-link>
          </li>
          <li class="nav-item" v-if="isTourist">
            <router-link class="nav-link" to="/purchases">
              <i class="fas fa-ticket-alt me-1"></i>My Tours
            </router-link>
          </li>
        </ul>

        <ul class="navbar-nav">
          <li class="nav-item" v-if="!isAuthenticated">
            <router-link class="nav-link" to="/login">Login</router-link>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <div class="dropdown-simple">
              <button class="btn btn-outline-light" @click="toggleDropdown">
                {{ userStore.user?.username || 'User' }} ▼
              </button>
              <div v-if="showDropdown" class="dropdown-content">
                <router-link to="/profile" @click="showDropdown = false">Profile</router-link>
                <button @click="logout">Logout</button>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useCartStore } from '../stores/cart'

export default {
  name: 'Navbar',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const cartStore = useCartStore()
    const showDropdown = ref(false)

    const isAuthenticated = computed(() => {
      return userStore.isAuthenticated
    })

    const isAdmin = computed(() => {
      return userStore.isAdmin
    })
    
    const isTourist = computed(() => {
      return userStore.isTourist
    })
    
    const cartItemCount = computed(() => {
      return cartStore.itemCount
    })
    
    onMounted(async () => {
      if (userStore.isAuthenticated && userStore.isTourist) {
        await cartStore.fetchCart()
      }
    })

    const toggleDropdown = () => {
      console.log('Dropdown clicked! Current state:', showDropdown.value)
      showDropdown.value = !showDropdown.value
      console.log('New state:', showDropdown.value)
    }

    const logout = () => {
      showDropdown.value = false
      userStore.logout()
      router.push('/login')
    }

    return {
      userStore,
      isAuthenticated,
      isAdmin,
      isTourist,
      cartItemCount,
      showDropdown,
      toggleDropdown,
      logout
    }
  }
}
</script>

<style scoped>
.dropdown-simple {
  position: relative;
  display: inline-block;
}

.dropdown-content {
  position: absolute;
  right: 0;
  background-color: white;
  min-width: 160px;
  box-shadow: 0px 8px 16px rgba(0,0,0,0.2);
  z-index: 1000;
  border-radius: 4px;
  padding: 8px 0;
  margin-top: 4px;
}

.dropdown-content a,
.dropdown-content button {
  color: black;
  padding: 8px 16px;
  text-decoration: none;
  display: block;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
}

.dropdown-content a:hover,
.dropdown-content button:hover {
  background-color: #f1f1f1;
}
</style>
