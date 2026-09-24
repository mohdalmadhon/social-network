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
const commentsLoading = ref(false)
const commentsLoadingMore = ref(false)
const commentsHasMore = ref(true)
const commentsOffset = ref(0)
const commentsOpen = ref(false)
const commentsLoaded = ref(false)
const newComment = ref('')
const commentError = ref('')
const likeLoading = ref(false)
const liked = ref(Boolean(props.post.liked))
const likeCount = ref(Number(props.post.like_count ?? props.post.likeCount ?? 0))
const commentCount = ref(
  Number(props.post.comment_count ?? props.post.commentCount ?? 0),
)

const commentsSentinel = ref(null)
let commentsObserver = null

const showLocationDialog = ref(false)

const imageUrl = computed(() => {
  if (!props.post.image_path && !props.post.imagePath) return ''

  const path = props.post.image_path ?? props.post.imagePath

  if (String(path).startsWith('http')) {
    return path
  }

  return `/uploads/${String(path).replace(/^\/+/, '')}`
})

const postLocation = computed(() => {
  if (!props.post.location) return null

  const parts = String(props.post.location).split(':')

  if (parts.length < 3) return null

  const label = parts[0].trim()
  const lat = Number(parts[1])
  const lon = Number(parts[2])

  if (!label || !Number.isFinite(lat) || !Number.isFinite(lon)) {
    return null
  }

  return {
    label,
    lat,
    lon,
  }
})

const googleMapsUrl = computed(() => {
  if (!postLocation.value) return ''

  const { lat, lon } = postLocation.value

  return `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(
    `${lat},${lon}`,
  )}`
})

const googleMapsEmbedUrl = computed(() => {
  if (!postLocation.value) return ''

  const { lat, lon } = postLocation.value

  return `https://www.google.com/maps?q=${encodeURIComponent(
    `${lat},${lon}`,
  )}&output=embed`
})

function openLocationDialog() {
  if (!postLocation.value) return

  showLocationDialog.value = true
}

function closeLocationDialog() {
  showLocationDialog.value = false
}

function handleLocationKeydown(event) {
  if (event.key === 'Escape') {
    closeLocationDialog()
  }
}

function initials(name) {
  if (!name) return '?'

  return String(name)
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part.charAt(0).toUpperCase())
    .join('')
}

function formatRelativeTime(dateValue) {
  if (!dateValue) return ''

  const date = new Date(dateValue)

  if (Number.isNaN(date.getTime())) {
    return String(dateValue)
  }

  const seconds = Math.floor((Date.now() - date.getTime()) / 1000)

  if (seconds < 60) {
    return 'just now'
  }

  const minutes = Math.floor(seconds / 60)

  if (minutes < 60) {
    return `${minutes}m`
  }

  const hours = Math.floor(minutes / 60)

  if (hours < 24) {
    return `${hours}h`
  }

  const days = Math.floor(hours / 24)

  if (days < 7) {
    return `${days}d`
  }

  const weeks = Math.floor(days / 7)

  if (weeks < 5) {
    return `${weeks}w`
  }

  const months = Math.floor(days / 30)

  if (months < 12) {
    return `${months}mo`
  }

  const years = Math.floor(days / 365)

  return `${years}y`
}

async function toggleLike() {
  if (likeLoading.value) return

  const previousLiked = liked.value
  const previousCount = likeCount.value

  liked.value = !previousLiked
  likeCount.value += liked.value ? 1 : -1

  if (likeCount.value < 0) {
    likeCount.value = 0
  }

  likeLoading.value = true

  try {
    const response = await setPostLike(props.post.id, liked.value)

    if (response?.liked !== undefined) {
      liked.value = Boolean(response.liked)
    }

    if (
      response?.like_count !== undefined ||
      response?.likeCount !== undefined
    ) {
      likeCount.value = Number(
        response.like_count ?? response.likeCount,
      )
    }
  } catch (error) {
    liked.value = previousLiked
    likeCount.value = previousCount
    console.error(error)
  } finally {
    likeLoading.value = false
  }
}

