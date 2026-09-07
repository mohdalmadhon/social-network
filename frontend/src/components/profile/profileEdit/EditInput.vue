<script setup>
defineProps({
    activeTab: String
})

defineEmits(['update:activeTab'])

const tabs = [
    { id: 'personal', label: 'Personal Info' },
    { id: 'additional', label: 'Additional Info' }
]
</script>

<template>
    <nav class="edit-tabs">
        <button
            v-for="tab in tabs"
            :key="tab.id"
            type="button"
            class="tab-button"
            :class="{ active: activeTab === tab.id }"
            @click="$emit('update:activeTab', tab.id)"
        >
            {{ tab.label }}
        </button>
    </nav>
</template>

<style scoped>
.edit-tabs {
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
    -webkit-overflow-scrolling: touch;
}

.tab-button {
    flex: 1 1 0;
    min-width: max-content;
    padding: var(--space-3) var(--space-4);
    border: none;
    border-radius: var(--radius-medium);
    background: transparent;
    color: var(--color-text-muted);
    text-align: center;
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 600;
    white-space: nowrap;
    cursor: pointer;
    transition: background 0.2s ease, color 0.2s ease;
}

.tab-button:hover {
    color: var(--color-text-soft);
    background: var(--color-surface-raised);
}

.tab-button.active {
    background: var(--gradient-action);
    color: var(--color-text);
}

.tab-button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

@media (max-width: 48rem) {
    .edit-tabs {
        top: 0;
    }
}
</style>
