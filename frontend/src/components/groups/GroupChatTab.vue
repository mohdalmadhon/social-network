<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';

import { addNotification } from '@/data/notifications';
import { sendWS } from '@/api/socket/socket';
import { getMessages } from '@/api/chats/chats';
import { activePage } from '@/data/chatState';

const props = defineProps({
    groupID: {
        type: Number,
        required: true
    },
    userID: {
        type: Number,
        default: null
    },
    userFirstName: {
        type: String,
        default: ''
    },
    userLastName: {
        type: String,
        default: ''
    },
    userAvatar: {
        type: String,
        default: ''
    },
    group: {
        type: Object,
        default: null
    }
});

const message = ref('');
const messages = ref([]);
const sending = ref(false);
const loading = ref(false);
const loadingMore = ref(false);
const hasMore = ref(true);
const offset = ref(0);
const messagesContainer = ref(null);

let fetchTimer = null;
let requestID = 0;

function formatMessage(msg) {
    return {
        id: msg.ID ?? msg.id,
        content: msg.Content ?? msg.content,
        createdAt: msg.CreatedAt ?? msg.createdAt,
        groupID: msg.GroupID ?? msg.groupID,
        sender: {
            id:
                msg.Sender?.ID ??
                msg.Sender?.id ??
                msg.sender?.ID ??
                msg.sender?.id,
            firstName:
                msg.Sender?.FirstName ??
                msg.Sender?.firstName ??
                msg.sender?.FirstName ??
                msg.sender?.firstName,
            lastName:
                msg.Sender?.LastName ??
                msg.Sender?.lastName ??
                msg.sender?.LastName ??
                msg.sender?.lastName,
            avatar:
                msg.Sender?.Avatar ??
                msg.Sender?.avatar ??
                msg.sender?.Avatar ??
                msg.sender?.avatar
        }
    };
}

function isOwnMessage(msg) {
    const senderID = Number(msg.sender?.id);
    const userID = Number(props.userID);

    return senderID === -1 || senderID === userID;
}

async function scrollToBottom() {
    await nextTick();

    if (messagesContainer.value) {
        messagesContainer.value.scrollTop =
            messagesContainer.value.scrollHeight;
    }
}

async function fetchChatMessages(groupID) {
    if (
        groupID === null ||
        groupID === undefined
    ) {
        messages.value = [];
        return;
    }

    const currentRequestID = ++requestID;

    offset.value = 0;
    hasMore.value = true;
    loading.value = true;
    loadingMore.value = false;
    messages.value = [];

    try {
        const result = await getMessages(groupID, 0);

        if (currentRequestID !== requestID) {
            return;
        }

        const data = Array.isArray(result)
            ? result
            : result.messages || result.data || [];

        messages.value = data
            .map(formatMessage)
            .reverse();

        offset.value = data.length;

        if (data.length < 20) {
            hasMore.value = false;
        }
    } catch (err) {
        if (currentRequestID !== requestID) {
            return;
        }

        messages.value = [];
        offset.value = 0;
        hasMore.value = false;

        addNotification(
            err.message ||
                'Error happened while fetching messages',
            'error'
        );
    } finally {
        if (currentRequestID === requestID) {
            loading.value = false;
            await scrollToBottom();
        }
    }
}

async function fetchOlderMessages() {
    if (
        props.groupID === null ||
        props.groupID === undefined ||
        loading.value ||
        loadingMore.value ||
        !hasMore.value
    ) {
        return;
    }

    const container = messagesContainer.value;

    if (!container) {
        return;
    }

    const currentRequestID = requestID;

    loadingMore.value = true;

    const oldScrollHeight =
        container.scrollHeight;

    const oldScrollTop =
        container.scrollTop;

    try {
        const result = await getMessages(
            props.groupID,
            offset.value
        );

        if (currentRequestID !== requestID) {
            return;
        }

        const data = Array.isArray(result)
            ? result
            : result.messages || result.data || [];

        if (data.length === 0) {
            hasMore.value = false;
            return;
        }

        const olderMessages = data
            .map(formatMessage)
            .reverse();

        messages.value = [
            ...olderMessages,
            ...messages.value
        ];

        offset.value += data.length;

        if (data.length < 20) {
            hasMore.value = false;
        }

        await nextTick();

        container.scrollTop =
            oldScrollTop +
            (container.scrollHeight -
                oldScrollHeight);
    } catch (err) {
        if (currentRequestID !== requestID) {
            return;
        }

        addNotification(
            err.message ||
                'Error happened while loading older messages',
            'error'
        );
    } finally {
        if (currentRequestID === requestID) {
            loadingMore.value = false;
        }
    }
}

