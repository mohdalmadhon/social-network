import { createRouter, createWebHistory } from 'vue-router'
import Auth from '@/views/Auth.vue'
import PersonalProfile from '@/views/profiles/PersonalProfile.vue'
import EditProfile from '@/views/profiles/EditProfile.vue'
import Profile from '@/views/profiles/Profile.vue'
import AddPostPage from '@/views/posts/AddPostPage.vue'
import HomePage from '@/views/home/HomePage.vue'

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
    }, 
    {
        path: '/post/new',
        component: AddPostPage
    },
    {
        path: '/home',
        component: HomePage
    }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})
