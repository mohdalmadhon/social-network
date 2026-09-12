import { checkSessionResponse } from '@/helpers/auth/auth'
import { router } from '@/router/router'

export async function getSearchResults(query, options = {}) {
  const params = new URLSearchParams({ q: query })
  if (options.excludeGroupId) params.set('excludeGroupId', options.excludeGroupId)

  const response = await fetch(`/api/search?${params.toString()}`, {
    credentials: 'include',
  })

  if (!checkSessionResponse(response)) {
    router.replace('/login')
    return
  }

  const result = await response.json()
  if (!response.ok) {
    throw new Error(result.message || 'Could not search Orbit')
  }

  return result
}
