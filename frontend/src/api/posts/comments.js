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

export async function getComments(postId) {
  return requestComments(`/api/posts/${postId}/comments`)
}

export async function createComment(postId, comment) {
  if (comment.file) {
    const formData = new FormData()
    formData.append('content', comment.content)
    formData.append('image', comment.file)

    return requestComments(`/api/posts/${postId}/comments`, {
      method: 'POST',
      body: formData,
    })
  }

  return requestComments(`/api/posts/${postId}/comments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ content: comment.content }),
  })
}
