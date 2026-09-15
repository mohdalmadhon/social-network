<script setup>
import { onMounted, ref } from 'vue'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import FeedSidebar from '@/components/posts/FeedSidebar.vue'
import PostCard from '@/components/posts/PostCard.vue'
import PostComposer from '@/components/posts/PostComposer.vue'
import { getPosts } from '@/api/posts/posts.js'

const posts = ref([])
const isLoading = ref(true)
const feedError = ref('')

const avatarColors = ['#3ee6b0', '#ff6b8a', '#7c5cff', '#ffb84d', '#4cc3ff']

function formatPostTime(value) {
  if (!value) return 'Just now'

  const date = new Date(value.replace(' ', 'T'))
  if (Number.isNaN(date.getTime())) return 'Recently'

  return date.toLocaleString([], {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

function privacyLabel(value) {
  const labels = {
    public: 'Public',
    followers: 'Followers only',
    selected: 'Selected followers',
  }

  return labels[value] || value
}

function toCardPost(post, index = 0) {
  return {
    id: post.id,
    author: post.author || 'Orbit member',
    avatarColor: avatarColors[index % avatarColors.length],
    avatarPath: post.avatarPath || '',
    time: formatPostTime(post.createdAt),
    privacy: privacyLabel(post.privacy),
    content: post.content,
    likes: post.likeCount || 0,
    liked: Boolean(post.liked),
    comments: post.commentCount || 0,
    imagePath: post.imagePath || '',
    hasMedia: false,
    mediaDescription: '',
  }
}

async function loadPosts() {
  isLoading.value = true
  feedError.value = ''

  try {
    const result = await getPosts()
    posts.value = (result?.posts || []).map(toCardPost)
  } catch (error) {
    feedError.value = error.message || 'Could not load your feed.'
  } finally {
    isLoading.value = false
  }
}

function addPost(post) {
  if (post) {
    posts.value.unshift(toCardPost(post))
  }
}

onMounted(loadPosts)
</script>

<template>
  <AuthenticatedLayout active-page="home">
    <div class="feed-layout">
      <div class="home-feed">
        <h1 class="visually-hidden">Home feed</h1>
        <PostComposer @post-created="addPost" />

        <p v-if="isLoading" class="feed-state orbit-surface">Loading your feed...</p>

        <div v-else-if="feedError" class="feed-state orbit-surface">
          <p>{{ feedError }}</p>
          <button type="button" @click="loadPosts">Try again</button>
        </div>

        <p v-else-if="posts.length === 0" class="feed-state orbit-surface">
          No posts yet. Share something with your orbit.
        </p>

        <PostCard v-for="post in posts" v-else :key="post.id" :post="post" />
      </div>
      <FeedSidebar />
    </div>
  </AuthenticatedLayout>
</template>

<style scoped>
.feed-layout {
  display: flex;
  align-items: flex-start;
  gap: var(--space-5);
}

.home-feed {
  display: grid;
  width: 100%;
  max-width: 48rem;
  gap: var(--space-4);
}

.feed-state {
  margin: 0;
  padding: var(--space-5);
  color: var(--color-text-muted);
  text-align: center;
}

.feed-state p {
  margin: 0;
}

.feed-state button {
  margin-top: var(--space-3);
  padding: var(--space-2) var(--space-4);
  border: 0;
  border-radius: 999px;
  background: var(--gradient-action);
  color: white;
  cursor: pointer;
  font-weight: 700;
}

@media (min-width: 48rem) {
  .home-feed {
    gap: var(--space-5);
  }
}

@media (min-width: 90rem) {
  .feed-layout {
    justify-content: center;
  }
}
</style>
