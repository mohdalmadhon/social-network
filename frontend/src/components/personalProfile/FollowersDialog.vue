<script setup>
import { getFriends } from '@/api/common/friends';
import { searchFollowing, searchFollows } from '@/api/users/profiles';
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const props = defineProps({
    type: {
        type: String,
        default: 'followers'
    },
    targetId: {
        type: [String, Number],
        required: true
    }
});

const emit = defineEmits(['close']);
const router = useRouter();

const PAGE_SIZE = 20;
const SCROLL_THROTTLE_MS = 250;

const list = ref([]);
const offset = ref(0);
const loading = ref(false);
const searching = ref(false);
const hasMore = ref(true);
const error = ref(null);
const searchResults = ref([]);
const searchQuery = ref('');
const scrollBox = ref(null);

let searchDebounce;

function throttle(fn, wait) {
    let lastCall = 0;
    let timeoutId = null;

    return function throttled(...args) {
        const now = Date.now();
        const remaining = wait - (now - lastCall);

        if (remaining <= 0) {
            if (timeoutId) {
                clearTimeout(timeoutId);
                timeoutId = null;
            }

            lastCall = now;
            fn.apply(this, args);
        } else if (!timeoutId) {
            timeoutId = setTimeout(() => {
                lastCall = Date.now();
                timeoutId = null;
                fn.apply(this, args);
            }, remaining);
        }
    };
}

const route = useRoute();
const targetID = route.query.id;

async function fetchPage() {
    if (loading.value || !hasMore.value || searchQuery.value.trim()) {
        return;
    }

    loading.value = true;
    error.value = null;

    const endpoint = props.type === 'following'
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
            throw new Error(body.message || 'failed to load');
        }

        const page = Object.entries(body.data).map(([key, value]) => ({
            ID: key,
            FirstName: value.FirstName,
            LastName: value.LastName,
            Avatar: value.Avatar
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
}

function checkAndFetch() {
    if (searchQuery.value.trim()) {
        return;
    }

    const el = scrollBox.value;

    if (!el) {
        return;
    }

    const distanceFromBottom =
        el.scrollHeight - el.scrollTop - el.clientHeight;

    if (distanceFromBottom < 120) {
        fetchPage();
    }
}

const throttledScroll = throttle(checkAndFetch, SCROLL_THROTTLE_MS);

async function runSearch(query) {
    if (!query) {
        searchResults.value = [];
        return;
    }

    searching.value = true;
    error.value = null;

    try {
        let result;
        if (props.type === 'following') {
            result = await searchFollowing(query, props.targetId)
        } else if (props.type === 'friends') {
            result = await getFriends(query, props.targetId)
        } else {
            result = await searchFollows(query, props.targetId);
        }

        searchResults.value = Object.entries(result.data).map(
            ([_, value]) => ({
                UserID: value.ID,
                ...value
            })
        );
        console.log(searchResults.value)

    } catch (err) {
        console.error(err);
        searchResults.value = [];
    } finally {
        searching.value = false;
    }
}

watch(searchQuery, (value) => {
    clearTimeout(searchDebounce);

    const query = value.trim();

    if (!query) {
        searchResults.value = [];
        searching.value = false;
        error.value = null;
        return;
    }

    searchDebounce = setTimeout(() => {
        runSearch(query);
    }, 300);
});

function handleKeydown(e) {
    if (e.key === 'Escape') {
        emit('close');
    }
}

async function goToProfile(id) {
    emit('close');
    await router.replace(`/user?id=${id}`);
    window.location.reload();
}

onMounted(() => {
    fetchPage();
    window.addEventListener('keydown', handleKeydown);
});

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown);
    clearTimeout(searchDebounce);
});
</script>

