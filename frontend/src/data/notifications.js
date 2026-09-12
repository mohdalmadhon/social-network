import { reactive } from 'vue';

export const notifications = reactive([]);

export function addNotification(message, type = 'success') {
    const id = Date.now();

    const displayMessage = message instanceof Error ? message.message : String(message);

    notifications.push({
        id,
        message: displayMessage,
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
