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
              <div v-if="!editing" class="profile-view">
                <div class="row mb-3">
                  <div class="col-sm-3">
                    <strong>Name:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ user?.firstName }} {{ user?.lastName }}
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
                    <strong>Email:</strong>
                  </div>
                  <div class="col-sm-9">
                    {{ user?.email }}
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
                
                <div class="d-flex gap-2">
                  <button class="btn btn-primary" @click="startEditing">
                    <i class="fas fa-edit me-2"></i>Edit Profile
                  </button>
                  <button class="btn btn-outline-danger" @click="confirmLogout">
                    <i class="fas fa-sign-out-alt me-2"></i>Logout
                  </button>
                </div>
              </div>
              
              <div v-else class="profile-edit">
                <form @submit.prevent="saveProfile">
                  <div class="row mb-3">
                    <div class="col-md-6">
                      <label class="form-label">First Name</label>
                      <input 
                        v-model="editForm.firstName" 
                        type="text" 
                        class="form-control"
                        required
                      />
                    </div>
                    <div class="col-md-6">
                      <label class="form-label">Last Name</label>
                      <input 
                        v-model="editForm.lastName" 
                        type="text" 
                        class="form-control"
                        required
                      />
                    </div>
                  </div>
                  
                  <div class="mb-3">
                    <label class="form-label">Email</label>
                    <input 
                      v-model="editForm.email" 
                      type="email" 
                      class="form-control"
                      required
                    />
                  </div>
                  
                  <div v-if="updateError" class="alert alert-danger">
                    {{ updateError }}
                  </div>
                  
                  <div v-if="updateSuccess" class="alert alert-success">
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
    const updateError = ref('')
    const updateSuccess = ref(false)
    
    const user = computed(() => userStore.user)
    const roleBadgeClass = computed(() => {
      return user.value?.role === 'guide' ? 'bg-success' : 'bg-info'
    })
    
    const editForm = ref({
      firstName: '',
      lastName: '',
      email: ''
    })
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
      }
    })
    
    const startEditing = () => {
      editForm.value = {
        firstName: user.value?.firstName || '',
        lastName: user.value?.lastName || '',
        email: user.value?.email || ''
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
        await userStore.updateProfile(editForm.value)
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
      updateError,
      updateSuccess,
      user,
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
