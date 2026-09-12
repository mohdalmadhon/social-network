<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import {
  createGroupPostComment,
  deleteGroupPost,
  deleteGroupPostComment,
  getGroupPostComments,
} from '@/api/groups/Groups.js'

const props = defineProps({
  groupId: {
    type: [String, Number],
    required: true,
  },
  post: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['post-deleted'])

const MAX_FILE_SIZE = 5 * 1024 * 1024
const comments = ref([])
const commentsVisible = ref(false)
const commentsLoaded = ref(false)
const commentsLoading = ref(false)
const commentsError = ref('')
const commentContent = ref('')
const commentFile = ref(null)
const commentPreviewUrl = ref('')
const commentFileInput = ref(null)
const isSubmittingComment = ref(false)
const isDeletingPost = ref(false)
const deletingCommentId = ref(null)
const postDeleteError = ref('')
const localCommentCount = ref(props.post.commentCount || 0)

const authorName = computed(() => `${props.post.firstName || ''} ${props.post.lastName || ''}`.trim() || props.post.username || 'Group member')
const canComment = computed(() => commentContent.value.trim() !== '' || commentFile.value !== null)

function assetUrl(path) {
  if (!path) return ''
  return path.startsWith('/') ? path : `/uploads/${path}`
}

function formatDate(value) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleString()
}

async function toggleComments() {
  commentsVisible.value = !commentsVisible.value
  if (!commentsVisible.value || commentsLoaded.value || commentsLoading.value) return

  commentsLoading.value = true
  commentsError.value = ''
  try {
    const result = await getGroupPostComments(props.groupId, props.post.id)
    comments.value = result?.comments || []
    commentsLoaded.value = true
  } catch (err) {
    commentsError.value = err.message || 'Could not load comments.'
  } finally {
    commentsLoading.value = false
  }
}

function selectCommentFile(event) {
  clearCommentPreview()
  const file = event.target.files?.[0] || null
  if (file && file.size > MAX_FILE_SIZE) {
    commentFile.value = null
    event.target.value = ''
    commentsError.value = 'Image must be smaller than 5 MB.'
    return
  }

  commentFile.value = file
  commentPreviewUrl.value = file ? URL.createObjectURL(file) : ''
  commentsError.value = ''
}

function clearCommentPreview() {
  if (commentPreviewUrl.value) {
    URL.revokeObjectURL(commentPreviewUrl.value)
    commentPreviewUrl.value = ''
  }
}

function clearCommentForm() {
  commentContent.value = ''
  removeCommentFile()
}

function removeCommentFile() {
  commentFile.value = null
  clearCommentPreview()
  if (commentFileInput.value) commentFileInput.value.value = ''
}

async function submitComment() {
  if (!canComment.value || isSubmittingComment.value) return

  const formData = new FormData()
  formData.append('content', commentContent.value)
  if (commentFile.value) formData.append('image', commentFile.value)

  isSubmittingComment.value = true
  commentsError.value = ''
  try {
    const result = await createGroupPostComment(props.groupId, props.post.id, formData)
    if (!result?.comment) throw new Error('Could not create comment')

    comments.value.push(result.comment)
    commentsLoaded.value = true
    localCommentCount.value += 1
    clearCommentForm()
  } catch (err) {
    commentsError.value = err.message || 'Could not create comment.'
  } finally {
    isSubmittingComment.value = false
  }
}

async function removePost() {
  if (isDeletingPost.value || !window.confirm('Delete this post?')) return

  isDeletingPost.value = true
  postDeleteError.value = ''
  try {
    await deleteGroupPost(props.groupId, props.post.id)
    emit('post-deleted', props.post.id)
  } catch (err) {
    postDeleteError.value = err.message || 'Could not delete group post.'
  } finally {
    isDeletingPost.value = false
  }
}

async function removeComment(comment) {
  if (deletingCommentId.value !== null || !window.confirm('Delete this comment?')) return

  deletingCommentId.value = comment.id
  commentsError.value = ''
  try {
    await deleteGroupPostComment(props.groupId, props.post.id, comment.id)
    comments.value = comments.value.filter(item => item.id !== comment.id)
    localCommentCount.value = Math.max(0, localCommentCount.value - 1)
  } catch (err) {
    commentsError.value = err.message || 'Could not delete comment.'
  } finally {
    deletingCommentId.value = null
  }
}

