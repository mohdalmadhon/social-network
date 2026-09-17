import { onMounted, onUnmounted, ref } from 'vue'
import { getNotifications } from '@/api/notifications.js'

const chatCount = ref(0)
const error = ref('')
let timer
let subscribers = 0
let pending = false

export async function refreshChats() {
  if (pending) return
  pending = true

  try {
    const result = await getNotifications('messages')
    chatCount.value = result?.unreadCount || 0
    error.value = ''
  } catch (failure) {
    error.value = failure.message || 'Could not load chats.'
  } finally {
    pending = false
  }
}

// Header and sidebar share one small polling loop so they show the same count.
export function useChatCount() {
  onMounted(() => {
    subscribers += 1
    if (subscribers === 1) {
      refreshChats()
      timer = window.setInterval(refreshChats, 5000)
    }
  })

  onUnmounted(() => {
    subscribers -= 1
    if (subscribers === 0) {
      window.clearInterval(timer)
      chatCount.value = 0
    }
  })

  return { chatCount, error, refreshChats }
}
