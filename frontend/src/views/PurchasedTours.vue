<template>
  <div class="purchased-tours">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-md-10">
          <div class="card shadow">
            <div class="card-header bg-success text-white">
              <h4 class="mb-0">
                <i class="fas fa-check-circle me-2"></i>My Purchased Tours
              </h4>
            </div>
            <div class="card-body">
              <!-- Loading State -->
              <div v-if="purchaseStore.loading" class="text-center py-4">
                <div class="spinner-border text-success" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-2">Loading purchased tours...</p>
              </div>
              
              <!-- No Purchases -->
              <div v-else-if="purchaseStore.purchasedTokens.length === 0" class="text-center py-5">
                <i class="fas fa-ticket-alt fa-3x text-muted mb-3"></i>
                <h5 class="text-muted">No purchased tours yet</h5>
                <p class="text-muted mb-4">Purchase tours to access their full content and key points!</p>
                <router-link to="/tours" class="btn btn-primary">
                  <i class="fas fa-search me-2"></i>Browse Tours
                </router-link>
              </div>
              
              <!-- Purchased Tours List -->
              <div v-else>
                <div class="row">
                  <div 
                    v-for="token in purchaseStore.purchasedTokens" 
                    :key="token.id"
                    class="col-lg-6 col-md-12 mb-4"
                  >
                    <div class="card h-100" :class="getTokenStatusClass(token)">
                      <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                          <h5 class="card-title">{{ token.tour_name }}</h5>
                          <span 
                            class="badge"
                            :class="getStatusBadgeClass(token.status)"
                          >
                            {{ token.status }}
                          </span>
                        </div>
                        
                        <div class="mb-3">
                          <small class="text-muted">
                            <strong>Purchase Token:</strong> 
                            <code class="text-break">{{ token.token.substring(0, 8) }}...</code>
                          </small>
                          <br>
                          <small class="text-muted">
                            <strong>Purchased:</strong> {{ formatDate(token.created_at) }}
                          </small>
                          <br>
                          <small class="text-muted">
                            <strong>Expires:</strong> {{ formatDate(token.expires_at) }}
                          </small>
                        </div>
                        
                        <div class="d-grid gap-2">
                          <router-link 
                            :to="`/tours?view=${token.tour_id}&token=${token.token}`"
                            class="btn btn-primary"
                            :class="{ disabled: token.status !== 'active' }"
                          >
                            <i class="fas fa-map-marked-alt me-2"></i>
                            {{ token.status === 'active' ? 'View Full Tour' : 'Token Expired' }}
                          </router-link>
                          
                          <button 
                            class="btn btn-outline-secondary btn-sm"
                            @click="copyToken(token.token)"
                          >
                            <i class="fas fa-copy me-1"></i>Copy Token
                          </button>
                        </div>
                      </div>
                      
                      <div class="card-footer bg-transparent">
                        <small class="text-muted">
                          Tour ID: {{ token.tour_id }}
                        </small>
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- Summary -->
                <div class="alert alert-info mt-4">
                  <h6 class="alert-heading">
                    <i class="fas fa-info-circle me-2"></i>Purchase Summary
                  </h6>
                  <p class="mb-0">
                    You have purchased <strong>{{ purchaseStore.purchasedTokens.length }}</strong> tour(s).
                    Active tokens: <strong>{{ activeTokensCount }}</strong>
                  </p>
                </div>
              </div>
              
              <!-- Error Messages -->
              <div v-if="error" class="alert alert-danger mt-3" role="alert">
                {{ error }}
              </div>
              
              <!-- Success Messages -->
              <div v-if="success" class="alert alert-success mt-3" role="alert">
                {{ success }}
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
import { usePurchaseStore } from '../stores/purchase'

export default {
  name: 'PurchasedTours',
  setup() {
    const purchaseStore = usePurchaseStore()
    
    const error = ref('')
    const success = ref('')
    
    const activeTokensCount = computed(() => {
      return purchaseStore.purchasedTokens.filter(token => token.status === 'active').length
    })
    
    onMounted(async () => {
      try {
        await purchaseStore.fetchPurchasedTours()
      } catch (err) {
        error.value = err.message
      }
    })
    
    const formatDate = (dateString) => {
      if (!dateString) return ''
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    
    const getTokenStatusClass = (token) => {
      switch (token.status) {
        case 'active':
          return 'border-success'
        case 'expired':
          return 'border-warning'
        default:
          return 'border-secondary'
      }
    }
    
    const getStatusBadgeClass = (status) => {
      switch (status) {
        case 'active':
          return 'bg-success'
        case 'expired':
          return 'bg-warning'
        default:
          return 'bg-secondary'
      }
    }
    
    const copyToken = async (token) => {
      try {
        await navigator.clipboard.writeText(token)
        success.value = 'Token copied to clipboard!'
        setTimeout(() => { success.value = '' }, 3000)
      } catch (err) {
        error.value = 'Failed to copy token to clipboard'
      }
    }
    
    return {
      purchaseStore,
      error,
      success,
      activeTokensCount,
      formatDate,
      getTokenStatusClass,
      getStatusBadgeClass,
      copyToken
    }
  }
}
</script>

<style scoped>
.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-2px);
}

.badge {
  font-size: 0.75em;
}

code {
  font-size: 0.85em;
  background-color: #f8f9fa;
  padding: 0.25rem;
  border-radius: 0.25rem;
}
</style>