<script setup>
defineProps({
    show: {
        type: Boolean,
        default: false
    },

    people: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['close']);

function close() {
    emit('close');
}
</script>

<template>
    <Teleport to="body">
        <div v-if="show" class="tagged-dialog-overlay" @click.self="close">
            <div class="tagged-dialog">
                <header class="tagged-dialog-header">
                    <span>Tagged people</span>

                    <button type="button" class="dialog-close" @click="close">
                        ✕
                    </button>
                </header>

                <ul class="tagged-list">
                    <li v-for="person in people" :key="person.id" class="tagged-list-item">
                        <a :href="`/user?id=${person.id}`" class="tagged-list-link">
                            <div class="tagged-list-avatar">
                                <img v-if="person.avatarPath" :src="`/uploads/${person.avatarPath}`"
                                    :alt="`${person.firstName} ${person.lastName}`">

                                <span v-else>
                                    {{ person.firstName?.charAt(0) }}
                                </span>
                            </div>

                            <span class="tagged-list-name">
                                {{ person.firstName }} {{ person.lastName }}
                            </span>
                        </a>
                    </li>
                </ul>
            </div>
        </div>
    </Teleport>
</template>

<style scoped>
.tagged-dialog-overlay {
    position: fixed;
    inset: 0;
    z-index: 1000;

    display: flex;
    align-items: center;
    justify-content: center;

    padding: 20px;

    background: rgba(0, 0, 0, 0.5);
}

.tagged-dialog {
    width: 100%;
    max-width: 400px;
    max-height: 70vh;

    display: flex;
    flex-direction: column;

    border: 2px solid var(--main-color);
    border-radius: 8px;
    overflow: hidden;

    background: var(--bg-color);
    box-shadow: 7px 7px var(--main-color);
}

.tagged-dialog-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    padding: 14px 16px;

    border-bottom: 2px solid var(--main-color);

    color: var(--main-color);

    font-family: "Liter", serif;
    font-size: 14px;
    font-weight: 600;
}

.dialog-close {
    padding: 2px 6px;

    border: none;
    background: transparent;

    color: var(--font-color-sub);

    font-size: 14px;
    font-weight: 600;
}

.dialog-close:hover {
    color: var(--main-color);
}

.tagged-list {
    margin: 0;
    padding: 8px;

    list-style: none;

    overflow-y: auto;
}

.tagged-list-item {
    border-bottom: 1px solid #eee;
}

.tagged-list-item:last-child {
    border-bottom: none;
}

.tagged-list-link {
    display: flex;
    align-items: center;
    gap: 12px;

    padding: 10px 8px;

    text-decoration: none;

    transition: background 0.15s;
}

.tagged-list-link:hover {
    background: var(--page-background);
}

.tagged-list-avatar {
    flex-shrink: 0;

    width: 40px;
    height: 40px;

    display: flex;
    align-items: center;
    justify-content: center;

    overflow: hidden;

    border: 2px solid var(--main-color);
    border-radius: 50%;

    background: var(--input-focus);

    color: white;

    font-family: "Liter", serif;
    font-size: 16px;
    font-weight: 600;
}

.tagged-list-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.tagged-list-name {
    color: var(--main-color);

    font-size: 13px;
    font-weight: 600;
}

@media (max-width: 650px) {
    .tagged-dialog {
        max-width: 100%;
    }
}
</style>