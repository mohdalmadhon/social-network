<script setup>
defineProps({
    chat: {
        type: Object,
        default: null
    }
});
</script>

<template>
    <section class="chat-window">
        <template v-if="chat">
            <header class="chat-window-header">
                <div class="avatar">
                    <img
                        v-if="chat.Avatar"
                        :src="`/uploads/${chat.Avatar}`"
                        alt=""
                    />
                </div>
            </header>

            <div class="messages">
            </div>

            <form class="composer" @submit.prevent>
                <input type="text" placeholder="Type a message..." />
                <button type="submit">Send</button>
            </form>
        </template>

        <div v-else class="empty-state">
            <p class="eyebrow">NO CHAT SELECTED</p>
            <h2>Pick a conversation</h2>
            <p class="hint">Choose a chat from the list to start messaging.</p>
        </div>
    </section>
</template>

<style scoped>
.chat-window {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: calc(100vh - 64px - 40px);
    border: 2px solid var(--main-color);
    border-radius: 8px;
    background: var(--bg-color);
    box-shadow: 6px 6px var(--main-color);
    overflow: hidden;
}

.chat-window-header {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 13px;
    padding: 16px 20px;
    border-bottom: 2px solid var(--page-background);
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

.chat-window-header strong {
    display: block;
    font-size: 14px;
}

.status {
    margin: 3px 0 0;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
}

.messages {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
    overflow-y: auto;
    background: var(--page-background);
}

.message {
    max-width: 60%;
    padding: 11px 15px;
    border: 2px solid var(--main-color);
    border-radius: 10px;
    font-size: 13px;
    line-height: 1.5;
}

.message p {
    margin: 0;
}

.message.received {
    align-self: flex-start;
    background: var(--bg-color);
    box-shadow: 3px 3px var(--main-color);
}

.message.sent {
    align-self: flex-end;
    background: var(--input-focus);
    color: white;
    box-shadow: 3px 3px var(--main-color);
}

.composer {
    flex-shrink: 0;
    display: flex;
    gap: 10px;
    padding: 14px 20px;
    border-top: 2px solid var(--page-background);
}

.composer input {
    flex: 1;
    height: 42px;
    padding: 0 14px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--page-background);
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
}

.composer button {
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

.composer button:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px var(--main-color);
}

.empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 20px;
}

.empty-state .eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
    letter-spacing: 2px;
}

.empty-state h2 {
    margin: 0 0 8px;
    font-family: "Liter", serif;
    font-size: 26px;
}

.empty-state .hint {
    margin: 0;
    color: var(--font-color-sub);
    font-size: 12px;
}

@media (max-width: 800px) {
    .chat-window {
        height: 480px;
    }

    .message {
        max-width: 80%;
    }
}
</style>