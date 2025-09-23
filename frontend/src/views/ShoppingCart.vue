<template>
  <div class="shopping-cart">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-md-8">
          <div class="card shadow">
            <div class="card-header bg-primary text-white">
              <h4 class="mb-0">
                <i class="fas fa-shopping-cart me-2"></i>Shopping Cart
                <span v-if="cartStore.itemCount > 0" class="badge bg-secondary ms-2">
                  {{ cartStore.itemCount }}
                </span>
              </h4>
            </div>
            <div class="card-body">
              <!-- Loading State -->
              <div v-if="cartStore.loading" class="text-center py-4">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-2">Loading cart...</p>
              </div>
              
              <!-- Empty Cart -->
              <div v-else-if="cartStore.itemCount === 0" class="text-center py-5">
                <i class="fas fa-shopping-cart fa-3x text-muted mb-3"></i>
                <h5 class="text-muted">Your cart is empty</h5>
                <p class="text-muted mb-4">Browse tours and add them to your cart to get started!</p>
                <router-link to="/tours" class="btn btn-primary">
                  <i class="fas fa-search me-2"></i>Browse Tours
                </router-link>
              </div>
              
              <!-- Cart Items -->
              <div v-else>
                <div class="row">
                  <div 
                    v-for="item in cartStore.cartItems" 
                    :key="item.id"
                    class="col-12 mb-3"
                  >
                    <div class="card">
                      <div class="card-body">
                        <div class="row align-items-center">
                          <div class="col-md-8">
                            <h5 class="card-title mb-1">{{ item.tour_name }}</h5>
                            <p class="text-muted mb-0">
                              <small>Tour ID: {{ item.tour_id }}</small>
                            </p>
                          </div>
                          <div class="col-md-2 text-center">
                            <h5 class="text-success mb-0">${{ item.price.toFixed(2) }}</h5>
                          </div>
                          <div class="col-md-2 text-end">
                            <button 
                              class="btn btn-outline-danger btn-sm"
                              @click="removeItem(item.tour_id, item.tour_name)"
                              :disabled="removing === item.tour_id"
                            >
                              <span v-if="removing === item.tour_id" class="spinner-border spinner-border-sm me-1"></span>
                              <i v-else class="fas fa-trash me-1"></i>
                              Remove
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- Cart Summary -->
                <div class="card bg-light mt-4">
                  <div class="card-body">
                    <div class="row align-items-center">
                      <div class="col-md-6">
                        <h5 class="mb-0">
                          Total: <span class="text-success">${{ cartStore.total.toFixed(2) }}</span>
                        </h5>
                        <small class="text-muted">{{ cartStore.itemCount }} item(s) in cart</small>
                      </div>
                      <div class="col-md-6 text-end">
                        <button 
                          class="btn btn-outline-secondary me-2"
                          @click="clearCartConfirm"
                          :disabled="clearing"
                        >
                          <span v-if="clearing" class="spinner-border spinner-border-sm me-1"></span>
                          <i v-else class="fas fa-trash me-1"></i>
                          Clear Cart
                        </button>
                        <button 
                          class="btn btn-success"
                          @click="proceedToCheckout"
                          :disabled="checkingOut || cartStore.itemCount === 0"
                        >
                          <span v-if="checkingOut" class="spinner-border spinner-border-sm me-1"></span>
                          <i v-else class="fas fa-credit-card me-1"></i>
                          Checkout
                        </button>
                      </div>
                    </div>
                  </div>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { usePurchaseStore } from '../stores/purchase'

export default {
  name: 'ShoppingCart',
  setup() {
    const router = useRouter()
    const cartStore = useCartStore()
    const purchaseStore = usePurchaseStore()
    
    const removing = ref(null)
    const clearing = ref(false)
    const checkingOut = ref(false)
    const error = ref('')
    const success = ref('')
    
    onMounted(async () => {
      await cartStore.fetchCart()
    })
    
    const removeItem = async (tourId, tourName) => {
      if (confirm(`Are you sure you want to remove "${tourName}" from your cart?`)) {
        removing.value = tourId
        error.value = ''
        
        try {
          await cartStore.removeFromCart(tourId)
          success.value = 'Tour removed from cart successfully!'
          setTimeout(() => { success.value = '' }, 3000)
        } catch (err) {
          error.value = err.message
        } finally {
          removing.value = null
        }
      }
    }
    
    const clearCartConfirm = async () => {
      if (confirm('Are you sure you want to clear your entire cart?')) {
        clearing.value = true
        error.value = ''
        
        try {
          await cartStore.clearCart()
          success.value = 'Cart cleared successfully!'
          setTimeout(() => { success.value = '' }, 3000)
        } catch (err) {
          error.value = err.message
        } finally {
          clearing.value = false
        }
      }
    }
    
    const proceedToCheckout = async () => {
      checkingOut.value = true
      error.value = ''
      
      try {
        const tokens = await cartStore.checkout()
        success.value = `Checkout successful! You've purchased ${tokens.length} tour(s).`
        
        // Refresh purchased tours
        await purchaseStore.fetchPurchasedTours()
        
        // Redirect to purchased tours after a delay
        setTimeout(() => {
          router.push('/purchases')
        }, 2000)
      } catch (err) {
        error.value = err.message
      } finally {
        checkingOut.value = false
      }
    }
    
    return {
      cartStore,
      removing,
      clearing,
      checkingOut,
      error,
      success,
      removeItem,
      clearCartConfirm,
      proceedToCheckout
    }
  }
}
</script>

<style scoped>
.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-1px);
}

.badge {
  font-size: 0.75em;
}
</style>