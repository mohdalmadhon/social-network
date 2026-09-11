<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import { getComments, addComment, deleteComment, voteComment } from '@/api/posts/comments';
import { Comment, createComment } from '@/models/posts';
import { router } from '@/router/router';

const props = defineProps({
    show: { type: Boolean, default: false },
    postId: { type: Number, required: true },
    currentUserId: { type: [Number, String], default: null },
    firstName: { type: String, default: '' },
    lastName: { type: String, default: '' },
    avatarPath: { type: String, default: '' },
    createdAt: { type: String, default: '' },
    content: { type: String, default: '' },
    imagePath: { type: String, default: '' }
});

const emit = defineEmits(['close']);

const comments = ref([]);
const loading = ref(false);
const submitting = ref(false);
const error = ref('');
const newComment = ref('');
const replyingTo = ref(null); 
const loadingReplies = ref({});
const replyErrors = ref({});
const menuComment = ref(null);
const commentInput = ref(null);

const displayName = computed(() => `${props.firstName} ${props.lastName}`.trim());
const replyingToName = computed(() => {
    const user = replyingTo.value?.user;
    return user ? `${user.firstName} ${user.lastName}`.trim() : '';
});
const inputPlaceholder = computed(() =>
    replyingTo.value ? `Reply to ${replyingToName.value}...` : 'Write a comment...'
);

function convertComments(data) {
    return data.map(comment => createComment(comment));
}

function formatDate(date) {
    const created = new Date(date);
    if (Number.isNaN(created.getTime())) return date;

    const diff = Date.now() - created.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (minutes < 1) return 'now';
    if (minutes < 60) return `${minutes}m`;
    if (hours < 24) return `${hours}h`;
    if (days < 7) return `${days}d`;
    return created.toLocaleDateString();
}

async function loadComments() {
    loading.value = true;
    error.value = '';
    try {
        const data = await getComments(props.postId);
        comments.value = convertComments(data);
        console.log(comments.value)
    } catch (err) {
        error.value = err.message || 'Failed to load comments';
    } finally {
        loading.value = false;
    }
}

function createTemporaryComment(content, replyTo = null) {
    const comment = new Comment(
        `temporary-${Date.now()}`, content, props.postId, replyTo, 0,
        new Date().toISOString(),
        { id: Number(props.currentUserId), firstName: 'You', lastName: '', avatarPath: '' },
        0
    );
    comment.pending = true;
    return comment;
}

async function submitTopLevel(content) {
    submitting.value = true;
    error.value = '';
    const temporaryComment = createTemporaryComment(content);
    comments.value.unshift(temporaryComment);

    try {
        const data = await addComment(props.postId, content);
        const realComment = createComment(data);
        const index = comments.value.findIndex(comment => comment.ID === temporaryComment.ID);
        if (index !== -1) comments.value[index] = realComment;
    } catch (err) {
        comments.value = comments.value.filter(comment => comment.ID !== temporaryComment.ID);
        newComment.value = content;
        error.value = err.message || 'Failed to add comment';
    } finally {
        submitting.value = false;
    }
}

async function submitReplyTo(parentComment, content) {
    submitting.value = true;
    replyErrors.value[parentComment.ID] = '';
    const temporaryReply = createTemporaryComment(content, parentComment.ID);

    if (!parentComment.loadedReplies) parentComment.loadedReplies = [];
    parentComment.loadedReplies.unshift(temporaryReply);
    parentComment.replies++;
    parentComment.showReplies = true;

    try {
        const data = await addComment(props.postId, content, parentComment.ID);
        const realReply = createComment(data);
        const index = parentComment.loadedReplies.findIndex(reply => reply.ID === temporaryReply.ID);
        if (index !== -1) parentComment.loadedReplies[index] = realReply;
        replyingTo.value = null;
    } catch (err) {
        parentComment.loadedReplies = parentComment.loadedReplies.filter(reply => reply.ID !== temporaryReply.ID);
        parentComment.replies--;
        newComment.value = content;
        replyErrors.value[parentComment.ID] = err.message || 'Failed to add reply';
    } finally {
        submitting.value = false;
    }
}

