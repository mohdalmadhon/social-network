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
.avatar-uploader {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-4);
}

.avatar-preview {
    flex-shrink: 0;
    width: 5rem;
    height: 5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border: 0.1875rem solid var(--color-surface);
    outline: 2px solid transparent;
    border-radius: 50%;
    background: var(--gradient-action);
    box-shadow: var(--shadow-raised);
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
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: 1.5625rem;
    background: var(--color-surface-raised);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 600;
    white-space: nowrap;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.upload-button:hover:not(:disabled) {
    filter: brightness(1.15);
}

.upload-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.hint {
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
}

.error {
    color: var(--color-coral);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
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
