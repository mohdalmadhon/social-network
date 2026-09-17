<script setup>
import { computed, reactive, ref } from 'vue'
import FormField from '@/components/ProfileEdit/FormField.vue'
import AvatarUploader from '@/components/ProfileEdit/AvatarUploader.vue'
import { updateUserInfo } from '@/api/users/editProfile'
import { addNotification } from '@/data/notifications'
import { profileData } from '@/data/usersData'

const props = defineProps({
  firstName: String,
  lastName: String,
  username: String,
  email: String,
  bio: String,
  avatarPath: String,
  isPrivate: Boolean,
})
const emit = defineEmits(['avatar-change'])
const form = reactive({
  FirstName: props.firstName || '',
  LastName: props.lastName || '',
  Username: props.username || '',
  Email: props.email || '',
  About: props.bio || '',
  Password: '',
  IsPrivate: Boolean(props.isPrivate),
})
const savedSnapshot = ref(snapshot())
const saving = ref(false)
const feedback = ref('')
const saveError = ref('')
const isDirty = computed(() => snapshot() !== savedSnapshot.value)

function snapshot() {
  return JSON.stringify({
    FirstName: form.FirstName,
    LastName: form.LastName,
    Username: form.Username,
    Email: form.Email,
    About: form.About,
    Password: form.Password,
    IsPrivate: form.IsPrivate,
  })
}

async function updateInfo() {
<<<<<<< HEAD
    if (form.IsPrivate) {
        form.IsPrivate = 1
    } else {
        form.IsPrivate = 0
    }
    try {
        const result = await updateUserInfo(form);
        if (!result.status) {
            addNotification(result.message, 'error')
            return;
        }
        addNotification(result.message, 'success')
    } catch (err) {
        addNotification(err.message, 'error')
    }


    console.log('Profile updated successfully');
=======
  if (saving.value || !isDirty.value) return
  saving.value = true
  feedback.value = ''
  saveError.value = ''
  try {
    const result = await updateUserInfo({ ...form, IsPrivate: form.IsPrivate ? 1 : 0 })
    profileData.userInfo.firstName = form.FirstName
    profileData.userInfo.lastName = form.LastName
    profileData.userInfo.userName = form.Username
    profileData.userInfo.email = form.Email
    profileData.userInfo.about = form.About
    profileData.userInfo.isPrivate = form.IsPrivate ? 1 : 0
    profileData.about.bio = form.About
    form.Password = ''
    savedSnapshot.value = snapshot()
    feedback.value = result.message || 'Profile saved.'
    addNotification(feedback.value, 'success')
  } catch (err) {
    saveError.value = err.message || 'Could not save profile changes.'
  } finally {
    saving.value = false
  }
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
}
</script>

<template>
  <section class="edit-section orbit-surface" aria-labelledby="personal-info-heading">
    <header>
      <p class="orbit-meta">Account</p>
      <h2 id="personal-info-heading">Personal information</h2>
    </header>

    <div class="edit-group">
      <div class="edit-group__heading"><h3>Profile picture</h3><p>Shown across posts, comments, and conversations.</p></div>
      <AvatarUploader :src="avatarPath" @change="$emit('avatar-change', $event)" />
    </div>

    <form @submit.prevent="updateInfo">
      <div class="edit-group">
        <div class="edit-group__heading"><h3>Basic information</h3><p>Use the name and username people know you by.</p></div>
        <div class="field-grid">
          <FormField id="firstName" v-model="form.FirstName" label="First name" placeholder="First name" />
          <FormField id="lastName" v-model="form.LastName" label="Last name" placeholder="Last name" />
          <FormField id="username" v-model="form.Username" label="Username" placeholder="Username" />
          <FormField id="email" v-model="form.Email" label="Email" type="email" placeholder="Email address" />
        </div>
        <FormField id="bio" v-model="form.About" label="About me" type="textarea" placeholder="Tell people about yourself" />
        <FormField id="password" v-model="form.Password" label="New password (optional)" type="password" placeholder="Leave blank to keep your current password" />
      </div>

<<<<<<< HEAD
        <div class="edit-card">
            <AvatarUploader :src="props.avatar_path" @change="onAvatarChange" />

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

                <button type="button" class="privacy-button" :class="{ private: form.IsPrivate }"
                    @click="togglePrivacy">
                    {{ form.IsPrivate ? 'Make public' : 'Make private' }}
                </button>
            </div>

            <form @submit.prevent="updateInfo">
                <div class="field-grid">
                    <FormField id="firstName" label="First name" v-model="form.FirstName" placeholder="First name" />

                    <FormField id="lastName" label="Last name" v-model="form.LastName" placeholder="Last name" />

                    <FormField id="username" label="Username" v-model="form.Username" placeholder="Username" />

                    <FormField id="email" label="Email" type="email" v-model="form.Email" placeholder="Email address" />

                    <FormField id="password" label="Password" type="password" v-model="form.Password"
                        placeholder="****************" />
                </div>

                <FormField id="bio" label="Bio" type="textarea" v-model="form.About"
                    placeholder="Tell people about yourself" />

                <button type="submit" class="confirm-button">
                    Confirm changes
                </button>
            </form>
        </div>
    </section>
