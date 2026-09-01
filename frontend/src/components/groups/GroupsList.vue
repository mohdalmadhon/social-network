<script setup>
import GroupCard from './GroupCard.vue'

defineProps({
    groups: {
        type: Array,
        default: () => []
    }
})

const emit = defineEmits(['join-group', 'view-group'])
</script>

<template>
    <div>
        <div v-if="groups.length" class="groups-list">
            <GroupCard v-for="group in groups" :key="group.id" :group="group" @join-group="emit('join-group', $event)"
                @view-group="emit('view-group', $event)" />
        </div>

        <p v-else>
            No groups found.
        </p>
    </div>
</template>

<style scoped>
.groups-list {
    display: grid;
    gap: 16px;
}
</style>