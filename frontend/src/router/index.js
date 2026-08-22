import { createRouter, createWebHistory } from 'vue-router'

import LoginPage from '@/pages/LoginPage.vue'
import ProfilePage from '@/pages/ProfilePage.vue'

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
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(async (to) => {
    if (!to.meta.requiresAuth) {
        return true
    }

    try {
        const resp = await fetch('/api/session', {
            method: 'POST'
        })

        console.log('Session status:', resp.status)

        if (!resp.ok) {
            console.log('No valid session')
            return '/login'
        }

        const result = await resp.json()

        console.log('Session result:', result)

        if (!result.status) {
            console.log('Session invalid')
            return '/login'
        }

        return true
    } catch (err) {
        console.error('Session check failed:', err)
        return '/login'
    }
})