async function submitComment() {
    const content = newComment.value.trim();
    if (!content || submitting.value) return;

    newComment.value = '';

    if (replyingTo.value) {
        await submitReplyTo(replyingTo.value, content);
    } else {
        await submitTopLevel(content);
    }
}

async function showReplies(comment) {
    if (comment.loadedReplies !== null) {
        comment.showReplies = !comment.showReplies;
        return;
    }

    loadingReplies.value[comment.ID] = true;
    replyErrors.value[comment.ID] = '';

    try {
        const data = await getComments(props.postId, comment.ID);
        comment.loadedReplies = convertComments(data);
        comment.showReplies = true;
    } catch (err) {
        replyErrors.value[comment.ID] = err.message || 'Failed to load replies';
    } finally {
        loadingReplies.value[comment.ID] = false;
    }
}

function startReply(comment) {
    replyingTo.value = comment;
    newComment.value = '';
    nextTick(() => commentInput.value?.focus());
}

function cancelReply() {
    replyingTo.value = null;
    newComment.value = '';
}

function isOwner(comment) {
    return Number(comment.user?.id) === Number(props.currentUserId);
}

async function removeComment(comment) {
    menuComment.value = null;
    const index = comments.value.findIndex(item => item.ID === comment.ID);
    if (index === -1) return;

    const removedComment = comments.value[index];
    comments.value.splice(index, 1);

    try {
        await deleteComment(comment.ID);
    } catch (err) {
        comments.value.splice(index, 0, removedComment);
        error.value = err.message || 'Failed to delete comment';
    }
}

async function removeReply(parentComment, reply) {
    menuComment.value = null;
    const index = parentComment.loadedReplies.findIndex(item => item.ID === reply.ID);
    if (index === -1) return;

    const removedReply = parentComment.loadedReplies[index];
    parentComment.loadedReplies.splice(index, 1);
    parentComment.replies--;

    try {
        await deleteComment(reply.ID);
    } catch (err) {
        parentComment.loadedReplies.splice(index, 0, removedReply);
        parentComment.replies++;
        replyErrors.value[parentComment.ID] = err.message || 'Failed to delete reply';
    }
}

async function likeComment(comment) {
    const oldVotes = comment.votes;
    comment.votes++;

    try {
        await voteComment(comment.ID, 1);
    } catch (err) {
        comment.votes = oldVotes;
        error.value = err.message || 'Failed to like comment';
    }
}

function close() {
    emit('close');
}

function takeToProfile(id) {
    router.push(`/user?id=${id}`)
    return;
}

watch(() => props.show, value => { if (value) loadComments(); });
</script>

