<script setup>
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import FollowersDialog from '@/components/personalProfile/FollowersDialog.vue';


const router = useRouter();

const props = defineProps({
    followers: {
        type: Object,
        default: () => ({})
    },
    type: {
        type: String,
        default: 'followers'
    },
    targetId: {
        type: [String, Number],
        default: null
    }
});

const showDialog = ref(false);

const followerList = computed(() => {
    return Object.entries(props.followers).map(([id, follower]) => ({
        id,
        ...follower
    }));
});

async function takeToProfile(id) {
    await router.push(`/user?id=${id}`);
    window.location.reload();
    return;
}

function openDialog() {
    showDialog.value = true;
}

function closeDialog() {
    showDialog.value = false;
}
</script>

<template>
    <section class="followers-section">
        <div class="section-heading">
            <div class="heading-row">
                <div>
                    <p class="eyebrow">SOCIAL</p>

                    <h2>
                        {{
                            type === 'following'
                                ? 'Following'
                                : type === 'friends'
                                    ? 'Friends'
                                    : 'Followers'
                        }}
                    </h2>
                </div>

                <button v-if="followerList.length" type="button" class="show-all-btn" @click="openDialog">
                    Show all
                </button>
            </div>
        </div>

        <div class="followers-card">
            <div v-if="followerList.length" class="followers-grid">
                <article v-for="follower in followerList" :key="follower.id" @click="takeToProfile(follower.id)"
                    class="follower-card">
                    <img :src="follower.Avatar ? `/uploads/${follower.Avatar}` : '/default-avatar.png'"
                        :alt="`${follower.FirstName} ${follower.LastName}`" class="follower-avatar">

                    <div class="follower-info">
                        <p class="follower-name">
                            {{ follower.FirstName }} {{ follower.LastName }}
                        </p>
                    </div>
                </article>
            </div>

            <p v-else class="empty">
                No followers yet.
            </p>
        </div>

        <FollowersDialog v-if="showDialog" :type="type" :target-id="targetId" @close="closeDialog"
            @navigate="takeToProfile" />
    </section>
</template>

<style scoped>
.followers-section {
    width: 100%;
    scroll-margin-top: 100px;
}

.section-heading {
    margin-bottom: clamp(14px, 2.5vw, 20px);
}

.heading-row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: var(--space-2);
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--color-cyan);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
    text-transform: uppercase;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: clamp(20px, 4vw, 29px);
}

.show-all-btn {
    flex: 0 0 auto;
    padding: var(--space-2) var(--space-4);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-aurora);
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: clamp(10px, 1.4vw, 12px);
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.show-all-btn:hover {
    transform: translateY(-1px);
    filter: brightness(1.08);
}

.show-all-btn:active {
    transform: translateY(0);
}

.show-all-btn:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.followers-card {
    width: 100%;
    padding: clamp(14px, 3vw, 25px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-teal);
    box-shadow: var(--shadow-raised);
    box-sizing: border-box;
}

.followers-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: clamp(12px, 2vw, 18px);
}

.follower-card {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    min-width: 0;
    padding: clamp(12px, 2vw, 16px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-surface-raised);
    box-sizing: border-box;
    transition: transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease;
    cursor: pointer;
}

.follower-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-raised);
    border-color: var(--color-cyan);
}

.follower-avatar {
    flex: 0 0 auto;
    width: clamp(45px, 6vw, 60px);
    height: clamp(45px, 6vw, 60px);
    border: 2px solid var(--color-border);
    border-radius: 50%;
    object-fit: cover;
}

.follower-info {
    min-width: 0;
}

.follower-name {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: clamp(12px, 1.8vw, 15px);
    font-weight: 600;
    line-height: 1.4;
    overflow-wrap: anywhere;
}

.empty {
    margin: 0;
    padding: var(--space-6) var(--space-2);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(9px, 1.4vw, 11px);
    text-align: center;
}

@media (max-width: 600px) {
    .followers-grid {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 400px) {
    .follower-card {
        padding: 11px;
    }

    .follower-avatar {
        width: 45px;
        height: 45px;
    }
}
</style>