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

export function profilePosts(posts, profile) {
  const fullName = `${profile.firstName || ''} ${profile.lastName || ''}`.trim()
  return (posts || [])
    .filter((post) => (
      (profile.id && Number(post.userId) === Number(profile.id))
      || (profile.userName && post.author === profile.userName)
      || (!profile.userName && fullName && post.author === fullName)
    ))
    .map((post, index) => ({
      id: post.id,
      authorId: post.userId,
      author: post.author || fullName || 'Orbit member',
      avatarColor: avatarColors[index % avatarColors.length],
      avatarPath: post.avatarPath || '',
      time: formatPostTime(post.createdAt),
      privacy: privacyLabel(post.privacy),
      content: post.content,
      likes: post.likeCount || 0,
      comments: post.commentCount || 0,
      imagePath: post.imagePath || '',
      hasMedia: false,
      location: post.location,
      mediaDescription: '',
    }))
}
