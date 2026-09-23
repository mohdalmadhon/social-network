<script setup>
import { computed, ref, watch } from 'vue'
import { createPost } from '@/api/posts/posts.js'
import { getFollowers } from '@/api/users/profiles.js'
import { normalizePostAudience } from '@/helpers/postAudience.js'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const emit = defineEmits(['post-created'])
const props = defineProps(["avatar"]);

const MAX_FILE_SIZE = 5 * 1024 * 1024
const content = ref('')
// This is the visibility chosen for the new post. The API still receives it
// as `privacy`, but the UI name makes its purpose easier to understand.
const postVisibility = ref('public')
const feeling = ref('')
const selectedFile = ref(null)
const selectedFollowerIds = ref([])
const followers = ref([])
const followersOffset = ref(0)
const hasMoreFollowers = ref(false)
const followersError = ref('')
const isLoadingFollowers = ref(false)
const previewUrl = ref('')
const fileInput = ref(null)
const showFeelings = ref(false)
const message = ref('')
const messageType = ref('success')
const isPosting = ref(false)

const feelings = ['Happy', 'Excited', 'Grateful', 'Thoughtful']

const canPost = computed(() => {
  const hasContent = content.value.trim() !== '' || selectedFile.value !== null
  const hasSelectedFollowers = postVisibility.value !== 'selected' || selectedFollowerIds.value.length > 0
  return hasContent && hasSelectedFollowers && !isLoadingFollowers.value
})

async function loadFollowers({ append = false } = {}) {
  if (isLoadingFollowers.value) return

  isLoadingFollowers.value = true
  followersError.value = ''

  try {
    const offset = append ? followersOffset.value : 0
    const result = await getFollowers('', 20, offset)
    if (!result?.status) {
      throw new Error(result?.message || 'Could not load your followers.')
    }
    const nextFollowers = normalizePostAudience(result?.data)

    followers.value = append
      ? [...followers.value, ...nextFollowers]
      : nextFollowers
    followersOffset.value = offset + nextFollowers.length
    hasMoreFollowers.value = nextFollowers.length === 20
  } catch (error) {
    followersError.value = error.message || 'Could not load your followers.'
  } finally {
    isLoadingFollowers.value = false
  }
}

function loadMoreFollowers() {
  loadFollowers({ append: true })
}

function selectFile(event) {
  const file = event.target.files[0] || null

  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }

  if (file && file.size > MAX_FILE_SIZE) {
    selectedFile.value = null
    message.value = 'File must be smaller than 5 MB.'
    messageType.value = 'error'
    event.target.value = ''
    return
  }

  selectedFile.value = file
  message.value = ''
  messageType.value = 'success'
  previewUrl.value = file ? URL.createObjectURL(file) : ''
}

function openFilePicker() {
  fileInput.value.click()
}

function removeFile() {
  selectedFile.value = null
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }

  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

function selectFeeling(value) {
  feeling.value = value
  showFeelings.value = false
}

async function preparePost() {
  if (!canPost.value) return

  const formData = new FormData()
  formData.append('content', content.value)
  formData.append('privacy', postVisibility.value)
  // Only send a private audience for the selected-followers option. This keeps
  // old checkbox choices from making a later public post fail validation.
  const selectedAudience = postVisibility.value === 'selected'
    ? selectedFollowerIds.value
    : []
  formData.append('selectedFollowerIds', JSON.stringify(selectedAudience))

  if (selectedFile.value) {
    formData.append('image', selectedFile.value)
  }

  isPosting.value = true
  message.value = ''

  try {
    const result = await createPost(formData)

    if (!result?.status) {
      throw new Error(result?.message || 'Could not create post')
    }

    content.value = ''
    feeling.value = ''
    postVisibility.value = 'public'
    selectedFollowerIds.value = []
    removeFile()
    emit('post-created', result.post)
    message.value = 'Post published successfully.'
    messageType.value = 'success'
  } catch (error) {
    message.value = error.message || 'Could not create post'
    messageType.value = 'error'
  } finally {
    isPosting.value = false
  }
}

watch(postVisibility, (value) => {
  if (value !== 'selected') {
    selectedFollowerIds.value = []
    return
  }

  if (followers.value.length === 0 && !followersError.value) {
    loadFollowers()
  }
})
</script>

