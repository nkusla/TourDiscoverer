import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Tours from '../views/Tours.vue'
import TourEditor from '../views/TourEditor.vue'
import Login from '../views/Login.vue'

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
    component: TourEditor
  },
  {
    path: '/tour/edit/:id',
    name: 'EditTour',
    component: TourEditor,
    props: true
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
