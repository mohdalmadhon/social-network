import { createApp, h } from 'vue'
import HomeFeedPage from './pages/HomeFeedPage.vue'
import PlaceholderPage from './pages/PlaceholderPage.vue'
import './styles/variables.css'
import './styles/global.css'

const pageName = document.body.dataset.page

const pages = {
  home: HomeFeedPage,
  profile: () => h(PlaceholderPage, { activePage: 'profile', title: 'Profile' }),
  groups: () => h(PlaceholderPage, { activePage: 'groups', title: 'Groups' }),
  chats: () => h(PlaceholderPage, { activePage: 'chats', title: 'Chats' }),
  notifications: () => h(PlaceholderPage, { activePage: 'notifications', title: 'Notifications' }),
}

const page = pages[pageName]

if (!page) {
  throw new Error(`Unknown Orbit page: ${pageName}`)
}

createApp(page).mount('#app')