<template>
  <form class="post-composer orbit-surface" @submit.prevent="preparePost">
    <div class="post-composer__input-row">
      <div class="post-composer__avatar">
        <img v-if="avatar" :src="avatar" alt="Profile avatar" />
        <IconGlyph v-else name="profile" :size="18" />
      </div>

      <label class="visually-hidden" for="post-content">Post content</label>
      <textarea id="post-content" v-model="content" maxlength="500" placeholder="What's happening in your orbit?"
        rows="2" @input="message = ''"></textarea>
    </div>

    <div v-if="selectedFile" class="selected-file">
      <span>{{ selectedFile.name }}</span>
      <button type="button" aria-label="Remove selected file" @click="removeFile">
        <IconGlyph name="close" :size="16" />
      </button>
    </div>

    <div v-if="previewUrl" class="selected-preview">
      <img :src="previewUrl" alt="Preview of the selected media" />
    </div>

    <div class="post-composer__toolbar">
      <div class="post-composer__tools">
        <button class="composer-action composer-action--media" type="button" @click="openFilePicker">
          <IconGlyph name="image" :size="17" />
          <span>Photo / GIF</span>
        </button>
        <input ref="fileInput" class="file-input" type="file" accept="image/jpeg,image/png,image/gif"
          @change="selectFile" />

      </div>

      <div class="post-composer__actions">
        <label class="privacy-control">
          <IconGlyph name="globe" :size="16" />
          <span class="visually-hidden">Post visibility</span>
          <select v-model="postVisibility" aria-label="Post visibility">
            <option value="public">Public</option>
            <option value="followers">Followers only</option>
            <option value="selected">Selected followers</option>
          </select>
        </label>

        <p class="privacy-description" aria-live="polite">
          <template v-if="postVisibility === 'public'">Anyone on Orbit can see this post.</template>
          <template v-else-if="postVisibility === 'followers'">People who follow you can see this post.</template>
          <template v-else>Only the followers you choose can see this post.</template>
        </p>

        <div v-if="postVisibility === 'selected'" class="selected-followers">
          <p class="selected-followers__label">Choose approved followers</p>
          <p v-if="isLoadingFollowers && followers.length === 0" class="selected-followers__state">Loading your followers...</p>
          <div v-else-if="followersError && followers.length === 0" class="selected-followers__state selected-followers__state--error">
            {{ followersError }}
            <button class="load-followers-button" type="button" :disabled="isLoadingFollowers" @click="loadFollowers()">
              Try again
            </button>
          </div>
          <p v-else-if="followers.length === 0" class="selected-followers__state">You have no approved followers yet.
          </p>
          <div v-else class="selected-followers__list">
            <p v-if="followersError" class="selected-followers__state selected-followers__state--error">
              {{ followersError }}
              <button class="load-followers-button" type="button" :disabled="isLoadingFollowers" @click="loadMoreFollowers">
                Try again
              </button>
            </p>
            <label v-for="person in followers" :key="person.id" class="selected-follower">
              <input v-model="selectedFollowerIds" type="checkbox" :value="person.id" />
              <span>{{ person.name }}</span>
            </label>
            <button v-if="hasMoreFollowers" class="load-followers-button" type="button"
              :disabled="isLoadingFollowers" @click="loadMoreFollowers">
              {{ isLoadingFollowers ? 'Loading…' : 'Show more followers' }}
            </button>
          </div>
        </div>

        <button class="post-button" type="submit" :disabled="!canPost || isPosting">
          {{ isPosting ? 'Posting…' : 'Post' }}
        </button>
      </div>
    </div>

    <p v-if="message" class="post-composer__message"
      :class="{ 'post-composer__message--error': messageType === 'error' }" role="status">
      {{ message }}
    </p>
  </form>
</template>

<style scoped>
.post-composer {
  position: sticky;
  top: 4rem;
  z-index: 15;
  align-self: start;
  width: 100%;
  max-height: calc(100vh - 8.5rem);
  max-height: calc(100dvh - 8.5rem);
  overflow-y: auto;
  padding: var(--space-4);
  box-shadow: 0 0.75rem 1.75rem rgb(0 0 0 / 24%);
  overscroll-behavior: contain;
}

.post-composer__input-row {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr);
  align-items: start;
  gap: var(--space-3);
}

.post-composer__avatar {
  width: 2.75rem;
  height: 2.75rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 50%;
  border: 2px solid var(--color-border);
  background: var(--gradient-action);
  color: white;
  box-shadow: var(--shadow-soft);
}

.post-composer__avatar img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

