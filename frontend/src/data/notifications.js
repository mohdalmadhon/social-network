import { reactive } from 'vue';

export const notifications = reactive([]);

export function addNotification(message, type = 'success') {
    const id = Date.now();

    notifications.push({
        id,
        message,
        type
    });

    setTimeout(() => {
        removeNotification(id);
    }, 3000);
}

export function removeNotification(id) {
    const index = notifications.findIndex(
        notification => notification.id === id
    );

    if (index !== -1) {
        notifications.splice(index, 1);
    }
}