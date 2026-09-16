<script setup>
import { reactive, ref, watch } from 'vue';

import { addNotification } from '@/data/notifications';
import AvatarPicker from './AvatarPicker.vue';
import AddMembers from './AddMembers.vue';

const props = defineProps({
    show: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['close', 'created']);

const NAME_LIMIT = 15;
const DESCRIPTION_LIMIT = 200;
const totalSlides = 3;

const group = reactive({
    name: '',
    description: '',
    avatar: null,
    members: []
});

const currentSlide = ref(1);
const submitting = ref(false);

const validation = ref({
    field: null,
    message: null
});

function validateGroup() {
    if (!group.name.trim()) {
        return {
            field: 'name',
            message: 'group name is required'
        };
    }

    if (group.name.length > NAME_LIMIT) {
        return {
            field: 'name',
            message: `group name must be ${NAME_LIMIT} characters or fewer`
        };
    }

    if (group.description.length > DESCRIPTION_LIMIT) {
        return {
            field: 'description',
            message: `description must be ${DESCRIPTION_LIMIT} characters or fewer`
        };
    }

    return {
        field: null,
        message: null
    };
}

watch(
    group,
    () => {
        validation.value = validateGroup();
    },
    { deep: true }
);

function nextSlide() {
    if (currentSlide.value === 1 && validation.value.field) {
        addNotification(validation.value.message, 'error');
        return;
    }

    if (currentSlide.value < totalSlides) {
        currentSlide.value++;
        validation.value = {
            field: null,
            message: null
        };
    }
}

function previousSlide() {
    if (currentSlide.value > 1) {
        currentSlide.value--;
        validation.value = {
            field: null,
            message: null
        };
    }
}

function resetForm() {
    currentSlide.value = 1;
    group.name = '';
    group.description = '';
    group.avatar = null;
    group.members = [];
    validation.value = {
        field: null,
        message: null
    };
    submitting.value = false;
}

function closeDialog() {
    emit('close');
    resetForm();
}

function handleOverlayClick() {
    if (!submitting.value) {
        closeDialog();
    }
}

async function handleSubmit() {
    const validationResult = validateGroup();

    if (validationResult.field) {
        validation.value = validationResult;
        addNotification(validationResult.message, 'error');
        currentSlide.value = 1;
        return;
    }

    if (submitting.value) {
        return;
    }

    submitting.value = true;

    try {
        const formData = new FormData();
        
        formData.append('title', group.name.trim());
        formData.append('description', group.description.trim());
        formData.append(
            'users',
            JSON.stringify(group.members.map(member => member.id))
        );

        if (group.avatar) {
            formData.append('avatar', group.avatar);
        }

        const response = await fetch('/api/groups', {
            method: 'POST',
            credentials: 'include',
            body: formData
        });

        const result = await response.json();

        if (!result.status) {
            addNotification(result.message || 'could not create group', 'error');
            return;
        }

        addNotification('group created', 'success');
        emit('created');
        closeDialog();
    } catch (err) {
        addNotification(err.message || 'could not create group', 'error');
    } finally {
        submitting.value = false;
    }
}
</script>

<template>
    <div v-if="show" class="dialog-overlay" @click.self="handleOverlayClick">
        <form class="create-group-card" @submit.prevent="handleSubmit">
            <button class="close-button" type="button" :disabled="submitting" @click="closeDialog">
                ×
            </button>

            <div class="progress">
                <div v-for="slide in totalSlides" :key="slide" class="progress-step" :class="{
                    active: slide === currentSlide,
                    completed: slide < currentSlide
                }">
                    {{ slide }}
                </div>
            </div>

            <div class="slide-title">
                <p class="slide-number">
                    STEP {{ currentSlide }} / {{ totalSlides }}
                </p>

                <h2 v-if="currentSlide === 1">
                    Name & Description
                </h2>

                <h2 v-else-if="currentSlide === 2">
                    Group Photo
                </h2>

                <h2 v-else>
                    Add Members
                </h2>
            </div>

            <div v-if="currentSlide === 1" class="slide">
                <div class="field">
                    <label>
                        <strong>Group name</strong>
                    </label>
                    
                    <label style="font-size: 8px;">
                        {{ group.name.length }}/{{ NAME_LIMIT }} characters
                    </label>

                    <input
                        v-model="group.name"
                        type="text"
                        name="title"
                        maxlength="15"
                        placeholder="goats"
                    >
                </div>

                <p v-if="validation.field === 'name'" class="validation-error">
                    {{ validation.message }}
                </p>

                <div class="field">
                    <label>
                        <strong>Description</strong>
                    </label>

                    <label style="font-size: 8px;">
                        {{ group.description.length }}/{{ DESCRIPTION_LIMIT }} characters
                    </label>

                    <textarea
                        v-model="group.description"
                        maxlength="200"
                        rows="4"
                        name="description"
                        placeholder="What's this group about?"
                    ></textarea>
                </div>

                <p v-if="validation.field === 'description'" class="validation-error">
                    {{ validation.message }}
                </p>
            </div>

            <div v-else-if="currentSlide === 2" class="slide">
                <AvatarPicker v-model="group.avatar" />
            </div>

            <div v-else class="slide">
                <AddMembers v-model="group.members" />
            </div>

            <div class="navigation-buttons">
                <button v-if="currentSlide > 1" class="previous-button" type="button" @click="previousSlide">
                    Previous
                </button>

                <button v-if="currentSlide < totalSlides" class="next-button" type="button" @click="nextSlide">
                    Next
                </button>

                <button v-else class="submit-button" type="submit" :disabled="submitting">
                    {{ submitting ? 'Creating...' : 'Create group' }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
.dialog-overlay {
    position: fixed;
    inset: 0;
    z-index: 100;

    display: flex;
    align-items: center;
    justify-content: center;

    padding: 20px;

    background: rgba(50, 50, 50, 0.5);
}

.create-group-card {
    position: relative;

    width: 100%;
    max-width: 480px;
    max-height: calc(100vh - 40px);
    overflow-y: auto;

    display: flex;
    flex-direction: column;
    gap: 20px;

    padding: 32px 34px;
    box-sizing: border-box;

    border: 2px solid var(--main-color);
    border-radius: 7px;

    background: var(--bg-color);
    box-shadow: 6px 6px var(--main-color);
}

.close-button {
    position: absolute;
    top: 16px;
    right: 16px;

    width: 28px;
    height: 28px;

    display: flex;
    align-items: center;
    justify-content: center;

    border: 2px solid var(--main-color);
    border-radius: 50%;

    background: var(--bg-color);
    color: var(--main-color);

    font-family: "JetBrains Mono", monospace;
    font-size: 15px;
    line-height: 1;

    transition:
        transform 0.1s,
        background 0.15s,
        color 0.15s;
}

.close-button:hover:not(:disabled) {
    background: var(--main-color);
    color: var(--bg-color);
}

.close-button:disabled {
    opacity: 0.5;
    cursor: default;
}

.progress {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
}

.progress-step {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
    font-weight: 700;
    background: var(--bg-color);
    color: var(--main-color);
}

.progress-step.active {
    background: var(--input-focus);
    color: white;
}

.progress-step.completed {
    background: var(--main-color);
    color: white;
}

.slide-title {
    border-bottom: 2px solid var(--main-color);
    padding-bottom: 16px;
}

.slide-number {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    letter-spacing: 2px;
}

.slide-title h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: 22px;
}

