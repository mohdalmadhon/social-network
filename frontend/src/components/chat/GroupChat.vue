<script setup>
import { onMounted, ref } from 'vue'
import { getGroupMessages, sendGroupMessage } from '@/api/chats.js'
import IconGlyph from '@/components/layout/IconGlyph.vue'
import MessageComposer from './MessageComposer.vue'
import MessageThread from './MessageThread.vue'

const props = defineProps({ groupId: { type: [String, Number], required: true } })
const messages = ref([])
const loading = ref(true)
const sending = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const result = await getGroupMessages(props.groupId)
    messages.value = result?.messages || []
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

async function send(content, clear) {
  if (sending.value) return
  sending.value = true
  error.value = ''
  try {
    const result = await sendGroupMessage(props.groupId, content)
    messages.value.push(result.message)
    clear()
  } catch (err) {
    error.value = err.message
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <section class="group-chat orbit-surface" aria-labelledby="group-chat-heading">
    <header>
      <div class="group-chat__icon"><IconGlyph name="chat" :size="19" /></div>
      <div>
        <p class="orbit-meta">Members</p>
        <h2 id="group-chat-heading">Group chat</h2>
      </div>
    </header>
    <p v-if="error" class="group-chat__error" role="alert">{{ error }}</p>
    <MessageThread :messages="messages" :loading="loading" empty-message="No group messages yet. Start the conversation." />
    <MessageComposer :sending="sending" @send="send" />
  </section>
</template>

<style scoped>
.group-chat { overflow: hidden; }
.group-chat > header { display: grid; grid-template-columns: 2.5rem minmax(0, 1fr); align-items: center; gap: var(--space-3); padding: var(--space-5); }
.group-chat__icon { display: grid; width: 2.5rem; height: 2.5rem; place-items: center; border-radius: 50%; background: var(--color-surface-teal); color: var(--color-mint); }
.group-chat .orbit-meta { margin: 0; }
.group-chat h2 { margin: var(--space-1) 0 0; font-family: var(--font-display); font-size: 1.4rem; letter-spacing: 0; }
.group-chat__error { margin: 0; padding: var(--space-3) var(--space-5); border-top: 1px solid var(--color-border); background: rgb(255 112 112 / 8%); color: var(--color-coral); }
@media (max-width: 520px) { .group-chat > header { padding: var(--space-4); } }
</style>
