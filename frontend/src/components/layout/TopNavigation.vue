<script setup>
import OrbitLogo from './OrbitLogo.vue'
import { useNotifications } from '@/helpers/useNotifications.js'
import { useChatCount } from '@/helpers/useChats.js'
import { logout } from '@/api/auth/auth.js'
import { ref } from 'vue'
import { router } from '@/router/router.js'
import IconGlyph from './IconGlyph.vue'
const logoutError = ref('')
const searchText = ref('')

function submitSearch() {
  const search = searchText.value.trim()
  router.push(search ? { path: '/search', query: { q: search } } : '/search')
}

async function signOut() {
  try { await logout() } catch { logoutError.value = 'Could not log out. Please try again.' }
}
const { unreadCount: notificationUnreadCount } = useNotifications()
const { chatCount } = useChatCount()
</script>

<template>
  <header class="top-navigation">
    <a class="brand" href="/home-feed" aria-label="Orbit home">
      <OrbitLogo />
      <span>orbit</span>
    </a>

    <form class="search" role="search" @submit.prevent="submitSearch">
      <IconGlyph name="search" :size="16" />
      <input v-model="searchText" type="search" placeholder="Search people, groups, posts…" aria-label="Search" />
    </form>

    <nav class="top-actions" aria-label="Account shortcuts">
      <a class="icon-link orbit-touch-target" href="/chats" aria-label="Messages" :title="`${chatCount} active chat${chatCount === 1 ? '' : 's'}`">
        <IconGlyph name="chat" :size="19" />
        <span v-if="chatCount" class="badge badge--message">{{ chatCount }}</span>
      </a>
      <a class="icon-link orbit-touch-target" href="/notifications" aria-label="Notifications">
        <IconGlyph name="bell" :size="19" />
        <span v-if="notificationUnreadCount" class="badge badge--notification">{{ notificationUnreadCount }}</span>
      </a>
      <a class="avatar orbit-touch-target" href="/profile" aria-label="My profile"><IconGlyph name="profile" :size="18" stroke-width="2" /></a>
      <button class="sign-out" type="button" @click="signOut">Log out</button>
    </nav>
    <p v-if="logoutError" role="alert">{{ logoutError }}</p>
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
  grid-template-columns: 1fr auto;
  align-items: center;
  min-height: 4rem;
  padding: var(--space-2) var(--space-3);
  background: rgb(11 13 23 / 92%);
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
  align-items: center;
  gap: var(--space-2);
  justify-self: center;
  width: min(100%, 28.5rem);
  min-height: var(--touch-target);
  padding-inline: var(--space-4);
  background: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  color: var(--color-text-faint);
}

.top-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--space-1);
}
.sign-out { min-height: 44px; border: 0; background: transparent; color: var(--color-text-muted); cursor: pointer; font-size: .8125rem; }
.sign-out:hover { color: var(--color-coral); }

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
  width: var(--touch-target);
}

.avatar {
  width: var(--touch-target);
  border-radius: 50%;
  background: var(--color-violet);
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
    grid-template-columns: minmax(0, 1fr) minmax(18rem, 28.5rem) minmax(0, 1fr);
    padding-inline: var(--space-5);
  }

  .search {
    display: flex;
  }

  .search input {
    width: 100%;
    border: 0;
    outline: 0;
    background: transparent;
    color: var(--color-text);
  }

  .search input::placeholder {
    color: var(--color-text-faint);
  }

  .search:focus-within {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .top-actions {
    justify-self: end;
    gap: var(--space-2);
  }
}
</style>
