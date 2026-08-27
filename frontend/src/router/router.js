import { createRouter, createWebHistory } from 'vue-router'
import Auth from '@/views/Auth.vue'

const routes = [
    {
        path: '/login',
        component: Auth
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})
