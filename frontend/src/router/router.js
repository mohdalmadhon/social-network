import LoginPage from '@/views/LoginPage.vue'
import EditProfile from '@/views/profiles/EditProfile.vue'
import PersonalProfile from '@/views/profiles/PersonalProfile.vue'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/login',
        component: LoginPage
    },
    {
        path: '/me',
        component: PersonalProfile
    }, 
    {
        path: '/me/edit',
        component: EditProfile
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})
