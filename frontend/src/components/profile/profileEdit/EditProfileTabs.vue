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
@import '../../../styles/global.css';

.edit-tabs {
    position: sticky;
    top: 64px;
    z-index: 50;
    display: flex;
    margin: var(--space-5) 0;
    overflow-x: auto;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    box-shadow: var(--shadow-raised);
    -webkit-overflow-scrolling: touch;
}

.tab-button {
    flex: 1 1 0;
    min-width: max-content;
    min-height: var(--touch-target);
    padding-inline: var(--space-4);
    border: none;
    border-right: 1px solid var(--color-border);
    background: transparent;
    color: var(--color-text-soft);
    text-align: center;
    font-family: var(--font-meta);
    font-size: 0.78rem;
    font-weight: 600;
    letter-spacing: 0.04em;
    white-space: nowrap;
    cursor: pointer;
    transition:
        background 0.15s ease,
        color 0.15s ease;
}

.tab-button:last-child {
    border-right: 0;
}

.tab-button:hover {
    background: var(--color-surface-raised);
    color: var(--color-text);
}

.tab-button.active {
    background: var(--gradient-action);
    color: #fff;
}

@media (max-width: 50rem) {
    .edit-tabs {
        top: 0;
    }
}
</style>
