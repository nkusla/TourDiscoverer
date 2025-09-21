<template>
  <div class="profile">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-md-8">
          <div class="card shadow">
            <div class="card-header bg-primary text-white">
              <h4 class="mb-0">
                <i class="fas fa-user me-2"></i>Profile
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
                <!-- Profile Picture Display -->
                <div v-if="fullProfile?.profile_picture" class="row mb-3">
                  <div class="col-sm-3">
                    <strong>Profile Picture:</strong>
                  </div>
                  <div class="col-sm-9">
                    <img
                      :src="fullProfile.profile_picture"
                      alt="Profile Picture"
                      class="img-thumbnail"
                      style="max-width: 200px; max-height: 200px; object-fit: cover;"
                    />
                  </div>
                </div>

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

                <!-- Followers and Following Section -->
                <div class="row mb-4">
                  <div class="col-12">
                    <h5 class="mb-3">
                      <i class="fas fa-users me-2"></i>Social Connections
                    </h5>

                    <!-- Clickable Stats Cards -->
                    <div class="row mb-3">
                      <div class="col-6">
                        <div
                          class="card bg-light text-center clickable-card"
                          data-bs-toggle="modal"
                          data-bs-target="#followersModal"
                        >
                          <div class="card-body py-3">
                            <h5 class="card-title mb-1 text-primary">{{ followers.length }}</h5>
                            <small class="text-muted">Followers</small>
                            <i class="fas fa-eye ms-2 text-muted"></i>
                          </div>
                        </div>
                      </div>
                      <div class="col-6">
                        <div
                          class="card bg-light text-center clickable-card"
                          data-bs-toggle="modal"
                          data-bs-target="#followingModal"
                        >
                          <div class="card-body py-3">
                            <h5 class="card-title mb-1 text-primary">{{ following.length }}</h5>
                            <small class="text-muted">Following</small>
                            <i class="fas fa-edit ms-2 text-muted"></i>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Empty state -->
                    <div v-if="followers.length === 0 && following.length === 0" class="text-muted text-center py-3">
                      <p>No social connections yet. Start following people to build your network!</p>
                      <router-link to="/users" class="btn btn-sm btn-outline-primary">
                        Browse Users
                      </router-link>
                    </div>
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
                    <label class="form-label">Profile Picture</label>
                    <input
                      ref="fileInput"
                      type="file"
                      class="form-control"
                      accept="image/*"
                      @change="handleFileChange"
                    />
                    <div class="form-text">Select an image file (JPEG, PNG, GIF, WebP)</div>

                    <!-- Preview current or selected image -->
                    <div v-if="imagePreview || fullProfile?.profile_picture" class="mt-2">
                      <img
                        :src="imagePreview || fullProfile?.profile_picture"
                        alt="Profile Preview"
                        class="img-thumbnail"
                        style="max-width: 150px; max-height: 150px; object-fit: cover;"
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

    <!-- Followers Modal -->
    <div class="modal fade" id="followersModal" tabindex="-1" aria-labelledby="followersModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="followersModalLabel">
              <i class="fas fa-users me-2"></i>Followers ({{ followers.length }})
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div v-if="followers.length === 0" class="text-center py-4">
              <i class="fas fa-user-friends fa-3x text-muted mb-3"></i>
              <p class="text-muted">No followers yet</p>
            </div>
            <div v-else class="row">
              <div
                v-for="follower in followers"
                :key="follower.username"
                class="col-12 col-md-6 mb-3"
              >
                <div class="card border-0 bg-light">
                  <div class="card-body py-2 px-3">
                    <div class="d-flex align-items-center">
                      <div class="user-avatar-small me-3">
                        <i class="fas fa-user"></i>
                      </div>
                      <div class="flex-grow-1">
                        <h6 class="mb-0">{{ follower.username }}</h6>
                        <small class="text-muted">{{ follower.role || 'User' }}</small>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Following Modal -->
    <div class="modal fade" id="followingModal" tabindex="-1" aria-labelledby="followingModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="followingModalLabel">
              <i class="fas fa-user-plus me-2"></i>Following ({{ following.length }})
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div v-if="following.length === 0" class="text-center py-4">
              <i class="fas fa-user-plus fa-3x text-muted mb-3"></i>
              <p class="text-muted">You're not following anyone yet</p>
              <router-link to="/users" class="btn btn-primary btn-sm" data-bs-dismiss="modal">
                Browse Users
              </router-link>
            </div>
            <div v-else class="row">
              <div
                v-for="followee in following"
                :key="followee.username"
                class="col-12 col-md-6 mb-3"
              >
                <div class="card border-0 bg-light">
                  <div class="card-body py-2 px-3">
                    <div class="d-flex align-items-center justify-content-between">
                      <div class="d-flex align-items-center">
                        <div class="user-avatar-small me-3">
                          <i class="fas fa-user"></i>
                        </div>
                        <div>
                          <h6 class="mb-0">{{ followee.username }}</h6>
                          <small class="text-muted">{{ followee.role || 'User' }}</small>
                        </div>
                      </div>
                      <button
                        class="btn btn-sm btn-outline-danger"
                        @click="unfollowUser(followee.username)"
                        :disabled="unfollowingUsers.has(followee.username)"
                        title="Unfollow"
                      >
                        <span v-if="unfollowingUsers.has(followee.username)" class="spinner-border spinner-border-sm me-1"></span>
                        {{ unfollowingUsers.has(followee.username) ? 'Unfollowing...' : 'Unfollow' }}
                      </button>
                    </div>
                  </div>
                </div>
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
import api from '../services/api'

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
    const selectedFile = ref(null)
    const imagePreview = ref('')
    const fileInput = ref(null)
    const followers = ref([])
    const following = ref([])
    const unfollowingUsers = ref(new Set())

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
        await fetchFollowersAndFollowing()
      } catch (error) {
        console.error('Error fetching profile:', error)
        updateError.value = 'Failed to load profile data'
      } finally {
        loading.value = false
      }
    }

    const fetchFollowersAndFollowing = async () => {
      if (!user.value?.username) return

      try {
        const [followersResponse, followingResponse] = await Promise.all([
          api.get(`/api/followers/user/${user.value.username}/followers`),
          api.get(`/api/followers/user/${user.value.username}/following`)
        ])

        followers.value = followersResponse.data || []
        following.value = followingResponse.data || []
      } catch (error) {
        console.error('Error fetching followers/following:', error)
        // Don't show error to user for this, just log it
      }
    }

    const unfollowUser = async (username) => {
      if (!user.value?.username) return

      unfollowingUsers.value.add(username)

      try {
        await api.delete('/api/followers/unfollow', {
          data: {
            follower: user.value.username,
            followee: username
          }
        })

        // Remove user from following list
        following.value = following.value.filter(user => user.username !== username)
      } catch (error) {
        console.error('Error unfollowing user:', error)
        updateError.value = error.response?.data?.message || `Failed to unfollow ${username}`
      } finally {
        unfollowingUsers.value.delete(username)
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
      selectedFile.value = null
      imagePreview.value = ''
      editing.value = true
      updateError.value = ''
      updateSuccess.value = false
    }

    const handleFileChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        selectedFile.value = file

        // Create preview
        const reader = new FileReader()
        reader.onload = (e) => {
          imagePreview.value = e.target.result
        }
        reader.readAsDataURL(file)
      } else {
        selectedFile.value = null
        imagePreview.value = ''
      }
    }

    const cancelEditing = () => {
      editing.value = false
      selectedFile.value = null
      imagePreview.value = ''
      updateError.value = ''
      updateSuccess.value = false

      // Reset file input
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }

    const saveProfile = async () => {
      updating.value = true
      updateError.value = ''
      updateSuccess.value = false

      try {
        await userStore.updateUserProfile(editForm.value, selectedFile.value)
        await fetchProfile() // Refresh the profile data
        updateSuccess.value = true
        editing.value = false
        selectedFile.value = null
        imagePreview.value = ''

        // Reset file input
        if (fileInput.value) {
          fileInput.value.value = ''
        }
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
      selectedFile,
      imagePreview,
      fileInput,
      followers,
      following,
      unfollowingUsers,
      startEditing,
      cancelEditing,
      saveProfile,
      handleFileChange,
      confirmLogout,
      unfollowUser
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

.user-badge {
  font-size: 0.8rem;
  padding: 0.4rem 0.6rem;
  cursor: default;
}

.card.bg-light {
  border: 1px solid #e9ecef;
}

.card.bg-light .card-body {
  padding: 0.75rem;
}

.card.bg-light h6 {
  margin-bottom: 0.25rem;
  font-weight: 600;
  color: #495057;
}

.clickable-card {
  cursor: pointer;
  transition: all 0.2s ease-in-out;
}

.clickable-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.clickable-card:hover .card-title {
  color: #0d6efd !important;
}

.user-avatar-small {
  width: 35px;
  height: 35px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.9rem;
}

.modal-body {
  max-height: 60vh;
  overflow-y: auto;
}

.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-thumb {
  background-color: #dee2e6;
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-track {
  background-color: #f8f9fa;
}
</style>
