import { checkSessionResponse } from '@/helpers/auth/auth'
import { router } from '@/router/router'

async function chatRequest(path, options = {}) {
  const response = await fetch(path, { credentials: 'include', ...options })
  if (!checkSessionResponse(response)) {
    router.replace('/login')
    return
  }
  const result = await response.json()
  if (!response.ok) throw new Error(result.message || 'Could not complete chat request')
  return result
}

export const getPrivateChats = () => chatRequest('/api/chats')
export const getPrivateChatUsers = () => chatRequest('/api/chats/private-users')
export const getPrivateMessages = (chatId, { limit = 20, offset = 0 } = {}) =>
  chatRequest(`/api/chats/${chatId}/messages?limit=${limit}&offset=${offset}`)
export const getGroupMessages = (groupId, { limit = 20, offset = 0 } = {}) =>
  chatRequest(`/api/groups/${groupId}/chat/messages?limit=${limit}&offset=${offset}`)

export function openPrivateChat(userId) {
  return chatRequest('/api/chats/private', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId }),
  })
}