onBeforeUnmount(clearCommentPreview)
</script>

<template>
  <article class="group-post-card">
    <header class="group-post-card__header">
      <img v-if="post.avatarPath" :src="assetUrl(post.avatarPath)" :alt="`${authorName}'s avatar`" />
      <span v-else class="group-post-card__avatar" aria-hidden="true">{{ authorName.charAt(0) }}</span>
      <div>
        <h3>{{ authorName }}</h3>
        <p>@{{ post.username }} <span aria-hidden="true">&middot;</span> {{ formatDate(post.createdAt) }}</p>
      </div>
      <button
        v-if="post.isOwner"
        type="button"
        class="delete-post-button"
        :disabled="isDeletingPost"
        @click="removePost"
      >
        {{ isDeletingPost ? 'Deleting...' : 'Delete' }}
      </button>
    </header>

    <p v-if="post.content" class="group-post-card__content">{{ post.content }}</p>
    <img v-if="post.imagePath" class="group-post-card__image" :src="assetUrl(post.imagePath)" alt="Image attached to this group post" />
    <p v-if="postDeleteError" class="comments-error" role="alert">{{ postDeleteError }}</p>

    <button type="button" class="comments-toggle" :aria-expanded="commentsVisible" @click="toggleComments">
      {{ commentsVisible ? 'Hide comments' : `Comments (${localCommentCount})` }}
    </button>

    <section v-if="commentsVisible" class="group-comments">
      <p v-if="commentsLoading" class="comments-state">Loading comments...</p>
      <p v-else-if="commentsLoaded && comments.length === 0" class="comments-state">No comments yet.</p>

      <div v-for="comment in comments" :key="comment.id" class="group-comment">
        <img v-if="comment.avatarPath" :src="assetUrl(comment.avatarPath)" :alt="`${comment.firstName}'s avatar`" />
        <span v-else class="group-comment__avatar" aria-hidden="true">{{ comment.firstName?.charAt(0) }}</span>
        <div class="group-comment__body">
          <div class="group-comment__meta">
            <div>
              <strong>{{ `${comment.firstName || ''} ${comment.lastName || ''}`.trim() || comment.username }}</strong>
              <small>@{{ comment.username }} <span aria-hidden="true">&middot;</span> {{ formatDate(comment.createdAt) }}</small>
            </div>
            <button
              v-if="comment.isOwner"
              type="button"
              :disabled="deletingCommentId !== null"
              @click="removeComment(comment)"
            >
              {{ deletingCommentId === comment.id ? 'Deleting...' : 'Delete' }}
            </button>
          </div>
          <p v-if="comment.content">{{ comment.content }}</p>
          <img v-if="comment.imagePath" class="group-comment__image" :src="assetUrl(comment.imagePath)" alt="Image attached to this comment" />
        </div>
      </div>

      <form class="comment-form" @submit.prevent="submitComment">
        <textarea v-model="commentContent" maxlength="200" rows="2" placeholder="Write a comment..."></textarea>
        <img v-if="commentPreviewUrl" class="comment-form__preview" :src="commentPreviewUrl" alt="Preview of the selected comment image" />
        <div class="comment-form__actions">
          <button type="button" @click="commentFileInput?.click()">Add image</button>
          <input
            ref="commentFileInput"
            class="file-input"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            @change="selectCommentFile"
          />
          <button v-if="commentFile" type="button" @click="removeCommentFile">Remove image</button>
          <button type="submit" :disabled="!canComment || isSubmittingComment">
            {{ isSubmittingComment ? 'Commenting...' : 'Comment' }}
          </button>
        </div>
      </form>
      <p v-if="commentsError" class="comments-error" role="alert">{{ commentsError }}</p>
    </section>
  </article>
</template>

<style scoped>
.group-post-card {
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  background: var(--color-surface);
}

