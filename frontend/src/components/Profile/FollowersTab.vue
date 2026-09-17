<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getFriends } from '@/api/common/friends'
import { searchFollowing, searchFollows } from '@/api/users/profiles'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const props = defineProps({
  type: { type: String, default: 'followers' },
  targetId: { type: [String, Number], default: null },
  followerList: { type: Object, default: () => ({}) },
})

const router = useRouter()
const route = useRoute()
const PAGE_SIZE = 20
const showDialog = ref(false)
const list = ref([])
const offset = ref(0)
const loading = ref(false)
const searching = ref(false)
const hasMore = ref(true)
const error = ref('')
const searchResults = ref([])
const searchQuery = ref('')
const scrollBox = ref(null)
let searchDebounce

const targetId = computed(() => props.targetId || route.query.id || '')
const title = computed(() => props.type === 'following' ? 'Following' : props.type === 'friends' ? 'Friends' : 'Followers')
const emptyText = computed(() => `No ${title.value.toLowerCase()} yet.`)
const previewList = computed(() => normalizeUsers(props.followerList).slice(0, 6))
const displayedList = computed(() => searchQuery.value.trim() ? searchResults.value : list.value)

function normalizeUsers(users = {}) {
  return Object.entries(users || {}).map(([id, user]) => ({
    id: Number(id),
    firstName: user.FirstName || user.firstName || '',
    lastName: user.LastName || user.lastName || '',
    username: user.UserName || user.username || '',
    avatar: user.Avatar || user.avatar || '',
  }))
}

function initials(user) {
  return `${user.firstName}${user.lastName}`.slice(0, 2).toUpperCase() || 'O'
}

function buildListUrl() {
  const endpoint = props.type === 'following'
    ? '/api/profile/following'
    : props.type === 'friends' ? '/api/friends/' : '/api/profile/follow'
  const params = new URLSearchParams({ offset: offset.value.toString() })
  if (props.type !== 'friends' && targetId.value) params.set('targetid', targetId.value)
  return `${endpoint}?${params.toString()}`
}

async function fetchPage() {
<<<<<<< HEAD
    if (
        loading.value ||
        !hasMore.value ||
        searchQuery.value.trim()
    ) {
        return;
    }

    loading.value = true;
    error.value = null;

    const endpoint =
        props.type === 'following'
            ? '/api/profile/following'
            : props.type === 'friends'
                ? '/api/friends'
                : '/api/profile/follow';

    try {
        const res = await fetch(
            `${endpoint}?targetid=${targetID}&offset=${offset.value}`,
            {
                method: 'GET',
                credentials: 'include'
            }
        );

        const body = await res.json();

        if (!res.ok || !body.status) {
            throw new Error(body.message || 'Failed to load');
        }

        const page = Object.entries(body.data || {}).map(([key, value]) => ({
            ID: Number(key),
            firstName: value.firstName,
            lastName: value.lastName,
            username: value.username,
            avatar: value.avatar
        }));

        list.value.push(...page);
        offset.value += PAGE_SIZE;

        if (page.length < PAGE_SIZE) {
            hasMore.value = false;
        }
    } catch (err) {
        console.error(err);
        error.value = 'Could not load more.';
    } finally {
        loading.value = false;
    }
=======
  if (loading.value || !hasMore.value || searchQuery.value.trim()) return
  loading.value = true
  error.value = ''
  try {
    const response = await fetch(buildListUrl(), { credentials: 'include' })
    const result = await response.json()
    if (!response.ok || !result.status) throw new Error(result.message || 'Could not load connections')
    const page = normalizeUsers(result.data)
    list.value.push(...page.filter(user => !list.value.some(existing => existing.id === user.id)))
    offset.value += PAGE_SIZE
    hasMore.value = page.length === PAGE_SIZE
  } catch (err) {
    error.value = err.message || 'Could not load connections.'
  } finally {
    loading.value = false
  }
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
}

function openDialog() {
  showDialog.value = true
  list.value = []
  offset.value = 0
  hasMore.value = true
  error.value = ''
  searchQuery.value = ''
  fetchPage()
}

