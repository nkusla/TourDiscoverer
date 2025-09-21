import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import Home from '../views/Home.vue'
import Tours from '../views/Tours.vue'
import TourEditor from '../views/TourEditor.vue'
import Login from '../views/Login.vue'
import Profile from '../views/Profile.vue'
import Users from '../views/Users.vue'
import BlogList from '../views/BlogList.vue'
import CreateBlog from '../views/CreateBlog.vue'
import PositionSimulator from '../views/PositionSimulator.vue'
import TourExecution from '../views/TourExecution.vue'
import Recommendations from '../views/Recommendations.vue'

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
    meta: { requiresAuth: true, requiresGuideOrAdmin: true }
  },
  {
    path: '/tour/edit/:id',
    name: 'EditTour',
    component: TourEditor,
    props: true,
    meta: { requiresAuth: true, requiresGuideOrAdmin: true }
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
  },
  {
    path: '/recommendations',
    name: 'Recommendations',
    component: Recommendations,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/blogs',
    name: 'BlogList',
    component: BlogList,
    meta: { requiresAuth: true }
  },
  {
    path: '/blog/create',
    name: 'CreateBlog',
    component: CreateBlog,
    meta: { requiresAuth: true }
  },
  {
    path: '/position-simulator',
    name: 'PositionSimulator',
    component: PositionSimulator,
    meta: { requiresAuth: true, requiresTourist: true }
  },
  {
    path: '/tour-execution',
    name: 'TourExecution',
    component: TourExecution,
    meta: { requiresAuth: true, requiresTourist: true }
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

  // Check if route requires admin privileges
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/')
    return
  }

  // Check if route requires guide or admin privileges
  if (to.meta.requiresGuideOrAdmin && !userStore.canCreateTours) {
    next('/')
    return
  }

  // Check if route requires guest (like login page)
  if (to.meta.requiresGuest && userStore.isAuthenticated) {
    next('/')
    return
  }

  // Check if route requires tourist privileges
  if (to.meta.requiresTourist && userStore.user?.role !== 'tourist') {
    next('/')
    return
  }

  next()
})

export default router
