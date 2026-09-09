<script setup>
const emit = defineEmits(['create-group'])
</script>

<template>
    <header class="groups-header">
        <div class="header-content">
            <p class="eyebrow">Community</p>
            <h1>Groups</h1>
            <p class="description">
                Discover communities, join discussions, and connect with other members.
            </p>
        </div>

        <button class="create-button" @click="emit('create-group')">
            Create Group
        </button>
    </header>
</template>

<style scoped>
.groups-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: var(--space-4);
}

.header-content {
    min-width: 0;
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
    font-weight: 600;
    letter-spacing: 2px;
    text-transform: uppercase;
}

h1 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: clamp(1.75rem, 4vw, 2.25rem);
    font-weight: 700;
    line-height: 1.2;
}

.description {
    max-width: 600px;
    margin: var(--space-2) 0 0;
    color: var(--color-text-muted);
    font-size: 0.875rem;
}

.create-button {
    min-height: var(--touch-target);
    padding: 0 var(--space-4);
    border: 0;
    border-radius: var(--radius-small);
    background: var(--gradient-action);
    color: #fff;
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 600;
    cursor: pointer;
    transition:
        transform 0.15s ease,
        opacity 0.15s ease;
}

.create-button:hover {
    transform: translateY(-1px);
    opacity: 0.92;
}

@media (max-width: 36rem) {
    .groups-header {
        align-items: stretch;
        flex-direction: column;
    }

    .create-button {
        width: 100%;
    }
}
</style>