import { ref } from 'vue';

import { getUnreadNotificationCount } from '@/api/common/notifications';

export const unreadNotificationCount = ref(0);

export async function refreshUnreadNotificationCount() {
    try {
        const result = await getUnreadNotificationCount();

        unreadNotificationCount.value = result.count || 0;
    } catch (err) {
        console.error(err);
    }
}

export function incrementUnreadNotificationCount() {
    unreadNotificationCount.value += 1;
}

export function clearUnreadNotificationCount() {
    unreadNotificationCount.value = 0;
}

function isDisplayableNotification(data) {
    return !data.message_user_id
        && !data.group_invite_user_id
        && !data.group_join_user_id
        && !data.group_accept_user_id;
}

export function handleIncomingNotification(data) {
    if (isDisplayableNotification(data)) {
        incrementUnreadNotificationCount();
    }
}
