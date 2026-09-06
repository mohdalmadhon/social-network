<script setup>
import { computed, ref } from 'vue'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'

const filters = [
  { id: 'all', label: 'All' },
  { id: 'requests', label: 'Requests' },
  { id: 'groups', label: 'Groups' },
  { id: 'events', label: 'Events' },
]

const notificationItems = ref([
  {
    id: 1,
    type: 'requests',
    icon: 'S',
    color: '#9b7cff',
    title: 'Sana Iqbal wants to follow you',
    detail: 'Your profile is private.',
    time: '12 min ago',
    unread: true,
    action: 'follow',
  },
  {
    id: 2,
    type: 'groups',
    icon: 'W',
    color: '#45d9d0',
    title: 'You were invited to Weekend Hikers',
    detail: 'Mara Voss invited you to join this group.',
    time: '1h ago',
    unread: true,
    action: 'join',
  },
  {
    id: 3,
    type: 'events',
    icon: 'E',
    color: '#ffb84d',
    title: 'Design meetup starts tomorrow',
    detail: 'You RSVP’d as going.',
    time: '3h ago',
    unread: false,
    action: 'rsvp',
  },
  {
    id: 4,
    type: 'groups',
    icon: 'K',
    color: '#ff8b5c',
    title: 'Kiko Tanaka commented on your post',
    detail: '“This view is incredible.”',
    time: 'Yesterday',
    unread: false,
    action: '',
  },
])

const activeFilter = ref('all')

const visibleNotifications = computed(() => {
  if (activeFilter.value === 'all') return notificationItems.value
  return notificationItems.value.filter((item) => item.type === activeFilter.value)
})

const unreadCount = computed(() => notificationItems.value.filter((item) => item.unread).length)

function markAsRead(item) {
  item.unread = false
}

function markAllAsRead() {
  notificationItems.value.forEach((item) => {
    item.unread = false
  })
}

function chooseAction(item, action) {
  item.action = action
  item.unread = false
}
</script>

<template>
  <AuthenticatedLayout active-page="notifications">
    <section class="notifications-page orbit-surface" aria-labelledby="notifications-title">
      <header class="notifications-page__header">
        <div>
          <p class="orbit-meta">Stay in the loop</p>
          <h1 id="notifications-title">Notifications</h1>
          <p class="notifications-page__summary">
            {{ unreadCount ? `${unreadCount} unread updates` : 'You are all caught up' }}
          </p>
        </div>
        <button
          class="mark-all-button"
          type="button"
          :disabled="unreadCount === 0"
          @click="markAllAsRead"
        >
          Mark all as read
        </button>
      </header>

      <nav class="notification-filters" aria-label="Notification filters">
        <button
          v-for="filter in filters"
          :key="filter.id"
          type="button"
          :class="{ 'notification-filter--active': activeFilter === filter.id }"
          @click="activeFilter = filter.id"
        >
          {{ filter.label }}
        </button>
      </nav>

      <div v-if="visibleNotifications.length" class="notification-list">
        <article
          v-for="item in visibleNotifications"
          :key="item.id"
          class="notification-item"
          :class="{ 'notification-item--unread': item.unread }"
          @click="markAsRead(item)"
        >
          <div class="notification-item__icon" :style="{ background: item.color }" aria-hidden="true">
            {{ item.icon }}
          </div>

          <div class="notification-item__body">
            <h2>{{ item.title }}</h2>
            <p>{{ item.detail }}</p>
            <time>{{ item.time }}</time>
          </div>

          <div v-if="item.action === 'follow'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary" @click.stop="chooseAction(item, 'accepted')">
              Accept
            </button>
            <button type="button" class="action-button" @click.stop="chooseAction(item, 'declined')">
              Decline
            </button>
          </div>
          <div v-else-if="item.action === 'join'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary" @click.stop="chooseAction(item, 'joined')">
              Join
            </button>
          </div>
          <span v-else-if="item.action" class="notification-item__result">{{ item.action }}</span>
        </article>
      </div>

      <p v-else class="notifications-empty">Nothing here yet.</p>
    </section>
  </AuthenticatedLayout>
</template>

<style scoped>
.notifications-page {
  width: 100%;
  max-width: 60rem;
  margin: 0 auto;
  padding: var(--space-4);
}

.notifications-page__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
}

.notifications-page h1,
.notifications-page__summary,
.notification-item h2,
.notification-item p,
.notification-item time {
  margin: 0;
}

.notifications-page h1 {
  margin-top: var(--space-1);
  font-family: var(--font-display);
  font-size: clamp(1.75rem, 5vw, 2.5rem);
}

.notifications-page__summary {
  margin-top: var(--space-2);
  color: var(--color-text-muted);
}

.mark-all-button,
.notification-filter,
.action-button {
  min-height: var(--touch-target);
  border-radius: 999px;
  cursor: pointer;
}

.mark-all-button {
  padding-inline: var(--space-3);
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.mark-all-button:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.notification-filters {
  display: flex;
  gap: var(--space-2);
  margin-top: var(--space-5);
  overflow-x: auto;
  padding-bottom: var(--space-1);
}

.notification-filters button {
  min-height: var(--touch-target);
  padding-inline: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  white-space: nowrap;
}

.notification-filters button:hover,
.notification-filter--active {
  border-color: var(--color-violet) !important;
  background: rgb(124 92 255 / 16%) !important;
  color: var(--color-text) !important;
}

.notification-list {
  display: grid;
  gap: var(--space-3);
  margin-top: var(--space-5);
}

.notification-item {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr);
  align-items: start;
  gap: var(--space-3);
  padding: var(--space-3);
  border: 1px solid transparent;
  border-radius: var(--radius-medium);
  background: var(--color-input);
  cursor: pointer;
}

.notification-item--unread {
  border-color: rgb(255 107 138 / 48%);
  background: linear-gradient(90deg, rgb(255 107 138 / 9%), var(--color-input));
}

.notification-item__icon {
  display: grid;
  width: var(--touch-target);
  height: var(--touch-target);
  place-items: center;
  border-radius: 50%;
  color: #0b0d17;
  font-weight: 700;
}

.notification-item__body {
  min-width: 0;
}

.notification-item__body h2 {
  color: var(--color-text);
  font-size: 1rem;
  font-weight: 600;
}

.notification-item__body p {
  margin-top: var(--space-1);
  color: var(--color-text-muted);
  overflow-wrap: anywhere;
}

.notification-item__body time {
  display: block;
  margin-top: var(--space-2);
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.75rem;
}

.notification-item__actions,
.notification-item__result {
  grid-column: 2;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.action-button {
  padding-inline: var(--space-3);
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
}

.action-button--primary {
  border-color: var(--color-mint);
  color: var(--color-mint);
}

.notification-item__result {
  color: var(--color-mint);
  font-size: 0.875rem;
  text-transform: capitalize;
}

.notifications-empty {
  margin: var(--space-6) 0 0;
  color: var(--color-text-muted);
  text-align: center;
}

@media (min-width: 48rem) {
  .notifications-page {
    padding: var(--space-6);
  }

  .notification-item {
    grid-template-columns: var(--touch-target) minmax(0, 1fr) auto;
    align-items: center;
    padding: var(--space-4);
  }

  .notification-item__actions,
  .notification-item__result {
    grid-column: auto;
  }
}
</style>
