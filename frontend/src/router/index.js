import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import Home from '../views/Home.vue'
import Tours from '../views/Tours.vue'
import TourEditor from '../views/TourEditor.vue'
import Login from '../views/Login.vue'
import Profile from '../views/Profile.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/tours',
    name: 'Tours',
    component: Tours
  },
  {
    path: '/tour/create',
    name: 'CreateTour',
    component: TourEditor,
    meta: { requiresAuth: true }
  },
  {
    path: '/tour/edit/:id',
    name: 'EditTour',
    component: TourEditor,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // Check if route requires authentication
  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    next('/login')
    return
  }
  
  // Check if route requires guest (like login page)
  if (to.meta.requiresGuest && userStore.isAuthenticated) {
    next('/')
    return
  }
  
  next()
})

export default router
