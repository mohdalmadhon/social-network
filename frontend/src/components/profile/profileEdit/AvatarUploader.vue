<script setup>
import { ref } from 'vue';
import { checkSessionResponse } from '@/helpers/auth/auth';
import { addNotification } from '@/data/notifications';
import { router } from '@/router/router';

const props = defineProps({
    src: {
        type: String,
        default: ''
    }
});

const emit = defineEmits(['change']);

const preview = ref(props.src);
const fileInput = ref(null);
const uploading = ref(false);
const error = ref('');

const allowedTypes = ['image/png', 'image/jpeg', 'image/gif'];
const maxFileSize = 5 * 1024 * 1024;

function triggerUpload() {
    if (!uploading.value) {
        fileInput.value?.click();
    }
}

function validateFile(file) {
    if (!allowedTypes.includes(file.type)) {
        return 'Only JPG, GIF and PNG images are allowed';
    }

    if (file.size > maxFileSize) {
        return 'Image must be smaller than 5MB';
    }

    return '';
}

async function uploadAvatar(file) {
    const formData = new FormData();
    formData.append('avatar', file);

    uploading.value = true;

    try {
        const response = await fetch('/api/profile/avatar', {
            method: 'PATCH',
            credentials: 'include',
            body: formData
        });

        const result = await response.json();

        if (!checkSessionResponse(response)) {
            router.replace('/login');
            return;
        }

        if (!response.ok) {
            throw new Error(result.message || 'Failed to update avatar');
        }

        emit('change', file);

        addNotification(
            result.message || 'Avatar updated successfully',
            'success'
        );
    } catch (err) {
        console.error(err);

        preview.value = props.src;

        addNotification(
            err.message || 'Failed to update avatar',
            'error'
        );
    } finally {
        uploading.value = false;
    }
}

async function onFileChange(event) {
    const file = event.target.files?.[0];

    if (!file) {
        return;
    }

    error.value = '';

    const validationError = validateFile(file);

    if (validationError) {
        error.value = validationError;
        event.target.value = '';
        return;
    }

    const previousPreview = preview.value;

    preview.value = URL.createObjectURL(file);

    try {
        await uploadAvatar(file);
    } catch (err) {
        preview.value = previousPreview;
    }

    event.target.value = '';
}
</script>

<template>
    <div class="avatar-uploader">
        <div class="avatar-preview">
            <img
                v-if="preview"
                :src="preview"
                alt="Profile avatar"
            >
        </div>

        <div class="avatar-actions">
            <button
                type="button"
                class="upload-button"
                :disabled="uploading"
                @click="triggerUpload"
            >
                {{ uploading ? 'Uploading...' : 'Change avatar' }}
            </button>

            <span class="hint">
                JPG, GIF or PNG, up to 5MB
            </span>

            <span
                v-if="error"
                class="error"
            >
                {{ error }}
            </span>
        </div>

        <input
            ref="fileInput"
            type="file"
            accept="image/png, image/jpeg, image/gif"
            class="hidden-input"
            @change="onFileChange"
        >
    </div>
</template>

<style scoped>
@import '../../../styles/global.css';

.avatar-uploader {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-4);
}

.avatar-preview {
    flex-shrink: 0;
    width: clamp(72px, 12vw, 96px);
    height: clamp(72px, 12vw, 96px);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background: var(--color-surface-raised);
    border: 4px solid var(--color-surface);
    border-radius: 50%;
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.2);
}

.avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.avatar-actions {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

.upload-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-target);
    padding-inline: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text-soft);
    font: inherit;
    font-weight: 600;
    white-space: nowrap;
    cursor: pointer;
    transition:
        border-color 0.15s ease,
        color 0.15s ease,
        background 0.15s ease;
}

.upload-button:hover:not(:disabled),
.upload-button:focus-visible:not(:disabled) {
    border-color: var(--color-violet);
    color: var(--color-text);
    background: var(--color-surface-raised);
    outline: none;
}

.upload-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.hint {
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.75rem;
}

.error {
    color: var(--color-danger, #ef4444);
    font-size: 0.75rem;
}

.hidden-input {
    display: none;
}

@media (max-width: 25rem) {
    .avatar-uploader {
        flex-direction: column;
        align-items: flex-start;
    }
}
</style>
