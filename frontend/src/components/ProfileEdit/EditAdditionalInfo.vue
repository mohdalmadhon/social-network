<script setup>
<<<<<<< HEAD
import { updateAbout } from '@/api/users/editProfile';
import { addNotification } from '@/data/notifications';
import { validateAboutField } from '@/helpers/validators/editProfile.js';
import { reactive } from 'vue';
=======
import { computed, reactive, ref } from 'vue'
import { updateAbout } from '@/api/users/editProfile'
import { profileData } from '@/data/usersData'
import { validateAboutField } from '@/helpers/validators/editProfile.js'
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53

const props = defineProps({
    about: {
        type: Object,
        default: () => ({})
    }
});

const form = reactive({
<<<<<<< HEAD
    work: props.about?.work || '',
    education: props.about?.education || '',
    hobbies: props.about?.hobbies || '',
    interests: props.about?.interests || '',
    travel: props.about?.travel || '',
    website: props.about?.website || '',
    linkedin: props.about?.linkedin || '',
    twitter: props.about?.twitter || '',
    instagram: props.about?.instagram || ''
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
=======
  work: profileData.about.work || '',
  education: profileData.about.education || '',
  hobbies: profileData.about.hobbies || '',
  interests: profileData.about.intrests || '',
  travel: profileData.about.travel || '',
  website: profileData.about.website || '',
  linkedin: profileData.about.linkedin || '',
  twitter: profileData.about.twitter || '',
  instagram: profileData.about.instgram || '',
})
const errors = reactive(Object.fromEntries(Object.keys(form).map((key) => [key, ''])))
const original = ref(JSON.stringify(form))
const saving = ref(false)
const feedback = ref('')
const saveError = ref('')
const isDirty = computed(() => JSON.stringify(form) !== original.value)
const informationFields = [
  { id: 'work', label: 'Work', placeholder: 'Where do you work?' },
  { id: 'education', label: 'Education', placeholder: 'Where did you study?' },
  { id: 'hobbies', label: 'Hobbies', placeholder: 'What do you enjoy doing?' },
  { id: 'interests', label: 'Interests', placeholder: 'What are you interested in?' },
  { id: 'travel', label: 'Travel', placeholder: 'Places you have visited' },
]
const linkFields = [
  { id: 'website', label: 'Website', placeholder: 'https://example.com' },
  { id: 'linkedin', label: 'LinkedIn', placeholder: 'https://linkedin.com/in/...' },
  { id: 'twitter', label: 'Twitter / X', placeholder: 'https://x.com/...' },
  { id: 'instagram', label: 'Instagram', placeholder: 'https://instagram.com/...' },
]
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53

function validateField(field) {
  errors[field] = validateAboutField(form[field])
}

<<<<<<< HEAD
async function confirmChanges() {
    Object.keys(form).forEach(validateField);

    const hasErrors = Object.values(errors).some(error => error !== '');

    if (hasErrors) {
        addNotification('Failed to connect to server', 'success');
        return;
    }

    try {
        const result = await updateAbout(form);
        if (!result.status) {
            addNotification(result.message, 'error')
            return;
        }
        addNotification(result.message, 'success')
    } catch (err) {
        addNotification(err.message, 'error')
    }
=======
async function save() {
  Object.keys(form).forEach(validateField)
  if (Object.values(errors).some(Boolean) || saving.value || !isDirty.value) return
  saving.value = true
  feedback.value = ''
  saveError.value = ''
  try {
    const result = await updateAbout(form)
    Object.assign(profileData.about, {
      work: form.work,
      education: form.education,
      hobbies: form.hobbies,
      intrests: form.interests,
      travel: form.travel,
      website: form.website,
      linkedin: form.linkedin,
      twitter: form.twitter,
      instgram: form.instagram,
    })
    original.value = JSON.stringify(form)
    feedback.value = result.message || 'Additional information saved.'
  } catch (err) {
    saveError.value = err.message || 'Could not save additional information.'
  } finally {
    saving.value = false
  }
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
}
</script>

<template>
  <section class="additional-section orbit-surface" aria-labelledby="additional-heading">
    <header>
      <p class="orbit-meta">Profile details</p>
      <h2 id="additional-heading">Additional information</h2>
      <p>Add context about your work, interests, and places people can find you.</p>
    </header>

    <form @submit.prevent="save">
      <fieldset>
        <legend>About you</legend>
        <div class="field-grid">
          <label v-for="field in informationFields" :key="field.id" :for="field.id">
            <span>{{ field.label }}</span>
            <textarea :id="field.id" v-model="form[field.id]" maxlength="200" :placeholder="field.placeholder" @input="validateField(field.id)" />
            <small v-if="errors[field.id]" role="alert">{{ errors[field.id] }}</small>
          </label>
        </div>
      </fieldset>

<<<<<<< HEAD
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
=======
      <fieldset>
        <legend>Social links</legend>
        <div class="field-grid">
          <label v-for="field in linkFields" :key="field.id" :for="field.id">
            <span>{{ field.label }}</span>
            <input :id="field.id" v-model="form[field.id]" type="url" maxlength="200" :placeholder="field.placeholder" @input="validateField(field.id)" />
            <small v-if="errors[field.id]" role="alert">{{ errors[field.id] }}</small>
          </label>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
        </div>
      </fieldset>

      <p v-if="feedback" class="form-feedback" role="status">{{ feedback }}</p>
      <p v-if="saveError" class="form-feedback form-feedback--error" role="alert">{{ saveError }}</p>
      <div class="form-actions">
        <RouterLink to="/me">Cancel</RouterLink>
        <button type="submit" :disabled="saving || !isDirty">{{ saving ? 'Saving...' : 'Save changes' }}</button>
      </div>
    </form>
  </section>
</template>

<style scoped>
<<<<<<< HEAD
.edit-section {
    scroll-margin-top: 100px;
    width: 100%;
    max-width: 100%;
}
.error {
    margin: 0;
    color: #d9534f;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(8px, 1.2vw, 9px);
}
.section-heading {
    margin-bottom: clamp(14px, 2.5vw, 20px);
}

.links-heading {
    margin-top: 5px;
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: clamp(20px, 4vw, 29px);
}

.edit-card {
    display: flex;
    flex-direction: column;
    gap: clamp(16px, 2.5vw, 22px);
    padding: clamp(14px, 3vw, 25px);
    border: 2px solid var(--main-color);
    border-radius: 7px;
    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
    max-width: 100%;
    box-sizing: border-box;
}

.field-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(220px, 100%), 1fr));
    gap: clamp(14px, 2.5vw, 20px);
}

