import LoginPage from "@/views/LoginPage.vue";
import PersonalProfile from "@/views/profiles/PersonalProfile.vue";
import GroupsPage from "@/views/group/GroupsPage.vue";
import GroupPage from "@/views/group/GroupPage.vue";
import HomeFeedPage from "@/views/HomeFeedPage.vue";
import ChatsPage from "@/views/ChatsPage.vue";
import NotificationsPage from "@/views/NotificationsPage.vue";
import SearchPage from "@/views/SearchPage.vue";

import { createRouter, createWebHistory } from "vue-router";
import EditProfile from "@/views/profiles/EditProfile.vue";
import UserProfile from "@/views/profiles/UserProfile.vue";

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/:pathMatch(.*)*', redirect: '/home' },
  {
    path: "/login",
    component: LoginPage,
  },
  {
    path: "/me",
    component: PersonalProfile,
    alias: "/profile",
  },
  {
    path: "/me/edit",
    component: EditProfile,
  },
  {
    path: "/user",
    component: UserProfile
  },
  {
    path: "/groups",
    component: GroupsPage,
  },
  {
    path: "/groups/:groupId",
    component: GroupPage,
  },

  {
    path: "/home",
    component: HomeFeedPage,
    alias: "/home-feed",
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/chats",
    component: ChatsPage,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/notifications",
    component: NotificationsPage,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/search",
    component: SearchPage,
    meta: {
      requiresAuth: true,
    },
  },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