.slide {
    min-height: 200px;
    display: flex;
    flex-direction: column;
    gap: 18px;
}

.field {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.field label {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.field input,
.field textarea {
    width: 100%;
    box-sizing: border-box;
    padding: 12px 14px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    outline: none;
    background: var(--bg-color);
    color: var(--font-color);
    font-family: "Hedvig Letters Sans", sans-serif;
    font-size: 14px;
    transition: box-shadow 0.15s ease;
}

.field textarea {
    resize: none;
    line-height: 1.5;
}

.field input:focus,
.field textarea:focus {
    box-shadow: 3px 3px var(--main-color);
}

.validation-error {
    margin: -12px 0 0;
    color: #d93025;
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
    font-weight: 700;
}

.navigation-buttons {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 15px;
}

.previous-button,
.next-button,
.submit-button {
    padding: 12px 26px;
    border: 2px solid var(--main-color);
    border-radius: 6px;
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
    font-weight: 700;
    transition: transform 0.1s ease, box-shadow 0.1s ease;
}

.previous-button:disabled,
.next-button:disabled,
.submit-button:disabled {
    opacity: 0.6;
    cursor: default;
}

.previous-button {
    background: var(--bg-color);
    color: var(--main-color);
    box-shadow: 4px 4px var(--main-color);
}

.next-button,
.submit-button {
    margin-left: auto;
    background: var(--input-focus);
    color: white;
    box-shadow: 4px 4px var(--main-color);
}

.previous-button:hover,
.next-button:hover,
.submit-button:hover {
    transform: translate(-1px, -1px);
    box-shadow: 5px 5px var(--main-color);
}

.previous-button:active,
.next-button:active,
.submit-button:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px var(--main-color);
}

@media (max-width: 550px) {
    .create-group-card {
        padding: 26px 20px;
        border-radius: 12px;
        gap: 16px;
    }

    .slide-title h2 {
        font-size: 19px;
    }

    .navigation-buttons {
        flex-direction: column-reverse;
        gap: 10px;
    }

    .previous-button,
    .next-button,
    .submit-button {
        width: 100%;
        margin-left: 0;
        text-align: center;
    }
}
</style>