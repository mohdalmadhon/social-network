<script setup>
import { computed, nextTick, onBeforeUnmount, ref } from 'vue'
import CommentInput from '@/components/comments/CommentInput.vue'
import CommentPreview from '@/components/comments/CommentPreview.vue'
import IconGlyph from '@/components/layout/IconGlyph.vue'

import { createComment, getComments } from '@/api/posts/comments.js'
import { setPostLike } from '@/api/posts/posts.js'

const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
})

const comments = ref([])
const commentInput = ref(null)
const commentsError = ref('')
const isLoadingComments = ref(false)
const isLoadingMoreComments = ref(false)
const areCommentsOpen = ref(false)
const commentsLoaded = ref(false)
const hasMoreComments = ref(false)
const commentsList = ref(null)
const commentsSentinel = ref(null)
const commentsOffset = ref(0)
const isSubmittingComment = ref(false)
const isLiked = ref(Boolean(props.post.liked))
const likeCount = ref(Number(props.post.likes) || 0)
const isLikePending = ref(false)
const likeError = ref('')
const COMMENTS_PAGE_SIZE = 20
let commentsObserver

const commentCount = computed(() => Math.max(props.post.comments, comments.value.length))

function formatRelativeTime(dateValue) {
  if (!dateValue) return ''

  const date = new Date(dateValue)

  if (Number.isNaN(date.getTime())) return dateValue

  const now = new Date()
  const difference = now.getTime() - date.getTime()

  const minutes = Math.floor(difference / 60000)
  const hours = Math.floor(difference / 3600000)
  const days = Math.floor(difference / 86400000)
  const weeks = Math.floor(difference / 604800000)

  if (difference < 0) return 'just now'

  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m`
  if (hours < 24) return `${hours}h`
  if (days < 7) return `${days}d`
  if (weeks < 5) return `${weeks}w`

  return date.toLocaleDateString(undefined, {
    day: 'numeric',
    month: 'short',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  })
}

const formattedTime = computed(() => formatRelativeTime(props.post.time))

function commentForPreview(comment) {
  return {
    ...comment,
    author: comment.author || 'Orbit member',
    avatarColor: 'var(--gradient-action)',
  }
}

async function loadComments({ append = false } = {}) {
  if (append) {
    if (isLoadingMoreComments.value || !hasMoreComments.value) return
    isLoadingMoreComments.value = true
  } else {
    if (commentsLoaded.value) return
    isLoadingComments.value = true
  }

  commentsError.value = ''
  const requestOffset = append ? commentsOffset.value : 0

  try {
    const result = await getComments(props.post.id, {
      limit: COMMENTS_PAGE_SIZE,
      offset: requestOffset,
    })
    const nextComments = (result?.comments || []).map(commentForPreview)
    const combinedComments = append ? [...comments.value, ...nextComments] : nextComments
    const uniqueComments = new Map(combinedComments.map(comment => [comment.id, comment]))
    comments.value = [...uniqueComments.values()].sort((first, second) => {
      const timeDifference = Date.parse(first.createdAt || '') - Date.parse(second.createdAt || '')
      return (Number.isFinite(timeDifference) ? timeDifference : 0) || Number(first.id) - Number(second.id)
    })
    commentsOffset.value = Number.isInteger(result?.nextOffset)
      ? result.nextOffset
      : requestOffset + nextComments.length
    hasMoreComments.value = Boolean(result?.hasMore)
    commentsLoaded.value = true
  } catch (error) {
    commentsError.value = error.message || 'Could not load comments.'
  } finally {
    isLoadingComments.value = false
    isLoadingMoreComments.value = false
  }
}

function observeCommentsEnd() {
  commentsObserver?.disconnect()

  if (!areCommentsOpen.value || !commentsList.value || !commentsSentinel.value || typeof IntersectionObserver === 'undefined') return

  commentsObserver = new IntersectionObserver(([entry]) => {
    if (entry.isIntersecting && areCommentsOpen.value && hasMoreComments.value) {
      loadComments({ append: true })
    }
  }, {
    root: commentsList.value,
    rootMargin: '0px 0px 120px',
  })

  commentsObserver.observe(commentsSentinel.value)
}

async function retryComments() {
  await loadComments({ append: commentsLoaded.value && comments.value.length > 0 })
  await nextTick()
  observeCommentsEnd()
}

function closeComments() {
  areCommentsOpen.value = false
  commentsObserver?.disconnect()
}

async function addComment(comment) {
  isSubmittingComment.value = true
  commentsError.value = ''

  try {
    const result = await createComment(props.post.id, comment)

    if (!result?.comment) {
      throw new Error('Could not create comment')
    }

    comments.value.push(commentForPreview(result.comment))
    commentInput.value?.reset()
  } catch (error) {
    commentsError.value = error.message || 'Could not create comment.'
  } finally {
    isSubmittingComment.value = false
  }
}

async function toggleLike() {
  if (isLikePending.value) return

  const nextLiked = !isLiked.value

  isLikePending.value = true
  likeError.value = ''

  try {
    const result = await setPostLike(props.post.id, nextLiked)

    if (!result?.status) {
      throw new Error('Could not update the like')
    }

    isLiked.value = Boolean(result.liked)
    likeCount.value = Number(result.likeCount) || 0
  } catch (error) {
    likeError.value = error.message || 'Could not update the like.'
  } finally {
    isLikePending.value = false
  }
}

function imageUrl(imagePath) {
  if (!imagePath) return ''

  return imagePath.startsWith('/') ? imagePath : `/uploads/${imagePath}`
}

async function toggleComments() {
  areCommentsOpen.value = !areCommentsOpen.value
  if (areCommentsOpen.value) {
    await loadComments()
    await nextTick()
    observeCommentsEnd()
  } else {
    closeComments()
  }
}

function initials(author) {
  return author.slice(0, 2).toUpperCase()
}

onBeforeUnmount(() => {
  commentsObserver?.disconnect()
})

</script>

<template>

  <article class="post-card orbit-surface">

    <header class="post-card__header">

      <RouterLink
        v-if="post.authorId"
        class="post-card__author-link"
        :to="{ path: '/user', query: { id: post.authorId } }"
        :aria-label="`View ${post.author}'s profile`"
      >

        <div class="post-card__avatar" :style="{ background: post.avatarColor }" aria-hidden="true">
          <img v-if="post.avatarPath" :src="imageUrl(post.avatarPath)" alt="" />
          <span v-else>{{ initials(post.author) }}</span>
        </div>

        <div class="post-card__author">
          <h2>{{ post.author }}</h2>
          <p>{{ formattedTime }} <span aria-hidden="true">•</span> {{ post.privacy }}</p>
        </div>

      </RouterLink>

      <template v-else>

        <div class="post-card__avatar" :style="{ background: post.avatarColor }" aria-hidden="true">
          <img v-if="post.avatarPath" :src="imageUrl(post.avatarPath)" alt="" />
          <span v-else>{{ initials(post.author) }}</span>
        </div>

        <div class="post-card__author">
          <h2>{{ post.author }}</h2>
          <p>{{ formattedTime }} <span aria-hidden="true">•</span> {{ post.privacy }}</p>
        </div>

      </template>

    </header>

    <p class="post-card__content">{{ post.content }}</p>

    <div v-if="post.imagePath" class="post-card__media post-card__media--uploaded">
      <img :src="imageUrl(post.imagePath)" alt="Image attached to this post" />
    </div>

    <div
      v-else-if="post.hasMedia"
      class="post-card__media"
      role="img"
      :aria-label="post.mediaDescription"
    >
      <span class="post-card__sun" aria-hidden="true"></span>
      <span class="post-card__mountain post-card__mountain--back" aria-hidden="true"></span>
      <span class="post-card__mountain post-card__mountain--front" aria-hidden="true"></span>
    </div>

    <footer class="post-card__actions">

      <button
        class="post-action"
        :class="{ 'post-action--liked': isLiked }"
        type="button"
        :aria-pressed="isLiked"
        :aria-label="isLiked ? 'Unlike this post' : 'Like this post'"
        :disabled="isLikePending"
        @click="toggleLike"
      >
        <IconGlyph name="heart" :size="21" />
        <span>{{ likeCount }} likes</span>
      </button>

      <button
        class="post-action post-action--comments"
        type="button"
        :aria-expanded="areCommentsOpen"
        :aria-controls="`comments-${post.id}`"
        @click="toggleComments"
      >
        <IconGlyph name="comment" :size="21" />
        <span>{{ commentCount }} comments</span>
        <span class="post-action__hint">{{ areCommentsOpen ? 'Hide' : 'View' }}</span>
      </button>

    </footer>

    <p v-if="likeError" class="post-action-error" role="alert">{{ likeError }}</p>

    <section
      v-if="areCommentsOpen"
      :id="`comments-${post.id}`"
      class="comments-panel"
      aria-label="Comments"
    >

      <header class="comments-panel__header">

        <div>
          <p class="comments-panel__eyebrow">The conversation</p>
          <h3>Comments <span>{{ commentCount }}</span></h3>
        </div>

        <button
          type="button"
          class="comments-panel__close"
          aria-label="Close comments"
          @click="closeComments"
        >
          <IconGlyph name="close" :size="17" />
        </button>

      </header>

      <p v-if="isLoadingComments" class="comments-state">
        Loading comments...
      </p>

      <div v-else-if="commentsError && !comments.length" class="comments-state comments-state--error">
        <span>{{ commentsError }}</span>
        <button type="button" @click="retryComments">Retry</button>
      </div>
      <p v-else-if="!comments.length" class="comments-state comments-state--empty">No comments yet. Start the conversation.</p>
      <div v-else class="comments-list" ref="commentsList">
        <CommentPreview v-for="comment in comments" :key="comment.id" :comment="comment" />
        <div ref="commentsSentinel" class="comments-sentinel" aria-hidden="true"></div>
      </div>

      <p v-if="isLoadingMoreComments" class="comments-load-state" role="status">Loading more comments...</p>

      <div v-if="commentsError && comments.length" class="comments-state comments-state--error">
        <span>{{ commentsError }}</span>
        <button type="button" @click="retryComments">Retry</button>
      </div>

      <CommentInput
        ref="commentInput"
        :input-id="`comment-${post.id}`"
        :disabled="isSubmittingComment"
        @submit="addComment"
      />

    </section>

  </article>

</template>

<style scoped>

.post-card {
  width: 100%;
  padding: var(--space-4);
}

.post-card__header {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  align-items: center;
  gap: var(--space-3);
}

.post-card__author-link {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr);
  align-items: center;
  gap: var(--space-3);
  min-width: 0;
  color: inherit;
  text-decoration: none;
}

