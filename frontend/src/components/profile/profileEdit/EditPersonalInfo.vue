<script setup>
import { updateUserInfo } from '@/api/users/editProfile';
import { reactive } from 'vue';
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
@import '../../../styles/global.css';

.edit-section {
    scroll-margin-top: 100px;
    width: 100%;
    max-width: 100%;
}

.section-heading {
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.4rem;
    font-weight: 700;
    letter-spacing: -0.02em;
}

.edit-card {
    position: relative;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
    padding: var(--space-6);
    background:
        linear-gradient(145deg, rgb(124 92 255 / 5%), transparent 32%),
        var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    box-shadow: var(--shadow-raised);
    max-width: 100%;
    box-sizing: border-box;
}

.edit-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 2px;
    background: var(--gradient-action);
}

form {
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
}

.privacy-setting {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-4);
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
}

.privacy-title {
    margin: 0 0 var(--space-1);
    color: var(--color-text);
    font-size: 0.9rem;
    font-weight: 700;
}

.privacy-description {
    margin: 0;
    color: var(--color-text-faint);
    font-size: 0.82rem;
}

.privacy-button {
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-target);
    padding-inline: var(--space-4);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-action);
    color: #fff;
    font: inherit;
    font-weight: 600;
    cursor: pointer;
    transition: filter 0.15s ease, transform 0.15s ease;
}

.privacy-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

.privacy-button.private {
    background: var(--color-input);
    border-color: var(--color-border);
    color: var(--color-text-soft);
}

.field-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--space-4);
}

.confirm-button {
    align-self: flex-start;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-target);
    padding-inline: var(--space-5);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-action);
    color: #fff;
    font: inherit;
    font-weight: 600;
    cursor: pointer;
    box-shadow: var(--shadow-raised);
    transition: filter 0.15s ease, transform 0.15s ease;
}

.confirm-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

@media (min-width: 40rem) {
    .field-grid {
        grid-template-columns: 1fr 1fr;
    }
}

@media (min-width: 64rem) {
    .edit-card {
        padding: var(--space-7);
    }
}
</style>
