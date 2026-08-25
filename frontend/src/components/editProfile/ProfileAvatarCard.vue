<script setup>
import { computed } from 'vue'
import { userData } from '@/stores/userData'

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
</script>

<template>
    <aside class="summary-card">
        <div class="summary-card__cover" aria-hidden="true"></div>

        <div class="avatar-wrapper">
            <div class="avatar">
                <img
                    v-if="userData.avatar_path"
                    :src="userData.avatar_path"
                    alt="Avatar"
                >

                <span v-else class="avatar__initial">
                    {{ initial }}
                </span>
            </div>

            <button class="avatar-camera" type="button" aria-label="Change profile photo">
                <svg
                    viewBox="0 0 24 24"
                    width="15"
                    height="15"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M3 7h4l2-2h6l2 2h4v13H3z" />
                    <circle cx="12" cy="13" r="4" />
                </svg>
            </button>
        </div>

        <div class="summary-card__identity">
            <p class="summary-card__name">
                {{ fullName || userData.username || 'User' }}
            </p>

            <p class="summary-card__handle">
                @{{ userData.username || 'username' }}
            </p>
        </div>

        <div class="summary-card__actions">
            <button class="btn btn--ghost btn--block" type="button">
                Upload new photo
            </button>

            <button
                v-if="userData.avatar_path"
                class="btn btn--text btn--block"
                type="button"
            >
                Remove photo
            </button>
        </div>
    </aside>
</template>

<style scoped>
@import '../../styles/variables.css';

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