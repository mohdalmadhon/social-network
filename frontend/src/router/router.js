import LoginPage from "@/views/LoginPage.vue";
import EditProfile from "@/views/profiles/EditProfile.vue";
import PersonalProfile from "@/views/profiles/PersonalProfile.vue";
import GroupsPage from "@/views/GroupsPage.vue";
import HomeFeedPage from "@/views/HomeFeedPage.vue";
import ChatsPage from "@/views/ChatsPage.vue";
import NotificationsPage from "@/views/NotificationsPage.vue";

import { createRouter, createWebHistory } from "vue-router";

const routes = [
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
    path: "/groups",
    component: GroupsPage,
    meta: {
      requiresAuth: true,
    },
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
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
