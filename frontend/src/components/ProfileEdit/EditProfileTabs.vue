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
@import '../../styles/global.css';
@import '../../styles/variables.css';
.edit-tabs {
    position: sticky;
    top: 64px;
    z-index: 50;
    display: flex;
    margin: clamp(16px, 3vw, 25px) 0;
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    -webkit-overflow-scrolling: touch;
}

.tab-button {
    flex: 1 1 0;
    min-width: max-content;
    padding: clamp(11px, 2vw, 15px) clamp(12px, 2.5vw, 18px);
    border: none;
    border-right: 1px solid var(--color-border);
    background: transparent;
    color: var(--color-text-muted);
    text-align: center;
    font-family: var(--font-meta);
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
    white-space: nowrap;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}

.tab-button:last-child {
    border-right: 0;
}

.tab-button:nth-child(1):hover {
    background: var(--color-surface-violet);
    color: var(--color-violet-soft);
}

.tab-button:nth-child(2):hover {
    background: var(--color-surface-teal);
    color: var(--color-cyan-soft);
}

.tab-button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.tab-button.active {
    position: relative;
    background: var(--gradient-cyber);
    color: var(--color-text);
}

.tab-button.active::after {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 2px;
    background: var(--color-amber);
}

@media (max-width: 800px) {
    .edit-tabs {
        top: 0;
    }
}
</style>