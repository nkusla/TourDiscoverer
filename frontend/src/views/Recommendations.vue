<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <h1 class="mb-4">
        	Friend Recommendations
        </h1>
        <p class="text-muted mb-4">
          Discover new people to follow based on who your friends are following!
        </p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-3">Finding recommendations for you...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>

    <!-- No Recommendations -->
    <div v-else-if="recommendations.length === 0" class="text-center my-5">
      <div class="card border-0 bg-light">
        <div class="card-body py-5">
          <h3 class="text-muted">No recommendations yet</h3>
          <p class="text-muted">
            Start following some people to get personalized recommendations!
          </p>
          <router-link to="/users" class="btn btn-primary">
            Browse Users
          </router-link>
        </div>
      </div>
    </div>

    <!-- Recommendations Grid -->
    <div v-else class="row">
      <div
        v-for="user in recommendations"
        :key="user.username"
        class="col-12 col-md-6 col-lg-4 mb-4"
      >
        <div class="card h-100 shadow-sm recommendation-card">
          <div class="card-body d-flex flex-column">
            <!-- User Avatar -->
            <div class="text-center mb-3">
              <div class="user-avatar mx-auto mb-2">
								👤
              </div>
              <h5 class="card-title mb-1">{{ user.username }}</h5>
              <span class="badge" :class="getRoleBadgeClass(user.role)">
                {{ user.role || 'User' }}
              </span>
            </div>

            <!-- User Info -->
            <div class="mb-3 flex-grow-1">
              <p class="text-muted small text-center">
                Suggested based on your connections
              </p>
            </div>

            <!-- Action Buttons -->
            <div class="d-grid gap-2">
              <button
                class="btn btn-primary btn-sm"
                @click="followUser(user.username)"
                :disabled="followingUsers.has(user.username)"
              >
                {{ followingUsers.has(user.username) ? 'Following...' : 'Follow' }}
              </button>
              <button
                class="btn btn-outline-secondary btn-sm"
                @click="viewProfile(user.username)"
              >
                View Profile
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Refresh Button -->
    <div class="text-center mt-4" v-if="recommendations.length > 0">
      <button
        class="btn btn-outline-primary"
        @click="fetchRecommendations"
        :disabled="loading"
      >
        Refresh Recommendations
      </button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import api from '../services/api'

export default {
  name: 'Recommendations',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const recommendations = ref([])
    const loading = ref(false)
    const error = ref(null)
    const followingUsers = ref(new Set())

    const currentUser = computed(() => userStore.user)

    const fetchRecommendations = async () => {
      if (!currentUser.value?.username) {
        error.value = 'Please log in to see recommendations'
        return
      }

      loading.value = true
      error.value = null

      try {
        const response = await api.get(`/api/followers/user/${currentUser.value.username}/recommendations`)
        recommendations.value = response.data || []
      } catch (err) {
        console.error('Error fetching recommendations:', err)
        error.value = err.response?.data?.message || 'Failed to load recommendations'
      } finally {
        loading.value = false
      }
    }

    const followUser = async (username) => {
      if (!currentUser.value?.username) {
        error.value = 'Please log in to follow users'
        return
      }

      followingUsers.value.add(username)

      try {
        await api.post('/api/followers/follow', {
          follower: currentUser.value.username,
          followee: username
        })

        // Remove user from recommendations after successful follow
        recommendations.value = recommendations.value.filter(user => user.username !== username)

        // Show success message
        // You might want to add a toast notification here
        console.log(`Successfully followed ${username}`)
      } catch (err) {
        console.error('Error following user:', err)
        error.value = err.response?.data?.message || `Failed to follow ${username}`
      } finally {
        followingUsers.value.delete(username)
      }
    }

    const viewProfile = (username) => {
      // Navigate to user profile - adjust route as needed
      router.push(`/profile/${username}`)
    }

    const getRoleBadgeClass = (role) => {
      switch (role?.toLowerCase()) {
        case 'guide':
          return 'badge-success'
        case 'tourist':
          return 'badge-info'
        default:
          return 'badge-secondary'
      }
    }

    onMounted(() => {
      fetchRecommendations()
    })

    return {
      recommendations,
      loading,
      error,
      followingUsers,
      fetchRecommendations,
      followUser,
      viewProfile,
      getRoleBadgeClass
    }
  }
}
</script>

<style scoped>
.recommendation-card {
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.recommendation-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.badge-success {
  background-color: #28a745;
  color: white;
}

.badge-info {
  background-color: #17a2b8;
  color: white;
}

.badge-secondary {
  background-color: #6c757d;
  color: white;
}

.card-title {
  color: #495057;
  font-weight: 600;
}

.btn-sm {
  font-size: 0.875rem;
  padding: 0.375rem 0.75rem;
}
</style>