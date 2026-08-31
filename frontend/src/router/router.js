import { createRouter, createWebHistory } from 'vue-router'
import Auth from '@/views/Auth.vue'
import PersonalProfile from '@/views/profiles/PersonalProfile.vue'
import EditProfile from '@/views/profiles/EditProfile.vue'
import Profile from '@/views/profiles/Profile.vue'

const routes = [
    {
        path: '/login',
        component: Auth
    },
    {
        path: '/me',
        component: PersonalProfile
    },
    {
        path: '/me/edit',
        component: EditProfile
    }, 
    {
        path: '/user',
        component: Profile
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})
