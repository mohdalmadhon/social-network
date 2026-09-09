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
    top: 64px;
    z-index: 50;
    display: flex;
    margin: var(--space-5) 0;
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.profile-tabs button {
    flex: 1;
    min-width: 105px;
    padding: var(--space-4);
    border: 0;
    border-right: 1px solid var(--color-border);
    background: transparent;
    color: var(--color-text-muted);
    text-align: center;
    text-transform: capitalize;
    font-family: var(--font-meta);
    font-size: 10px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}

.profile-tabs button:last-child {
    border-right: 0;
}

.profile-tabs button:nth-child(1):hover {
    background: var(--color-surface-teal);
    color: var(--color-cyan-soft);
}

.profile-tabs button:nth-child(2):hover {
    background: var(--color-surface-violet);
    color: var(--color-violet-soft);
}

.profile-tabs button:nth-child(3):hover {
    background: var(--color-surface-coral);
    color: var(--color-coral-soft);
}

.profile-tabs button:nth-child(4):hover {
    background: var(--color-surface-amber);
    color: var(--color-amber-soft);
}

.profile-tabs button:nth-child(5):hover {
    background: var(--color-surface-violet);
    color: var(--color-magenta-soft);
}

.profile-tabs button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.profile-tabs button.active {
    background: var(--gradient-cyber);
    color: var(--color-text);
}

@media (max-width: 800px) {
    .profile-tabs {
        top: 0;
    }
}
</style>