function closeDialog() {
  showDialog.value = false
  searchQuery.value = ''
  searchResults.value = []
}

async function runSearch(query) {
  searching.value = true
  error.value = ''
  try {
    const result = props.type === 'following'
      ? await searchFollowing(query, targetId.value)
      : props.type === 'friends'
        ? await getFriends(query, targetId.value)
        : await searchFollows(query, targetId.value)
    searchResults.value = normalizeUsers(result.data)
  } catch (err) {
    searchResults.value = []
    error.value = err.message || 'Could not search connections.'
  } finally {
    searching.value = false
  }
}

watch(searchQuery, value => {
  clearTimeout(searchDebounce)
  const query = value.trim()
  if (!query) {
    searchResults.value = []
    searching.value = false
    error.value = ''
    return
  }
  searchDebounce = setTimeout(() => runSearch(query), 300)
})

function handleScroll() {
  const box = scrollBox.value
  if (box && box.scrollHeight - box.scrollTop - box.clientHeight < 120) fetchPage()
}

function handleKeydown(event) {
  if (event.key === 'Escape' && showDialog.value) closeDialog()
}

async function goToProfile(id) {
  closeDialog()
  await router.push({ path: '/user', query: { id } })
}

<<<<<<< HEAD
onMounted(() => {
    window.addEventListener('keydown', handleKeydown);
});

=======
onMounted(() => window.addEventListener('keydown', handleKeydown))
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  clearTimeout(searchDebounce)
})
</script>

<template>
  <section class="connections orbit-surface">
    <header class="connections-header">
      <div>
        <p class="orbit-meta">Connections</p>
        <h2>{{ title }}</h2>
      </div>
      <button v-if="previewList.length" type="button" class="text-button" @click="openDialog">Show all</button>
    </header>

    <div v-if="previewList.length" class="preview-grid">
      <button v-for="user in previewList" :key="user.id" type="button" class="person" @click="goToProfile(user.id)">
        <span class="avatar">
          <img v-if="user.avatar" :src="`/uploads/${user.avatar}`" alt="" />
          <span v-else>{{ initials(user) }}</span>
        </span>
        <span class="person-copy">
          <strong>{{ user.firstName }} {{ user.lastName }}</strong>
          <small v-if="user.username">@{{ user.username }}</small>
        </span>
      </button>
    </div>
    <p v-else class="empty-state">{{ emptyText }}</p>
  </section>

  <Teleport to="body">
    <div v-if="showDialog" class="dialog-backdrop" @click.self="closeDialog">
      <section class="dialog" role="dialog" aria-modal="true" :aria-label="title">
        <header class="dialog-header">
          <div><p class="orbit-meta">Connections</p><h2>{{ title }}</h2></div>
          <button type="button" class="icon-button" aria-label="Close" title="Close" @click="closeDialog">
            <IconGlyph name="close" :size="18" />
          </button>
        </header>

        <label class="search-field">
          <IconGlyph name="search" :size="17" />
          <input v-model="searchQuery" type="search" placeholder="Search by name..." />
        </label>

        <div ref="scrollBox" class="dialog-list" @scroll="handleScroll">
          <button v-for="user in displayedList" :key="user.id" type="button" class="person person--row" @click="goToProfile(user.id)">
            <span class="avatar">
              <img v-if="user.avatar" :src="`/uploads/${user.avatar}`" alt="" />
              <span v-else>{{ initials(user) }}</span>
            </span>
            <span class="person-copy">
              <strong>{{ user.firstName }} {{ user.lastName }}</strong>
              <small v-if="user.username">@{{ user.username }}</small>
            </span>
            <IconGlyph name="arrowRight" :size="17" />
          </button>

          <p v-if="loading || searching" class="dialog-status">{{ searching ? 'Searching...' : 'Loading...' }}</p>
          <p v-else-if="error" class="dialog-status dialog-status--error" role="alert">{{ error }}</p>
          <p v-else-if="!displayedList.length" class="dialog-status">{{ searchQuery.trim() ? 'No users found.' : emptyText }}</p>
          <p v-else-if="!searchQuery.trim() && !hasMore" class="dialog-status">You have reached the end.</p>
        </div>
