<script setup>
import { nextTick, ref, watch } from 'vue'

const props = defineProps({
  messages: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  emptyMessage: { type: String, default: 'No messages yet. Say hello.' },
})
const thread = ref(null)

watch(() => [props.loading, props.messages.length], async () => {
  await nextTick()
  if (thread.value) thread.value.scrollTop = thread.value.scrollHeight
})

function senderName(message) {
  return `${message.firstName || ''} ${message.lastName || ''}`.trim() || message.username || 'Orbit member'
}

function messageTime(value) {
  if (!value) return ''
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' })
}
</script>

<template>
  <div ref="thread" class="message-thread" aria-live="polite">
    <p v-if="loading" class="thread-state">Loading messages...</p>
    <p v-else-if="!messages.length" class="thread-state">{{ emptyMessage }}</p>
    <article v-for="message in messages" v-else :key="message.id" class="message-row" :class="{ 'message-row--own': message.isOwn }">
      <span v-if="!message.isOwn" class="message-avatar" aria-hidden="true">
        <img v-if="message.avatarPath" :src="`/uploads/${message.avatarPath}`" alt="" />
        <span v-else>{{ senderName(message).slice(0, 2).toUpperCase() }}</span>
      </span>
      <div class="message-content">
        <span v-if="!message.isOwn" class="message-sender">{{ senderName(message) }}</span>
        <p>{{ message.content }}</p>
        <time :datetime="message.createdAt">{{ messageTime(message.createdAt) }}</time>
      </div>
    </article>
  </div>
</template>

<style scoped>
.message-thread { display: flex; min-height: 20rem; max-height: 34rem; flex-direction: column; gap: var(--space-3); padding: var(--space-5); overflow-y: auto; background: var(--color-input); }
.thread-state { margin: auto; padding: var(--space-5); color: var(--color-text-muted); text-align: center; }
.message-row { display: flex; max-width: min(78%, 38rem); align-items: end; gap: var(--space-2); }
.message-row--own { align-self: flex-end; }
.message-avatar { display: grid; width: 2rem; height: 2rem; flex: 0 0 2rem; place-items: center; overflow: hidden; border-radius: 50%; background: var(--gradient-action); color: white; font-size: .6875rem; font-weight: 700; }
.message-avatar img { width: 100%; height: 100%; object-fit: cover; }
.message-content { display: grid; min-width: 0; gap: var(--space-1); }
.message-sender { color: var(--color-text-muted); font-size: .75rem; font-weight: 600; }
.message-content p { max-height: 12rem; margin: 0; padding: var(--space-3) var(--space-4); border: 1px solid var(--color-border); border-radius: var(--radius-medium) var(--radius-medium) var(--radius-medium) var(--radius-small); background: var(--color-surface); color: var(--color-text-soft); line-height: 1.55; overflow-y: auto; overflow-wrap: anywhere; overscroll-behavior: contain; scrollbar-color: var(--color-violet) transparent; scrollbar-width: thin; white-space: pre-wrap; }
.message-row--own .message-content p { border-color: rgb(72 217 193 / 28%); border-radius: var(--radius-medium) var(--radius-medium) var(--radius-small) var(--radius-medium); background: var(--color-surface-teal); color: var(--color-text); }
.message-content time { color: var(--color-text-faint); font-family: var(--font-meta); font-size: .6875rem; }
.message-row--own time { text-align: right; }
@media (max-width: 520px) { .message-thread { min-height: 18rem; padding: var(--space-4); } .message-row { max-width: 92%; } }
</style>
