<script setup>
defineProps({
    posts: {
        type: Array,
        default: () => []
    }
});
</script>

<template>
    <section class="tab-panel">
        <div v-if="!posts.length" class="empty-state">
            No posts in this group yet.
        </div>

        <div v-else class="posts-list">
            <article
                v-for="post in posts"
                :key="post.id"
                class="post-placeholder"
            >
                {{ post.content }}
            </article>
        </div>
    </section>
</template>

<style scoped>
.tab-panel {
    width: 100%;
}

.empty-state {
    padding: 40px 20px;
    text-align: center;

    border: 2px solid var(--main-color);
    border-radius: 8px;

    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);

    color: var(--font-color-sub);

    font-family: "JetBrains Mono", monospace;
    font-size: 12px;
}

.posts-list {
    display: flex;
    flex-direction: column;
    gap: 18px;
}

.post-placeholder {
    padding: 18px 20px;

    border: 2px solid var(--main-color);
    border-radius: 8px;

    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);

    color: var(--font-color);

    font-size: 14px;
    line-height: 1.55;
}
</style>
