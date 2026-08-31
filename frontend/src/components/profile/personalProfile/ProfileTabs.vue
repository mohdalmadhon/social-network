<script setup>
const emit = defineEmits(['changeTab']);

const tabs = [
    'posts',
    'friends',
    'groups',
    'following',
    'followers',
    'about'
];

function selectTab(tab) {
    emit('changeTab', tab);
}
</script>

<template>
    <nav class="profile-tabs">
        <button
            v-for="tab in tabs"
            :key="tab"
            type="button"
            @click="selectTab(tab)"
        >
            {{ tab }}
        </button>
    </nav>
</template>

<style scoped>
.profile-tabs {
    position: sticky;
    top: 64px;
    z-index: 50;
    display: flex;
    gap: 28px;
    margin: 25px 0;
    overflow-x: auto;
    border-bottom: 1px solid #232332;
    background: #0a0a12;
    padding: 0 4px;
}

.profile-tabs button {
    flex: none;
    min-width: 80px;
    padding: 12px 0;
    border: 0;
    background: transparent;
    color: #8b8b9e;
    text-align: center;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    text-transform: capitalize;
}

.profile-tabs button:hover {
    color: #fff;
}

.profile-tabs button.active {
    color: #a855f7;
    font-weight: 600;
    border-bottom: 2px solid #a855f7;
}

@media (max-width: 800px) {
    .profile-tabs {
        top: 0;
    }
}
</style>