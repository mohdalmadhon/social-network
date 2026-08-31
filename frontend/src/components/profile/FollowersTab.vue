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
    scroll-margin-top: 100px;
}

.section-heading {
    margin-bottom: clamp(14px, 2.5vw, 20px);
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: clamp(20px, 4vw, 29px);
}

.followers-card {
    width: 100%;
    padding: clamp(14px, 3vw, 25px);
    border: 2px solid var(--main-color);
    border-radius: 7px;
    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
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
    gap: 14px;
    min-width: 0;
    padding: clamp(12px, 2vw, 16px);
    border: 2px solid var(--main-color);
    border-radius: 6px;
    background: var(--bg-color);
    box-shadow: 3px 3px var(--main-color);
    box-sizing: border-box;
    transition: transform 0.15s ease;
}

.follower-card:hover {
    transform: translate(-2px, -2px);
}

.follower-avatar {
    flex: 0 0 auto;
    width: clamp(45px, 6vw, 60px);
    height: clamp(45px, 6vw, 60px);
    border: 2px solid var(--main-color);
    border-radius: 50%;
    object-fit: cover;
}

.follower-info {
    min-width: 0;
}

.follower-name {
    margin: 0;
    color: var(--font-color);
    font-family: "Hedvig Letters Sans", sans-serif;
    font-size: clamp(12px, 1.8vw, 15px);
    font-weight: 600;
    line-height: 1.4;
    overflow-wrap: anywhere;
}

.empty {
    margin: 0;
    padding: 25px 10px;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
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
