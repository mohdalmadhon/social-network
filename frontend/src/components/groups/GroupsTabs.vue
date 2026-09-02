<script setup>
defineProps({
    activeTab: {
        type: String,
        required: true
    }
})

const emit = defineEmits(['change-tab'])
</script>

<template>
    <div class="groups-tabs">
        <button :class="{ active: activeTab === 'discover' }" @click="emit('change-tab', 'discover')">
            Discover
        </button>

        <button :class="{ active: activeTab === 'my-groups' }" @click="emit('change-tab', 'my-groups')">
            My Groups
        </button>
    </div>
</template>

<style scoped>
.groups-tabs {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-1);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-surface);
    width: fit-content;
}

.groups-tabs button {
    min-height: 38px;
    padding: 0 var(--space-4);
    border: 0;
    border-radius: 6px;
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 500;
    cursor: pointer;
    transition:
        color 0.15s ease,
        background 0.15s ease;
}

.groups-tabs button:hover {
    color: var(--color-text);
}

.groups-tabs button.active {
    background: var(--color-surface-raised);
    color: var(--color-violet);
    font-weight: 600;
}

@media (max-width: 30rem) {
    .groups-tabs {
        width: 100%;
    }

    .groups-tabs button {
        flex: 1;
    }
}
</style>