<template>
    <Teleport to="body">
        <div class="dialog-overlay" @click="$emit('close')">
            <aside class="dialog-panel" @click.stop>
                <header class="dialog-header">
                    <h2>
                        {{ type === 'following' ? 'Following' : 'Followers' }}
                    </h2>

                    <button type="button" class="close-btn" aria-label="Close" @click="$emit('close')">
                        &times;
                    </button>
                </header>

                <input v-model="searchQuery" type="text" placeholder="Search by name..." class="group-search-input" />

                <div ref="scrollBox" class="dialog-body" @scroll="throttledScroll">
                    <template v-if="searchQuery.trim()">
                        <article v-for="user in searchResults" :key="user.ID" class="follower-row"
                            @click="goToProfile(user.ID)">
                            <img :src="user.Avatar
                                ? `/uploads/${user.Avatar}`
                                : '/default-avatar.png'
                                " :alt="`${user.FirstName} ${user.LastName}`" class="follower-avatar" />

                            <p class="follower-name">
                                {{ user.FirstName }} {{ user.LastName }}
                            </p>
                        </article>

                        <p v-if="searching" class="status-text">
                            Searching…
                        </p>

                        <p v-if="
                            !searching &&
                            !searchResults.length
                        " class="status-text">
                            No users found.
                        </p>
                    </template>

                    <template v-else>
                        <article v-for="follower in list" :key="follower.ID" class="follower-row"
                            @click="goToProfile(follower.ID)">
                            <img :src="follower.Avatar
                                ? `/uploads/${follower.Avatar}`
                                : '/default-avatar.png'
                                " :alt="`${follower.FirstName} ${follower.LastName}`" class="follower-avatar" />

                            <p class="follower-name">
                                {{ follower.FirstName }}
                                {{ follower.LastName }}
                            </p>
                        </article>

                        <p v-if="loading" class="status-text">
                            Loading…
                        </p>

                        <p v-if="error" class="status-text error">
                            {{ error }}
                        </p>

                        <p v-if="
                            !hasMore &&
                            !list.length &&
                            !loading
                        " class="status-text">
                            No
                            {{
                                type === 'following'
                                    ? 'following'
                                    : 'followers'
                            }}
                            yet.
                        </p>

                        <p v-if="
                            !hasMore &&
                            list.length &&
                            !loading
                        " class="status-text">
                            That's everyone.
                        </p>
                    </template>
                </div>
            </aside>
        </div>
    </Teleport>
</template>

<style scoped>
.dialog-overlay {
    position: fixed;
    inset: 0;
    background: rgb(11 13 31 / 65%);
    display: flex;
    justify-content: flex-end;
    z-index: 1000;
}

.group-search-input {
    padding: var(--space-3) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 15px;
    max-width: 320px;
    margin-top: var(--space-4);
    margin-left: var(--space-4);
    transition: border-color 0.15s;
}

.group-search-input::placeholder {
    color: var(--color-text-faint);
}

.group-search-input:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
    border-color: var(--color-cyan);
}

.dialog-panel {
    width: min(380px, 92vw);
    height: 100%;
    background: var(--color-surface-violet);
    border-left: 1px solid var(--color-border);
    box-shadow: var(--shadow-raised);
    display: flex;
    flex-direction: column;
    animation: slide-in 0.2s ease-out;
}

@keyframes slide-in {
    from {
        transform: translateX(100%);
    }

    to {
        transform: translateX(0);
    }
}

.dialog-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-5) var(--space-5);
    border-bottom: 1px solid var(--color-border);
    flex: 0 0 auto;
}

.dialog-header h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 22px;
}

.close-btn {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    width: var(--touch-target);
    height: var(--touch-target);
    line-height: 1;
    font-size: 18px;
    cursor: pointer;
    color: var(--color-text);
    transition: transform 0.1s, background 0.15s, color 0.15s;
}

.close-btn:hover {
    transform: translateY(-1px);
    background: var(--color-surface-coral);
    color: var(--color-coral-soft);
}

.close-btn:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.dialog-body {
    flex: 1 1 auto;
    overflow-y: auto;
    padding: var(--space-4) var(--space-4) var(--space-6);
}

.follower-row {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-2);
    border-radius: var(--radius-small);
    cursor: pointer;
    transition: background 0.15s ease;
}

.follower-row:hover {
    background: var(--color-surface-raised);
}

.follower-avatar {
    flex: 0 0 auto;
    width: 44px;
    height: 44px;
    border: 2px solid var(--color-cyan);
    border-radius: 50%;
    object-fit: cover;
}

.follower-name {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 14px;
    font-weight: 600;
    overflow-wrap: anywhere;
}

.status-text {
    margin: 0;
    padding: var(--space-4) 0;
    text-align: center;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 11px;
}

.status-text.error {
    color: var(--color-danger);
}
</style>