<template>
    <Teleport to="body">
        <div v-if="show" class="comments-overlay" @click.self="close">
            <section class="comments-dialog">

                <header class="comments-header">
                    <h2>Comments</h2>
                    <button type="button" class="close-button" @click="close">×</button>
                </header>

                <div class="comments-content">
                    <article class="dialog-post">
                        <div class="post-user">
                            <img v-if="avatarPath" :src="`/uploads/${avatarPath}`" class="post-avatar">
                            <div v-else class="post-avatar avatar-fallback">{{ firstName.charAt(0) }}</div>
                            <div class="post-user-info">
                                <strong>{{ displayName }}</strong>
                                <span>{{ formatDate(createdAt) }}</span>
                            </div>
                        </div>
                        <div v-if="content" class="post-text">{{ content }}</div>
                        <img v-if="imagePath" :src="`/uploads/${imagePath}`" class="post-image">
                    </article>

                    <div v-if="error" class="comments-error">{{ error }}</div>
                    <div v-if="loading" class="comments-loading">Loading comments...</div>
                    <div v-else-if="!comments.length" class="no-comments">No comments yet.</div>

                    <div v-for="comment in comments" :key="comment.ID" class="comment"
                        :class="{ 'comment-active': replyingTo && replyingTo.ID === comment.ID }">
                        <div class="comment-row">
                            <img style="cursor: pointer;" v-if="comment.user?.avatar" :src="`/uploads/${comment.user.avatar}`"
                                class="comment-avatar" @click="takeToProfile(comment.user.ID)">
                            <div v-else class="comment-avatar avatar-fallback">{{ comment.user?.firstName?.charAt(0) }}
                            </div>
                            <div class="comment-main">

                                <div class="comment-header">
                                    <strong>{{ comment.user?.firstName }} {{ comment.user?.lastName }}</strong>
                                    <button v-if="isOwner(comment)" type="button" class="dots-button"
                                        @click="menuComment = menuComment === comment.ID ? null : comment.ID">⋯</button>
                                </div>

                                <div class="comment-bubble">{{ comment.content }}</div>

                                <div class="comment-actions">
                                    <span>{{ formatDate(comment.createdAt) }}</span>
                                    <button type="button" @click="likeComment(comment)">Like</button>
                                    <span class="vote-count">{{ comment.votes }}</span>
                                    <button type="button" @click="startReply(comment)">Reply</button>
                                    <button v-if="comment.replies > 0" type="button" @click="showReplies(comment)">
                                        {{ loadingReplies[comment.ID] ? 'Loading...' : comment.showReplies ? 'Hide replies' : `Show ${comment.replies} replies` }}
                                    </button>
                                </div>

                                <div v-if="menuComment === comment.ID" class="comment-menu">
                                    <button type="button" @click="removeComment(comment)">Delete</button>
                                </div>

                                <div v-if="replyErrors[comment.ID]" class="reply-error">{{ replyErrors[comment.ID] }}
                                </div>

                                <div v-if="comment.showReplies" class="replies">
                                    <div v-for="reply in comment.loadedReplies" :key="reply.ID" class="comment reply">
                                        <div class="comment-row">
                                            <img v-if="reply.user?.avatarPath"
                                                :src="`/uploads/${reply.user.avatarPath}`" class="comment-avatar">
                                            <div v-else class="comment-avatar avatar-fallback">{{
                                                reply.user?.firstName?.charAt(0) }}</div>
                                            <div class="comment-main">

                                                <div class="comment-header">
                                                    <strong>{{ reply.user?.firstName }} {{ reply.user?.lastName
                                                        }}</strong>
                                                    <button v-if="isOwner(reply)" type="button" class="dots-button"
                                                        @click="menuComment = menuComment === reply.ID ? null : reply.ID">⋯</button>
                                                </div>

                                                <div class="comment-bubble">{{ reply.content }}</div>

                                                <div class="comment-actions">
                                                    <span>{{ formatDate(reply.createdAt) }}</span>
                                                    <button type="button" @click="likeComment(reply)">Like</button>
                                                    <span class="vote-count">{{ reply.votes }}</span>
                                                </div>

                                                <div v-if="menuComment === reply.ID" class="comment-menu">
                                                    <button type="button"
                                                        @click="removeReply(comment, reply)">Delete</button>
                                                </div>

                                            </div>
                                        </div>
                                    </div>
                                </div>

                            </div>
                        </div>
                    </div>

                </div>

                <div v-if="replyingTo" class="reply-banner">
                    <span>Replying to <strong>{{ replyingToName }}</strong></span>
                    <button type="button" @click="cancelReply">✕</button>
                </div>

                <form class="add-comment" @submit.prevent="submitComment">
                    <input ref="commentInput" v-model="newComment" maxlength="200" :placeholder="inputPlaceholder">
                    <button type="submit" :disabled="submitting || !newComment.trim()">{{ replyingTo ? 'Reply' : 'Post'
                        }}</button>
                </form>

            </section>
        </div>
    </Teleport>
</template>

<style scoped>
.comments-dialog {
    --cd-bg: #ffffff;
    --cd-surface: #f4f4f4;
    --cd-border: #000000;
    --cd-text: #111114;
    --cd-text-muted: #6b6b70;
    --cd-accent: #2f6fed;
    --cd-accent-ink: #ffffff;
    --cd-danger: #e5484d;
    --cd-radius: 14px;
    --cd-shadow: 4px 4px 0 var(--cd-border);
    --cd-shadow-sm: 2px 2px 0 var(--cd-border);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.comments-overlay {
    position: fixed;
    inset: 0;
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
    background: rgb(0 0 0 / 55%);
}

.comments-dialog {
    width: 100%;
    max-width: 620px;
    height: min(760px, 90vh);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border: 2px solid var(--cd-border);
    border-radius: var(--cd-radius);
    background: var(--cd-bg);
    box-shadow: 6px 6px 0 var(--cd-border);
}

.comments-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 22px;
    border-bottom: 2px solid var(--cd-border);
}