function throttleFetchOlder() {
    if (fetchTimer) {
        return;
    }

    fetchTimer = setTimeout(() => {
        fetchTimer = null;

        if (
            messagesContainer.value &&
            messagesContainer.value.scrollTop <= 100
        ) {
            fetchOlderMessages();
        }
    }, 200);
}

function handleScroll() {
    const container = messagesContainer.value;

    if (!container) {
        return;
    }

    if (
        container.scrollTop <= 100 &&
        !loading.value &&
        !loadingMore.value &&
        hasMore.value
    ) {
        throttleFetchOlder();
    }
}

function receiveMessage(event) {
    const incoming = event.detail;

    if (!incoming) {
        return;
    }

    const incomingGroupID =
        incoming.GroupID ??
        incoming.groupID;

    if (
        Number(incomingGroupID) !==
        Number(props.groupID)
    ) {
        return;
    }

    const formatted =
        formatMessage(incoming);

    if (
        formatted.id &&
        messages.value.some(
            msg =>
                msg.id === formatted.id
        )
    ) {
        return;
    }

    const senderID =
        Number(formatted.sender.id);

    messages.value.push(formatted);

    nextTick(() => {
        const container =
            messagesContainer.value;

        if (!container) {
            return;
        }

        const distanceFromBottom =
            container.scrollHeight -
            container.scrollTop -
            container.clientHeight;

        const ownMessage =
            senderID ===
            Number(props.userID);

        if (
            ownMessage ||
            distanceFromBottom < 150
        ) {
            container.scrollTop =
                container.scrollHeight;
        }
    });
}

function send() {
    const content = message.value.trim();

    if (
        !content ||
        sending.value ||
        !props.groupID
    ) {
        return;
    }

    sending.value = true;

    try {
        sendWS({
            type: 'privateMessage',
            data: {
                userID: props.userID,
                groupID: props.groupID,
                content,
                private: 1
            }
        });

        messages.value.push({
            id: `local-${Date.now()}`,
            content,
            createdAt:
                new Date().toISOString(),
            groupID: props.groupID,
            sender: {
                id: -1,
                firstName:
                    props.userFirstName,
                lastName:
                    props.userLastName,
                avatar:
                    props.userAvatar
            }
        });

        message.value = '';

        scrollToBottom();
    } catch (err) {
        addNotification(
            err.message ||
                'Error happened while sending message',
            'error'
        );
    } finally {
        sending.value = false;
    }
}

function handleKeydown(event) {
    if (
        event.key === 'Enter' &&
        !event.shiftKey
    ) {
        event.preventDefault();
        send();
    }
}

watch(
    () => props.groupID,
    newGroupID => {
        if (fetchTimer) {
            clearTimeout(fetchTimer);
            fetchTimer = null;
        }

        requestID++;

        messages.value = [];
        offset.value = 0;
        hasMore.value = true;
        loadingMore.value = false;

        if (
            newGroupID === null ||
            newGroupID === undefined
        ) {
            loading.value = false;
            return;
        }

        activePage.value =
            'group:' + newGroupID;

        fetchChatMessages(newGroupID);
    },
    {
        immediate: true
    }
);

watch(
    messagesContainer,
    (newEl, oldEl) => {
        if (oldEl) {
            oldEl.removeEventListener(
                'scroll',
                handleScroll
            );
        }

        if (newEl) {
            newEl.addEventListener(
                'scroll',
                handleScroll
            );
        }
    },
    {
        immediate: true
    }
);

onMounted(() => {
    activePage.value =
        'group:' + props.groupID;

    window.addEventListener(
        'chat-message',
        receiveMessage
    );
});

onUnmounted(() => {
    window.removeEventListener(
        'chat-message',
        receiveMessage
    );

    if (messagesContainer.value) {
        messagesContainer.value.removeEventListener(
            'scroll',
            handleScroll
        );
    }

    if (fetchTimer) {
        clearTimeout(fetchTimer);
        fetchTimer = null;
    }

    requestID++;
});
</script>

