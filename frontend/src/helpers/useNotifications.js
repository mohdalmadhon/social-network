import { ref, onMounted, onUnmounted } from 'vue'
import { getNotifications } from '@/api/notifications.js'

const items = ref([])
const unreadCount = ref(0)
const error = ref('')
let timer
let subscribers = 0
let pending = false

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

// One shared timer, even when the header, sidebar and page all subscribe.
export function useNotifications() {
  onMounted(() => {
    subscribers += 1
    if (subscribers === 1) {
      refreshNotifications()
      timer = window.setInterval(refreshNotifications, 5000)
    }
  })
  onUnmounted(() => {
    subscribers -= 1
    if (subscribers === 0) {
      window.clearInterval(timer)
      items.value = []
      unreadCount.value = 0
    }
  })
  return { items, unreadCount, error, refreshNotifications }
}