.comments-header h2 {
    margin: 0;
    color: var(--cd-text);
    font-size: 16px;
    font-weight: 800;
    letter-spacing: -0.01em;
}

.close-button {
    border: 2px solid var(--cd-border);
    width: 30px;
    height: 30px;
    border-radius: 8px;
    background: var(--cd-bg);
    color: var(--cd-text);
    font-size: 18px;
    line-height: 1;
    font-weight: 700;
    cursor: pointer;
    box-shadow: var(--cd-shadow-sm);
    transition: transform 0.1s, box-shadow 0.1s, background 0.1s, color 0.1s;
}

.close-button:hover {
    background: var(--cd-text);
    color: #fff;
    transform: translate(1px, 1px);
    box-shadow: 1px 1px 0 var(--cd-border);
}

.comments-content {
    flex: 1;
    overflow-y: auto;
    padding: 20px 22px;
}

.dialog-post {
    margin-bottom: 18px;
    padding: 14px;
    border: 2px solid var(--cd-border);
    border-radius: var(--cd-radius);
    box-shadow: var(--cd-shadow-sm);
}

.post-user {
    display: flex;
    align-items: center;
    gap: 10px;
}

.post-user-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.post-user-info strong {
    color: var(--cd-text);
    font-size: 13px;
    font-weight: 700;
}

.post-user-info span {
    color: var(--cd-text-muted);
    font-size: 12px;
    font-weight: 600;
}

.post-avatar,
.comment-avatar {
    width: 34px;
    height: 34px;
    flex-shrink: 0;
    border-radius: 50%;
    object-fit: cover;
    background: var(--cd-surface);
    border: 2px solid var(--cd-border);
}

.avatar-fallback {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--cd-accent);
    color: #fff;
    font-weight: 800;
    font-size: 14px;
    text-transform: uppercase;
}

.post-text {
    padding-top: 10px;
    color: var(--cd-text);
    font-size: 14px;
    font-weight: 600;
    line-height: 1.55;
    white-space: pre-wrap;
    word-break: break-word;
}

.post-image {
    width: 100%;
    max-height: 340px;
    margin-top: 12px;
    border: 2px solid var(--cd-border);
    border-radius: 10px;
    object-fit: contain;
    background: var(--cd-surface);
}

.comment {
    position: relative;
    padding: 14px;
    margin-bottom: 12px;
    border: 2px solid var(--cd-border);
    border-radius: var(--cd-radius);
    box-shadow: var(--cd-shadow-sm);
    transition: box-shadow 0.12s, transform 0.12s;
}

.comment.comment-active {
    box-shadow: 3px 3px 0 var(--cd-accent);
    border-color: var(--cd-accent);
}

.comment-row {
    display: flex;
    gap: 10px;
}

.comment-main {
    position: relative;
    min-width: 0;
    flex: 1;
}

.comment-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.comment-header strong {
    color: var(--cd-text);
    font-size: 13px;
    font-weight: 700;
}

.dots-button {
    border: 0;
    padding: 2px 6px;
    border-radius: 6px;
    background: transparent;
    color: var(--cd-text-muted);
    font-size: 16px;
    font-weight: 700;
    line-height: 1;
    cursor: pointer;
}

.dots-button:hover {
    background: var(--cd-surface);
    color: var(--cd-text);
}

.comment-bubble {
    margin-top: 4px;
    color: var(--cd-text);
    font-size: 14px;
    font-weight: 500;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
}

.comment-actions {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: 9px;
    color: var(--cd-text-muted);
    font-size: 12px;
}

.comment-actions .vote-count {
    font-weight: 700;
    color: var(--cd-text);
}

.comment-actions button {
    border: 2px solid var(--cd-border);
    padding: 3px 10px;
    border-radius: 999px;
    background: var(--cd-bg);
    color: var(--cd-text);
    font-size: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.12s, color 0.12s, transform 0.1s;
}