.form-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
}

.form-field label {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(8px, 1.2vw, 9px);
    font-weight: 600;
    letter-spacing: 1px;
}

.form-field textarea {
    width: 100%;
    min-height: 90px;
    padding: clamp(10px, 1.8vw, 12px) clamp(10px, 2vw, 14px);
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--bg-color);
    color: var(--font-color);
    font-family: "Hedvig Letters Sans", sans-serif;
    font-size: clamp(12px, 1.8vw, 13px);
    resize: vertical;
    box-sizing: border-box;
}

.form-field textarea:focus {
    outline: none;
    border-color: var(--input-focus);
    box-shadow: 3px 3px var(--input-focus);
}

.confirm-button {
    align-self: flex-start;
    padding: clamp(11px, 2vw, 13px) clamp(16px, 3vw, 22px);
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--input-focus);
    box-shadow: 4px 4px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
}

.confirm-button:hover {
    transform: translate(-1px, -1px);
}
</style>
=======
.additional-section { padding: var(--space-5); }
.additional-section > header { margin-bottom: var(--space-5); }
.additional-section .orbit-meta { margin: 0; color: var(--color-violet-soft); }
.additional-section h2 { margin: var(--space-1) 0; font-family: var(--font-display); font-size: 1.5rem; letter-spacing: 0; }
.additional-section header > p:last-child { margin: 0; color: var(--color-text-muted); }
form { display: grid; gap: var(--space-5); }
fieldset { margin: 0; padding: var(--space-5) 0 0; border: 0; border-top: 1px solid var(--color-border); }
legend { padding: 0; font-weight: 700; }
.field-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: var(--space-4); margin-top: var(--space-4); }
label { display: grid; min-width: 0; gap: var(--space-2); color: var(--color-text-muted); font-size: .8125rem; font-weight: 600; }
input, textarea { width: 100%; min-width: 0; padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-small); outline: none; background: var(--color-input); color: var(--color-text); }
textarea { min-height: 6.5rem; line-height: 1.5; }
input:focus, textarea:focus { border-color: var(--color-violet); box-shadow: var(--focus-ring); }
label small { color: var(--color-coral-soft); font-weight: 400; }
.form-feedback { margin: 0; padding: var(--space-3) var(--space-4); border-left: 3px solid var(--color-mint); background: var(--color-surface-teal); color: var(--color-text-soft); }
.form-feedback--error { border-color: var(--color-coral); background: var(--color-surface-coral); color: var(--color-coral-soft); }
.form-actions { display: flex; justify-content: flex-end; gap: var(--space-3); padding-top: var(--space-4); border-top: 1px solid var(--color-border); }
.form-actions a, .form-actions button { display: inline-flex; min-height: var(--touch-target); align-items: center; justify-content: center; padding: 0 var(--space-4); border-radius: var(--radius-small); font-weight: 700; text-decoration: none; }
.form-actions a { border: 1px solid var(--color-border); color: var(--color-text-soft); }
.form-actions button { border: 0; background: var(--gradient-action); color: white; cursor: pointer; }
@media (max-width: 620px) {
  .additional-section { padding: var(--space-4); }
  .field-grid { grid-template-columns: 1fr; }
  .form-actions { display: grid; }
}
</style>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