<<<<<<< HEAD

        <div class="followers-card">
            <div
                v-if="previewList.length"
                class="followers-grid"
            >
                <article
                    v-for="follower in previewList"
                    :key="follower.ID"
                    class="follower-card"
                    @click="goToProfile(follower.ID)"
                >
                    <img
                        :src="
                            follower.avatar
                                ? `/uploads/${follower.avatar}`
                                : '/default-avatar.png'
                        "
                        :alt="`${follower.firstName} ${follower.lastName}`"
                        class="follower-avatar"
                    >

                    <div class="follower-info">
                        <p class="follower-name">
                            {{ follower.firstName }}
                            {{ follower.lastName }}
                        </p>
                    </div>
                </article>
            </div>

            <p v-else class="empty">
                {{ emptyText }}
            </p>
        </div>
    </section>

    <Teleport to="body">
        <div
            v-if="showDialog"
            class="dialog-overlay"
            @click="closeDialog"
        >
            <aside
                class="dialog-panel"
                @click.stop
            >
                <header class="dialog-header">
                    <div class="header-title">
                        <span class="header-accent"></span>

                        <div>
                            <p class="eyebrow">SOCIAL</p>
                            <h2>{{ dialogTitle }}</h2>
                        </div>
                    </div>

                    <button
                        type="button"
                        class="close-btn"
                        aria-label="Close"
                        @click="closeDialog"
                    >
                        ×
                    </button>
                </header>

                <div class="search-wrap">
                    <div class="search-box">
                        <span class="search-icon">⌕</span>

                        <input
                            v-model="searchQuery"
                            type="text"
                            placeholder="Search by name..."
                            class="group-search-input"
                        >
                    </div>
                </div>

                <div
                    ref="scrollBox"
                    class="dialog-body"
                    @scroll="throttledScroll"
                >
                    <template v-if="searchQuery.trim()">
                        <article
                            v-for="user in searchResults"
                            :key="user.ID"
                            class="follower-row"
                            @click="goToProfile(user.ID)"
                        >
                            <div class="row-avatar-wrap">
                                <img
                                    :src="
                                        user.avatar
                                            ? `/uploads/${user.avatar}`
                                            : '/default-avatar.png'
                                    "
                                    :alt="`${user.firstName} ${user.lastName}`"
                                    class="follower-avatar"
                                >
                            </div>

                            <div class="row-info">
                                <p class="follower-name">
                                    {{ user.firstName }}
                                    {{ user.lastName }}
                                </p>

                                <span class="row-label">
                                    VIEW PROFILE
                                </span>
                            </div>

                            <span class="row-arrow">↗</span>
                        </article>

                        <p
                            v-if="searching"
                            class="status-text"
                        >
                            Searching…
                        </p>

                        <p
                            v-if="
                                !searching &&
                                !searchResults.length
                            "
                            class="status-text empty-status"
                        >
                            No users found.
                        </p>
                    </template>

                    <template v-else>
                        <article
                            v-for="follower in list"
                            :key="follower.ID"
                            class="follower-row"
                            @click="goToProfile(follower.ID)"
                        >
                            <div class="row-avatar-wrap">
                                <img
                                    :src="
                                        follower.avatar
                                            ? `/uploads/${follower.avatar}`
                                            : '/default-avatar.png'
                                    "
                                    :alt="`${follower.firstName} ${follower.lastName}`"
                                    class="follower-avatar"
                                >
                            </div>

                            <div class="row-info">
                                <p class="follower-name">
                                    {{ follower.firstName }}
                                    {{ follower.lastName }}
                                </p>

                                <span class="row-label">
                                    VIEW PROFILE
                                </span>
                            </div>

                            <span class="row-arrow">↗</span>
                        </article>

                        <p
                            v-if="loading"
                            class="status-text"
                        >
                            Loading…
                        </p>

                        <p
                            v-if="error"
                            class="status-text error"
                        >
                            {{ error }}
                        </p>

                        <p
                            v-if="
                                !hasMore &&
                                !list.length &&
                                !loading
                            "
                            class="status-text empty-status"
                        >
                            No
                            {{
                                type === 'following'
                                    ? 'following'
                                    : type === 'friends'
                                        ? 'friends'
                                        : 'followers'
                            }}
                            yet.
                        </p>

                        <p
                            v-if="
                                !hasMore &&
                                list.length &&
                                !loading
                            "
                            class="status-text end-status"
                        >
                            — That's everyone —
                        </p>
                    </template>
                </div>
            </aside>
        </div>
    </Teleport>
