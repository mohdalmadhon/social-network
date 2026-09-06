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

function submitComment() {
  const cleanContent = content.value.trim()
  if (!cleanContent) return

  emit('submit', cleanContent)
  content.value = ''
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
    <button type="submit" :disabled="disabled || !content.trim()" aria-label="Send comment">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="m4 12 16-8-6 16-2-6-8-2Zm8 2 3-3" />
      </svg>
    </button>
  </form>
</template>

<style scoped>
.comment-input {
  display: grid;
  grid-template-columns: minmax(0, 1fr) var(--touch-target);
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
