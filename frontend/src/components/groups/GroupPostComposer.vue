<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { createGroupPost } from '@/api/groups/Groups.js'

const props = defineProps({
  groupId: {
    type: [String, Number],
    required: true,
  },
})

const emit = defineEmits(['post-created'])
const MAX_FILE_SIZE = 5 * 1024 * 1024

const content = ref('')
const selectedFile = ref(null)
const previewUrl = ref('')
const fileInput = ref(null)
const isPosting = ref(false)
const error = ref('')

const canPost = computed(() => content.value.trim() !== '' || selectedFile.value !== null)

function chooseImage() {
  fileInput.value?.click()
}

function selectFile(event) {
  clearPreview()
  const file = event.target.files?.[0] || null

  if (file && file.size > MAX_FILE_SIZE) {
    selectedFile.value = null
    event.target.value = ''
    error.value = 'Image must be smaller than 5 MB.'
    return
  }

  selectedFile.value = file
  previewUrl.value = file ? URL.createObjectURL(file) : ''
  error.value = ''
}

function clearPreview() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }
}

function removeFile() {
  selectedFile.value = null
  clearPreview()
  if (fileInput.value) fileInput.value.value = ''
}

async function submitPost() {
  if (!canPost.value || isPosting.value) return

  const formData = new FormData()
  formData.append('content', content.value)
  if (selectedFile.value) formData.append('image', selectedFile.value)

  isPosting.value = true
  error.value = ''
  try {
    const result = await createGroupPost(props.groupId, formData)
    if (!result?.post) throw new Error('Could not create group post')

    content.value = ''
    removeFile()
    emit('post-created', result.post)
  } catch (err) {
    error.value = err.message || 'Could not create group post.'
  } finally {
    isPosting.value = false
  }
}

onBeforeUnmount(clearPreview)
</script>

<template>
  <form class="group-post-composer" @submit.prevent="submitPost">
    <label for="group-post-content">Share with the group</label>
    <textarea
      id="group-post-content"
      v-model="content"
      maxlength="500"
      rows="3"
      placeholder="Write a post..."
      @input="error = ''"
    ></textarea>

    <div v-if="previewUrl" class="group-post-composer__preview">
      <img :src="previewUrl" alt="Preview of the selected image" />
    </div>

    <div class="group-post-composer__actions">
      <button type="button" class="image-button" @click="chooseImage">
        Add image
      </button>
      <input
        ref="fileInput"
        class="file-input"
        type="file"
        accept="image/jpeg,image/png,image/gif"
        @change="selectFile"
      />
      <button v-if="selectedFile" type="button" class="remove-button" @click="removeFile">
        Remove image
      </button>
      <button type="submit" class="publish-button" :disabled="!canPost || isPosting">
        {{ isPosting ? 'Posting...' : 'Post' }}
      </button>
    </div>

    <p v-if="error" class="group-post-composer__error" role="alert">{{ error }}</p>
  </form>
</template>

<style scoped>
.group-post-composer {
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  background: var(--color-surface);
}

.group-post-composer label {
  display: block;
  margin-bottom: var(--space-2);
  color: var(--color-text);
  font-weight: 600;
}

.group-post-composer textarea {
  width: 100%;
  min-height: 6rem;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  outline: none;
  background: var(--color-input);
  color: var(--color-text);
  font: inherit;
  line-height: 1.5;
  resize: vertical;
}

.group-post-composer textarea:focus {
  border-color: var(--color-mint);
}

.group-post-composer__preview {
  margin-top: var(--space-3);
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: var(--color-input);
}

.group-post-composer__preview img {
  display: block;
  width: 100%;
  max-height: 22rem;
  object-fit: contain;
}

.group-post-composer__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-3);
}

.group-post-composer button {
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

.group-post-composer .publish-button {
  margin-left: auto;
  border-color: var(--color-mint);
  color: var(--color-mint);
}

.group-post-composer button:hover:not(:disabled) {
  background: var(--color-input);
  color: var(--color-text);
}

.group-post-composer button:disabled {
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

.group-post-composer__error {
  margin: var(--space-2) 0 0;
  color: var(--color-coral);
  font-size: 0.875rem;
}

@media (max-width: 520px) {
  .group-post-composer__actions {
    flex-wrap: wrap;
  }

  .group-post-composer .publish-button {
    width: 100%;
    margin-left: 0;
  }
}
</style>