async function loadComments(reset = false) {
  if (commentsLoading.value || commentsLoadingMore.value) return

  if (!reset && !commentsHasMore.value) return

  if (reset) {
    commentsOffset.value = 0
    commentsHasMore.value = true
    comments.value = []
    commentsLoading.value = true
  } else {
    commentsLoadingMore.value = true
  }

  try {
    const response = await getComments(
      props.post.id,
      commentsOffset.value,
      10,
    )

    const receivedComments =
      response?.comments ??
      response?.data ??
      (Array.isArray(response) ? response : [])

    const normalizedComments = Array.isArray(receivedComments)
      ? receivedComments
      : []

    if (reset) {
      comments.value = normalizedComments
    } else {
      comments.value.push(...normalizedComments)
    }

    commentsOffset.value += normalizedComments.length

    if (
      normalizedComments.length < 10 ||
      response?.hasMore === false ||
      response?.has_more === false
    ) {
      commentsHasMore.value = false
    }
  } catch (error) {
    console.error(error)
  } finally {
    commentsLoading.value = false
    commentsLoadingMore.value = false
  }
}

async function openComments() {
  commentsOpen.value = !commentsOpen.value

  if (commentsOpen.value && !commentsLoaded.value) {
    commentsLoaded.value = true
    await loadComments(true)
    await nextTick()
    setupCommentsObserver()
  }
}

async function addComment(content, clearInput) {
  const value = String(content ?? '').trim()

  if (!value) return

  commentError.value = ''

  try {
    const response = await createComment(props.post.id, value)

    const createdComment =
      response?.comment ??
      response?.data ??
      response

    if (createdComment && typeof createdComment === 'object') {
      comments.value.unshift(createdComment)
    }

    commentCount.value += 1

    if (typeof clearInput === 'function') {
      clearInput()
    }
  } catch (error) {
    console.error(error)
    commentError.value = 'Could not add comment.'
  }
}

function setupCommentsObserver() {
  if (!commentsSentinel.value) return

  if (commentsObserver) {
    commentsObserver.disconnect()
  }

  commentsObserver = new IntersectionObserver(
    (entries) => {
      if (!entries[0]?.isIntersecting) return

      loadComments(false)
    },
    {
      root: null,
      rootMargin: '200px',
      threshold: 0,
    },
  )

  commentsObserver.observe(commentsSentinel.value)
}

onBeforeUnmount(() => {
  if (commentsObserver) {
    commentsObserver.disconnect()
  }
})
</script>

