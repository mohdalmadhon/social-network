<script setup>
import { ref } from 'vue';

defineProps({
    comments: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['submit-comment']);

const commentText = ref('');

function formatRelativeDate(dateString) {
    if (!dateString) {
        return '';
    }

    const created = new Date(dateString);

    if (Number.isNaN(created.getTime())) {
        return dateString;
    }

    const diffMs = Date.now() - created.getTime();
    const diffHours = diffMs / (1000 * 60 * 60);
    const diffDays = diffHours / 24;

    if (diffHours < 24) {
        const hours = Math.max(1, Math.floor(diffHours));
        return `${hours}h ago`;
    }

    if (diffDays <= 6) {
        const days = Math.floor(diffDays);
        return `${days}d ago`;
    }

    const weeks = Math.floor(diffDays / 7);
    return `${weeks}w ago`;
}

function submitComment() {
    const text = commentText.value.trim();

    if (!text) {
        return;
    }

    emit('submit-comment', text);
    commentText.value = '';
}
</script>

<template>
    <div class="comments-section">
        <form class="comment-form" @submit.prevent="submitComment">
            <input v-model="commentText" type="text" placeholder="Write a comment...">

            <button type="submit">
                Post
            </button>
        </form>

        <div v-if="comments.length" class="comments-list">
            <div v-for="comment in comments" :key="comment.id" class="comment">
                <div class="comment-avatar">
                    <img v-if="comment.avatarPath" :src="comment.avatarPath" alt="">

                    <span v-else>
                        {{ comment.firstName?.charAt(0) }}
                    </span>
                </div>

                <div class="comment-body">
                    <div class="comment-author">
                        {{ comment.firstName }}
                        {{ comment.lastName }}
                    </div>

                    <div class="comment-content">
                        {{ comment.content }}
                    </div>

                    <div v-if="comment.createdAt" class="comment-date">
                        {{ formatRelativeDate(comment.createdAt) }}
                    </div>
                </div>
            </div>
        </div>

        <p v-else class="no-comments">
            No comments yet.
        </p>
    </div>
</template>

<style scoped>
.comments-section {
    padding: 15px 18px 18px;

    border-top: 2px solid var(--main-color);

    background: #fafafa;
}

.comment-form {
    display: flex;
    gap: 9px;
}

.comment-form input {
    width: 100%;

    padding: 11px 12px;

    border: 2px solid var(--main-color);
    border-radius: 5px;

    outline: none;

    background: var(--bg-color);

    color: var(--main-color);

    font-family: "Hedvig Letters Sans", sans-serif;
    font-size: 12px;

    box-shadow: 3px 3px var(--main-color);
}

.comment-form input:focus {
    border-color: var(--input-focus);
}

.comment-form button {
    flex-shrink: 0;

    padding: 10px 15px;

    border: 2px solid var(--main-color);
    border-radius: 5px;

    background: var(--input-focus);
    box-shadow: 3px 3px var(--main-color);

    color: white;

    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
}

.comment-form button:active {
    transform: translate(2px, 2px);
    box-shadow: 1px 1px var(--main-color);
}

.comments-list {
    margin-top: 20px;

    display: flex;
    flex-direction: column;
    gap: 15px;
}

.comment {
    display: flex;
    align-items: flex-start;
    gap: 10px;
}

.comment-avatar {
    flex-shrink: 0;

    width: 35px;
    height: 35px;

    display: flex;
    align-items: center;
    justify-content: center;

    overflow: hidden;

    border: 2px solid var(--main-color);
    border-radius: 50%;

    background: var(--main-color);

    color: white;

    font-family: "Liter", serif;
    font-size: 13px;
}

.comment-avatar img {
    width: 100%;
    height: 100%;

    object-fit: cover;
}

.comment-body {
    min-width: 0;

    padding: 9px 12px;

    border: 1.5px solid var(--main-color);
    border-radius: 6px;

    background: var(--bg-color);
}

.comment-author {
    margin-bottom: 3px;

    color: var(--main-color);

    font-size: 11px;
    font-weight: 600;
}

.comment-content {
    color: var(--font-color);

    font-size: 12px;
    line-height: 1.4;

    word-break: break-word;
}

.comment-date {
    display: inline-block;
    margin-top: 5px;

    padding: 1px 5px;

    border-radius: 3px;

    background: var(--input-focus);

    color: white;

    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
    font-weight: 600;
}

.no-comments {
    margin: 18px 0 0;

    text-align: center;

    color: var(--font-color-sub);

    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
}

@media (max-width: 650px) {
    .comments-section {
        padding: 12px;
    }
}
</style>