textarea {
  width: 100%;
  min-height: 4.5rem;
  padding: var(--space-2) 0;
  overflow: hidden;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--color-text);
  line-height: 1.5;
  resize: vertical;
}

textarea::placeholder {
  color: var(--color-text-faint);
}

.selected-file {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin: var(--space-3) 0 0 calc(var(--touch-target) + var(--space-3));
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-small);
  background: var(--color-input);
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.selected-preview {
  margin: var(--space-3) 0 0 calc(var(--touch-target) + var(--space-3));
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: var(--color-input);
}

.selected-preview img {
  display: block;
  width: 100%;
  max-height: 18rem;
  object-fit: contain;
}

.selected-file span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.selected-file button {
  min-width: var(--touch-target);
  min-height: var(--touch-target);
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 1.4rem;
}

.post-composer__toolbar {
  display: flex;
  align-items: stretch;
  flex-direction: column;
  gap: var(--space-3);
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border);
}

.post-composer__tools,
.post-composer__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.post-composer__actions {
  justify-content: space-between;
  flex-wrap: wrap;
}

.composer-action,
.privacy-control,
.post-button {
  min-height: var(--touch-target);
}

.composer-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding-inline: var(--space-2);
  border: 0;
  border-radius: var(--radius-small);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 0.875rem;
}

.composer-action:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.composer-action svg {
  width: 1.25rem;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.8;
}

.composer-action--media {
  color: var(--color-mint);
}

.composer-action--media span {
  color: var(--color-text-muted);
}

.file-input {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
}

.composer-action--feeling>span:first-child {
  color: var(--color-amber);
  font-size: 1.4rem;
}

.feeling-picker {
  position: relative;
}

.feeling-menu {
  position: absolute;
  top: calc(100% + var(--space-2));
  left: 0;
  z-index: 5;
  display: grid;
  min-width: 10rem;
  padding: var(--space-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: var(--color-surface-raised);
  box-shadow: var(--shadow-raised);
}

.feeling-menu button {
  min-height: var(--touch-target);
  padding-inline: var(--space-3);
  border: 0;
  border-radius: var(--radius-small);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  text-align: left;
}

.feeling-menu button:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.privacy-control {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding-inline: var(--space-3);
  border: 1px solid var(--color-blue);
  border-radius: 999px;
  color: var(--color-blue);
}

.privacy-control select {
  max-width: 9.5rem;
  border: 0;
  outline: 0;
  background: transparent;
  color: inherit;
  cursor: pointer;
  font-size: 0.875rem;
}

.privacy-control option {
  background: var(--color-surface-raised);
  color: var(--color-text);
}

.privacy-description {
  flex: 1 1 100%;
  margin: 0;
  color: var(--color-text-faint);
  font-size: 0.75rem;
  line-height: 1.4;
}

.selected-followers {
  flex: 1 1 100%;
  display: grid;
  gap: var(--space-2);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: var(--color-input);
}

.selected-followers__list {
  display: grid;
  max-height: 12rem;
  gap: var(--space-1);
  overflow-y: auto;
  overscroll-behavior: contain;
}

.selected-followers__label,
.selected-followers__state {
  margin: 0;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.selected-followers__state--error {
  color: var(--color-coral);
}

.selected-follower {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  min-height: var(--touch-target);
  color: var(--color-text-soft);
  cursor: pointer;
}

.selected-follower input {
  width: 1.1rem;
  height: 1.1rem;
  accent-color: var(--color-violet);
}

.load-followers-button {
  min-height: var(--touch-target);
  border: 0;
  background: transparent;
  color: var(--color-violet-soft);
  cursor: pointer;
  font: inherit;
  text-align: left;
}

.load-followers-button:disabled {
  opacity: 0.6;
  cursor: wait;
}

.post-button {
  min-width: 5rem;
  padding-inline: var(--space-4);
  border: 0;
  border-radius: 999px;
  background: var(--gradient-action);
  color: white;
  cursor: pointer;
  font-weight: 700;
}

.post-button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.post-composer__message {
  margin: var(--space-3) 0 0;
  color: var(--color-mint);
  font-size: 0.875rem;
}

.post-composer__message--error {
  color: var(--color-coral);
}

@media (min-width: 48rem) {
  .post-composer {
    padding: var(--space-5);
  }

  .post-composer__toolbar {
    align-items: center;
    flex-direction: row;
    justify-content: space-between;
  }

  .post-composer__actions {
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .selected-followers {
    order: -1;
  }
}
</style>
