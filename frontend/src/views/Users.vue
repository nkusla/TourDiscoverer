<template>
  <div class="container mt-4">
    <h1>User Management</h1>
    <p class="text-muted">Manage all registered users in the system</p>

    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-content-center">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="alert alert-danger">
      {{ error }}
    </div>

    <!-- Users table -->
    <div v-if="!loading && !error" class="card">
      <div class="card-header">
        <h5 class="mb-0">All Users ({{ users.length }})</h5>
      </div>
      <div class="card-body">
        <div v-if="users.length === 0" class="text-center text-muted">
          No users found.
        </div>
        <div v-else class="table-responsive">
          <table class="table table-striped">
            <thead>
              <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>{{ user.id }}</td>
                <td>
                  <strong>{{ user.username }}</strong>
                </td>
                <td>{{ user.email }}</td>
                <td>
                  <span class="badge" :class="getRoleBadgeClass(user.role)">
                    {{ user.role }}
                  </span>
                </td>
                <td>
                  <span class="badge" :class="getStatusBadgeClass(user.is_blocked)">
                    {{ user.is_blocked ? 'Banned' : 'Active' }}
                  </span>
                </td>
                <td>
                  <button
                    v-if="user.role !== 'admin' && user.username !== currentUser?.username"
                    @click="toggleUserBlock(user)"
                    :disabled="actionLoading[user.id]"
                    class="btn btn-sm"
                    :class="user.is_blocked ? 'btn-success' : 'btn-warning'"
                  >
                    <span v-if="actionLoading[user.id]" class="spinner-border spinner-border-sm me-1"></span>
                    {{ user.is_blocked ? 'Unblock' : 'Block' }}
                  </button>
                  <span v-else class="text-muted">
                    {{ user.role === 'admin' ? 'Admin' : 'Self' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../stores/user'
import api from '../services/api'

export default {
  name: 'Users',
  setup() {
    const userStore = useUserStore()
    const users = ref([])
    const loading = ref(false)
    const error = ref('')
    const actionLoading = ref({})

    const currentUser = computed(() => userStore.user)

    const fetchUsers = async () => {
      loading.value = true
      error.value = ''

      try {
        const response = await api.get('/api/auth/users')
        users.value = response.data || []
      } catch (err) {
        console.error('Error fetching users:', err)
        error.value = err.response?.data?.message || 'Failed to fetch users'
      } finally {
        loading.value = false
      }
    }

    const toggleUserBlock = async (user) => {
      // Set loading state for this specific user
      actionLoading.value[user.id] = true

      try {
        await api.post('/api/auth/block', {
          username: user.username
        })

        // Toggle the banned status locally
        user.is_blocked = !user.is_blocked

        // Show success message
        console.log(`User ${user.username} ${user.is_blocked ? 'blocked' : 'unblocked'} successfully`)

      } catch (err) {
        console.error('Error toggling user block status:', err)
        error.value = err.response?.data?.message || 'Failed to update user status'
      } finally {
        actionLoading.value[user.id] = false
      }
    }

    const getRoleBadgeClass = (role) => {
      switch (role) {
        case 'admin':
          return 'bg-danger'
        case 'tourist':
          return 'bg-primary'
        case 'author':
          return 'bg-success'
        default:
          return 'bg-secondary'
      }
    }

    const getStatusBadgeClass = (isBlocked) => {
      return isBlocked ? 'bg-danger' : 'bg-success'
    }

    // Check if user is admin before mounting
    onMounted(() => {
      if (!userStore.isAdmin) {
        error.value = 'Access denied. Admin privileges required.'
        return
      }
      fetchUsers()
    })

    return {
      users,
      loading,
      error,
      actionLoading,
      currentUser,
      toggleUserBlock,
      getRoleBadgeClass,
      getStatusBadgeClass
    }
  }
}
</script>

<style scoped>
.table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.badge {
  font-size: 0.75em;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}

.btn:disabled {
  opacity: 0.6;
}
</style>