<template>
    <section class="chat-window">
        <header
            v-if="group"
            class="chat-window-header"
        >
            <div class="avatar">
                <img
                    v-if="group.avatar || group.Avatar"
                    :src="`/uploads/${group.avatar || group.Avatar}`"
                    alt=""
                />
            </div>

            <div class="chat-user-info">
                <strong>
                    {{ group.name || group.Name }}
                </strong>
            </div>
        </header>

        <div
            ref="messagesContainer"
            class="messages"
        >
            <div
                v-if="loadingMore"
                class="loading-more"
            >
                <div class="small-loader"></div>
                <span>
                    Loading older messages...
                </span>
            </div>

            <div
                v-if="loading"
                class="loading-state"
            >
                <div class="loader"></div>
                <p>Loading messages...</p>
            </div>

            <template v-else>
                <div
                    v-for="(msg, index) in messages"
                    :key="msg.id ?? index"
                    class="message"
                    :class="
                        isOwnMessage(msg)
                            ? 'sent'
                            : 'received'
                    "
                >
                    <div
                        v-if="!isOwnMessage(msg)"
                        class="message-avatar"
                    >
                        <img
                            v-if="msg.sender?.avatar"
                            :src="`/uploads/${msg.sender.avatar}`"
                            alt=""
                        />

                        <span v-else>
                            {{ msg.sender?.firstName?.[0] }}
                            {{ msg.sender?.lastName?.[0] }}
                        </span>
                    </div>

                    <div class="message-body">
                        <span
                            v-if="!isOwnMessage(msg)"
                            class="message-sender-name"
                        >
                            {{ msg.sender?.firstName }}
                            {{ msg.sender?.lastName }}
                        </span>

                        <p>
                            {{ msg.content }}
                        </p>
                    </div>
                </div>

                <div
                    v-if="messages.length === 0"
                    class="no-messages"
                >
                    <p>No messages yet</p>
                </div>
            </template>
        </div>

        <form
            class="composer"
            @submit.prevent="send"
        >
            <input
                v-model="message"
                type="text"
                placeholder="Type a message..."
                :disabled="loading"
                @keydown="handleKeydown"
            />

            <button
                type="submit"
                :disabled="
                    sending ||
                    loading ||
                    !message.trim()
                "
            >
                {{ sending ? 'Sending...' : 'Send' }}
            </button>
        </form>
    </section>
</template>

<style scoped>
.chat-window {
    width: 100%;
    height: calc(100vh - 100px);
    min-height: 0;
    display: flex;
    flex-direction: column;
    border: 2px solid var(--main-color);
    border-radius: 8px;
    background: var(--bg-color);
    box-shadow: 6px 6px var(--main-color);
    overflow: hidden;
    box-sizing: border-box;
}

.chat-window-header {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    gap: 13px;
    padding: 16px 20px;
    border-bottom: 2px solid var(--page-background);
}

.chat-user-info {
    display: flex;
    flex-direction: column;
}

.chat-window-header strong {
    display: block;
    font-size: 14px;
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

.messages {
    position: relative;
    flex: 1 1 0;
    min-height: 0;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
    overflow-y: auto;
    overflow-x: hidden;
    background: var(--page-background);
    box-sizing: border-box;
}

.loading-state {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    background: var(--page-background);
    z-index: 2;
}

.loading-state p {
    margin: 0;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
}

.loader {
    width: 28px;
    height: 28px;
    border: 3px solid var(--main-color);
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
}

.loading-more {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    min-height: 24px;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
}

.small-loader {
    width: 14px;
    height: 14px;
    border: 2px solid var(--main-color);
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.message {
    flex-shrink: 0;
    display: flex;
    align-items: flex-end;
    gap: 8px;
    max-width: 60%;
    font-size: 13px;
    line-height: 1.5;
}

.message-avatar {
    flex-shrink: 0;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--main-color);
    color: white;
    font-size: 10px;
    overflow: hidden;
}

.message-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.message-body {
    min-width: 0;
    padding: 11px 15px;
    border: 2px solid var(--main-color);
    border-radius: 10px;
    overflow-wrap: anywhere;
}

.message-body p {
    margin: 0;
}

.message-sender-name {
    display: block;
    margin-bottom: 3px;
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
    opacity: 0.8;
}

.message.received {
    align-self: flex-start;
}

.message.received .message-body {
    background: var(--input-focus);
    color: white;
    box-shadow: 3px 3px var(--main-color);
}

.message.sent {
    align-self: flex-end;
    flex-direction: row-reverse;
}

.message.sent .message-body {
    background: var(--input-focus);
    color: white;
    box-shadow: 3px 3px var(--main-color);
}

.composer {
    flex: 0 0 auto;
    display: flex;
    gap: 10px;
    padding: 14px 20px;
    border-top: 2px solid var(--page-background);
}

.composer input {
    flex: 1;
    min-width: 0;
    height: 42px;
    padding: 0 14px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--page-background);
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
    box-sizing: border-box;
}

.composer input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.composer button {
    flex-shrink: 0;
    padding: 0 18px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--input-focus);
    box-shadow: 4px 4px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
    font-weight: 600;
}

.composer button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.composer button:active:not(:disabled) {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px var(--main-color);
}

.no-messages {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
}

.no-messages p {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
}

@media (max-width: 800px) {
    .chat-window {
        height: 480px;
        max-height: 480px;
    }

    .message {
        max-width: 80%;
    }
}
</style>