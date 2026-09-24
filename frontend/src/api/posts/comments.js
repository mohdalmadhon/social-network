import { checkSessionResponse } from '@/helpers/auth/auth'
import { router } from '@/router/router'

async function requestComments(url, options = {}) {
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
    throw new Error(result.message || 'Could not load comments')
  }

  return result
}

export async function getComments(postId, { limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  })

  return requestComments(`/api/posts/${postId}/comments?${params}`)
}

export async function createComment(postId, comment) {
  return requestComments(`/api/posts/${postId}/comments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ content: comment.content }),
  })
}

export async function deleteComment(postId, commentId) {
  return requestComments(`/api/posts/${postId}/comments/${commentId}`, {
    method: 'DELETE',
  })
}