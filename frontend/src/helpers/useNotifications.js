import { ref, onMounted, onUnmounted } from 'vue'
import { getNotifications } from '@/api/notifications.js'
import { subscribeRealtime } from '@/services/realtime.js'

const items = ref([])
const unreadCount = ref(0)
const error = ref('')
let subscribers = 0
let pending = false
let stopNotificationListener
let stopConnectionListener

function receiveNotification(event) {
  const notification = event?.notification
  if (!notification?.id) return
  const index = items.value.findIndex(item => item.id === notification.id)
  if (index >= 0) {
    items.value[index] = notification
    return
  }
  items.value.unshift(notification)
  if (!notification.isRead) unreadCount.value += 1
}

export async function refreshNotifications() {
  if (pending) return
  pending = true
  try {
    const result = await getNotifications()
    items.value = result?.notifications || []
    unreadCount.value = result?.unreadCount || 0
    error.value = ''
  } catch (failure) {
    error.value = failure.message || 'Could not load notifications.'
  } finally {
    pending = false
  }
}

// The header and sidebar share one initial HTTP load and one realtime listener.
export function useNotifications() {
  onMounted(() => {
    subscribers += 1
    if (subscribers === 1) {
      refreshNotifications()
      stopNotificationListener = subscribeRealtime('notification', receiveNotification)
      stopConnectionListener = subscribeRealtime('connection', event => {
        if (event.status === 'connected') refreshNotifications()
      })
    }
  })
  onUnmounted(() => {
    subscribers -= 1
    if (subscribers === 0) {
      stopNotificationListener?.()
      stopConnectionListener?.()
      stopNotificationListener = null
      stopConnectionListener = null
      items.value = []
      unreadCount.value = 0
    }
  })
  return { items, unreadCount, error, refreshNotifications }
}
