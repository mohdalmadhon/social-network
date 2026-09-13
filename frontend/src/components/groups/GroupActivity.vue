<script setup>
import { computed, onMounted, ref } from 'vue'
import {
  createGroupEvent,
  deleteGroupEvent,
  getGroupEvents,
  getInviteUsers,
  inviteUserToGroup,
  removeEventRSVP,
  setEventRSVP,
  undoGroupInvitation,
} from '@/api/groups/Groups.js'

const props = defineProps({ groupId: { type: [String, Number], required: true } })
const events = ref([])
const title = ref('')
const description = ref('')
const startsAt = ref('')
const inviteUsers = ref([])
const inviteSearch = ref('')
const inviteBusyUserIDs = ref(new Set())
const eventBusy = ref(false)
const busyEventActions = ref({})
const message = ref('')

const filteredInviteUsers = computed(() => {
  const search = inviteSearch.value.trim().toLowerCase()
  if (!search) return inviteUsers.value

  return inviteUsers.value.filter((user) => {
    const username = user.username?.toLowerCase() || ''
    const firstName = user.firstName?.toLowerCase() || ''
    const lastName = user.lastName?.toLowerCase() || ''
    const fullName = `${firstName} ${lastName}`.trim()
    return username.includes(search) || firstName.includes(search) || lastName.includes(search) || fullName.includes(search)
  })
})

async function load() {
  try {
    const [eventResult, inviteResult] = await Promise.all([
      getGroupEvents(props.groupId),
      getInviteUsers(props.groupId),
    ])
    events.value = eventResult?.events || []
    inviteUsers.value = inviteResult?.users || []
  } catch (error) {
    message.value = error.message
  }
}

async function createEvent() {
  if (eventBusy.value) return
  eventBusy.value = true
  try {
    const result = await createGroupEvent(props.groupId, { title: title.value, description: description.value, startsAt: new Date(startsAt.value).toISOString() })
    title.value = ''; description.value = ''; startsAt.value = ''
    message.value = 'Event created. Other members have been notified.'
    events.value = [...events.value, result.event].sort((a, b) => a.startsAt.localeCompare(b.startsAt) || a.id - b.id)
  } catch (error) { message.value = error.message } finally { eventBusy.value = false }
}

async function invite(user) {
  if (inviteBusyUserIDs.value.has(user.id)) return
  inviteBusyUserIDs.value = new Set(inviteBusyUserIDs.value).add(user.id)
  try {
    const result = await inviteUserToGroup(props.groupId, user.id)
    user.isInvited = true
    user.invitationId = result.invitationId
    message.value = 'Invitation sent.'
  } catch (error) { message.value = error.message } finally {
    const busyUsers = new Set(inviteBusyUserIDs.value)
    busyUsers.delete(user.id)
    inviteBusyUserIDs.value = busyUsers
  }
}

async function undoInvitation(user) {
  if (inviteBusyUserIDs.value.has(user.id) || !user.invitationId) return
  inviteBusyUserIDs.value = new Set(inviteBusyUserIDs.value).add(user.id)
  try {
    await undoGroupInvitation(props.groupId, user.invitationId)
    user.isInvited = false
    user.invitationId = null
    message.value = 'Invitation cancelled.'
  } catch (error) { message.value = error.message } finally {
    const busyUsers = new Set(inviteBusyUserIDs.value)
    busyUsers.delete(user.id)
    inviteBusyUserIDs.value = busyUsers
  }
}

onMounted(load)
async function respond(event, response) {
  if (busyEventActions.value[event.id]) return
  setEventBusy(event.id, 'rsvp')
  try {
    if (event.response === response) {
      const result = await removeEventRSVP(props.groupId, event.id)
      updateEventRSVP(event, '', { id: result.userId })
      message.value = 'Response removed.'
    } else {
      const result = await setEventRSVP(event.id, response)
      updateEventRSVP(event, response, result.user)
      message.value = 'Response saved.'
    }
  } catch (error) { message.value = error.message } finally { setEventBusy(event.id, '') }
}

