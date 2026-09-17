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
const error = ref('')
let objectUrl = ''

watch(() => props.src, (value) => { if (!objectUrl) preview.value = value })

function triggerUpload() {
  if (!uploading.value) fileInput.value?.click()
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
      <button type="button" :disabled="uploading" @click="triggerUpload">
        <IconGlyph name="upload" :size="16" />
        {{ uploading ? 'Uploading...' : 'Change avatar' }}
      </button>
      <span>JPG, GIF or PNG, up to 5MB</span>
      <span v-if="error" class="avatar-error" role="alert">{{ error }}</span>
    </div>
    <input ref="fileInput" class="visually-hidden" type="file" accept="image/png, image/jpeg, image/gif" @change="onFileChange" />
  </div>
</template>

<style scoped>
<<<<<<< HEAD
.avatar-uploader {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: clamp(12px, 3vw, 20px);
}

.avatar-preview {
    flex-shrink: 0;
    width: clamp(72px, 12vw, 100px);
    height: clamp(72px, 12vw, 100px);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border: 3px solid var(--bg-color);
    outline: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--main-color);
    box-shadow: 4px 4px var(--main-color);
}

.avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.avatar-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
}

.upload-button {
    padding: clamp(8px, 1.5vw, 10px) clamp(12px, 2vw, 16px);
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--input-focus);
    box-shadow: 3px 3px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(8px, 1.4vw, 9px);
    font-weight: 600;
    white-space: nowrap;
    cursor: pointer;
}

.upload-button:hover:not(:disabled) {
    transform: translate(-1px, -1px);
}

.upload-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.hint {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(7px, 1.2vw, 8px);
}

.error {
    color: #d9534f;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(7px, 1.2vw, 8px);
}

.hidden-input {
    display: none;
}

@media (max-width: 400px) {
    .avatar-uploader {
        flex-direction: column;
        align-items: flex-start;
    }
}
=======
.avatar-uploader { display: flex; flex-wrap: wrap; align-items: center; gap: var(--space-4); }
.avatar-preview { display: grid; width: 6rem; height: 6rem; flex: 0 0 6rem; place-items: center; overflow: hidden; border: 3px solid var(--color-surface); outline: 1px solid var(--color-border); border-radius: 50%; background: var(--color-input); color: var(--color-text-muted); }
.avatar-preview img { width: 100%; height: 100%; object-fit: cover; }
.avatar-actions { display: grid; justify-items: start; gap: var(--space-2); color: var(--color-text-faint); font-size: .75rem; }
.avatar-actions button { display: inline-flex; min-height: var(--touch-target); align-items: center; gap: var(--space-2); padding: 0 var(--space-4); border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-input); color: var(--color-text); cursor: pointer; font-weight: 700; }
.avatar-actions button:hover:not(:disabled) { border-color: var(--color-violet); }
.avatar-error { color: var(--color-coral-soft); }
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
</style>
