<script setup>
import { reactive } from 'vue';

import FormField from '@/components/ProfileEdit/FormField.vue';
import AvatarUploader from '@/components/ProfileEdit/AvatarUploader.vue';
import { updateUserInfo } from '@/api/users/editProfile';
import { addNotification } from '@/data/notifications';

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
        addNotification(result.message);
        return;
    }

    addNotification(result.message);
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
                        <span
                            class="status-dot"
                            :class="{ private: form.IsPrivate }"
                        ></span>
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
@import '../../styles/global.css';
@import '../../styles/variables.css';
.edit-section {
    scroll-margin-top: 100px;
    width: 100%;
    max-width: 100%;
}

.section-heading {
    margin-bottom: clamp(14px, 2.5vw, 20px);
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: clamp(20px, 4vw, 29px);
}

.edit-card {
    display: flex;
    flex-direction: column;
    gap: clamp(16px, 2.5vw, 22px);
    padding: clamp(14px, 3vw, 25px);
    border: 1px solid var(--color-border);
    border-top: 3px solid var(--color-violet);
    border-radius: var(--radius-medium);
    background: var(--color-surface-violet);
    box-shadow: var(--shadow-raised);
    max-width: 100%;
    box-sizing: border-box;
}

.privacy-setting {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: clamp(12px, 2.5vw, 20px);
    padding: clamp(12px, 2.5vw, 18px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-surface);
}

.privacy-title {
    margin: 0 0 5px;
    display: flex;
    align-items: center;
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: clamp(11px, 1.6vw, 12px);
    font-weight: 700;
}

.status-dot {
    width: 7px;
    height: 7px;
    margin-right: 7px;
    border-radius: 50%;
    background: var(--color-mint);
    box-shadow: 0 0 6px var(--color-mint);
}

.status-dot.private {
    background: var(--color-coral);
    box-shadow: 0 0 6px var(--color-coral);
}

.privacy-description {
    margin: 0;
    color: var(--color-text-muted);
    font-size: clamp(11px, 1.6vw, 12px);
}

.privacy-button {
    flex-shrink: 0;
    padding: clamp(9px, 1.8vw, 11px) clamp(14px, 2.5vw, 18px);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-aurora);
    box-shadow: var(--shadow-raised);
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.1s;
}

.privacy-button:hover {
    transform: translateY(-1px);
}

.privacy-button.private {
    background: var(--color-surface-coral);
    border-color: var(--color-coral);
    color: var(--color-coral-soft);
}

.field-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(220px, 100%), 1fr));
    gap: clamp(14px, 2.5vw, 20px);
}

.confirm-button {
    align-self: flex-start;
    padding: clamp(11px, 2vw, 13px) clamp(16px, 3vw, 22px);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-cyber);
    box-shadow: var(--shadow-raised);
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.1s;
}

.confirm-button:hover {
    transform: translateY(-1px);
}
</style>