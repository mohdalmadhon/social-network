<script setup>
import GroupCard from './GroupCard.vue'

defineProps({
    groups: {
        type: Array,
        default: () => []
    }
})

const emit = defineEmits(['toggle-join-request', 'view-group'])
</script>

<template>
    <div class="groups-list-wrapper">
        <div v-if="groups.length" class="groups-grid">
            <GroupCard v-for="group in groups" :key="group.id" :group="group"
                @toggle-join-request="emit('toggle-join-request', $event)" @view-group="emit('view-group', $event)" />
        </div>

        <div v-else class="empty-state">
            No groups found.
        </div>
    </div>
</template>

<style scoped>
.groups-list-wrapper {
    width: 100%;
}

.groups-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: var(--space-4);
}

.empty-state {
    padding: var(--space-7) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.75rem;
    text-align: center;
}

@media (max-width: 48rem) {
    .groups-grid {
        grid-template-columns: 1fr;
    }
}
</style>