async function deleteEvent(event) {
  if (busyEventActions.value[event.id]) return
  setEventBusy(event.id, 'delete')
  try {
    await deleteGroupEvent(props.groupId, event.id)
    events.value = events.value.filter((item) => item.id !== event.id)
    message.value = 'Event deleted.'
  } catch (error) { message.value = error.message } finally { setEventBusy(event.id, '') }
}

function updateEventRSVP(event, response, user) {
  event.goingUsers = (event.goingUsers || []).filter((voter) => voter.id !== user.id)
  event.notGoingUsers = (event.notGoingUsers || []).filter((voter) => voter.id !== user.id)
  if (response === 'going') event.goingUsers.push(user)
  if (response === 'declined') event.notGoingUsers.push(user)
  event.goingCount = event.goingUsers.length
  event.notGoingCount = event.notGoingUsers.length
  event.response = response
}

function setEventBusy(eventId, action) {
  const actions = { ...busyEventActions.value }
  if (action) actions[eventId] = action
  else delete actions[eventId]
  busyEventActions.value = actions
}

function voterName(user) {
  return `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.username || 'Group member'
}
</script>

<template>
  <section class="activity orbit-surface">
    <p v-if="message" role="status">{{ message }}</p>
    <section>
      <h2>Invite people</h2>
      <label for="invite-user">Search by username or name</label>
      <input id="invite-user" v-model="inviteSearch" type="search" autocomplete="off" placeholder="Search by username or name..." />
      <div v-if="filteredInviteUsers.length" class="invite-suggestions" aria-label="People to invite">
        <div v-for="user in filteredInviteUsers" :key="user.id" class="invite-candidate">
          <span class="invite-avatar">
            <img v-if="user.avatarPath" :src="`/uploads/${user.avatarPath}`" alt="" />
            <span v-else>{{ `${user.firstName}${user.lastName}`.slice(0, 2).toUpperCase() }}</span>
          </span>
          <span><strong>{{ user.firstName }} {{ user.lastName }}</strong><small>@{{ user.username || 'orbit member' }}</small></span>
          <span v-if="user.isInvited" class="invite-actions">
            <span class="invite-selected">Invited</span>
            <button type="button" :disabled="inviteBusyUserIDs.has(user.id)" @click="undoInvitation(user)">
              {{ inviteBusyUserIDs.has(user.id) ? 'Undoing...' : 'Undo' }}
            </button>
          </span>
          <button v-else type="button" :disabled="inviteBusyUserIDs.has(user.id)" @click="invite(user)">
            {{ inviteBusyUserIDs.has(user.id) ? 'Sending...' : 'Invite' }}
          </button>
        </div>
      </div>
      <p v-else class="invite-hint">No users found.</p>
    </section>
    <form @submit.prevent="createEvent">
      <h2>Plan an event</h2>
      <label for="event-title">Title</label>
      <input id="event-title" v-model="title" maxlength="50" required />
      <label for="event-description">Description</label>
      <textarea id="event-description" v-model="description" maxlength="500" required />
      <label for="event-time">Date and time (your local time)</label>
      <input id="event-time" v-model="startsAt" type="datetime-local" required />
      <button :disabled="eventBusy">{{ eventBusy ? 'Creating...' : 'Create event' }}</button>
    </form>
    <section>
      <h2>Events</h2>
      <p v-if="!events.length">No events yet.</p>
      <article v-for="event in events" :key="event.id">
        <div class="event-heading">
          <h3>{{ event.title }}</h3>
          <button v-if="event.isCreator" type="button" class="delete-event" :disabled="Boolean(busyEventActions[event.id])" @click="deleteEvent(event)">
            {{ busyEventActions[event.id] === 'delete' ? 'Deleting...' : 'Delete event' }}
          </button>
        </div>
        <p>{{ event.description }}</p>
        <time>{{ event.startsAt ? new Date(event.startsAt).toLocaleString() : 'Date not set' }}</time>
        <p v-if="event.response">Your response: {{ event.response === 'going' ? 'Going' : 'Not going' }}</p>
        <div class="responses">
          <button type="button" :disabled="Boolean(busyEventActions[event.id])" :aria-pressed="event.response === 'going'" @click="respond(event, 'going')">Going</button>
          <button type="button" :disabled="Boolean(busyEventActions[event.id])" :aria-pressed="event.response === 'declined'" @click="respond(event, 'declined')">Not going</button>
        </div>
        <div class="event-voters">
          <section>
            <h4>Going ({{ event.goingCount || 0 }})</h4>
            <p v-if="!event.goingUsers?.length" class="event-voters-empty">No responses yet.</p>
            <ul v-else>
              <li v-for="user in event.goingUsers" :key="user.id">
                <span class="event-voter-avatar">
                  <img v-if="user.avatarPath" :src="`/uploads/${user.avatarPath}`" alt="" />
                  <span v-else>{{ voterName(user).slice(0, 2).toUpperCase() }}</span>
                </span>
                <span><strong>{{ voterName(user) }}</strong><small v-if="user.username">@{{ user.username }}</small></span>
              </li>
            </ul>
          </section>
          <section>
            <h4>Not going ({{ event.notGoingCount || 0 }})</h4>
            <p v-if="!event.notGoingUsers?.length" class="event-voters-empty">No responses yet.</p>
            <ul v-else>
              <li v-for="user in event.notGoingUsers" :key="user.id">
                <span class="event-voter-avatar">
                  <img v-if="user.avatarPath" :src="`/uploads/${user.avatarPath}`" alt="" />
                  <span v-else>{{ voterName(user).slice(0, 2).toUpperCase() }}</span>
                </span>
                <span><strong>{{ voterName(user) }}</strong><small v-if="user.username">@{{ user.username }}</small></span>
              </li>
            </ul>
          </section>
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
.invite-suggestions { display: grid; gap: .35rem; max-height: 14rem; overflow: auto; padding: .35rem; border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-surface-raised); }
.invite-candidate { display: grid; grid-template-columns: 2rem minmax(0, 1fr) auto; align-items: center; gap: .6rem; padding: .55rem; }
.invite-candidate > span:nth-child(2) { display: grid; gap: .1rem; }
.invite-suggestions small, .invite-hint, .invite-selected { margin: 0; color: var(--color-text-muted); font-size: .75rem; }
.invite-avatar { display: grid; flex: 0 0 2rem; width: 2rem; height: 2rem; place-items: center; overflow: hidden; border-radius: 50%; background: var(--gradient-action); color: white; font-size: .7rem; font-weight: 700; }
.invite-avatar img { width: 100%; height: 100%; object-fit: cover; }
.invite-actions { display: flex; align-items: center; gap: .5rem; }
.invite-selected { color: var(--color-mint); }
button { justify-self: start; min-height: 44px; margin-top: .5rem; padding: .6rem 1.2rem; border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; cursor: pointer; }
article { padding-block: 1rem; border-top: 1px solid var(--color-border); }
.event-heading { display: flex; align-items: center; justify-content: space-between; gap: var(--space-3); }
.event-heading h3 { min-width: 0; }
.delete-event { border: 1px solid var(--color-coral); background: transparent; color: var(--color-coral); }
.responses { display: flex; flex-wrap: wrap; gap: .75rem; }
.responses button[aria-pressed="true"] { outline: 2px solid var(--color-mint); outline-offset: 2px; }
.event-voters { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: var(--space-4); margin-top: var(--space-4); }
.event-voters h4 { margin: 0 0 var(--space-2); font-size: .9rem; }
.event-voters ul { display: grid; gap: .4rem; margin: 0; padding: 0; list-style: none; }
.event-voters li { display: flex; min-width: 0; align-items: center; gap: .5rem; }
.event-voters li > span:last-child { display: grid; min-width: 0; }
.event-voters small, .event-voters-empty { color: var(--color-text-muted); font-size: .75rem; }
.event-voter-avatar { display: grid; flex: 0 0 1.75rem; width: 1.75rem; height: 1.75rem; place-items: center; overflow: hidden; border-radius: 50%; background: var(--color-input); font-size: .65rem; font-weight: 700; }
.event-voter-avatar img { width: 100%; height: 100%; object-fit: cover; }
@media (max-width: 560px) { .event-voters { grid-template-columns: 1fr; } }
a { display: inline-block; margin-top: .75rem; color: var(--color-mint); }
</style>
