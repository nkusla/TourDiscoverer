<template>
  <div class="profile">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-md-8">
          <div class="card shadow">
            <div class="card-header bg-primary text-white">
              <h4 class="mb-0">
                <i class="fas fa-user me-2"></i>User Profile
              </h4>
            </div>
            <div class="card-body">
              <!-- Loading State -->
              <div v-if="loading" class="text-center py-4">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-2">Loading profile...</p>
              </div>

              <!-- Error State -->
              <div v-else-if="updateError && !editing" class="alert alert-danger" role="alert">
                {{ updateError }}
              </div>

              <!-- Profile View -->
              <div v-else-if="!editing" class="profile-view">
                <div class="row mb-3">
                  <div class="col-sm-3">
                    <strong>Name:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ fullProfile?.first_name || 'Not set' }} {{ fullProfile?.last_name || '' }}
                  </div>
                </div>

                <div class="row mb-3">
                  <div class="col-sm-3">
                    <strong>Username:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ user?.username }}
                  </div>
                </div>

                <div class="row mb-3">
                  <div class="col-sm-3">
                    <strong>Role:</strong>
                  </div>
                  <div class="col-sm-9">
                    <span class="badge" :class="roleBadgeClass">
                      {{ user?.role }}
                    </span>
                  </div>
                </div>

                <div class="row mb-3" v-if="fullProfile?.biography">
                  <div class="col-sm-3">
                    <strong>Biography:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ fullProfile.biography }}
                  </div>
                </div>

                <div class="row mb-3" v-if="fullProfile?.motto">
                  <div class="col-sm-3">
                    <strong>Motto:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ fullProfile.motto }}
                  </div>
                </div>

                <div class="d-flex gap-2">
                  <button class="btn btn-primary" @click="startEditing">
                    <i class="fas fa-edit me-2"></i>Edit Profile
                  </button>
                  <button class="btn btn-outline-danger" @click="confirmLogout">
                    <i class="fas fa-sign-out-alt me-2"></i>Logout
                  </button>
                </div>
              </div>

              <!-- Edit Form -->
              <div v-else class="profile-edit">
                <form @submit.prevent="saveProfile">
                  <div class="row mb-3">
                    <div class="col-md-6">
                      <label class="form-label">First Name</label>
                      <input
                        v-model="editForm.first_name"
                        type="text"
                        class="form-control"
                        placeholder="Enter your first name"
                      />
                    </div>
                    <div class="col-md-6">
                      <label class="form-label">Last Name</label>
                      <input
                        v-model="editForm.last_name"
                        type="text"
                        class="form-control"
                        placeholder="Enter your last name"
                      />
                    </div>
                  </div>

                  <div class="mb-3">
                    <label class="form-label">Biography</label>
                    <textarea
                      v-model="editForm.biography"
                      class="form-control"
                      rows="4"
                      placeholder="Tell others about yourself..."
                    ></textarea>
                  </div>

                  <div class="mb-3">
                    <label class="form-label">Motto</label>
                    <input
                      v-model="editForm.motto"
                      type="text"
                      class="form-control"
                      placeholder="Your personal motto or tagline"
                    />
                  </div>

                  <div v-if="updateError" class="alert alert-danger" role="alert">
                    {{ updateError }}
                  </div>

                  <div v-if="updateSuccess" class="alert alert-success" role="alert">
                    Profile updated successfully!
                  </div>

                  <div class="d-flex gap-2">
                    <button type="submit" class="btn btn-success" :disabled="updating">
                      <span v-if="updating" class="spinner-border spinner-border-sm me-2"></span>
                      {{ updating ? 'Saving...' : 'Save Changes' }}
                    </button>
                    <button type="button" class="btn btn-secondary" @click="cancelEditing">
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

export default {
  name: 'Profile',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()

    const editing = ref(false)
    const updating = ref(false)
    const loading = ref(false)
    const updateError = ref('')
    const updateSuccess = ref(false)
    const fullProfile = ref(null)

    const user = computed(() => userStore.user)
    const roleBadgeClass = computed(() => {
      return user.value?.role === 'guide' ? 'bg-success' : 'bg-info'
    })

    const editForm = ref({
      first_name: '',
      last_name: '',
      biography: '',
      motto: ''
    })

    const fetchProfile = async () => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }

      try {
        loading.value = true
        fullProfile.value = await userStore.fetchUserProfile()
      } catch (error) {
        console.error('Error fetching profile:', error)
        updateError.value = 'Failed to load profile data'
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchProfile()
    })

    const startEditing = () => {
      editForm.value = {
        first_name: fullProfile.value?.first_name || '',
        last_name: fullProfile.value?.last_name || '',
        biography: fullProfile.value?.biography || '',
        motto: fullProfile.value?.motto || ''
      }
      editing.value = true
      updateError.value = ''
      updateSuccess.value = false
    }

    const cancelEditing = () => {
      editing.value = false
      updateError.value = ''
      updateSuccess.value = false
    }

    const saveProfile = async () => {
      updating.value = true
      updateError.value = ''
      updateSuccess.value = false

      try {
        await userStore.updateUserProfile(editForm.value)
        await fetchProfile() // Refresh the profile data
        updateSuccess.value = true
        editing.value = false
      } catch (error) {
        updateError.value = error.message || 'Failed to update profile'
      } finally {
        updating.value = false
      }
    }

    const confirmLogout = () => {
      if (confirm('Are you sure you want to logout?')) {
        userStore.logout()
        router.push('/login')
      }
    }

    return {
      editing,
      updating,
      loading,
      updateError,
      updateSuccess,
      user,
      fullProfile,
      roleBadgeClass,
      editForm,
      startEditing,
      cancelEditing,
      saveProfile,
      confirmLogout
    }
  }
}
</script>

<style scoped>
.profile {
  min-height: calc(100vh - 56px);
  background-color: #f8f9fa;
  padding-top: 2rem;
}

.card {
  border: none;
  border-radius: 10px;
}

.badge {
  font-size: 0.875em;
  text-transform: capitalize;
}
</style>