.comment-actions button:hover {
    background: var(--cd-text);
    color: #fff;
}

.comment-actions button:active {
    transform: translate(1px, 1px);
}

.comment-menu {
    position: absolute;
    top: 30px;
    right: 0;
    z-index: 2;
    padding: 4px;
    border: 2px solid var(--cd-border);
    border-radius: 10px;
    background: var(--cd-bg);
    box-shadow: var(--cd-shadow-sm);
}

.comment-menu button {
    border: 0;
    width: 100%;
    padding: 7px 12px;
    border-radius: 6px;
    background: transparent;
    color: var(--cd-danger);
    font-size: 13px;
    font-weight: 700;
    text-align: left;
    cursor: pointer;
}

.comment-menu button:hover {
    background: var(--cd-surface);
}

.replies {
    margin-top: 12px;
    margin-left: 14px;
    padding-left: 14px;
    border-left: 2px solid var(--cd-border);
}

.reply {
    box-shadow: none;
    border-style: dashed;
}

.reply-banner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 0 22px;
    padding: 8px 14px;
    border: 2px solid var(--cd-accent);
    border-bottom: 0;
    border-radius: 12px 12px 0 0;
    background: color-mix(in srgb, var(--cd-accent) 10%, white);
    color: var(--cd-text);
    font-size: 12px;
    font-weight: 600;
}

.reply-banner strong {
    font-weight: 800;
}

.reply-banner button {
    border: 0;
    background: transparent;
    color: var(--cd-text-muted);
    font-size: 14px;
    font-weight: 800;
    cursor: pointer;
}

.reply-banner button:hover {
    color: var(--cd-danger);
}

.add-comment {
    display: flex;
    gap: 10px;
    padding: 16px 22px;
    border-top: 2px solid var(--cd-border);
}

.add-comment input {
    min-width: 0;
    flex: 1;
    border: 2px solid var(--cd-border);
    border-radius: 999px;
    padding: 10px 16px;
    outline: none;
    background: var(--cd-surface);
    color: var(--cd-text);
    font-size: 13px;
    font-weight: 600;
    box-shadow: var(--cd-shadow-sm);
    transition: background 0.12s, box-shadow 0.12s;
}

.add-comment input:focus {
    background: var(--cd-bg);
    box-shadow: 3px 3px 0 var(--cd-accent);
}

.add-comment button {
    border: 2px solid var(--cd-border);
    border-radius: 999px;
    padding: 10px 18px;
    background: var(--cd-text);
    color: #fff;
    font-size: 13px;
    font-weight: 800;
    cursor: pointer;
    box-shadow: var(--cd-shadow-sm);
    transition: transform 0.1s, box-shadow 0.1s, background 0.1s;
}

.add-comment button:hover:not(:disabled) {
    background: var(--cd-accent);
}

.add-comment button:active:not(:disabled) {
    transform: translate(1px, 1px);
    box-shadow: 1px 1px 0 var(--cd-border);
}

.add-comment button:disabled {
    cursor: not-allowed;
    opacity: 0.4;
}

.comments-error,
.reply-error {
    margin-bottom: 14px;
    padding: 9px 12px;
    border: 2px solid var(--cd-danger);
    border-radius: 10px;
    background: color-mix(in srgb, var(--cd-danger) 8%, white);
    color: var(--cd-danger);
    font-size: 13px;
    font-weight: 700;
}

.comments-loading,
.no-comments {
    padding: 40px 0;
    text-align: center;
    color: var(--cd-text-muted);
    font-size: 13px;
    font-weight: 600;
}

@media (max-width: 650px) {
    .comments-overlay {
        padding: 0;
    }

    .comments-dialog {
        width: 100%;
        height: 100%;
        max-width: none;
        border-radius: 0;
        border-width: 0 0 2px 0;
        box-shadow: none;
    }

    .comments-content {
        padding: 16px 16px;
    }

    .reply-banner {
        margin: 0 16px;
    }

    .add-comment {
        padding: 12px 16px;
    }
}
</style>