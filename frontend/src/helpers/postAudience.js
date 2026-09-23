// The followers endpoint returns an array for the signed-in user's followers,
// but an ID-keyed object for another profile. Read the ID from the user record
// first so an array position can never be mistaken for a real account ID.
export function normalizePostAudience(data) {
  const rows = Array.isArray(data)
    ? data.map(user => [null, user])
    : Object.entries(data || {})

  return rows
    .map(([fallbackId, user]) => {
      const id = Number(user?.ID ?? user?.id ?? fallbackId)
      const firstName = user?.FirstName ?? user?.firstName ?? ''
      const lastName = user?.LastName ?? user?.lastName ?? ''
      const username = user?.UserName ?? user?.username ?? ''
      const name = `${firstName} ${lastName}`.trim() || username || 'Orbit member'

      return { id, name }
    })
    .filter(user => Number.isSafeInteger(user.id) && user.id > 0)
    .sort((left, right) => left.id - right.id)
}
