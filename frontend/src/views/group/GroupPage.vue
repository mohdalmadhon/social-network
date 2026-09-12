<script setup>
import { ref, onMounted } from 'vue'
import { router } from '@/router/router.js'
import { useRoute } from 'vue-router'
import { getUsers } from '@/api/users/users.js'

import {
    getGroup,
    deleteGroupApi,
    inviteUserToGroup,
} from '@/api/groups/Groups.js'

const route = useRoute()
const groupId = route.params.groupId

const group = ref(null)
const users = ref([])

const loading = ref(true)
const error = ref(null)

const inviteLoadingUserId = ref(null)
const inviteMessage = ref('')
const inviteError = ref('')

onMounted(async () => {
    try {
        const groupResult = await getGroup(groupId)

        group.value = groupResult.group

        // Only members are allowed to invite other users
        if (group.value.isMember) {
            const usersResult = await getUsers()

            users.value = usersResult.users
        }
    } catch (err) {
        console.error(err)
        error.value = 'Could not load group'
    } finally {
        loading.value = false
    }
})

async function deleteGroup() {
    try {
        const result = await deleteGroupApi(groupId)

        if (result?.status) {
            router.replace('/groups')
        }
    } catch (err) {
        console.error(err)
    }
}

async function inviteUser(userId) {
    inviteMessage.value = ''
    inviteError.value = ''

    inviteLoadingUserId.value = userId

    try {
        const result = await inviteUserToGroup(
            groupId,
            userId
        )

        if (result?.status) {
            inviteMessage.value = 'Invitation sent successfully'
        }
    } catch (err) {
        console.error(err)

        inviteError.value =
            err.message || 'Could not send invitation'
    } finally {
        inviteLoadingUserId.value = null
    }
}
</script>

<template>
    <main class="group-page">
        <section class="group-container">

            <p v-if="loading">
                Loading group...
            </p>

            <p v-else-if="error">
                {{ error }}
            </p>

            <template v-else-if="group">

                <!-- Group Header -->
                <header class="group-header">

                    <div>
                        <h1>
                            {{ group.title }}
                        </h1>

                        <p>
                            {{ group.description }}
                        </p>
                    </div>

                    <button v-if="group.isCreator" class="delete-button" @click="deleteGroup">
                        Delete Group
                    </button>

                </header>

                <!-- Group Information -->
                <section class="group-content">

                    <p>
                        Members: {{ group.memberCount }}
                    </p>

                    <p v-if="group.isMember" class="member-status">
                        You are a member of this group.
                    </p>

                </section>

                <!-- Invite Section -->
                <section v-if="group.isMember" class="invite-section">

                    <div class="invite-section__header">

                        <h2>
                            Invite people
                        </h2>

                        <p>
                            Invite other users to join this group.
                        </p>

                    </div>

                    <!-- Users List -->
                    <div v-if="users.length" class="invite-users">

                        <div v-for="user in users" :key="user.id" class="invite-user">

                            <div class="invite-user__info">

                                <div class="invite-user__avatar">
                                    {{ user.firstName?.charAt(0) }}
                                </div>

                                <div>

                                    <p class="invite-user__name">
                                        {{ user.firstName }}
                                        {{ user.lastName }}
                                    </p>

                                    <p class="invite-user__username">
                                        @{{ user.username }}
                                    </p>

                                </div>

                            </div>

                            <button @click="inviteUser(user.id)" :disabled="inviteLoadingUserId === user.id">
                                {{
                                    inviteLoadingUserId === user.id
                                        ? 'Sending...'
                                        : 'Invite'
                                }}
                            </button>

                        </div>

                    </div>

                    <p v-else class="no-users">
                        No users found.
                    </p>

                    <!-- Success Message -->
                    <p v-if="inviteMessage" class="invite-message invite-message--success">
                        {{ inviteMessage }}
                    </p>

                    <!-- Error Message -->
                    <p v-if="inviteError" class="invite-message invite-message--error">
                        {{ inviteError }}
                    </p>

                </section>

            </template>

        </section>
    </main>
</template>

<style scoped>
.group-page {
    width: 100%;
    min-height: 100vh;
    padding: var(--space-6);

    background: var(--color-background);
    color: var(--color-text);
}

.group-container {
    width: 100%;
    max-width: 960px;

    margin: 0 auto;
}

/* =========================
   Group Header
   ========================= */

.group-header {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);

    padding: var(--space-6);
    margin-bottom: var(--space-5);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);

    background: var(--color-surface);

    box-shadow: var(--shadow-raised);
}

.group-header h1 {
    margin: 0;

    color: var(--color-text);

    font-family: var(--font-display);
    font-size: clamp(2rem, 5vw, 3rem);
    font-weight: 700;

    line-height: 1.1;
}

.group-header p {
    max-width: 700px;

    margin: 0;

    color: var(--color-text-muted);

    font-family: var(--font-body);
    font-size: 1rem;

    line-height: 1.7;
}

