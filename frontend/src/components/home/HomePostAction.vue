<script setup>
defineProps({
    userReaction: {
        type: String,
        default: ''
    }
});

const emit = defineEmits(['like', 'dislike', 'toggle-comments']);

function handleLike() {
    emit('like');
}

function handleDislike() {
    emit('dislike');
}

function toggleComments() {
    emit('toggle-comments');
}
</script>

<template>
    <div class="post-actions">
        <button class="action-button" :class="{ active: userReaction === 'like' }" type="button"
            @click="handleLike">
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path
                    d="M7 10v10H3V10h4Zm3 10h7.2a2 2 0 0 0 1.9-1.4l2.3-7A2 2 0 0 0 19.5 9H15l.7-3.4A2.2 2.2 0 0 0 13.5 3L9 9v11h1Z" />
            </svg>

            Like
        </button>

        <button class="action-button dislike" :class="{ active: userReaction === 'dislike' }" type="button"
            @click="handleDislike">
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path
                    d="M7 14V4H3v10h4Zm3-10h7.2a2 2 0 0 1 1.9 1.4l2.3 7a2 2 0 0 1-1.9 2.6H15l.7 3.4a2.2 2.2 0 0 1-2.2 2.6L9 15V4h1Z" />
            </svg>

            Dislike
        </button>

        <button class="action-button" type="button" @click="toggleComments">
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M21 11.5a8.4 8.4 0 0 1-9 8.4 9.6 9.6 0 0 1-4-.9L3 21l1.5-4.2A8.5 8.5 0 1 1 21 11.5Z" />
            </svg>

            Comment
        </button>
    </div>
</template>

<style scoped>
.post-actions {
    display: grid;
    grid-template-columns: repeat(3, 1fr);

    padding: 7px;
}

.action-button {
    min-height: 42px;

    display: flex;
    align-items: center;
    justify-content: center;
    gap: 7px;

    border: 2px solid transparent;
    border-radius: 5px;

    background: transparent;

    color: var(--font-color-sub);

    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;

    transition:
        background 0.15s,
        color 0.15s,
        border 0.15s;
}

.action-button:hover {
    background: var(--page-background);
    color: var(--main-color);
}

.action-button svg {
    width: 18px;
    height: 18px;

    fill: none;
    stroke: currentColor;
    stroke-width: 1.8;
    stroke-linecap: round;
    stroke-linejoin: round;
}

.action-button.active {
    border-color: var(--main-color);

    background: var(--input-focus);

    color: white;
}

.action-button.dislike.active {
    background: var(--main-color);
    color: white;
}

@media (max-width: 650px) {
    .post-actions {
        padding: 5px;
    }

    .action-button {
        gap: 4px;

        font-size: 8px;
    }

    .action-button svg {
        width: 16px;
        height: 16px;
    }
}
</style>