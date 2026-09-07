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
.edit-section {
    scroll-margin-top: 6.25rem;
    width: 100%;
    max-width: 100%;
}

.error {
    margin: 0;
    color: var(--color-coral);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

.section-heading {
    margin-bottom: var(--space-4);
}

.links-heading {
    margin-top: var(--space-1);
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

.field-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(13.75rem, 100%), 1fr));
    gap: var(--space-4);
}

.form-field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

.form-field label {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
}

.form-field textarea {
    width: 100%;
    min-height: 5.625rem;
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    resize: vertical;
    box-sizing: border-box;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-field textarea::placeholder {
    color: var(--color-text-faint);
}

.form-field textarea:focus {
    outline: none;
    border-color: var(--color-violet);
    box-shadow: var(--focus-ring);
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