.post-card__author-link:hover .post-card__author h2,
.post-card__author-link:focus-visible .post-card__author h2 {
  color: var(--color-violet-soft);
}

.post-card__avatar {
  display: grid;
  flex: 0 0 var(--touch-target);
  width: var(--touch-target);
  height: var(--touch-target);
  place-items: center;
  aspect-ratio: 1;
  overflow: hidden;
  border-radius: 50%;
  color: #0b0d17;
  font-weight: 700;
  background: var(--color-input);
}

.post-card__avatar img {
  width: 100%;
  height: 100%;
  border-radius: inherit;
  object-fit: cover;
  object-position: center;
}

.post-card__author {
  min-width: 0;
}

.post-card__author h2,
.post-card__author p,
.post-card__content {
  margin: 0;
}

.post-card__author h2 {
  overflow: hidden;
  color: var(--color-text);
  font-size: 1rem;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-card__author p {
  color: var(--color-text-faint);
  font-size: 0.8125rem;
}

.post-card__menu,
.post-action {
  min-width: var(--touch-target);
  min-height: var(--touch-target);
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.post-action--comments {
  align-items: center;
  font: inherit;
}

.post-action__hint {
  color: var(--color-violet-soft);
  font-size: 0.75rem;
  opacity: 0;
  transform: translateX(-0.25rem);
  transition: opacity 160ms ease, transform 160ms ease;
}

.post-action--comments:hover .post-action__hint,
.post-action--comments:focus-visible .post-action__hint {
  opacity: 1;
  transform: translateX(0);
}

.post-card__menu {
  border-radius: 50%;
  font-weight: 700;
  letter-spacing: 0.1em;
}

.post-card__menu:hover,
.post-action:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.post-card__content {
  margin-top: var(--space-4);
  color: var(--color-text-soft);
  line-height: 1.55;
  overflow-wrap: anywhere;
}

.post-card__media {
  position: relative;
  min-height: 12rem;
  margin-top: var(--space-4);
  overflow: hidden;
  border-radius: var(--radius-medium);
  background: linear-gradient(110deg, #44538e 0%, #aa5e9d 54%, #ff8e8b 100%);
}

.post-card__media--uploaded {
  min-height: 0;
  background: var(--color-input);
}

.post-card__media--uploaded img {
  display: block;
  width: 100%;
  max-height: 30rem;
  object-fit: contain;
}

.post-card__sun {
  position: absolute;
  top: 18%;
  right: 18%;
  width: clamp(2.75rem, 9vw, 4.5rem);
  aspect-ratio: 1;
  border-radius: 50%;
  background: #ffe1a3;
}

.post-card__mountain {
  position: absolute;
  right: -4%;
  bottom: -1px;
  left: -4%;
  height: 55%;
  background: #1c2447;
  clip-path: polygon(0 85%, 18% 30%, 35% 76%, 51% 20%, 66% 72%, 84% 12%, 100% 80%, 100% 100%, 0 100%);
}

.post-card__mountain--front {
  height: 42%;
  background: #11182f;
  clip-path: polygon(0 100%, 31% 30%, 49% 78%, 63% 22%, 79% 83%, 100% 34%, 100% 100%);
}

.post-card__actions {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-top: var(--space-2);
}

.comments-state {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin: var(--space-3) 0 0;
  color: var(--color-text-faint);
  font-size: 0.875rem;
}

.comments-panel {
  margin-top: var(--space-4);
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  background: rgb(11 13 23 / 46%);
  animation: comments-panel-in 180ms ease-out;
}

.comments-panel__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.comments-panel__eyebrow {
  margin: 0 0 var(--space-1);
  color: var(--color-mint);
  font-family: var(--font-meta);
  font-size: 0.6875rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.comments-panel h3 {
  margin: 0;
  color: var(--color-text);
  font-size: 1rem;
}

.comments-panel h3 span {
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.75rem;
  font-weight: 400;
}

.comments-panel__close {
  display: grid;
  min-width: var(--touch-target);
  min-height: var(--touch-target);
  place-items: center;
  border: 1px solid var(--color-border);
  border-radius: 50%;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.comments-panel__close:hover,
.comments-panel__close:focus-visible {
  border-color: var(--color-violet);
  color: var(--color-text);
}

.comments-list {
  display: grid;
  gap: var(--space-2);
  max-height: min(28rem, 55vh);
  overflow-y: auto;
  padding-right: var(--space-2);
  overscroll-behavior: contain;
}

.comments-sentinel {
  width: 100%;
  height: 1px;
  pointer-events: none;
}

.comments-load-state {
  margin: var(--space-2) 0 0;
  color: var(--color-text-faint);
  font-size: 0.8125rem;
  text-align: center;
}

.comments-state--empty {
  display: block;
  padding: var(--space-3);
  border-radius: var(--radius-small);
  background: var(--color-surface);
}

@keyframes comments-panel-in {
  from {
    opacity: 0;
    transform: translateY(-0.35rem);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.comments-state--error {
  color: var(--color-coral);
}

.comments-state button {
  min-height: var(--touch-target);
  padding-inline: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: transparent;
  color: inherit;
  cursor: pointer;
}

.post-action {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding-inline: var(--space-2);
  border-radius: var(--radius-small);
  font-size: 0.875rem;
}

.post-action :deep(.icon-glyph) {
  width: 1.35rem;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.7;
}

.post-action--liked {
  color: var(--color-coral);
}

.post-action--liked :deep(.icon-glyph) {
  fill: currentColor;
}

.post-action:disabled {
  cursor: wait;
  opacity: 0.65;
}

.post-action-error {
  margin: var(--space-2) 0 0;
  color: var(--color-coral);
  font-size: 0.8125rem;
}

@media (min-width: 48rem) {
  .post-card {
    padding: var(--space-5);
  }

  .post-card__media {
    min-height: 15rem;
  }
}

</style>
