<script setup>
import { computed, onMounted, ref } from 'vue'
import CommentInput from '@/components/comments/CommentInput.vue'
import CommentPreview from '@/components/comments/CommentPreview.vue'
import { createComment, getComments } from '@/api/posts/comments.js'

const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
})

const liked = ref(false)
const comments = ref([])
const commentsError = ref('')
const isLoadingComments = ref(true)
const isSubmittingComment = ref(false)

const likeCount = computed(() => props.post.likes + (liked.value ? 1 : 0))
const commentCount = computed(() => Math.max(props.post.comments, comments.value.length))

function toggleLike() {
  liked.value = !liked.value
}

function commentForPreview(comment) {
  return {
    ...comment,
    author: comment.author || 'Orbit member',
    avatarColor: 'var(--gradient-action)',
  }
}

async function loadComments() {
  isLoadingComments.value = true
  commentsError.value = ''

  try {
    const result = await getComments(props.post.id)
    comments.value = (result?.comments || []).map(commentForPreview)
  } catch (error) {
    commentsError.value = error.message || 'Could not load comments.'
  } finally {
    isLoadingComments.value = false
  }
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
  } catch (error) {
    commentsError.value = error.message || 'Could not create comment.'
  } finally {
    isSubmittingComment.value = false
  }
}

function imageUrl(imagePath) {
  if (!imagePath) return ''

  return imagePath.startsWith('/') ? imagePath : `/uploads/${imagePath}`
}

onMounted(loadComments)
</script>

<template>
  <article class="post-card orbit-surface">
    <header class="post-card__header">
      <div class="post-card__avatar" :style="{ background: post.avatarColor }" aria-hidden="true">
        {{ post.author.charAt(0) }}
      </div>

      <div class="post-card__author">
        <h2>{{ post.author }}</h2>
        <p>{{ post.time }} <span aria-hidden="true">•</span> {{ post.privacy }}</p>
      </div>

      <button class="post-card__menu" type="button" :aria-label="`More options for ${post.author}'s post`">
        <span aria-hidden="true">•••</span>
      </button>
    </header>

    <p class="post-card__content">{{ post.content }}</p>

    <div v-if="post.imagePath" class="post-card__media post-card__media--uploaded">
      <img :src="imageUrl(post.imagePath)" alt="Image attached to this post" />
    </div>

    <div v-else-if="post.hasMedia" class="post-card__media" role="img" :aria-label="post.mediaDescription">
      <span class="post-card__sun" aria-hidden="true"></span>
      <span class="post-card__mountain post-card__mountain--back" aria-hidden="true"></span>
      <span class="post-card__mountain post-card__mountain--front" aria-hidden="true"></span>
    </div>

    <p v-if="isLoadingComments" class="comments-state">Loading comments...</p>
    <div v-else-if="commentsError" class="comments-state comments-state--error">
      <span>{{ commentsError }}</span>
      <button type="button" @click="loadComments">Retry</button>
    </div>
    <CommentPreview v-for="comment in comments" v-else :key="comment.id" :comment="comment" />

    <footer class="post-card__actions">
      <button
        class="post-action"
        :class="{ 'post-action--liked': liked }"
        type="button"
        :aria-pressed="liked"
        @click="toggleLike"
      >
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8l1.1 1.1L12 21l7.8-7.5 1.1-1.1a5.5 5.5 0 0 0-.1-7.8Z" />
        </svg>
        <span>{{ likeCount }}</span>
      </button>

      <button class="post-action" type="button" :aria-label="`${commentCount} comments`">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M21 11.5a8.4 8.4 0 0 1-9 8.5 9.8 9.8 0 0 1-3.8-.8L3 21l1.8-4.6A8.4 8.4 0 1 1 21 11.5Z" />
        </svg>
        <span>{{ commentCount }} comments</span>
      </button>
    </footer>

    <CommentInput
      :input-id="`comment-${post.id}`"
      :disabled="isSubmittingComment"
      @submit="addComment"
    />
  </article>
</template>

<style scoped>
.post-card {
  width: 100%;
  padding: var(--space-4);
}

.post-card__header {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr) var(--touch-target);
  align-items: center;
  gap: var(--space-3);
}

.post-card__avatar {
  display: grid;
  width: var(--touch-target);
  height: var(--touch-target);
  place-items: center;
  border-radius: 50%;
  color: #0b0d17;
  font-weight: 700;
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

.post-action svg {
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

.post-action--liked svg {
  fill: currentColor;
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
