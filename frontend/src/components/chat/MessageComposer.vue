<script setup>
import { ref } from 'vue'
import IconGlyph from '@/components/layout/IconGlyph.vue'

defineProps({ sending: { type: Boolean, default: false } })
const emit = defineEmits(['send'])
const content = ref('')

function submit() {
  if (!content.value.trim()) return
  emit('send', content.value, () => { content.value = '' })
}

function handleKeydown(event) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    submit()
  }
}
</script>

<template>
  <form class="message-composer" @submit.prevent="submit">
    <label class="sr-only" for="chat-message">Message</label>
    <textarea id="chat-message" v-model="content" rows="2" maxlength="2000" placeholder="Write a message..." :disabled="sending" @keydown="handleKeydown" />
    <button type="submit" :disabled="sending || !content.trim()">
      <span>{{ sending ? 'Sending...' : 'Send' }}</span>
      <IconGlyph name="arrowRight" :size="17" />
    </button>
  </form>
</template>

<style scoped>
.sr-only { position: absolute; width: 1px; height: 1px; padding: 0; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }
.message-composer { display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: end; gap: var(--space-3); padding: var(--space-4); border-top: 1px solid var(--color-border); background: var(--color-surface); }
textarea { width: 100%; min-height: var(--touch-target); max-height: 10rem; resize: vertical; padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-small); outline: none; background: var(--color-input); color: var(--color-text); font: inherit; line-height: 1.45; }
textarea:focus { border-color: var(--color-violet); box-shadow: var(--focus-ring); }
button { display: inline-flex; min-height: var(--touch-target); align-items: center; gap: var(--space-2); padding: 0 var(--space-4); border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; cursor: pointer; font-weight: 700; }
@media (max-width: 520px) { .message-composer { grid-template-columns: 1fr; } button { justify-content: center; } }
</style>
