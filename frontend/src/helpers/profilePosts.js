const avatarColors = ['#3ee6b0', '#ff6b8a', '#7c5cff', '#ffb84d', '#4cc3ff']

function formatPostTime(value) {
  if (!value) return 'Just now'
  const date = new Date(String(value).replace(' ', 'T'))
  return Number.isNaN(date.getTime())
    ? 'Recently'
    : date.toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' })
}

function privacyLabel(value) {
  return { public: 'Public', followers: 'Followers only', selected: 'Selected followers' }[value] || value
}


