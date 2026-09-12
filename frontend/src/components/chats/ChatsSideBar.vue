<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { getPrivateChatsLists, searchChats } from '@/api/chats/chats';
import { addNotification } from '@/data/notifications';

const emit = defineEmits(['select-chat']);

const chats = ref([]);
const offset = ref(0);
const loading = ref(false);
const hasMore = ref(true);
const searchValue = ref('');
const activeChatId = ref(null);

const listEl = ref(null);

function throttle(fn, wait = 300) {
    let lastCallTime = 0;
    let pendingTimeout = null;

    return function throttled(...args) {
        const now = Date.now();
        const remaining = wait - (now - lastCallTime);

        if (remaining <= 0) {
            if (pendingTimeout) {
                clearTimeout(pendingTimeout);
                pendingTimeout = null;
            }
            lastCallTime = now;
            fn.apply(this, args);
        } else if (!pendingTimeout) {
            pendingTimeout = setTimeout(() => {
                lastCallTime = Date.now();
                pendingTimeout = null;
                fn.apply(this, args);
            }, remaining);
        }
    };
}

function debounce(fn, wait = 300) {
    let pendingTimeout = null;

    return function debounced(...args) {
        if (pendingTimeout) {
            clearTimeout(pendingTimeout);
        }

        pendingTimeout = setTimeout(() => {
            pendingTimeout = null;
            fn.apply(this, args);
        }, wait);
    };
}

function normalizeList(result) {
    if (Array.isArray(result)) return result;
    return result.data || result.chats || result.groups || [];
}

function normalizeHasMore(result, list) {
    if (typeof result?.hasMore === 'boolean') return result.hasMore;
    return list.length > 0;
}

function mergeByUserId(existing, incoming) {
    const merged = new Map(existing.map((chat) => [chat.UserID, chat]));

    for (const chat of incoming) {
        if (!merged.has(chat.UserID)) {
            merged.set(chat.UserID, chat);
        }
    }

    return Array.from(merged.values());
}

async function loadChats({ reset = false } = {}) {
    if (loading.value) return;
    if (!reset && !hasMore.value) return;

    loading.value = true;
    const nextOffset = reset ? 0 : offset.value;

    try {
        const result = searchValue.value
            ? await searchChats(nextOffset, searchValue.value)
            : await getPrivateChatsLists(nextOffset);

        console.log(result)
        const list = normalizeList(result);

        chats.value = mergeByUserId(reset ? [] : chats.value, list);


        offset.value = nextOffset + list.length;
        hasMore.value = normalizeHasMore(result, list);

        if (reset && !activeChatId.value && chats.value.length) {
            selectChat(chats.value[0]);
        }
    } catch (err) {
        addNotification(err.message || "could not load chats", 'error');
        console.error(err);
    } finally {
        loading.value = false;
    }
}

function selectChat(chat) {
    activeChatId.value = chat.UserID;
    emit('select-chat', chat);
}

const handleScroll = throttle(() => {
    const el = listEl.value;
    if (!el) return;

    const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 120;
    if (nearBottom) {
        loadChats();
    }
}, 250);

const handleSearchInput = debounce(() => {
    hasMore.value = true;
    loadChats({ reset: true });
}, 500);

onMounted(() => {
    loadChats({ reset: true });
    listEl.value?.addEventListener('scroll', handleScroll);
});

onBeforeUnmount(() => {
    listEl.value?.removeEventListener('scroll', handleScroll);
});
</script>

<template>
    <aside class="chat-sidebar">
        <div class="sidebar-heading">
            <p class="eyebrow">MESSAGES</p>
            <h2>Chats</h2>
        </div>

        <div class="search">
            <span>⌕</span>
            <input type="text" placeholder="Search friends" v-model="searchValue" @input="handleSearchInput" />
        </div>

        <div class="chat-list" ref="listEl">
            <button v-for="chat in chats" :key="chat.UserID" type="button" class="chat-item"
                :class="{ active: activeChatId === chat.UserID }" @click="selectChat(chat)">
                <div class="avatar">
                    <img v-if="chat.Avatar" :src="`/uploads/${chat.Avatar}`" alt="" />
                    <span v-else>
                        {{ (chat.FirstName || '?').charAt(0).toUpperCase() }}
                    </span>
                </div>

                <div class="chat-info">
                    <div class="chat-info-top">
                        <strong>{{ chat.FirstName + ' ' + chat.LastName }}</strong>
                    </div>
                </div>
            </button>

            <p v-if="loading" class="status-text">Loading...</p>
            <p v-else-if="!chats.length" class="status-text">No chats found</p>
            <p v-else-if="!hasMore" class="status-text">No more chats</p>
        </div>
    </aside>
</template>

<style scoped>
.chat-sidebar {
    display: flex;
    flex-direction: column;
    width: 320px;
    flex-shrink: 0;
    height: calc(100vh - 64px - 40px);
    border: 2px solid var(--main-color);
    border-radius: 8px;
    background: var(--bg-color);
    box-shadow: 6px 6px var(--main-color);
    overflow: hidden;
}

.sidebar-heading {
    padding: 20px 20px 12px;
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: 24px;
}

.search {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 0 16px 15px;
    padding: 0 12px;
    height: 40px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--page-background);
}

.search span {
    flex-shrink: 0;
    font-size: 18px;
    color: var(--font-color-sub);
}

.search input {
    width: 100%;
    border: 0;
    outline: 0;
    background: transparent;
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
}

.chat-list {
    flex: 1;
    overflow-y: auto;
    border-top: 2px solid var(--page-background);
}

.chat-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border: 0;
    border-bottom: 2px solid var(--page-background);
    background: transparent;
    text-align: left;
    cursor: pointer;
}

.chat-item:hover {
    background: var(--page-background);
}

.chat-item.active {
    background: var(--input-focus);
}

.chat-item.active strong,
.chat-item.active .preview,
.chat-item.active .time {
    color: white;
}

.avatar {
    flex-shrink: 0;
    width: 44px;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--main-color);
    color: white;
    font-family: "Liter", serif;
    font-size: 18px;
    overflow: hidden;
}

.avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.chat-info {
    flex: 1;
    min-width: 0;
}

.chat-info-top {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 8px;
}

.chat-info-top strong {
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.time {
    flex-shrink: 0;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
}

.preview {
    margin: 4px 0 0;
    color: var(--font-color-sub);
    font-size: 11px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.badge {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    border: 2px solid var(--main-color);
    border-radius: 999px;
    background: #d9534f;
    color: #fff;
    font-size: 9px;
    font-weight: 700;
}

.status-text {
    padding: 16px;
    margin: 0;
    text-align: center;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
}

@media (max-width: 800px) {
    .chat-sidebar {
        width: 100%;
        height: 320px;
    }
}
</style>