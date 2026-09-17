<script setup>
import { onMounted, ref } from 'vue'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import FeedSidebar from '@/components/posts/FeedSidebar.vue'
import PostCard from '@/components/posts/PostCard.vue'
import PostComposer from '@/components/posts/PostComposer.vue'
import { getPosts } from '@/api/posts/posts.js'

const posts = ref([])
const isLoading = ref(true)
const isLoadingMore = ref(false)
const hasMorePosts = ref(false)
const feedError = ref('')
const FEED_PAGE_SIZE = 20

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
    authorId: post.userId,
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

async function loadPosts({ append = false } = {}) {
  if (append) {
    if (isLoadingMore.value || !hasMorePosts.value) return
    isLoadingMore.value = true
  } else {
    isLoading.value = true
  }
  feedError.value = ''

  try {
    const result = await getPosts({
      limit: FEED_PAGE_SIZE,
      offset: append ? posts.value.length : 0,
    })
    const nextPosts = (result?.posts || []).map(toCardPost)
    posts.value = append ? [...posts.value, ...nextPosts] : nextPosts
    hasMorePosts.value = Boolean(result?.hasMore)
  } catch (error) {
    feedError.value = error.message || 'Could not load your feed.'
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
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

        <template v-else>
          <PostCard v-for="post in posts" :key="post.id" :post="post" />

          <div v-if="hasMorePosts" class="feed-pagination">
            <button type="button" :disabled="isLoadingMore" @click="loadPosts({ append: true })">
              {{ isLoadingMore ? 'Loading more posts...' : 'Load more posts' }}
            </button>
          </div>
        </template>
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

.feed-pagination {
  display: flex;
  justify-content: center;
  padding: var(--space-2) 0 var(--space-3);
}

.feed-pagination button {
  min-height: var(--touch-target);
  padding: 0 var(--space-5);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: var(--color-surface);
  color: var(--color-text-muted);
  cursor: pointer;
  font: inherit;
  font-weight: 600;
}

.feed-pagination button:hover:not(:disabled),
.feed-pagination button:focus-visible {
  border-color: var(--color-violet);
  color: var(--color-text);
}

.feed-pagination button:disabled {
  cursor: wait;
  opacity: 0.6;
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