=======
      </section>
    </div>
  </Teleport>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
</template>

<style scoped>
.connections { padding: var(--space-5); }
.connections-header { display: flex; align-items: center; justify-content: space-between; gap: var(--space-4); margin-bottom: var(--space-4); }
.connections h2, .dialog h2 { margin: var(--space-1) 0 0; font-family: var(--font-display); font-size: 1.2rem; letter-spacing: 0; }
.text-button { min-height: var(--touch-target); padding: 0 var(--space-3); border: 0; background: transparent; color: var(--color-violet-soft); cursor: pointer; font-weight: 700; }
.preview-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: var(--space-2); }
.person { display: flex; min-width: 0; min-height: 4.25rem; align-items: center; gap: var(--space-3); padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-input); color: var(--color-text); cursor: pointer; text-align: left; }
.person:hover { border-color: var(--color-violet); background: var(--color-surface-raised); }
.avatar { display: grid; width: 2.75rem; height: 2.75rem; flex: 0 0 2.75rem; place-items: center; overflow: hidden; border-radius: 50%; background: var(--gradient-action); color: white; font-size: .75rem; font-weight: 800; }
.avatar img { width: 100%; height: 100%; object-fit: cover; }
.person-copy { display: grid; min-width: 0; flex: 1; gap: .2rem; }
.person-copy strong, .person-copy small { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.person-copy strong { font-size: .875rem; }
.person-copy small { color: var(--color-text-muted); }
.empty-state, .dialog-status { margin: 0; padding: var(--space-6); color: var(--color-text-muted); text-align: center; }
.dialog-backdrop { position: fixed; inset: 0; display: flex; justify-content: flex-end; background: rgb(0 0 0 / 62%); z-index: 1000; }
.dialog { display: flex; width: min(27rem, 100%); height: 100%; flex-direction: column; border-left: 1px solid var(--color-border); background: var(--color-surface); box-shadow: -1rem 0 3rem rgb(0 0 0 / 25%); }
.dialog-header { display: flex; align-items: center; justify-content: space-between; gap: var(--space-4); padding: var(--space-5); border-bottom: 1px solid var(--color-border); }
.icon-button { display: grid; width: var(--touch-target); height: var(--touch-target); place-items: center; border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-input); color: var(--color-text); cursor: pointer; }
.search-field { display: flex; align-items: center; gap: var(--space-2); margin: var(--space-4); padding: 0 var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-input); color: var(--color-text-muted); }
.search-field:focus-within { border-color: var(--color-violet); }
.search-field input { width: 100%; min-height: var(--touch-target); border: 0; outline: 0; background: transparent; color: var(--color-text); font: inherit; }
.dialog-list { display: flex; min-height: 0; flex: 1; flex-direction: column; gap: var(--space-2); overflow-y: auto; padding: 0 var(--space-4) var(--space-5); }
.person--row { flex: 0 0 auto; width: 100%; }
.person--row > svg { color: var(--color-text-faint); }
.dialog-status--error { color: var(--color-coral-soft); }
@media (max-width: 600px) {
  .connections { padding: var(--space-4); }
  .preview-grid { grid-template-columns: 1fr; }
  .dialog { border-left: 0; }
}
</style>