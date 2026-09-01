<script setup>
import { updateAbout } from '@/api/users/editProfile';
import { addNotification } from '@/data/notifications';
import { profileData } from '@/data/usersData';
import { validateAboutField } from '@/helpers/validators/editProfile.js';
import { reactive } from 'vue';

const form = reactive({
    work: profileData.about.work || '',
    education: profileData.about.education || '',
    hobbies: profileData.about.hobbies || '',
    interests: profileData.about.intrests || '',
    travel: profileData.about.travel || '',
    website: profileData.about.website || '',
    linkedin: profileData.about.linkedin || '',
    twitter: profileData.about.twitter || '',
    instagram: profileData.about.instgram || ''
});

const errors = reactive({
    work: '',
    education: '',
    hobbies: '',
    interests: '',
    travel: '',
    website: '',
    linkedin: '',
    twitter: '',
    instagram: ''
});

function validateField(field) {
    errors[field] = validateAboutField(form[field]);
}

async function confirmChanges() {
    Object.keys(form).forEach(validateField);

    const hasErrors = Object.values(errors).some(error => error !== '');

    if (hasErrors) {
        addNotification('Failed to connect to server', 'success');
        return;
    }

    const result = await updateAbout(form);

    console.log(result);
}
</script>

<template>
    <section class="edit-section">
        <div class="section-heading">
            <p class="eyebrow">PROFILE</p>
            <h2>Additional Info</h2>
        </div>

        <div class="edit-card">
            <div class="field-grid">
                <div class="form-field">
                    <label for="work">Work</label>
                    <textarea
                        id="work"
                        v-model="form.work"
                        maxlength="200"
                        placeholder="Where do you work"
                        @input="validateField('work')"
                    ></textarea>

                    <p v-if="errors.work" class="error">
                        {{ errors.work }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="education">Education</label>
                    <textarea
                        id="education"
                        v-model="form.education"
                        maxlength="200"
                        placeholder="Where did you study"
                        @input="validateField('education')"
                    ></textarea>

                    <p v-if="errors.education" class="error">
                        {{ errors.education }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="hobbies">Hobbies</label>
                    <textarea
                        id="hobbies"
                        v-model="form.hobbies"
                        maxlength="200"
                        placeholder="Your hobbies"
                        @input="validateField('hobbies')"
                    ></textarea>

                    <p v-if="errors.hobbies" class="error">
                        {{ errors.hobbies }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="interests">Interests</label>
                    <textarea
                        id="interests"
                        v-model="form.interests"
                        maxlength="200"
                        placeholder="Your interests"
                        @input="validateField('interests')"
                    ></textarea>

                    <p v-if="errors.interests" class="error">
                        {{ errors.interests }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="travel">Travel</label>
                    <textarea
                        id="travel"
                        v-model="form.travel"
                        maxlength="200"
                        placeholder="Places you've been"
                        @input="validateField('travel')"
                    ></textarea>

                    <p v-if="errors.travel" class="error">
                        {{ errors.travel }}
                    </p>
                </div>
            </div>

            <div class="details-divider"></div>

            <div class="section-heading links-heading">
                <p class="eyebrow">SOCIAL</p>
                <h2>Links</h2>
            </div>

            <div class="field-grid">
                <div class="form-field">
                    <label for="website">Website</label>
                    <textarea
                        id="website"
                        v-model="form.website"
                        maxlength="200"
                        placeholder="https://"
                        @input="validateField('website')"
                    ></textarea>

                    <p v-if="errors.website" class="error">
                        {{ errors.website }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="linkedin">LinkedIn</label>
                    <textarea
                        id="linkedin"
                        v-model="form.linkedin"
                        maxlength="200"
                        placeholder="https://linkedin.com/in/"
                        @input="validateField('linkedin')"
                    ></textarea>

                    <p v-if="errors.linkedin" class="error">
                        {{ errors.linkedin }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="twitter">Twitter / X</label>
                    <textarea
                        id="twitter"
                        v-model="form.twitter"
                        maxlength="200"
                        placeholder="https://x.com/"
                        @input="validateField('twitter')"
                    ></textarea>

                    <p v-if="errors.twitter" class="error">
                        {{ errors.twitter }}
                    </p>
                </div>

                <div class="form-field">
                    <label for="instagram">Instagram</label>
                    <textarea
                        id="instagram"
                        v-model="form.instagram"
                        maxlength="200"
                        placeholder="https://instagram.com/"
                        @input="validateField('instagram')"
                    ></textarea>

                    <p v-if="errors.instagram" class="error">
                        {{ errors.instagram }}
                    </p>
                </div>
            </div>

            <button
                type="button"
                class="confirm-button"
                @click="confirmChanges"
            >
                Confirm changes
            </button>
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

.error {
    margin: 0;
    color: var(--color-danger, #ef4444);
    font-size: 0.75rem;
}

.section-heading {
    margin-bottom: var(--space-5);
}

.links-heading {
    margin-top: 0;
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

.details-divider {
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--color-border), transparent);
}

.field-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--space-4);
}

.form-field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

.form-field label {
    color: var(--color-text-soft);
    font-size: 0.82rem;
    font-weight: 600;
}

.form-field textarea {
    width: 100%;
    min-height: 6.5rem;
    padding: var(--space-3);
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    color: var(--color-text);
    font: inherit;
    line-height: 1.55;
    resize: vertical;
    box-sizing: border-box;
    transition:
        border-color 0.18s ease,
        background 0.18s ease,
        box-shadow 0.18s ease;
}

.form-field textarea::placeholder {
    color: var(--color-text-faint);
}

.form-field textarea:hover {
    border-color: var(--color-text-faint);
    background: var(--color-surface-raised);
}

.form-field textarea:focus {
    outline: none;
    border-color: var(--color-violet);
    background: var(--color-surface-raised);
    box-shadow: var(--focus-ring);
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