/* =========================
   Delete Button
   ========================= */

.delete-button {
    align-self: flex-start;

    min-height: var(--touch-target);

    padding: 0 var(--space-4);

    border: 1px solid var(--color-coral);
    border-radius: var(--radius-small);

    background: transparent;
    color: var(--color-coral);

    font-family: var(--font-body);
    font-size: 0.9rem;
    font-weight: 600;

    cursor: pointer;

    transition:
        background 0.15s ease,
        color 0.15s ease,
        transform 0.15s ease;
}

.delete-button:hover {
    background: var(--color-coral);
    color: var(--color-background);

    transform: translateY(-1px);
}

.delete-button:active {
    transform: translateY(0);
}

/* =========================
   Group Content
   ========================= */

.group-content {
    display: grid;
    gap: var(--space-4);

    padding: var(--space-5);
    margin-bottom: var(--space-5);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);

    background: var(--color-surface);
}

.group-content p {
    margin: 0;

    color: var(--color-text-soft);

    font-family: var(--font-body);
    font-size: 0.95rem;

    line-height: 1.6;
}

.group-content p:first-child {
    color: var(--color-text);
    font-weight: 600;
}

.member-status {
    display: inline-flex;
    align-items: center;

    width: fit-content;

    padding: var(--space-2) var(--space-3);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);

    background: var(--color-surface-raised);

    color: var(--color-mint) !important;

    font-size: 0.875rem;
    font-weight: 600;
}

/* =========================
   Invite Section
   ========================= */

.invite-section {
    padding: var(--space-5);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);

    background: var(--color-surface);
}

.invite-section__header h2 {
    margin: 0;

    color: var(--color-text);

    font-family: var(--font-display);
    font-size: 1.3rem;
}

.invite-section__header p {
    margin: var(--space-2) 0 0;

    color: var(--color-text-muted);

    font-family: var(--font-body);
    font-size: 0.9rem;
}

/* =========================
   Users List
   ========================= */

.invite-users {
    display: grid;
    gap: var(--space-3);

    margin-top: var(--space-4);
}

.invite-user {
    display: flex;
    align-items: center;
    justify-content: space-between;

    gap: var(--space-4);

    padding: var(--space-3);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);

    background: var(--color-surface-raised);
}

.invite-user__info {
    display: flex;
    align-items: center;

    gap: var(--space-3);
}

/* =========================
   Avatar
   ========================= */

.invite-user__avatar {
    display: flex;
    align-items: center;
    justify-content: center;

    width: 42px;
    height: 42px;

    flex-shrink: 0;

    border-radius: 50%;

    background: var(--color-input);
    color: var(--color-text);

    font-family: var(--font-body);
    font-size: 1rem;
    font-weight: 700;
}

/* =========================
   User Information
   ========================= */

.invite-user__name {
    margin: 0;

    color: var(--color-text);

    font-family: var(--font-body);
    font-size: 0.95rem;
    font-weight: 600;
}

.invite-user__username {
    margin: var(--space-1) 0 0;

    color: var(--color-text-muted);

    font-family: var(--font-body);
    font-size: 0.85rem;
}

/* =========================
   Invite Button
   ========================= */

.invite-user button {
    min-height: var(--touch-target);

    padding: 0 var(--space-4);

    border: 1px solid var(--color-mint);
    border-radius: var(--radius-small);

    background: transparent;
    color: var(--color-mint);

    font-family: var(--font-body);
    font-size: 0.875rem;
    font-weight: 600;

    cursor: pointer;

    transition:
        background 0.15s ease,
        color 0.15s ease,
        transform 0.15s ease;
}

.invite-user button:hover:not(:disabled) {
    background: var(--color-mint);
    color: var(--color-background);

    transform: translateY(-1px);
}

.invite-user button:active:not(:disabled) {
    transform: translateY(0);
}

.invite-user button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
}

/* =========================
   Messages
   ========================= */

.invite-message {
    margin: var(--space-3) 0 0;

    font-family: var(--font-body);
    font-size: 0.875rem;
}

.invite-message--success {
    color: var(--color-mint);
}

.invite-message--error {
    color: var(--color-coral);
}

.no-users {
    margin: var(--space-4) 0 0;

    color: var(--color-text-muted);

    font-family: var(--font-body);
    font-size: 0.9rem;
}

/* =========================
   Loading / Error
   ========================= */

.group-container>p {
    margin: var(--space-6) 0;

    text-align: center;

    color: var(--color-text-muted);

    font-family: var(--font-body);
}

/* =========================
   Mobile
   ========================= */

@media (max-width: 700px) {
    .group-page {
        padding: var(--space-4);
    }

    .group-header,
    .group-content,
    .invite-section {
        padding: var(--space-4);
    }

    .delete-button {
        width: 100%;
    }

    .invite-user {
        align-items: flex-start;
    }

    .invite-user button {
        flex-shrink: 0;
    }
}
</style>