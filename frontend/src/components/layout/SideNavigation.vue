<script setup>
defineProps({
  activePage: {
    type: String,
    required: true,
  },
})

const links = [
  { name: 'home', label: 'Home', href: '/home-feed', icon: '⌂' },
  { name: 'profile', label: 'Profile', href: '/profile', icon: '◎' },
  { name: 'groups', label: 'Groups', href: '/groups', icon: '▱' },
  { name: 'chats', label: 'Chats', href: '/chats', icon: '◌', badge: 5, badgeType: 'message' },
  { name: 'notifications', label: 'Notifications', href: '/notifications', icon: '♢', badge: 3, badgeType: 'notification' },
]

async function logOut() {
  try {
    await fetch('/api/logout', {
      method: 'POST',
      credentials: 'include',
    })
  } finally {
    window.location.href = '/login'
  }
}
</script>

<template>
  <aside class="side-navigation">
    <nav aria-label="Main navigation">
      <a
        v-for="link in links"
        :key="link.name"
        class="navigation-link"
        :class="{ 'navigation-link--active': activePage === link.name }"
        :href="link.href"
        :aria-current="activePage === link.name ? 'page' : undefined"
      >
        <span class="navigation-link__icon" aria-hidden="true">{{ link.icon }}</span>
        <span class="navigation-link__label">{{ link.label }}</span>
        <span v-if="link.badge" class="navigation-link__badge" :class="`navigation-link__badge--${link.badgeType}`">
          {{ link.badge }}
        </span>
      </a>
    </nav>

    <button class="logout-link" type="button" @click="logOut">↪ <span>Log out</span></button>
  </aside>

  <nav class="mobile-navigation" aria-label="Mobile navigation">
    <a
      v-for="link in links"
      :key="link.name"
      class="mobile-link"
      :class="{ 'mobile-link--active': activePage === link.name }"
      :href="link.href"
      :aria-label="link.label"
      :aria-current="activePage === link.name ? 'page' : undefined"
    >
      <span aria-hidden="true">{{ link.icon }}</span>
      <span>{{ link.label }}</span>
    </a>
  </nav>
</template>

<style scoped>
.side-navigation {
  display: none;
}

.mobile-navigation {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 20;
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  min-height: 4.25rem;
  background: rgb(15 18 34 / 97%);
  border-top: 1px solid var(--color-border);
  backdrop-filter: blur(1rem);
}

.mobile-link {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 0.125rem;
  min-width: 0;
  color: var(--color-text-faint);
  font-size: 0.7rem;
  text-decoration: none;
}

.mobile-link > span:first-child {
  font-size: 1.3rem;
}

.mobile-link--active {
  color: var(--color-violet);
}

@media (min-width: 64rem) {
  .side-navigation {
    position: sticky;
    top: 4rem;
    display: flex;
    align-self: start;
    flex-direction: column;
    justify-content: space-between;
    width: 14.5rem;
    height: calc(100vh - 4rem);
    padding: var(--space-4) var(--space-3) var(--space-5);
    background: var(--color-sidebar);
    border-right: 1px solid var(--color-border);
  }

  .navigation-link {
    display: grid;
    grid-template-columns: 1.75rem 1fr auto;
    align-items: center;
    min-height: var(--touch-target);
    padding-inline: var(--space-3);
    border-radius: var(--radius-small);
    color: var(--color-text-muted);
    text-decoration: none;
  }

  .navigation-link:hover,
  .navigation-link--active {
    background: rgb(124 92 255 / 16%);
    color: var(--color-violet);
  }

  .navigation-link__icon {
    font-size: 1.25rem;
  }

  .navigation-link__badge {
    display: grid;
    min-width: 1.35rem;
    height: 1.35rem;
    place-items: center;
    border-radius: 999px;
    color: #081018;
    font-size: 0.75rem;
    font-weight: 700;
  }

  .navigation-link__badge--message {
    background: var(--color-mint);
  }

  .navigation-link__badge--notification {
    background: var(--color-coral);
  }

  .logout-link {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    min-height: var(--touch-target);
    padding-inline: var(--space-3);
    border: 0;
    background: transparent;
    color: var(--color-text-faint);
    cursor: pointer;
    font: inherit;
    text-decoration: none;
  }

  .mobile-navigation {
    display: none;
  }
}
</style>
