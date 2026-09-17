import { onMounted, onUnmounted, ref } from 'vue'
import { getNotifications } from '@/api/notifications.js'
import { subscribeRealtime } from '@/services/realtime.js'

const chatCount = ref(0)
const error = ref('')
let subscribers = 0
let pending = false
let stopNotificationListener
let stopConnectionListener
const knownNotificationIDs = new Set()

export async function refreshChats() {
  if (pending) return
  pending = true

  try {
    const result = await getNotifications('messages')
    chatCount.value = result?.unreadCount || 0
    knownNotificationIDs.clear()
    for (const notification of result?.notifications || []) knownNotificationIDs.add(notification.id)
    error.value = ''
  } catch (failure) {
    error.value = failure.message || 'Could not load chats.'
  } finally {
    pending = false
  }
}

export function useChatCount() {
  onMounted(() => {
    subscribers += 1
    if (subscribers === 1) {
      refreshChats()
      stopNotificationListener = subscribeRealtime('notification', event => {
        const notification = event?.notification
        if (notification?.category !== 'messages' || knownNotificationIDs.has(notification.id)) return
        knownNotificationIDs.add(notification.id)
        if (!notification.isRead) chatCount.value += 1
      })
      stopConnectionListener = subscribeRealtime('connection', event => {
        if (event.status === 'connected') refreshChats()
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
      knownNotificationIDs.clear()
      chatCount.value = 0
    }
  })

  return { chatCount, error, refreshChats }
}
