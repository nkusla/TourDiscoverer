<template>
  <div class="login">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-6 col-lg-4">
          <div class="card shadow">
            <div class="card-body p-4">
              <div class="text-center mb-4">
                <h3>Login</h3>
                <p class="text-muted">Sign in to your account</p>
              </div>
              
              <form @submit.prevent="handleLogin">
                <div class="mb-3">
                  <label class="form-label">Username</label>
                  <input 
                    v-model="loginForm.username" 
                    type="text" 
                    class="form-control"
                    :class="{ 'is-invalid': errors.username }"
                    required
                    placeholder="Enter your username"
                  />
                  <div v-if="errors.username" class="invalid-feedback">
                    {{ errors.username }}
                  </div>
                </div>
                
                <div class="mb-3">
                  <label class="form-label">Password</label>
                  <input 
                    v-model="loginForm.password" 
                    type="password" 
                    class="form-control"
                    :class="{ 'is-invalid': errors.password }"
                    required
                    placeholder="Enter your password"
                  />
                  <div v-if="errors.password" class="invalid-feedback">
                    {{ errors.password }}
                  </div>
                </div>
                
                <div v-if="successMessage" class="alert alert-success" role="alert">
                  <i class="fas fa-check-circle me-2"></i>{{ successMessage }}
                </div>
                
                <div v-if="errors.general" class="alert alert-danger" role="alert">
                  {{ errors.general }}
                </div>
                
                <div class="d-grid">
                  <button 
                    type="submit" 
                    class="btn btn-primary"
                    :disabled="loading"
                  >
                    <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                    {{ loading ? 'Signing in...' : 'Sign In' }}
                  </button>
                </div>
              </form>
              
              <div class="text-center mt-3">
                <p class="text-muted">
                  Don't have an account? 
                  <a href="#" @click.prevent="showRegister = true">Register here</a>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Register Modal -->
    <div 
      class="modal fade" 
      id="registerModal" 
      tabindex="-1"
      ref="registerModal"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Create Account</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleRegister">
              <div class="mb-3">
                <label class="form-label">Username</label>
                <input 
                  v-model="registerForm.username" 
                  type="text" 
                  class="form-control"
                  :class="{ 'is-invalid': registerErrors.username }"
                  required
                  placeholder="Choose a username"
                />
                <div v-if="registerErrors.username" class="invalid-feedback">
                  {{ registerErrors.username }}
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Email</label>
                <input 
                  v-model="registerForm.email" 
                  type="email" 
                  class="form-control"
                  :class="{ 'is-invalid': registerErrors.email }"
                  required
                  placeholder="Enter your email"
                />
                <div v-if="registerErrors.email" class="invalid-feedback">
                  {{ registerErrors.email }}
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Password</label>
                <input 
                  v-model="registerForm.password" 
                  type="password" 
                  class="form-control"
                  :class="{ 'is-invalid': registerErrors.password }"
                  required
                  placeholder="Choose a password"
                />
                <div v-if="registerErrors.password" class="invalid-feedback">
                  {{ registerErrors.password }}
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Confirm Password</label>
                <input 
                  v-model="registerForm.confirmPassword" 
                  type="password" 
                  class="form-control"
                  :class="{ 'is-invalid': registerErrors.confirmPassword }"
                  required
                  placeholder="Confirm your password"
                />
                <div v-if="registerErrors.confirmPassword" class="invalid-feedback">
                  {{ registerErrors.confirmPassword }}
                </div>
              </div>
              
              <div class="mb-3">
                <label class="form-label">Role</label>
                <select 
                  v-model="registerForm.role" 
                  class="form-select"
                  :class="{ 'is-invalid': registerErrors.role }"
                  required
                >
                  <option value="">Select your role</option>
                  <option value="guide">Guide</option>
                  <option value="tourist">Tourist</option>
                </select>
                <div v-if="registerErrors.role" class="invalid-feedback">
                  {{ registerErrors.role }}
                </div>
              </div>
              
              <div v-if="registerErrors.general" class="alert alert-danger" role="alert">
                {{ registerErrors.general }}
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
              Cancel
            </button>
            <button 
              type="button" 
              class="btn btn-primary"
              @click="handleRegister"
              :disabled="registerLoading"
            >
              <span v-if="registerLoading" class="spinner-border spinner-border-sm me-2"></span>
              {{ registerLoading ? 'Creating...' : 'Create Account' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Modal } from 'bootstrap'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const registerModal = ref(null)
    const registerModalInstance = ref(null)
    const showRegister = ref(false)
    
    const loading = ref(false)
    const registerLoading = ref(false)
    const successMessage = ref('')
    
    const loginForm = ref({
      username: '',
      password: ''
    })
    
    const registerForm = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      role: ''
    })
    
    const errors = ref({})
    const registerErrors = ref({})
    
    onMounted(() => {
      // Initialize Bootstrap modal
      if (registerModal.value) {
        registerModalInstance.value = new Modal(registerModal.value)
      }
    })
    
    watch(showRegister, (newValue) => {
      if (newValue && registerModalInstance.value) {
        registerModalInstance.value.show()
        showRegister.value = false
      }
    })
    
    const validateLoginForm = () => {
      errors.value = {}
      
      if (!loginForm.value.username.trim()) {
        errors.value.username = 'Username is required'
      }
      
      if (!loginForm.value.password) {
        errors.value.password = 'Password is required'
      }
      
      return Object.keys(errors.value).length === 0
    }
    
    const validateRegisterForm = () => {
      registerErrors.value = {}
      
      if (!registerForm.value.username.trim()) {
        registerErrors.value.username = 'Username is required'
      } else if (registerForm.value.username.length < 3) {
        registerErrors.value.username = 'Username must be at least 3 characters'
      }
      
      if (!registerForm.value.email.trim()) {
        registerErrors.value.email = 'Email is required'
      } else if (!/\S+@\S+\.\S+/.test(registerForm.value.email)) {
        registerErrors.value.email = 'Please enter a valid email'
      }
      
      if (!registerForm.value.password) {
        registerErrors.value.password = 'Password is required'
      } else if (registerForm.value.password.length < 6) {
        registerErrors.value.password = 'Password must be at least 6 characters'
      }
      
      if (registerForm.value.password !== registerForm.value.confirmPassword) {
        registerErrors.value.confirmPassword = 'Passwords do not match'
      }
      
      if (!registerForm.value.role) {
        registerErrors.value.role = 'Please select a role'
      }
      
      return Object.keys(registerErrors.value).length === 0
    }
    
    const handleLogin = async () => {
      if (!validateLoginForm()) return
      
      loading.value = true
      errors.value = {}
      successMessage.value = ''
      
      try {
        await userStore.login(loginForm.value)
        successMessage.value = 'Login successful! Redirecting...'
        setTimeout(() => {
          router.push('/')
        }, 1000)
      } catch (error) {
        errors.value.general = error.message || 'Login failed. Please try again.'
      } finally {
        loading.value = false
      }
    }
    
    const handleRegister = async () => {
      if (!validateRegisterForm()) return
      
      registerLoading.value = true
      registerErrors.value = {}
      
      try {
        console.log('userStore:', userStore)
        console.log('userStore.register:', userStore.register)
        
        await userStore.register({
          username: registerForm.value.username,
          email: registerForm.value.email,
          password: registerForm.value.password,
          role: registerForm.value.role
        })
        
        // Close modal
        registerModalInstance.value?.hide()
        
        // Registration automatically logs the user in via the store
        router.push('/')
      } catch (error) {
        console.error('Registration error:', error)
        registerErrors.value.general = error.message || 'Registration failed. Please try again.'
      } finally {
        registerLoading.value = false
      }
    }
    
    return {
      registerModal,
      showRegister,
      loading,
      registerLoading,
      loginForm,
      registerForm,
      errors,
      registerErrors,
      successMessage,
      handleLogin,
      handleRegister
    }
  }
}
</script>

<style scoped>
.login {
  min-height: calc(100vh - 56px);
  display: flex;
  align-items: center;
  background-color: #f8f9fa;
}

.card {
  border: none;
  border-radius: 10px;
}

.btn-primary {
  border-radius: 20px;
}

.form-control {
  border-radius: 8px;
}

.form-control:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}
</style>
