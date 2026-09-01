<script setup>
defineProps({
    group: {
        type: Object,
        required: true
    }
})

const emit = defineEmits(['join-group', 'view-group'])
</script>

<template>
    <div class="group-card">
        <h3>{{ group.title }}</h3>

        <p>{{ group.description }}</p>

        <span>
            {{ group.memberCount }} members
        </span>

        <div class="actions">
            <button @click="emit('view-group', group.id)">
                View Group
            </button>

            <button v-if="!group.isMember && !group.requested" @click="emit('join-group', group.id)">
                Join
            </button>

            <button v-else-if="group.requested" disabled>
                Requested
            </button>
        </div>
    </div>
</template>

<style scoped>
.group-card {
    padding: 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
}

.group-card h3 {
    margin-top: 0;
}

.actions {
    margin-top: 15px;
    display: flex;
    gap: 10px;
}

button {
    padding: 8px 12px;
    cursor: pointer;
}
</style>