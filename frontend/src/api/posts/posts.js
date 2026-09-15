import { checkSessionResponse } from '@/helpers/auth/auth'
import { router } from '@/router/router'

export async function createPost(formData) {
    const response = await fetch('/api/posts', {
        method: 'POST',
        credentials: 'include',
        body: formData,
    })

    if (!checkSessionResponse(response)) {
        router.replace('/login')
        return
    }

    const result = await response.json()

    if (!response.ok) {
        throw new Error(result.message || 'Failed to create post')
    }

    return result
}

export async function getPosts() {
    const response = await fetch('/api/posts', {
        method: 'GET',
        credentials: 'include',
    })

    if (!checkSessionResponse(response)) {
        router.replace('/login')
        return
    }

    const result = await response.json()

    if (!response.ok) {
        throw new Error(result.message || 'Could not load posts')
    }

    return result
}

export async function setPostLike(postId, liked) {
    const response = await fetch(`/api/posts/${postId}/like`, {
        method: liked ? 'PUT' : 'DELETE',
        credentials: 'include',
    })

    if (!checkSessionResponse(response)) {
        router.replace('/login')
        return
    }

    const result = await response.json()
    if (!response.ok) {
        throw new Error(result.message || 'Could not update the like')
    }

    return result
}
