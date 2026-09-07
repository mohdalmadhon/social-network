<script setup>
import { ref } from 'vue'

defineProps({
  inputId: {
    type: String,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['submit'])

const content = ref('')
const selectedFile = ref(null)

function chooseFile(event) {
  selectedFile.value = event.target.files[0] || null
}

function removeFile() {
  selectedFile.value = null
}

function submitComment() {
  const cleanContent = content.value.trim()
  if (!cleanContent && !selectedFile.value) return

  emit('submit', { content: cleanContent, file: selectedFile.value })
  content.value = ''
  selectedFile.value = null
}
</script>

<template>
  <form class="comment-input" @submit.prevent="submitComment">
    <label class="visually-hidden" :for="inputId">Write a comment</label>
    <input
      :id="inputId"
      v-model="content"
      :disabled="disabled"
      maxlength="200"
      placeholder="Write a comment..."
      type="text"
    />
    <label class="comment-input__file" :for="`${inputId}-file`" aria-label="Attach an image to your comment">
      <span aria-hidden="true">＋</span>
      <input
        :id="`${inputId}-file`"
        :disabled="disabled"
        accept="image/jpeg,image/png,image/gif"
        type="file"
        @change="chooseFile"
      />
    </label>
    <button type="submit" :disabled="disabled || (!content.trim() && !selectedFile)" aria-label="Send comment">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="m4 12 16-8-6 16-2-6-8-2Zm8 2 3-3" />
      </svg>
    </button>
  </form>
  <div v-if="selectedFile" class="comment-input__selected-file">
    <span>{{ selectedFile.name }}</span>
    <button type="button" aria-label="Remove comment image" @click="removeFile">×</button>
  </div>
</template>

<style scoped>
.comment-input {
  display: grid;
  grid-template-columns: minmax(0, 1fr) var(--touch-target) var(--touch-target);
  align-items: center;
  margin-top: var(--space-3);
  padding-left: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: var(--color-input);
}

.comment-input input {
  width: 100%;
  min-width: 0;
  min-height: var(--touch-target);
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--color-text);
}

.comment-input input::placeholder {
  color: var(--color-text-faint);
}

.comment-input input:disabled {
  cursor: wait;
  opacity: 0.7;
}

.comment-input button {
  display: grid;
  min-width: var(--touch-target);
  min-height: var(--touch-target);
  place-items: center;
  border: 0;
  border-radius: 50%;
  background: transparent;
  color: var(--color-violet);
  cursor: pointer;
}

.comment-input__file {
  display: grid;
  min-width: var(--touch-target);
  min-height: var(--touch-target);
  place-items: center;
  color: var(--color-text-muted);
  cursor: pointer;
}

.comment-input__file:hover {
  color: var(--color-text);
}

.comment-input__file input {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
}

.comment-input__selected-file {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  margin-top: var(--space-2);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-small);
  background: var(--color-input);
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

.comment-input__selected-file span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-input__selected-file button {
  min-width: 2rem;
  min-height: 2rem;
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
}

.comment-input button:disabled {
  color: var(--color-text-faint);
  cursor: not-allowed;
}

.comment-input svg {
  width: 1.25rem;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.8;
}
</style>
