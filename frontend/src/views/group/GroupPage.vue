<script setup>
import { ref, onMounted } from 'vue'
import { router } from '@/router/router.js'
import { useRoute } from 'vue-router'
import { getGroup, deleteGroupApi, getGroupPosts } from '@/api/groups/Groups.js'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import GroupActivity from '@/components/groups/GroupActivity.vue'
import GroupPostCard from '@/components/groups/GroupPostCard.vue'
import GroupPostComposer from '@/components/groups/GroupPostComposer.vue'

const route = useRoute()
const groupId = route.params.groupId

const group = ref(null)
const groupPosts = ref([])
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
    try {
        const result = await getGroup(groupId)

        group.value = result.group
        if (group.value?.isMember) {
            const postsResult = await getGroupPosts(groupId)
            groupPosts.value = postsResult?.posts || []
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

function addGroupPost(post) {
    groupPosts.value.unshift(post)
}

function removeGroupPost(postId) {
    groupPosts.value = groupPosts.value.filter((post) => post.id !== postId)
}
</script>

<template>
    <AuthenticatedLayout active-page="groups">
    <div class="group-page">
        <section class="group-container">

            <p v-if="loading">
                Loading group...
            </p>

            <p v-else-if="error">
                {{ error }}
            </p>

            <template v-else-if="group">
                <header class="group-header">
                    <div>
                        <h1>{{ group.title }}</h1>
                        <p>{{ group.description }}</p>
                    </div>

                    <button v-if="group.isCreator" @click="deleteGroup">
                        Delete Group
                    </button>
                </header>

                <section class="group-content">
                    <p>
                        Members: {{ group.memberCount }}
                    </p>

                    <p v-if="group.isMember">
                        You are a member of this group.
                    </p>
                </section>

                <section v-if="group.isMember" class="group-feed orbit-surface">
                    <div class="group-feed__heading">
                        <div>
                            <p class="orbit-meta">Shared space</p>
                            <h2>Group posts</h2>
                        </div>
                        <span>{{ groupPosts.length }} {{ groupPosts.length === 1 ? 'post' : 'posts' }}</span>
                    </div>
                    <GroupPostComposer :group-id="groupId" @post-created="addGroupPost" />
                    <div v-if="groupPosts.length" class="group-posts">
                        <GroupPostCard
                            v-for="post in groupPosts"
                            :key="post.id"
                            :group-id="groupId"
                            :post="post"
                            @post-deleted="removeGroupPost"
                        />
                    </div>
                    <p v-else class="group-feed__state">No posts yet. Start the group conversation.</p>
                </section>
                <GroupActivity v-if="group.isMember" :group-id="groupId" />
            </template>

        </section>
    </div>
    </AuthenticatedLayout>
</template>

<style scoped>
.group-page {
    width: 100%;
    min-height: 100vh;
    padding: clamp(.25rem, 2vw, 1.5rem);
    background: var(--color-background);
    color: var(--color-text);
}

.group-container {
    width: 100%;
    max-width: 960px;
    margin: 0 auto;
}

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

.group-header button {
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

.group-header button:hover {
    background: var(--color-coral);
    color: var(--color-background);
    transform: translateY(-1px);
}

.group-header button:active {
    transform: translateY(0);
}

.group-content {
    display: grid;
    gap: var(--space-4);

    padding: var(--space-5);

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

.group-content p:last-child {
    display: inline-flex;
    align-items: center;
    width: fit-content;

    padding: var(--space-2) var(--space-3);

    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);

    background: var(--color-surface-raised);
    color: var(--color-mint);

    font-size: 0.875rem;
    font-weight: 600;
}

.group-feed {
    display: grid;
    gap: var(--space-4);
    margin-top: var(--space-5);
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
}

.group-feed__heading {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: var(--space-3);
}

.group-feed__heading h2 {
    margin: var(--space-1) 0 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: clamp(1.3rem, 3vw, 1.7rem);
}

.group-feed__heading > span {
    color: var(--color-text-muted);
    font-size: .875rem;
}

.group-posts {
    display: grid;
    gap: var(--space-4);
}

.group-feed__state {
    margin: 0;
    padding: var(--space-5);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-small);
    color: var(--color-text-muted);
    text-align: center;
}

.group-container>p {
    margin: var(--space-6) 0;
    text-align: center;

    color: var(--color-text-muted);
    font-family: var(--font-body);
}

@media (max-width: 700px) {
    .group-page {
        padding: var(--space-4);
    }

    .group-header {
        padding: var(--space-5);
    }

    .group-content {
        padding: var(--space-4);
    }

    .group-feed {
        padding: var(--space-4);
    }

    .group-header button {
        width: 100%;
    }
}
</style>
