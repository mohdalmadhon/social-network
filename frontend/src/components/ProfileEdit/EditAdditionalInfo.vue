<script setup>
import { computed, reactive, ref } from 'vue'
import { updateAbout } from '@/api/users/editProfile'
import { profileData } from '@/data/usersData'
import { validateAboutField } from '@/helpers/validators/editProfile.js'

const form = reactive({
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

function validateField(field) {
  errors[field] = validateAboutField(form[field])
}

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

      <fieldset>
        <legend>Social links</legend>
        <div class="field-grid">
          <label v-for="field in linkFields" :key="field.id" :for="field.id">
            <span>{{ field.label }}</span>
            <input :id="field.id" v-model="form[field.id]" type="url" maxlength="200" :placeholder="field.placeholder" @input="validateField(field.id)" />
            <small v-if="errors[field.id]" role="alert">{{ errors[field.id] }}</small>
          </label>
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
