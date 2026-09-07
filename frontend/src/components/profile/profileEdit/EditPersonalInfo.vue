<script setup>
import { reactive } from 'vue';

import { updateUserInfo } from '@/api/users/editProfile';
import AvatarUploader from './AvatarUploader.vue';
import FormField from './FormField.vue';

const props = defineProps({
    firstName: String,
    lastName: String,
    username: String,
    email: String,
    bio: String,
    avatar_path: String,
    isPrivate: Boolean
});

const form = reactive({
    FirstName: props.firstName || '',
    LastName: props.lastName || '',
    Username: props.username || '',
    Email: props.email || '',
    About: props.bio || '',
    Password: '',
    IsPrivate: props.isPrivate || false
});

function togglePrivacy() {
    form.IsPrivate = !form.IsPrivate;
}

async function updateInfo() {
    if(form.IsPrivate) {
        form.IsPrivate = 1
    } else {
        form.IsPrivate = 0
    }
    const result = await updateUserInfo(form);

    if (!result.status) {
        console.error('failed to update data');
        return;
    }

    console.log('Profile updated successfully');
}
</script>

<template>
    <section class="edit-section">
        <div class="section-heading">
            <p class="eyebrow">PROFILE</p>
            <h2>Personal Info</h2>
        </div>

        <div class="edit-card">
            <AvatarUploader
                :src="props.avatar_path"
                @change="onAvatarChange"
            />

            <div class="privacy-setting">
                <div>
                    <p class="privacy-title">
                        {{ form.IsPrivate ? 'Private account' : 'Public account' }}
                    </p>

                    <p class="privacy-description">
                        {{
                            form.IsPrivate
                                ? 'Only approved followers can see your posts.'
                                : 'Anyone can see your posts and profile.'
                        }}
                    </p>
                </div>

                <button
                    type="button"
                    class="privacy-button"
                    :class="{ private: form.IsPrivate }"
                    @click="togglePrivacy"
                >
                    {{ form.IsPrivate ? 'Make public' : 'Make private' }}
                </button>
            </div>

            <form @submit.prevent="updateInfo">
                <div class="field-grid">
                    <FormField
                        id="firstName"
                        label="First name"
                        v-model="form.FirstName"
                        placeholder="First name"
                    />

                    <FormField
                        id="lastName"
                        label="Last name"
                        v-model="form.LastName"
                        placeholder="Last name"
                    />

                    <FormField
                        id="username"
                        label="Username"
                        v-model="form.Username"
                        placeholder="Username"
                    />

                    <FormField
                        id="email"
                        label="Email"
                        type="email"
                        v-model="form.Email"
                        placeholder="Email address"
                    />

                    <FormField
                        id="password"
                        label="Password"
                        type="password"
                        v-model="form.Password"
                        placeholder="****************"
                    />
                </div>

                <FormField
                    id="bio"
                    label="Bio"
                    type="textarea"
                    v-model="form.About"
                    placeholder="Tell people about yourself"
                />

                <button type="submit" class="confirm-button">
                    Confirm changes
                </button>
            </form>
        </div>
    </section>
</template>

<style scoped>
.edit-section {
    scroll-margin-top: 6.25rem;
    width: 100%;
    max-width: 100%;
}

.section-heading {
    margin-bottom: var(--space-4);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: clamp(1.25rem, 4vw, 1.75rem);
    color: var(--color-text);
}

.edit-card {
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    max-width: 100%;
    box-sizing: border-box;
}

.edit-card form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
}

.privacy-setting {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-raised);
}

.privacy-title {
    margin: 0 0 var(--space-1);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 700;
    color: var(--color-text);
}

.privacy-description {
    margin: 0;
    font-family: var(--font-body);
    font-size: 0.75rem;
    color: var(--color-text-muted);
}

.privacy-button {
    flex-shrink: 0;
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: 1.5625rem;
    background: var(--color-surface);
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.privacy-button:hover {
    filter: brightness(1.1);
}

.privacy-button.private {
    background: var(--gradient-action);
    border-color: transparent;
    color: var(--color-text);
}

.field-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(13.75rem, 100%), 1fr));
    gap: var(--space-4);
}

.confirm-button {
    align-self: flex-start;
    padding: var(--space-3) var(--space-5);
    border: none;
    border-radius: 1.5625rem;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    font-weight: 700;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.confirm-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}
</style>