<template>
  <article class="post-card">
    <header class="post-card__header">
      <div class="post-card__author-info">
        <div class="post-card__author-name-row">
          <strong class="post-card__author-name">
            {{ post.author || 'Unknown user' }}
          </strong>

          <button v-if="postLocation" type="button" class="post-card__location-small"
            :title="`View location: ${postLocation.label}`" @click="openLocationDialog">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 21s7-6.1 7-12a7 7 0 1 0-14 0c0 5.9 7 12 7 12Z" />
              <circle cx="12" cy="9" r="2.25" />
            </svg>

            <span>{{ postLocation.label }}</span>
          </button>
        </div>

        <span class="post-card__time">
          {{ formatRelativeTime(post.created_at) }}
        </span>
      </div>

      <button type="button" class="post-card__more" aria-label="Post options">
        <IconGlyph name="more-horizontal" :size="18" />
      </button>
    </header>

    <div class="post-card__body">
      <p v-if="post.content" class="post-card__content">
        {{ post.content }}
      </p>

      <div v-if="postLocation" class="post-card__location">
        <button type="button" class="post-card__location-button" @click="openLocationDialog">
          <span class="post-card__location-icon">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 21s7-6.1 7-12a7 7 0 1 0-14 0c0 5.9 7 12 7 12Z" />
              <circle cx="12" cy="9" r="2.25" />
            </svg>
          </span>

          <span class="post-card__location-content">
            <small>Location</small>
            <strong>{{ postLocation.label }}</strong>
          </span>

          <span class="post-card__location-arrow">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="m9 18 6-6-6-6" />
            </svg>
          </span>
        </button>
      </div>

      <div v-if="imageUrl" class="post-card__media">
        <img :src="imageUrl" alt="Post image" loading="lazy" />
      </div>
    </div>

    <footer class="post-card__footer">
      <button type="button" class="post-card__action" :class="{ 'post-card__action--liked': liked }"
        :disabled="likeLoading" @click="toggleLike">
        <IconGlyph :name="liked ? 'heart-filled' : 'heart'" :size="18" />

        <span>{{ likeCount }}</span>
      </button>

      <button type="button" class="post-card__action" @click="openComments">
        <IconGlyph name="message-circle" :size="18" />

        <span>{{ commentCount }}</span>
      </button>
    </footer>

    <section v-if="commentsOpen" class="post-card__comments">
      <div v-if="commentsLoading" class="post-card__comments-loading">
        Loading comments...
      </div>

      <template v-else>
        <div v-if="comments.length" class="post-card__comments-list">
          <CommentPreview v-for="comment in comments" :key="comment.id" :comment="comment" />
        </div>

        <div v-else class="post-card__comments-empty">
          No comments yet.
        </div>

        <div ref="commentsSentinel" class="post-card__comments-sentinel">
          <span v-if="commentsLoadingMore">
            Loading more...
          </span>

          <span v-else-if="!commentsHasMore && comments.length">
            No more comments
          </span>
        </div>

        <p v-if="commentError" class="post-card__comment-error">
          {{ commentError }}
        </p>

        <CommentInput @send="addComment" />
      </template>
    </section>

    <Teleport to="body">
      <div v-if="showLocationDialog && postLocation" class="location-map-dialog" @keydown="handleLocationKeydown">
        <button type="button" class="location-map-dialog__backdrop" aria-label="Close map"
          @click="closeLocationDialog" />

        <section class="location-map-dialog__panel" role="dialog" aria-modal="true"
          aria-labelledby="location-map-title">
          <header class="location-map-dialog__header">
            <div class="location-map-dialog__title">
              <p class="location-map-dialog__eyebrow">
                Post location
              </p>

              <h2 id="location-map-title">
                {{ postLocation.label }}
              </h2>
            </div>

            <button type="button" class="location-map-dialog__close" aria-label="Close map"
              @click="closeLocationDialog">
              <IconGlyph name="close" :size="18" />
            </button>
          </header>

          <div class="location-map-dialog__map">
            <iframe :src="googleMapsEmbedUrl" title="Google Maps location" loading="lazy"
              referrerpolicy="no-referrer-when-downgrade" />
          </div>

          <footer class="location-map-dialog__footer">
            <div class="location-map-dialog__coordinates">
              {{ postLocation.lat.toFixed(6) }},
              {{ postLocation.lon.toFixed(6) }}
            </div>

            <a class="location-map-dialog__open" :href="googleMapsUrl" target="_blank" rel="noopener noreferrer">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M14 5h5v5" />
                <path d="m19 5-8 8" />
                <path d="M19 13v5a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h5" />
              </svg>

              Open in Google Maps
            </a>
          </footer>
        </section>
      </div>
    </Teleport>
  </article>
</template>

<style scoped>
.post-card {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: 1rem;
  background: var(--color-surface-raised);
  box-shadow: var(--shadow-soft);
}

.post-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.1rem 0.8rem;
}

.post-card__author {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 0.75rem;
}

.post-card__avatar {
  display: flex;
  width: 2.7rem;
  height: 2.7rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 50%;
  background: var(--color-violet-soft);
  color: var(--color-text);
  font-size: 0.78rem;
  font-weight: 700;
}

