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
    gap: 10px;
}

button {
    padding: 8px 14px;
    border: none;
    background: transparent;
    cursor: pointer;
}

button.active {
    font-weight: bold;
    border-bottom: 2px solid black;
}
</style>