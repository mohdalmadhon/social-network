<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import { useNotifications } from '@/helpers/useNotifications.js'
import {
  applyNotificationAction,
  getNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from '@/api/notifications.js'

const filters = [
  { id: 'all', label: 'All' },
  { id: 'requests', label: 'Requests' },
  { id: 'groups', label: 'Groups' },
  { id: 'events', label: 'Events' },
]

const notificationItems = ref([])
const isLoading = ref(true)
const loadError = ref('')

const activeFilter = ref('all')
const { items: liveItems, refreshNotifications } = useNotifications()
watch(liveItems, (items) => { notificationItems.value = items.map(notificationForDisplay) })

const visibleNotifications = computed(() => {
  if (activeFilter.value === 'all') return notificationItems.value
  return notificationItems.value.filter((item) => item.type === activeFilter.value)
})

const unreadCount = computed(() => notificationItems.value.filter((item) => item.unread).length)

function formatNotificationTime(value) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Recently'

  return date.toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' })
}

function actionLabel(action) {
  return {
    accept: 'Accepted',
    decline: 'Declined',
    reject: 'Rejected',
    join: 'Joined',
    rsvp: 'Going',
  }[action] || action
}

function notificationForDisplay(notification) {
  const categoryStyles = {
    requests: { icon: 'R', color: '#9b7cff' },
    groups: { icon: 'G', color: '#45d9d0' },
    events: { icon: 'E', color: '#ffb84d' },
  }

  const style =
    categoryStyles[notification.category] || categoryStyles.groups

  let action = ''

  if (notification.category === 'requests' && notification.type === 'follow_request') {
    action = notification.followStatus === 0 ? 'follow' : notification.followStatus === 1 ? 'accept' : ''
  }

  if (notification.category === 'groups' && notification.type === 'join_request') {
    if (notification.requestStatus === 'accepted') {
      action = 'accept'
    } else if (notification.requestStatus === 'rejected') {
      action = 'reject'
    } else {
      action = 'join_request'
    }
  }

  if (notification.category === 'groups' && notification.type === 'invitation') {
    if (notification.invitationStatus === 'accepted') {
      action = 'join'
    } else if (notification.invitationStatus === 'declined') {
      action = 'decline'
    } else {
      action = 'invitation'
    }
  }

  if (notification.category === 'events' && notification.type === 'event_created') {
    action = notification.eventResponse === 'going' ? 'going' : notification.eventResponse === 'declined' ? 'decline' : 'rsvp'
  }

  return {
    ...notification,
    type: notification.category,
    icon: style.icon,
    color: style.color,
    title: notification.message,
    detail: notification.type.replaceAll('_', ' '),
    time: formatNotificationTime(notification.createdAt),
    unread: !notification.isRead,
    action,
  }
}

async function loadNotifications() {
  isLoading.value = true
  loadError.value = ''

  try {
    const result = await getNotifications('all')
    notificationItems.value = (result?.notifications || []).map(notificationForDisplay)
  } catch (error) {
    loadError.value = error.message || 'Could not load notifications.'
  } finally {
    isLoading.value = false
  }
}

async function markAsRead(item) {
  if (!item.unread) return

  try {
    await markNotificationRead(item.id)
    item.unread = false
    await refreshNotifications()
  } catch (error) {
    loadError.value = error.message || 'Could not mark notification as read.'
  }
}

async function markAllAsRead() {
  try {
    await markAllNotificationsRead()
    await refreshNotifications()
    notificationItems.value.forEach((item) => {
      item.unread = false
    })
  } catch (error) {
    loadError.value = error.message || 'Could not mark notifications as read.'
  }
}

async function chooseAction(item, action) {
  if (item.busy) return
  item.busy = true
  const previousAction = item.action

  try {
    await applyNotificationAction(item.id, action)
    item.action = action
    item.unread = false
    await refreshNotifications()
  } catch (error) {
    item.action = previousAction
    loadError.value = error.message || 'Could not complete notification action.'
  } finally {
    item.busy = false
  }
}

onMounted(loadNotifications)
</script>

<template>
  <AuthenticatedLayout active-page="notifications">
    <section class="notifications-page orbit-surface" aria-labelledby="notifications-title">
      <header class="notifications-page__header">
        <div>
          <p class="orbit-meta">Stay in the loop</p>

          <h1 id="notifications-title">
            Notifications
          </h1>

          <p class="notifications-page__summary">
            {{ unreadCount ? `${unreadCount} unread updates` : 'You are all caught up' }}
          </p>
        </div>

        <button class="mark-all-button" type="button" :disabled="unreadCount === 0" @click="markAllAsRead">
          Mark all as read
        </button>
      </header>

      <nav class="notification-filters" aria-label="Notification filters">
        <button v-for="filter in filters" :key="filter.id" type="button"
          :class="{ 'notification-filter--active': activeFilter === filter.id }" @click="activeFilter = filter.id">
          {{ filter.label }}
        </button>
      </nav>

      <p v-if="isLoading" class="notifications-state">
        Loading notifications...
      </p>

      <div v-else-if="loadError" class="notifications-state notifications-state--error">
        <span>{{ loadError }}</span>

        <button type="button" @click="loadNotifications">
          Try again
        </button>
      </div>

      <div v-else-if="visibleNotifications.length" class="notification-list">
        <article v-for="item in visibleNotifications" :key="item.id" class="notification-item"
          :class="{ 'notification-item--unread': item.unread }" @click="markAsRead(item)">
          <div class="notification-item__icon" :style="{ background: item.color }" aria-hidden="true">
            {{ item.icon }}
          </div>

          <div class="notification-item__body">
            <h2>{{ item.title }}</h2>
            <p>{{ item.detail }}</p>
            <time>{{ item.time }}</time>
          </div>

          <!-- Follow request -->
          <div v-if="item.action === 'follow'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary"
              @click.stop="chooseAction(item, 'accept')">
              Accept
            </button>

            <button type="button" class="action-button" @click.stop="chooseAction(item, 'decline')">
              Decline
            </button>
          </div>

          <!-- Someone requested to join my group -->
          <div v-else-if="item.action === 'join_request'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary"
              @click.stop="chooseAction(item, 'accept')">
              Accept
            </button>

            <button type="button" class="action-button" @click.stop="chooseAction(item, 'reject')">
              Reject
            </button>
          </div>

          <!-- I received a group invitation -->
          <div v-else-if="item.action === 'invitation'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary" @click.stop="chooseAction(item, 'join')">
              Join
            </button>

            <button type="button" class="action-button" @click.stop="chooseAction(item, 'decline')">
              Decline
            </button>
          </div>

          <!-- Event -->
          <div v-else-if="item.action === 'rsvp'" class="notification-item__actions">
            <button type="button" class="action-button action-button--primary" @click.stop="chooseAction(item, 'rsvp')">
              RSVP
            </button>

            <button type="button" class="action-button" @click.stop="chooseAction(item, 'decline')">
              Decline
            </button>
          </div>

          <span v-else-if="item.action" class="notification-item__result" :class="{
            'notification-item__result--negative':
              item.action === 'reject' || item.action === 'decline'
          }">
            {{ actionLabel(item.action) }}
          </span>
        </article>
      </div>

      <p v-else class="notifications-empty">
        Nothing here yet.
      </p>
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

.notifications-state {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin-top: var(--space-6);
  color: var(--color-text-muted);
}

.notifications-state--error {
  color: var(--color-coral);
}

.notifications-state button {
  min-height: var(--touch-target);
  padding-inline: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: transparent;
  color: inherit;
  cursor: pointer;
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

.notification-item__result--negative {
  color: var(--color-coral);
}
</style>