.group-post-card__header,
.group-comment {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr);
  align-items: start;
  gap: var(--space-3);
}

.group-post-card__header {
  grid-template-columns: var(--touch-target) minmax(0, 1fr) auto;
}

.group-post-card__header > img,
.group-post-card__avatar {
  width: var(--touch-target);
  height: var(--touch-target);
  border-radius: 50%;
  object-fit: cover;
}

.group-post-card__avatar,
.group-comment__avatar {
  display: grid;
  place-items: center;
  background: var(--color-input);
  color: var(--color-text);
  font-weight: 700;
}

.group-post-card h3,
.group-post-card__header p {
  margin: 0;
}

.group-post-card h3 {
  color: var(--color-text);
  font-size: 1rem;
}

.delete-post-button,
.group-comment__meta button {
  min-height: var(--touch-target);
  padding: 0 var(--space-3);
  border: 1px solid var(--color-coral);
  border-radius: var(--radius-small);
  background: transparent;
  color: var(--color-coral);
  cursor: pointer;
  font: inherit;
  font-size: 0.8125rem;
  font-weight: 600;
}

.delete-post-button:hover:not(:disabled),
.group-comment__meta button:hover:not(:disabled) {
  background: var(--color-coral);
  color: var(--color-background);
}

.delete-post-button:disabled,
.group-comment__meta button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.group-post-card__header p,
.group-comment small {
  color: var(--color-text-faint);
  font-size: 0.8125rem;
}

.group-post-card__content {
  margin: var(--space-4) 0 0;
  color: var(--color-text-soft);
  line-height: 1.55;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.group-post-card__image {
  display: block;
  width: 100%;
  max-height: 32rem;
  margin-top: var(--space-4);
  border-radius: var(--radius-small);
  background: var(--color-input);
  object-fit: contain;
}

.comments-toggle {
  min-height: var(--touch-target);
  margin-top: var(--space-3);
  padding: 0 var(--space-2);
  border: 0;
  border-radius: var(--radius-small);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font: inherit;
}

.comments-toggle:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.group-comments {
  display: grid;
  gap: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border);
}

.group-comment {
  grid-template-columns: 2.25rem minmax(0, 1fr);
  padding: var(--space-3);
  border-radius: var(--radius-small);
  background: var(--color-input);
}

.group-comment > img,
.group-comment__avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  object-fit: cover;
}

.group-comment strong,
.group-comment small {
  display: block;
}

.group-comment__body {
  min-width: 0;
}

.group-comment__meta {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-3);
}

.group-comment p {
  margin: var(--space-1) 0 0;
  color: var(--color-text-soft);
  font-size: 0.9rem;
  line-height: 1.45;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.group-comment .group-comment__image {
  display: block;
  width: min(100%, 22rem);
  height: auto;
  max-height: 18rem;
  margin-top: var(--space-2);
  border-radius: var(--radius-small);
  object-fit: contain;
}

.comment-form textarea {
  width: 100%;
  min-height: 4rem;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  outline: none;
  background: var(--color-input);
  color: var(--color-text);
  font: inherit;
  resize: vertical;
}

.comment-form textarea:focus {
  border-color: var(--color-mint);
}

.comment-form__preview {
  display: block;
  width: min(100%, 22rem);
  max-height: 18rem;
  margin-top: var(--space-2);
  border-radius: var(--radius-small);
  object-fit: contain;
}

.comment-form__actions {
  display: flex;
  gap: var(--space-2);
  margin-top: var(--space-2);
}

.comment-form button {
  min-height: var(--touch-target);
  padding: 0 var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font: inherit;
  font-weight: 600;
}

.comment-form button:last-child {
  margin-left: auto;
  border-color: var(--color-mint);
  color: var(--color-mint);
}

.comment-form button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.file-input {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
}

.comments-state,
.comments-error {
  margin: 0;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.comments-error {
  color: var(--color-coral);
}

@media (max-width: 520px) {
  .comment-form__actions {
    flex-wrap: wrap;
  }

  .comment-form button:last-child {
    width: 100%;
    margin-left: 0;
  }
}
</style>
