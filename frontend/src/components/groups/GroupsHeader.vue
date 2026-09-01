<script setup>
const emit = defineEmits(['create-group'])
</script>

<template>
    <div class="groups-header">
        <div>
            <h1>Groups</h1>
            <p>Discover groups and connect with other users.</p>
        </div>

        <button @click="emit('create-group')">
            Create Group
        </button>
    </div>
</template>

<style scoped>
.groups-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.groups-header h1 {
    margin: 0;
}

.groups-header p {
    margin-top: 5px;
    color: #666;
}

button {
    padding: 10px 16px;
    cursor: pointer;
}
</style>