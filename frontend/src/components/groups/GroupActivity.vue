<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import IconGlyph from '@/components/layout/IconGlyph.vue'
import { getSearchResults } from '@/api/search.js'
const props = defineProps({ groupId: { type: [String, Number], required: true } })
const events = ref([])
const title = ref('')
const description = ref('')
const startsAt = ref('')
const inviteeQuery = ref('')
const inviteeUser = ref(null)
const inviteSuggestions = ref([])
const inviteSearchLoading = ref(false)
const message = ref('')
const busy = ref(false)
let inviteSearchTimer

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
  if (!inviteeUser.value) {
    message.value = 'Choose a person from the search results first.'
    return
  }
  busy.value = true
  try {
    await request('invitations', { userId: inviteeUser.value.id })
    inviteeQuery.value = ''
    inviteeUser.value = null
    inviteSuggestions.value = []
    message.value = 'Invitation sent.'
  } catch (error) { message.value = error.message } finally { busy.value = false }
}

function chooseInvitee(user) {
  inviteeUser.value = user
  inviteeQuery.value = `${user.firstName} ${user.lastName}`.trim()
  inviteSuggestions.value = []
}

watch(inviteeQuery, (value) => {
  const selectedName = inviteeUser.value ? `${inviteeUser.value.firstName} ${inviteeUser.value.lastName}`.trim() : ''
  if (selectedName && value.trim() === selectedName) {
    inviteSuggestions.value = []
    return
  }

  inviteeUser.value = null
  inviteSuggestions.value = []
  clearTimeout(inviteSearchTimer)

  const cleanQuery = value.trim()
  if (cleanQuery.length < 2) return

  inviteSearchTimer = setTimeout(async () => {
    inviteSearchLoading.value = true
    try {
      const result = await getSearchResults(cleanQuery, { excludeGroupId: props.groupId })
      inviteSuggestions.value = result?.users || []
    } catch (error) {
      message.value = error.message || 'Could not search for a member.'
    } finally {
      inviteSearchLoading.value = false
    }
  }, 250)
})

onUnmounted(() => clearTimeout(inviteSearchTimer))
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
      <label for="invite-user">Search by name or username</label>
      <div class="invite-picker">
        <IconGlyph name="search" :size="16" />
        <input id="invite-user" v-model="inviteeQuery" type="search" autocomplete="off" placeholder="Try Mira or @mira" />
      </div>
      <p v-if="inviteSearchLoading" class="invite-hint">Looking for people…</p>
      <div v-else-if="inviteSuggestions.length" class="invite-suggestions" role="listbox" aria-label="People to invite">
        <button v-for="user in inviteSuggestions" :key="user.id" type="button" role="option" @click="chooseInvitee(user)">
          <span class="invite-avatar">
            <img v-if="user.avatarPath" :src="`/uploads/${user.avatarPath}`" alt="" />
            <span v-else>{{ `${user.firstName}${user.lastName}`.slice(0, 2).toUpperCase() }}</span>
          </span>
          <span><strong>{{ user.firstName }} {{ user.lastName }}</strong><small>@{{ user.username || 'orbit member' }}</small></span>
        </button>
      </div>
      <p v-if="inviteeUser" class="invite-selected"><IconGlyph name="check" :size="14" /> {{ inviteeQuery }} selected</p>
      <button :disabled="busy || !inviteeUser">Send invitation</button>
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
.invite-picker { display: flex; align-items: center; gap: .5rem; padding: 0 .75rem; border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-input); color: var(--color-text-faint); }
input, textarea { min-width: 0; width: 100%; padding: .75rem; color: var(--color-text); background: var(--color-input); border: 1px solid var(--color-border); border-radius: var(--radius-small); }
.invite-picker input { padding-inline: 0; border: 0; background: transparent; outline: 0; }
.invite-suggestions { display: grid; gap: .35rem; max-height: 14rem; overflow: auto; padding: .35rem; border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-surface-raised); }
.invite-suggestions button { display: flex; align-items: center; gap: .6rem; width: 100%; margin: 0; padding: .55rem; border: 0; border-radius: var(--radius-small); background: transparent; color: var(--color-text); text-align: left; cursor: pointer; }
.invite-suggestions button:hover { background: var(--color-input); }
.invite-suggestions button > span:last-child { display: grid; gap: .1rem; }
.invite-suggestions small, .invite-hint, .invite-selected { margin: 0; color: var(--color-text-muted); font-size: .75rem; }
.invite-avatar { display: grid; flex: 0 0 2rem; width: 2rem; height: 2rem; place-items: center; overflow: hidden; border-radius: 50%; background: var(--gradient-action); color: white; font-size: .7rem; font-weight: 700; }
.invite-avatar img { width: 100%; height: 100%; object-fit: cover; }
.invite-selected { display: flex; align-items: center; gap: .35rem; color: var(--color-mint); }
button { justify-self: start; min-height: 44px; margin-top: .5rem; padding: .6rem 1.2rem; border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; cursor: pointer; }
article { padding-block: 1rem; border-top: 1px solid var(--color-border); }
.responses { display: flex; flex-wrap: wrap; gap: .75rem; }
.responses button[aria-pressed="true"] { outline: 2px solid var(--color-mint); outline-offset: 2px; }
a { display: inline-block; margin-top: .75rem; color: var(--color-mint); }
</style>
