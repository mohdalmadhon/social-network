<script setup>
import { notifications, removeNotification } from '@/data/notifications';
</script>
    
<template>
    <div class="notification-container">
        <div
            v-for="notification in notifications"
            :key="notification.id"
            class="notification"
            :class="notification.type"
        >
            <span>{{ notification.message }}</span>

            <button
                type="button"
                @click="removeNotification(notification.id)"
            >
                ×
            </button>
        </div>
    </div>
</template>

<style scoped>
.notification-container {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 9999;

    display: flex;
    flex-direction: column;
    gap: 10px;

    width: min(350px, calc(100vw - 40px));
}

.notification {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 15px;

    padding: 14px 16px;

    border: 2px solid var(--main-color);
    border-radius: 5px;

    background: var(--bg-color);
    color: var(--font-color);

    box-shadow: 4px 4px var(--main-color);

    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
}

.notification.success {
    border-color: var(--input-focus);
}

.notification.error {
    border-color: #d9534f;
}

.notification button {
    border: none;
    background: none;
    color: inherit;
    font-size: 18px;
    cursor: pointer;
}
</style>