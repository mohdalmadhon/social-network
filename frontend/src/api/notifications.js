import { checkSessionResponse } from '@/helpers/auth/auth'
import { router } from '@/router/router'

async function request(url, options = {}) {
  const response = await fetch(url, {
    credentials: 'include',
    ...options,
  })

  if (!checkSessionResponse(response)) {
    router.replace('/login')
    return
  }

  const result = await response.json()
  if (!response.ok) {
    throw new Error(result.message || 'Could not update notifications')
  }

  return result
}

export function getNotifications(category = 'all') {
  return request(`/api/notifications?category=${encodeURIComponent(category)}`)
}

export function markNotificationRead(notificationId) {
  return request(`/api/notifications/${notificationId}/read`, { method: 'PATCH' })
}

export function markAllNotificationsRead() {
  return request('/api/notifications/read-all', { method: 'PATCH' })
}

export function applyNotificationAction(notificationId, action) {
  return request(`/api/notifications/${notificationId}/action`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ action }),
  })
}
