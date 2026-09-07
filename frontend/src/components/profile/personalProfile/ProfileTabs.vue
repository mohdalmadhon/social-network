<script setup>
const emit = defineEmits(['changeTab']);

const tabs = [
    'posts',
    'friends',
    'following',
    'followers',
    'about'
];

function selectTab(tab) {
    emit('changeTab', tab);
}
</script>

<template>
    <nav class="profile-tabs">
        <button
            v-for="tab in tabs"
            :key="tab"
            type="button"
            @click="selectTab(tab)"
        >
            {{ tab }}
        </button>
    </nav>
</template>

<style scoped>
.profile-tabs {
    position: sticky;
    top: var(--space-4);
    z-index: 50;
    display: flex;
    gap: var(--space-1);
    margin: var(--space-5) 0;
    padding: var(--space-1);
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.profile-tabs button {
    flex: 1;
    min-width: 6.5625rem;
    padding: var(--space-3) var(--space-4);
    border: 0;
    border-radius: var(--radius-medium);
    background: transparent;
    color: var(--color-text-muted);
    text-align: center;
    text-transform: capitalize;
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s ease, color 0.2s ease;
}

.profile-tabs button:hover {
    color: var(--color-text-soft);
    background: var(--color-surface-raised);
}

.profile-tabs button.active {
    background: var(--gradient-action);
    color: var(--color-text);
}

@media (max-width: 48rem) {
    .profile-tabs {
        top: 0;
    }
}
</style>
