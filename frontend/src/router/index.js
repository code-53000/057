import { createRouter, createWebHistory } from 'vue-router'
import CustomerService from '../views/CustomerService.vue'
import Dispatcher from '../views/Dispatcher.vue'
import WorkerView from '../views/WorkerView.vue'

const routes = [
  {
    path: '/',
    redirect: '/customer-service'
  },
  {
    path: '/customer-service',
    name: 'CustomerService',
    component: CustomerService
  },
  {
    path: '/dispatcher',
    name: 'Dispatcher',
    component: Dispatcher
  },
  {
    path: '/worker',
    name: 'WorkerView',
    component: WorkerView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
