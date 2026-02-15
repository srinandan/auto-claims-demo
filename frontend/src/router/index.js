import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import ClaimDetail from '../views/ClaimDetail.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/claims/:id',
    name: 'ClaimDetail',
    component: ClaimDetail
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
