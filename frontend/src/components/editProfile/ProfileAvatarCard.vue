<script setup>
import { computed, ref } from 'vue'

import { userData } from '@/stores/userData'

import {
    notification,
    showNotification,
    hideNotification
} from '@/helpers/errors'

const avatarVersion = ref(Date.now())

const initial = computed(() => {
    if (userData.value.firstName) {
        return userData.value.firstName.charAt(0).toUpperCase()
    }

    if (userData.value.username) {
        return userData.value.username.charAt(0).toUpperCase()
    }

    return '?'
})

const fullName = computed(() => {
    return `${userData.value.firstName || ''} ${userData.value.lastName || ''}`.trim()
})

async function uploadImage(e) {
    try {
        const file = e.target.files[0]

        if (!file) {
            showNotification('error', 'No image uploaded')
            return
        }

        const formData = new FormData()
        formData.append('avatar', file)

        const resp = await fetch('/api/user/avatar', {
            method: 'PUT',
            credentials: 'include',
            body: formData
        })

        if (!resp.ok) {
            showNotification('error', 'Error uploading image')
            return
        }

        const result = await resp.json()

        if (!result.status) {
            showNotification(
                'error',
                result.message || 'Error uploading image'
            )
            return
        }

        userData.value.avatar_path = result.avatar_path
        window.location.reload();
        showNotification('success', 'Avatar Updated!')
    } catch (err) {
        console.error(err)
        showNotification('error', 'Error uploading image')
    } finally {
        e.target.value = ''
    }
}

async function deleteAvatar() {
    try {
        const resp = await fetch("/api/user/avatar", {
            method: "DELETE",
            credentials: 'include'
        });

        if (!resp.ok) {
            showNotification('error', 'could not delete avatar');
            return
        }

        const result = await resp.json();
        if (!result.status) {
            showNotification('error', 'could not delete avatar');
            return
        }

        showNotification('success', 'Avatar Deleted!');
        window.location.reload();
        return
    } catch (err) {

    }
}
</script>

<template>
    <Transition name="notification">
        <div v-if="notification.show" class="notification" :class="`notification--${notification.type}`" role="alert">
            <span class="notification__icon">
                {{ notification.type === 'success' ? '✓' : '!' }}
            </span>

            <span class="notification__message">
                {{ notification.message }}
            </span>

            <button class="notification__close" type="button" aria-label="Close notification" @click="hideNotification">
                ×
            </button>
        </div>
    </Transition>

    <aside class="summary-card">

        <div class="summary-card__cover" aria-hidden="true"></div>

        <div class="avatar-wrapper">
            <div class="avatar">

                <img v-if="userData.avatar_path" :src="`${userData.avatar_path}?v=${avatarVersion}`" alt="Avatar">

                <span v-else class="avatar__initial">
                    {{ initial }}
                </span>

            </div>
        </div>

        <div class="summary-card__identity">
            <p class="summary-card__name">
                {{ fullName || userData.username || 'User' }}
            </p>
        </div>

        <div class="summary-card__actions">

            <form>
                <label for="avatar-upload" class="btn btn--ghost btn--block">
                    Upload new photo
                </label>

                <input id="avatar-upload" type="file" name="avatar" accept="image/*" hidden @change="uploadImage">
            </form>

            <button @click="deleteAvatar" v-if="userData.avatar_path" class="btn btn--text btn--block" type="button">
                Remove photo
            </button>

        </div>

    </aside>
</template>

<style scoped>
@import '../../styles/variables.css';


.summary-card__actions form {
    width: 100%;
}


.notification {
    position: fixed;
    top: var(--space-5);
    right: var(--space-5);
    z-index: 1000;

    display: flex;
    align-items: center;
    gap: var(--space-3);

    min-width: 18rem;
    max-width: 28rem;

    padding: var(--space-3) var(--space-4);

    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);

    color: var(--color-text);
    box-shadow: var(--shadow-raised);
}

