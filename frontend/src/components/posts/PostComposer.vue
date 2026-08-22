<script setup>
import { computed, ref } from 'vue'

const content = ref('')
const privacy = ref('public')
const feeling = ref('')
const selectedFile = ref(null)
const fileInput = ref(null)
const showFeelings = ref(false)
const message = ref('')

const feelings = ['Happy', 'Excited', 'Grateful', 'Thoughtful']

const canPost = computed(() => content.value.trim() !== '' || selectedFile.value !== null)

function selectFile(event) {
  selectedFile.value = event.target.files[0] || null
  message.value = ''
}

function openFilePicker() {
  fileInput.value.click()
}

function removeFile() {
  selectedFile.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function selectFeeling(value) {
  feeling.value = value
  showFeelings.value = false
}

function preparePost() {
  if (!canPost.value) return

  message.value = 'Your draft is ready. Saving it comes with the posts API task.'
}
</script>

<template>
  <form class="post-composer orbit-surface" @submit.prevent="preparePost">
    <div class="post-composer__input-row">
      <div class="post-composer__avatar" aria-hidden="true">N</div>

      <label class="visually-hidden" for="post-content">Post content</label>
      <textarea
        id="post-content"
        v-model="content"
        maxlength="1000"
        placeholder="What's happening in your orbit, Noa?"
        rows="2"
        @input="message = ''"
      ></textarea>
    </div>

    <div v-if="selectedFile" class="selected-file">
      <span>{{ selectedFile.name }}</span>
      <button type="button" aria-label="Remove selected file" @click="removeFile">×</button>
    </div>

    <div class="post-composer__toolbar">
      <div class="post-composer__tools">
        <button class="composer-action composer-action--media" type="button" @click="openFilePicker">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 5.5h16v13H4zM7 15l3-3 2.5 2.5 2-2L18 16M8 9h.01" />
          </svg>
          <span>Photo / GIF</span>
        </button>
        <input
          ref="fileInput"
          class="file-input"
          type="file"
          accept="image/jpeg,image/png,image/gif"
          @change="selectFile"
        />

        <div class="feeling-picker">
          <button
            class="composer-action composer-action--feeling"
            type="button"
            :aria-expanded="showFeelings"
            @click="showFeelings = !showFeelings"
          >
            <span aria-hidden="true">☺</span>
            <span>{{ feeling || 'Feeling' }}</span>
          </button>

          <div v-if="showFeelings" class="feeling-menu">
            <button v-for="option in feelings" :key="option" type="button" @click="selectFeeling(option)">
              {{ option }}
            </button>
          </div>
        </div>
      </div>

      <div class="post-composer__actions">
        <label class="privacy-control">
          <span aria-hidden="true">◉</span>
          <span class="visually-hidden">Post privacy</span>
          <select v-model="privacy">
            <option value="public">Public</option>
            <option value="followers">Followers only</option>
            <option value="private">Selected followers</option>
          </select>
        </label>

        <button class="post-button" type="submit" :disabled="!canPost">Post</button>
      </div>
    </div>

    <p v-if="message" class="post-composer__message" role="status">{{ message }}</p>
  </form>
</template>

<style scoped>
.post-composer {
  width: 100%;
  padding: var(--space-4);
}

.post-composer__input-row {
  display: grid;
  grid-template-columns: var(--touch-target) minmax(0, 1fr);
  align-items: start;
  gap: var(--space-3);
}

.post-composer__avatar {
  display: grid;
  width: var(--touch-target);
  height: var(--touch-target);
  place-items: center;
  border-radius: 50%;
  background: var(--gradient-action);
  color: white;
  font-weight: 700;
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

.composer-action--feeling > span:first-child {
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
    justify-content: flex-end;
  }
}
</style>
