import { ref } from 'vue'

export const notification = ref({
    show: false,
    type: '',
    message: ''
})

let notificationTimer = null

export function showNotification(type, message) {
    notification.value = {
        show: true,
        type,
        message
    }

    clearTimeout(notificationTimer)

    notificationTimer = setTimeout(() => {
        notification.value.show = false
    }, 4000)
}

export function hideNotification() {
    notification.value.show = false
    clearTimeout(notificationTimer)
}