.post-card__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-card__author-info {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 0.15rem;
}

.post-card__author-name {
  overflow: hidden;
  color: var(--color-text);
  font-size: 0.85rem;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-card__time {
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.65rem;
}

.post-card__more {
  display: flex;
  width: 2.25rem;
  height: 2.25rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 0;
  border-radius: 50%;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.post-card__more:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.post-card__body {
  padding: 0 1.1rem 1rem;
}

.post-card__content {
  margin: 0;
  color: var(--color-text);
  font-size: 0.88rem;
  line-height: 1.6;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.post-card__media {
  margin-top: 1rem;
  overflow: hidden;
  border-radius: 0.85rem;
  background: var(--color-input);
}

.post-card__media img {
  display: block;
  width: 100%;
  max-height: 38rem;
  object-fit: cover;
}

.post-card__location {
  margin-top: var(--space-3);
}

.post-card__location-button {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.75rem;
  padding: 0.7rem 0.8rem;
  border: 1px solid color-mix(in srgb,
      var(--color-violet) 28%,
      var(--color-border));
  border-radius: 0.8rem;
  background: color-mix(in srgb,
      var(--color-violet) 6%,
      var(--color-input));
  color: var(--color-text);
  cursor: pointer;
  text-align: left;
  transition:
    background 160ms ease,
    border-color 160ms ease,
    transform 160ms ease;
}

.post-card__location-button:hover {
  border-color: var(--color-violet);
  background: color-mix(in srgb,
      var(--color-violet) 11%,
      var(--color-input));
  transform: translateY(-1px);
}

.post-card__location-icon {
  display: flex;
  width: 2.35rem;
  height: 2.35rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 0.65rem;
  background: color-mix(in srgb,
      var(--color-violet) 14%,
      var(--color-input));
  color: var(--color-violet-soft);
}

.post-card__location-icon svg {
  width: 1.1rem;
  height: 1.1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.post-card__location-content {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 0.15rem;
}

.post-card__location-content small {
  color: var(--color-text-faint);
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.post-card__location-content strong {
  overflow: hidden;
  color: var(--color-text);
  font-size: 0.82rem;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-card__location-arrow {
  display: flex;
  width: 2rem;
  height: 2rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  color: var(--color-text-faint);
}

.post-card__location-arrow svg {
  width: 1rem;
  height: 1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.post-card__footer {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.7rem 1rem;
  border-top: 1px solid var(--color-border);
}

.post-card__action {
  display: inline-flex;
  min-width: 4rem;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  min-height: 2.2rem;
  padding: 0 0.7rem;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 0.75rem;
  font-weight: 650;
  transition:
    background 160ms ease,
    color 160ms ease;
}

.post-card__action:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.post-card__action--liked {
  color: var(--color-violet-soft);
}

.post-card__action:disabled {
  cursor: default;
  opacity: 0.65;
}

.post-card__comments {
  padding: 0 1rem 1rem;
  border-top: 1px solid var(--color-border);
}

.post-card__comments-list {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  padding-top: 0.8rem;
}

.post-card__comments-loading,
.post-card__comments-empty {
  padding: 1rem 0;
  color: var(--color-text-faint);
  font-size: 0.75rem;
  text-align: center;
}

.post-card__comments-sentinel {
  min-height: 2rem;
  padding: 0.5rem;
  color: var(--color-text-faint);
  font-size: 0.68rem;
  text-align: center;
}

.post-card__comment-error {
  margin: 0.5rem 0;
  color: var(--color-danger, #ff5c7a);
  font-size: 0.72rem;
}

.location-map-dialog {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.location-map-dialog__backdrop {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  padding: 0;
  border: 0;
  background: rgb(0 0 0 / 65%);
  backdrop-filter: blur(4px);
  cursor: default;
}

.location-map-dialog__panel {
  position: relative;
  z-index: 1;
  display: flex;
  width: min(100%, 48rem);
  max-height: 90vh;
  overflow: hidden;
  flex-direction: column;
  border: 1px solid var(--color-border);
  border-radius: 1rem;
  background: var(--color-surface-raised);
  box-shadow: 0 1.5rem 4rem rgb(0 0 0 / 40%);
}

.location-map-dialog__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.1rem 1.25rem;
  border-bottom: 1px solid var(--color-border);
}

.location-map-dialog__title {
  min-width: 0;
}

.location-map-dialog__eyebrow {
  margin: 0 0 0.2rem;
  color: var(--color-violet-soft);
  font-family: var(--font-meta);
  font-size: 0.65rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.location-map-dialog__header h2 {
  margin: 0;
  overflow: hidden;
  color: var(--color-text);
  font-size: 1rem;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.location-map-dialog__close {
  display: flex;
  width: 2.4rem;
  height: 2.4rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 1px solid var(--color-border);
  border-radius: 50%;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.location-map-dialog__close:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.location-map-dialog__map {
  width: 100%;
  height: min(55vh, 28rem);
  background: var(--color-input);
}

.location-map-dialog__map iframe {
  display: block;
  width: 100%;
  height: 100%;
  border: 0;
}

.location-map-dialog__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.25rem;
  border-top: 1px solid var(--color-border);
}

.location-map-dialog__coordinates {
  min-width: 0;
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.7rem;
}

.location-map-dialog__open {
  display: inline-flex;
  min-height: 2.75rem;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0 1rem;
  border-radius: 999px;
  background: var(--gradient-action);
  color: white;
  font-size: 0.8rem;
  font-weight: 700;
  text-decoration: none;
  white-space: nowrap;
  transition:
    transform 160ms ease,
    box-shadow 160ms ease;
}

.location-map-dialog__open:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-soft);
}

.location-map-dialog__open svg {
  width: 1rem;
  height: 1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

@media (max-width: 36rem) {
  .post-card__header {
    padding: 0.85rem 0.85rem 0.7rem;
  }

  .post-card__body {
    padding: 0 0.85rem 0.85rem;
  }

  .post-card__footer {
    padding: 0.6rem 0.75rem;
  }

  .post-card__location-button {
    padding: 0.65rem;
  }

  .location-map-dialog {
    align-items: flex-end;
    padding: 0;
  }

  .location-map-dialog__panel {
    width: 100%;
    max-height: 92vh;
    border-radius: 1rem 1rem 0 0;
  }

  .location-map-dialog__map {
    height: 55vh;
  }

  .location-map-dialog__footer {
    align-items: stretch;
    flex-direction: column;
  }

  .location-map-dialog__open {
    width: 100%;
  }
}

.post-card__author-name-row {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 0.45rem;
}

.post-card__location-small {
  display: inline-flex;
  min-width: 0;
  max-width: 12rem;
  align-items: center;
  gap: 0.25rem;
  padding: 0.18rem 0.45rem;
  border: 1px solid color-mix(
    in srgb,
    var(--color-violet) 25%,
    var(--color-border)
  );
  border-radius: 999px;
  background: color-mix(
    in srgb,
    var(--color-violet) 7%,
    var(--color-input)
  );
  color: var(--color-violet-soft);
  cursor: pointer;
  font-size: 0.62rem;
  font-weight: 650;
  line-height: 1;
  transition:
    background 160ms ease,
    border-color 160ms ease,
    transform 160ms ease;
}

.post-card__location-small:hover {
  border-color: var(--color-violet);
  background: color-mix(
    in srgb,
    var(--color-violet) 14%,
    var(--color-input)
  );
  transform: translateY(-1px);
}

.post-card__location-small svg {
  width: 0.75rem;
  height: 0.75rem;
  flex-shrink: 0;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.post-card__location-small span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>