</template>

<style scoped>
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

.privacy-setting {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: clamp(12px, 2.5vw, 20px);
    padding: clamp(12px, 2.5vw, 18px);
    border: 2px solid var(--main-color);
    border-radius: 5px;
}

.privacy-title {
    margin: 0 0 5px;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(11px, 1.6vw, 12px);
    font-weight: 700;
}

.privacy-description {
    margin: 0;
    font-size: clamp(11px, 1.6vw, 12px);
    opacity: 0.7;
}

.privacy-button {
    flex-shrink: 0;
    padding: clamp(9px, 1.8vw, 11px) clamp(14px, 2.5vw, 18px);
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--input-focus);
    box-shadow: 3px 3px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
    cursor: pointer;
}

.privacy-button:hover {
    transform: translate(-1px, -1px);
}

.privacy-button.private {
    background: var(--main-color);
}

.field-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(220px, 100%), 1fr));
    gap: clamp(14px, 2.5vw, 20px);
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
      <fieldset class="visibility-group">
        <legend>Profile visibility</legend>
        <label :class="{ selected: !form.IsPrivate }">
          <input v-model="form.IsPrivate" type="radio" :value="false" />
          <span><strong>Public</strong><small>Anyone can view your profile and public activity.</small></span>
        </label>
        <label :class="{ selected: form.IsPrivate }">
          <input v-model="form.IsPrivate" type="radio" :value="true" />
          <span><strong>Private</strong><small>Only approved followers can view your full profile.</small></span>
        </label>
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
.edit-section { padding: var(--space-5); }
.edit-section > header { margin-bottom: var(--space-5); }
.edit-section .orbit-meta { margin: 0; color: var(--color-violet-soft); }
.edit-section h2 { margin: var(--space-1) 0 0; font-family: var(--font-display); font-size: 1.5rem; letter-spacing: 0; }
form, .edit-group { display: grid; gap: var(--space-4); }
.edit-group { padding: var(--space-5) 0; border-top: 1px solid var(--color-border); }
.edit-group__heading h3 { margin: 0; font-size: 1rem; }
.edit-group__heading p { margin: var(--space-1) 0 0; color: var(--color-text-muted); font-size: .875rem; }
.field-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: var(--space-4); }
.visibility-group { display: grid; gap: var(--space-3); margin: 0; padding: var(--space-5) 0; border: 0; border-top: 1px solid var(--color-border); }
.visibility-group legend { margin-bottom: var(--space-3); padding: 0; font-weight: 700; }
.visibility-group label { display: grid; grid-template-columns: auto minmax(0, 1fr); align-items: start; gap: var(--space-3); padding: var(--space-4); border: 1px solid var(--color-border); border-radius: var(--radius-small); cursor: pointer; }
.visibility-group label.selected { border-color: var(--color-mint); background: var(--color-surface-teal); }
.visibility-group input { margin-top: .25rem; accent-color: var(--color-mint); }
.visibility-group span { display: grid; gap: var(--space-1); }
.visibility-group small { color: var(--color-text-muted); line-height: 1.5; }
.form-feedback { margin: 0; padding: var(--space-3) var(--space-4); border-left: 3px solid var(--color-mint); background: var(--color-surface-teal); color: var(--color-text-soft); }
.form-feedback--error { border-color: var(--color-coral); background: var(--color-surface-coral); color: var(--color-coral-soft); }
.form-actions { display: flex; justify-content: flex-end; gap: var(--space-3); padding-top: var(--space-4); border-top: 1px solid var(--color-border); }
.form-actions a, .form-actions button { display: inline-flex; min-height: var(--touch-target); align-items: center; justify-content: center; padding: 0 var(--space-4); border-radius: var(--radius-small); font-weight: 700; text-decoration: none; }
.form-actions a { border: 1px solid var(--color-border); color: var(--color-text-soft); }
.form-actions button { border: 0; background: var(--gradient-action); color: white; cursor: pointer; }
@media (max-width: 600px) {
  .edit-section { padding: var(--space-4); }
  .field-grid { grid-template-columns: 1fr; }
  .form-actions { display: grid; }
}
</style>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
