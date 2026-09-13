<script setup>
import { ref, onMounted } from 'vue'
import { router } from '@/router/router.js'
import { useRoute } from 'vue-router'
import { getGroup, deleteGroupApi, getGroupPosts } from '@/api/groups/Groups.js'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import GroupActivity from '@/components/groups/GroupActivity.vue'
import GroupPostCard from '@/components/groups/GroupPostCard.vue'
import GroupPostComposer from '@/components/groups/GroupPostComposer.vue'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const route = useRoute()
const groupId = route.params.groupId

const group = ref(null)
const groupPosts = ref([])
const loading = ref(true)
const error = ref(null)
const isDeletingGroup = ref(false)
const deleteError = ref('')

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
    if (isDeletingGroup.value || !window.confirm('Delete this group and all of its content?')) return

    isDeletingGroup.value = true
    deleteError.value = ''
    try {
        const result = await deleteGroupApi(groupId)

        if (result?.status) {
            router.replace('/groups')
        }
    } catch (err) {
        console.error(err)
        deleteError.value = err.message || 'Could not delete group.'
    } finally {
        isDeletingGroup.value = false
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
        <div class="group-container">

            <p v-if="loading" class="group-page__state">
                Loading group...
            </p>

            <p v-else-if="error" class="group-page__state group-page__state--error" role="alert">
                {{ error }}
            </p>

            <template v-else-if="group">
                <header class="group-header orbit-surface">
                    <div class="group-header__copy">
                        <p class="orbit-meta">Community</p>
                        <h1>{{ group.title }}</h1>
                        <p>{{ group.description }}</p>
                    </div>

                    <button v-if="group.isCreator" class="group-header__delete" :disabled="isDeletingGroup" @click="deleteGroup">
                        {{ isDeletingGroup ? 'Deleting...' : 'Delete group' }}
                    </button>
                    <p v-if="deleteError" class="group-header__error" role="alert">{{ deleteError }}</p>
                </header>

                <section class="group-membership" aria-label="Group membership">
                    <div>
                        <IconGlyph name="groups" :size="19" />
                        <strong>{{ group.memberCount }}</strong>
                        <span>{{ group.memberCount === 1 ? 'member' : 'members' }}</span>
                    </div>
                    <p v-if="group.isMember" class="group-membership__status">
                        <IconGlyph name="check" :size="16" />
                        Member
                    </p>
                </section>

                <GroupActivity v-if="group.isMember" :group-id="groupId" />

                <section v-if="group.isMember" class="group-feed orbit-surface">
                    <div class="group-feed__heading">
                        <div>
                            <p class="orbit-meta">Conversation</p>
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
            </template>

        </div>
    </div>
    </AuthenticatedLayout>
</template>

<style scoped>
.group-page {
    width: 100%;
    padding: 0;
    background: var(--color-background);
    color: var(--color-text);
}

.group-container {
    display: grid;
    width: 100%;
    max-width: 64rem;
    margin: 0 auto;
    gap: var(--space-7);
}

.group-header {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: end;
    gap: var(--space-5);
    padding: var(--space-6);
}

.group-header__copy { min-width: 0; }
.group-header .orbit-meta { margin: 0 0 var(--space-2); color: var(--color-violet-soft); }

.group-header h1 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 2.25rem;
    font-weight: 700;
    line-height: 1.15;
    letter-spacing: 0;
}

.group-header__copy > p:last-child {
    max-width: 700px;
    margin: var(--space-3) 0 0;
    color: var(--color-text-muted);
    font-size: 1rem;
    line-height: 1.65;
}

.group-header__delete {
    min-height: var(--touch-target);
    padding: 0 var(--space-4);
    border: 1px solid var(--color-coral);
    border-radius: var(--radius-small);
    background: transparent;
    color: var(--color-coral);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s ease, color 0.15s ease;
}

.group-header__delete:hover:not(:disabled) {
    background: var(--color-coral);
    color: var(--color-background);
}

.group-header__error {
    grid-column: 1 / -1;
    margin: 0;
    color: var(--color-coral);
    font-size: .875rem;
}

.group-membership {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
    margin-top: calc(var(--space-5) * -1);
    padding: 0 var(--space-2) var(--space-5);
    border-bottom: 1px solid var(--color-border);
    color: var(--color-text-muted);
}
.group-membership > div,
.group-membership__status {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
}
.group-membership strong { color: var(--color-text); }
.group-membership__status {
    margin: 0;
    color: var(--color-mint);
    font-size: 0.875rem;
    font-weight: 600;
}

.group-feed {
    display: grid;
    gap: var(--space-5);
    padding: var(--space-6) 0 0;
    border-width: 1px 0 0;
    border-radius: 0;
    background: transparent;
    box-shadow: none;
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
    font-size: 1.75rem;
    letter-spacing: 0;
}

.group-feed__heading > span {
    color: var(--color-text-muted);
    font-size: .875rem;
}

.group-posts {
    display: grid;
    gap: var(--space-5);
}

.group-feed__state {
    margin: 0;
    padding: var(--space-5);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-small);
    color: var(--color-text-muted);
    text-align: center;
}

.group-page__state {
    margin: var(--space-7) 0;
    padding: var(--space-6);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-medium);
    text-align: center;
    color: var(--color-text-muted);
}
.group-page__state--error { color: var(--color-coral); }

@media (max-width: 700px) {
    .group-container { gap: var(--space-6); }
    .group-header {
        grid-template-columns: 1fr;
        align-items: start;
        padding: var(--space-5);
    }
    .group-header h1 { font-size: 1.875rem; }
    .group-header__delete {
        width: 100%;
    }
    .group-membership { margin-top: calc(var(--space-4) * -1); }
    .group-feed { padding-top: var(--space-5); }
    .group-feed__heading { align-items: start; }
    .group-feed__heading h2 { font-size: 1.5rem; }
}
</style>
