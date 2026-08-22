import { createApp } from 'vue'
import ChatsPage from './pages/ChatsPage.vue'
import GroupsPage from './pages/GroupsPage.vue'
import HomeFeedPage from './pages/HomeFeedPage.vue'
import NotificationsPage from './pages/NotificationsPage.vue'
import ProfilePage from './pages/ProfilePage.vue'
import './styles/variables.css'
import './styles/global.css'

const pageName = document.body.dataset.page

const pages = {
  home: HomeFeedPage,
  profile: ProfilePage,
  groups: GroupsPage,
  chats: ChatsPage,
  notifications: NotificationsPage,
}

const page = pages[pageName]

if (!page) {
  throw new Error(`Unknown Orbit page: ${pageName}`)
}

createApp(page).mount('#app')
