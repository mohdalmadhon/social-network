<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const props = defineProps({
    followers: {
        type: Object,
        default: () => ({})
    }
});

const followerList = computed(() => {
    return Object.entries(props.followers).map(([id, follower]) => ({
        id,
        ...follower
    }));
});

function takeToProfile(id) {
    router.push(`/user?id=${id}`);
    window.location.reload();
    return;
}
</script>

<template>
    <section class="followers-section">
        <div class="section-heading">
            <p class="eyebrow">SOCIAL</p>
            <h2>Followers</h2>
        </div>

        <div class="followers-card">
            <div 
                v-if="followerList.length"
                class="followers-grid"
            >
                <article
                    v-for="follower in followerList"
                    :key="follower.id"
                    @click="takeToProfile(follower.id)"
                    class="follower-card"
                >
                    <img
                        :src="follower.Avatar ? `/uploads/${follower.Avatar}` : '/default-avatar.png'"
                        :alt="`${follower.firstName} ${follower.LastName}`"
                        class="follower-avatar"
                    >

                    <div class="follower-info">
                        <p class="follower-name">
                            {{ follower.FirstName }} {{ follower.LastName }}
                        </p>
                    </div>
                </article>
            </div>

            <p
                v-else
                class="empty"
            >
                No followers yet.
            </p>
        </div>
    </section>
</template>

<style scoped>
.followers-section {
    width: 100%;
    scroll-margin-top: 6.25rem;
}

.section-heading {
    margin-bottom: var(--space-4);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.followers-card {
    width: 100%;
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    box-sizing: border-box;
}

.followers-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(13.75rem, 1fr));
    gap: var(--space-4);
}

.follower-card {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    min-width: 0;
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-raised);
    box-sizing: border-box;
    cursor: pointer;
    transition: border-color 0.15s ease, transform 0.15s ease;
}

.follower-card:hover {
    transform: translateY(-1px);
    border-color: var(--color-violet);
}

.follower-avatar {
    flex: 0 0 auto;
    width: 3.125rem;
    height: 3.125rem;
    border-radius: 50%;
    object-fit: cover;
}

.follower-info {
    min-width: 0;
}

.follower-name {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.875rem;
    font-weight: 600;
    line-height: 1.4;
    overflow-wrap: anywhere;
}

.empty {
    margin: 0;
    padding: var(--space-5) var(--space-2);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
    text-align: center;
}

@media (max-width: 37.5rem) {
    .followers-grid {
        grid-template-columns: 1fr;
    }
}
</style>
