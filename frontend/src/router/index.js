import { createRouter, createWebHistory } from 'vue-router'

import LoginPage from '@/pages/LoginPage.vue'
import ProfilePage from '@/pages/ProfilePage.vue'
import EditProfile from '@/pages/EditProfile.vue'

const routes = [
    {
        path: '/login',
        component: LoginPage
    },
    {
        path: '/profile',
        component: ProfilePage,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: '/profile/edit',
        component: EditProfile,
        meta: {
            requiresAuth: true
        }
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})
