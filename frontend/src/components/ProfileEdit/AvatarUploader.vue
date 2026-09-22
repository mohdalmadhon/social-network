<script setup>
import { onUnmounted, ref, watch } from 'vue'
import { checkSessionResponse } from '@/helpers/auth/auth'
import { addNotification } from '@/data/notifications'
import { router } from '@/router/router'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const props = defineProps({ src: { type: String, default: '' } })
const emit = defineEmits(['change'])
const preview = ref(props.src)
const fileInput = ref(null)
const uploading = ref(false)
const deleting = ref(false)
const error = ref('')
let objectUrl = ''

watch(() => props.src, (value) => { if (!objectUrl) preview.value = value })

function triggerUpload() {
  if (!uploading.value && !deleting.value) fileInput.value?.click()
}

async function onFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return
  error.value = ''
  if (!['image/png', 'image/jpeg', 'image/gif'].includes(file.type)) error.value = 'Only JPG, GIF and PNG images are allowed.'
  else if (file.size > 5 * 1024 * 1024) error.value = 'Image must be smaller than 5MB.'
  if (error.value) {
    event.target.value = ''
    return
  }

  if (objectUrl) URL.revokeObjectURL(objectUrl)
  objectUrl = URL.createObjectURL(file)
  preview.value = objectUrl
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('avatar', file)
    const response = await fetch('/api/profile/avatar', { method: 'PATCH', credentials: 'include', body: formData })
    const result = await response.json()
    if (!checkSessionResponse(response)) {
      router.replace('/login')
      return
    }
    if (!response.ok) throw new Error(result.message || 'Failed to update avatar')
    emit('change', result.avatar)
    addNotification(result.message || 'Avatar updated successfully', 'success')
  } catch (err) {
    preview.value = props.src
    error.value = err.message || 'Failed to update avatar.'
  } finally {
    uploading.value = false
    event.target.value = ''
  }
}

async function deleteAvatar() {
  if (uploading.value || deleting.value || !preview.value) return
  error.value = ''
  deleting.value = true
  try {
    const response = await fetch('/api/profile/avatar', { method: 'DELETE', credentials: 'include' })
    const result = await response.json()
    if (!checkSessionResponse(response)) {
      router.replace('/login')
      return
    }
    if (!response.ok) throw new Error(result.message || 'Failed to delete avatar')
    if (objectUrl) {
      URL.revokeObjectURL(objectUrl)
      objectUrl = ''
    }
    preview.value = ''
    emit('change', null)
    addNotification(result.message || 'Avatar deleted successfully', 'success')
  } catch (err) {
    error.value = err.message || 'Failed to delete avatar.'
  } finally {
    deleting.value = false
  }
}

onUnmounted(() => {
  if (objectUrl) URL.revokeObjectURL(objectUrl)
})
</script>

<template>
  <div class="avatar-uploader">
    <div class="avatar-preview">
      <img v-if="preview" :src="preview" alt="Current profile avatar" />
      <IconGlyph v-else name="profile" :size="30" />
    </div>
    <div class="avatar-actions">
      <div class="avatar-buttons">
        <button type="button" :disabled="uploading || deleting" @click="triggerUpload">
          <IconGlyph name="upload" :size="16" />
          {{ uploading ? 'Uploading...' : 'Change avatar' }}
        </button>
        <button v-if="preview" type="button" class="avatar-delete" :disabled="uploading || deleting"
          @click="deleteAvatar">
          <IconGlyph name="trash" :size="16" />
          {{ deleting ? 'Deleting...' : 'Delete avatar' }}
        </button>
      </div>
      <span>JPG, GIF or PNG, up to 5MB</span>
      <span v-if="error" class="avatar-error" role="alert">{{ error }}</span>
    </div>
    <input ref="fileInput" class="visually-hidden" type="file" accept="image/png, image/jpeg, image/gif"
      @change="onFileChange" />
  </div>
</template>

<style scoped>
.avatar-uploader {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--space-4);
}

.avatar-preview {
  display: grid;
  width: 6rem;
  height: 6rem;
  flex: 0 0 6rem;
  place-items: center;
  overflow: hidden;
  border: 3px solid var(--color-surface);
  outline: 1px solid var(--color-border);
  border-radius: 50%;
  background: var(--color-input);
  color: var(--color-text-muted);
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-actions {
  display: grid;
  justify-items: start;
  gap: var(--space-2);
  color: var(--color-text-faint);
  font-size: .75rem;
}

.avatar-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.avatar-actions button {
  display: inline-flex;
  min-height: var(--touch-target);
  align-items: center;
  gap: var(--space-2);
  padding: 0 var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  background: var(--color-input);
  color: var(--color-text);
  cursor: pointer;
  font-weight: 700;
}

.avatar-actions button:hover:not(:disabled) {
  border-color: var(--color-violet);
}

.avatar-actions button:disabled {
  cursor: not-allowed;
  opacity: .6;
}

.avatar-delete:hover:not(:disabled) {
  border-color: var(--color-coral-soft);
}

.avatar-error {
  color: var(--color-coral-soft);
}
</style>