.notification--success {
    border-color: rgb(34 197 94 / 40%);
}

.notification--error {
    border-color: rgb(239 68 68 / 40%);
}

.notification__icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;

    flex: 0 0 auto;

    width: 1.5rem;
    height: 1.5rem;

    border-radius: 50%;

    font-size: 0.75rem;
    font-weight: 700;
}

.notification--success .notification__icon {
    background: rgb(34 197 94 / 12%);
    color: #22c55e;
}

.notification--error .notification__icon {
    background: rgb(239 68 68 / 12%);
    color: #ef4444;
}

.notification__message {
    flex: 1;
    font-size: 0.82rem;
    line-height: 1.4;
}

.notification__close {
    display: inline-flex;
    align-items: center;
    justify-content: center;

    width: 1.5rem;
    height: 1.5rem;

    padding: 0;

    border: 0;
    background: transparent;

    color: var(--color-text-muted);

    font-size: 1.2rem;
    cursor: pointer;
}

.notification__close:hover {
    color: var(--color-text);
}

.notification-enter-active,
.notification-leave-active {
    transition:
        opacity 0.2s ease,
        transform 0.2s ease;
}

.notification-enter-from,
.notification-leave-to {
    opacity: 0;
    transform: translateY(-10px);
}

.summary-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    overflow: hidden;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    box-shadow: var(--shadow-raised);
}

.summary-card__cover {
    width: 100%;
    height: 4.5rem;
    background: var(--gradient-action);
}

.avatar-wrapper {
    position: relative;
    width: 6.5rem;
    height: 6.5rem;
    margin-top: -3.25rem;
}

.avatar {
    width: 100%;
    height: 100%;
    display: grid;
    place-items: center;
    overflow: hidden;
    background: var(--color-surface-raised);
    border: 4px solid var(--color-surface);
    border-radius: 50%;
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.2);
}

.avatar img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.avatar__initial {
    display: grid;
    width: 100%;
    height: 100%;
    place-items: center;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 2.25rem;
    font-weight: 700;
    line-height: 1;
    user-select: none;
}

.avatar-camera {
    position: absolute;
    right: -0.2rem;
    bottom: -0.1rem;
    display: grid;
    width: 2.25rem;
    height: 2.25rem;
    place-items: center;
    padding: 0;
    background: var(--color-input);
    border: 2px solid var(--color-surface);
    border-radius: 50%;
    color: var(--color-text-soft);
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    transition:
        background 0.15s ease,
        color 0.15s ease,
        transform 0.15s ease;
}

.avatar-camera:hover {
    background: var(--color-violet);
    color: white;
    transform: translateY(-1px);
}

.avatar-camera:focus-visible {
    outline: 2px solid var(--color-violet);
    outline-offset: 2px;
}

.summary-card__identity {
    margin-top: var(--space-3);
    text-align: center;
}

.summary-card__name {
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.1rem;
    font-weight: 700;
}

.summary-card__handle {
    margin-top: 0.2rem;
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.8rem;
}

.summary-card__actions {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    width: 100%;
    padding: var(--space-4);
    margin-top: var(--space-4);
}

.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-target);
    padding-inline: var(--space-4);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    font: inherit;
    font-weight: 600;
    cursor: pointer;
    transition:
        border-color 0.15s ease,
        background 0.15s ease,
        color 0.15s ease;
}

.btn--block {
    width: 100%;
}

.btn--ghost {
    background: var(--color-input);
    border-color: var(--color-border);
    color: var(--color-text-soft);
}

.btn--ghost:hover,
.btn--ghost:focus-visible {
    border-color: var(--color-violet);
    color: var(--color-text);
    outline: none;
}

.btn--text {
    background: transparent;
    color: var(--color-coral);
}

.btn--text:hover,
.btn--text:focus-visible {
    text-decoration: underline;
    outline: none;
}

@media (min-width: 64rem) {
    .summary-card {
        position: relative;
        top: calc(4rem + var(--space-5));
    }
}
</style>