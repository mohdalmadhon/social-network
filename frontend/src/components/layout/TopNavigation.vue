<script setup>
import { onMounted, ref } from 'vue'
import OrbitLogo from './OrbitLogo.vue'
import { getNotifications } from '@/api/notifications.js'

const notificationUnreadCount = ref(0)

onMounted(async () => {
  try {
    const result = await getNotifications('all')
    notificationUnreadCount.value = result?.unreadCount || 0
  } catch {
    notificationUnreadCount.value = 0
  }
})
</script>

<template>
  <header class="top-navigation">
    <a class="brand" href="/home-feed" aria-label="Orbit home">
      <OrbitLogo />
      <span>orbit</span>
    </a>

    <form class="search" role="search" @submit.prevent>
      <label class="visually-hidden" for="orbit-search">Search Orbit</label>
      <span aria-hidden="true">⌕</span>
      <input id="orbit-search" type="search" placeholder="Search people, groups, posts..." />
    </form>

    <nav class="top-actions" aria-label="Account shortcuts">
      <a class="icon-link orbit-touch-target" href="/chats" aria-label="Messages">
        <span aria-hidden="true">◌</span>
        <span class="badge badge--message">5</span>
      </a>
      <a class="icon-link orbit-touch-target" href="/notifications" aria-label="Notifications">
        <span aria-hidden="true">♢</span>
        <span v-if="notificationUnreadCount" class="badge badge--notification">{{ notificationUnreadCount }}</span>
      </a>
      <a class="avatar orbit-touch-target" href="/profile" aria-label="My profile">N</a>
    </nav>
  </header>
</template>

<style scoped>
 @import '../../styles/global.css';
 @import '../../styles/variables.css';
.top-navigation {
  position: sticky;
  top: 0;
  z-index: 20;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  min-height: 4rem;
  padding: var(--space-2) var(--space-3);
  background: rgb(15 18 34 / 96%);
  border-bottom: 1px solid var(--color-border);
  backdrop-filter: blur(1rem);
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 700;
  text-decoration: none;
}

.search {
  display: none;
}

.top-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--space-1);
}

.icon-link,
.avatar {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  text-decoration: none;
}

.icon-link {
  font-size: 1.5rem;
}

.avatar {
  width: var(--touch-target);
  border-radius: 50%;
  background: var(--gradient-action);
  color: white;
  font-weight: 700;
}

.avatar::after {
  position: absolute;
  right: 0;
  bottom: 0.125rem;
  width: 0.55rem;
  height: 0.55rem;
  border: 2px solid var(--color-sidebar);
  border-radius: 50%;
  background: var(--color-mint);
  content: "";
}

.badge {
  position: absolute;
  top: 0.125rem;
  right: 0.125rem;
  display: grid;
  min-width: 1.1rem;
  height: 1.1rem;
  padding-inline: 0.25rem;
  place-items: center;
  border-radius: 999px;
  color: #081018;
  font-size: 0.7rem;
  font-weight: 700;
}

.badge--message {
  background: var(--color-mint);
}

.badge--notification {
  background: var(--color-coral);
}

@media (min-width: 48rem) {
  .top-navigation {
    padding-inline: var(--space-5);
  }

  .search {
    display: flex;
    align-items: center;
    justify-self: center;
    width: min(100%, 26rem);
    min-height: var(--touch-target);
    padding-inline: var(--space-4);
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 999px;
    color: var(--color-text-faint);
  }

  .search input {
    width: 100%;
    border: 0;
    outline: 0;
    background: transparent;
    color: var(--color-text);
  }

  .top-actions {
    gap: var(--space-2);
  }
}
</style>
