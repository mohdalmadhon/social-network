<script setup>
import { onMounted, ref } from 'vue'
const props = defineProps({ groupId: { type: [String, Number], required: true } })
const events = ref([])
const title = ref('')
const description = ref('')
const startsAt = ref('')
const invitee = ref('')
const message = ref('')
const busy = ref(false)

async function request(path, data) {
  const response = await fetch(`/api/groups/${props.groupId}/${path}`, {
    method: data ? 'POST' : 'GET', credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: data ? JSON.stringify(data) : undefined,
  })
  const result = await response.json()
  if (!response.ok) throw new Error(result.message || 'Could not complete request.')
  return result
}
async function load() {
  try { events.value = (await request('events')).events } catch (error) { message.value = error.message }
}
async function createEvent() {
  if (busy.value) return
  busy.value = true
  try {
    await request('events', { title: title.value, description: description.value, startsAt: new Date(startsAt.value).toISOString() })
    title.value = ''; description.value = ''; startsAt.value = ''
    message.value = 'Event created. Other members have been notified.'
    await load()
  } catch (error) { message.value = error.message } finally { busy.value = false }
}
async function invite() {
  if (busy.value) return
  busy.value = true
  try {
    await request('invitations', { userId: Number(invitee.value) })
    invitee.value = ''
    message.value = 'Invitation sent.'
  } catch (error) { message.value = error.message } finally { busy.value = false }
}
onMounted(load)
async function respond(event, response) {
  if (busy.value) return
  busy.value = true
  try {
    const result = await fetch(`/api/events/${event.id}/rsvp`, {
      method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ response }),
    })
    if (!result.ok) throw new Error((await result.json()).message || 'Could not save response.')
    event.response = response
    message.value = 'Response saved.'
  } catch (error) { message.value = error.message } finally { busy.value = false }
}
</script>

<template>
  <section class="activity orbit-surface">
    <p v-if="message" role="status">{{ message }}</p>
    <form @submit.prevent="invite">
      <h2>Invite a member</h2>
      <label for="invite-user">User ID (shown in their profile address)</label>
      <input id="invite-user" v-model="invitee" type="number" min="1" required />
      <button :disabled="busy">Send invitation</button>
    </form>
    <form @submit.prevent="createEvent">
      <h2>Plan an event</h2>
      <label for="event-title">Title</label>
      <input id="event-title" v-model="title" maxlength="50" required />
      <label for="event-description">Description</label>
      <textarea id="event-description" v-model="description" maxlength="500" required />
      <label for="event-time">Date and time (your local time)</label>
      <input id="event-time" v-model="startsAt" type="datetime-local" required />
      <button :disabled="busy">Create event</button>
    </form>
    <section>
      <h2>Events</h2>
      <p v-if="!events.length">No events yet.</p>
      <article v-for="event in events" :key="event.id">
        <h3>{{ event.title }}</h3><p>{{ event.description }}</p>
        <time>{{ event.startsAt ? new Date(event.startsAt).toLocaleString() : 'Date not set' }}</time>
        <p v-if="event.response">Your response: {{ event.response === 'going' ? 'Going' : 'Not going' }}</p>
        <div class="responses">
          <button type="button" :disabled="busy" :aria-pressed="event.response === 'going'" @click="respond(event, 'going')">Going</button>
          <button type="button" :disabled="busy" :aria-pressed="event.response === 'declined'" @click="respond(event, 'declined')">Not going</button>
        </div>
      </article>
    </section>
  </section>
</template>

<style scoped>
.activity { padding: var(--space-5); display: grid; gap: var(--space-6); }
form { display: grid; gap: var(--space-2); }
h2, h3, p { margin: 0 0 var(--space-2); }
h2 { font-size: 1.25rem; }
label, time { color: var(--color-text-muted); }
input, textarea { min-width: 0; width: 100%; padding: .75rem; color: var(--color-text); background: var(--color-input); border: 1px solid var(--color-border); border-radius: var(--radius-small); }
button { justify-self: start; min-height: 44px; margin-top: .5rem; padding: .6rem 1.2rem; border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; cursor: pointer; }
article { padding-block: 1rem; border-top: 1px solid var(--color-border); }
.responses { display: flex; flex-wrap: wrap; gap: .75rem; }
.responses button[aria-pressed="true"] { outline: 2px solid var(--color-mint); outline-offset: 2px; }
a { display: inline-block; margin-top: .75rem; color: var(--color-mint); }
</style>
