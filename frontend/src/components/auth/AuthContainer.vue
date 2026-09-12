<script setup>
import { reactive, ref } from 'vue'
import { loggingSession, registerUser } from '@/api/auth/auth.js'
import {
  validateName,
  validateUsername,
  validateEmail,
  validateDOB,
  validatePassword,
  validateAbout,
  handleEmailInput,
  handleUsernameInput,
  handleNameInput,
  validateAvatar,
} from '@/helpers/validators/registration.js'
import InputHolder from './InputHolder.vue'
import { router } from '@/router/router.js'
import { addNotification } from '@/data/notifications.js'

const signingUp = ref(false)
const avatarInput = ref(null)
const avatarName = ref('')

const loginErrors = reactive({
  identifier: '',
  password: '',
  form: '',
})

const errors = reactive({
  firstName: '',
  lastName: '',
  username: '',
  dob: '',
  email: '',
  about: '',
  password: '',
  avatar: '',
})

const touched = reactive({
  firstName: false,
  lastName: false,
  username: false,
  dob: false,
  email: false,
  about: false,
  password: false,
  avatar: false,
})

function validateField(field, value) {
  touched[field] = true

  const validators = {
    avatar: validateAvatar,
    firstName: validateName,
    lastName: validateName,
    username: validateUsername,
    email: validateEmail,
    dob: validateDOB,
    password: validatePassword,
    about: validateAbout,
  }

  errors[field] = validators[field](value)
}

function handleName(field, event) {
  const value = handleNameInput(event.target.value)
  event.target.value = value
  validateField(field, value)
}

function handleUsername(event) {
  const value = handleUsernameInput(event.target.value)
  event.target.value = value
  validateField('username', value)
}

function handleEmail(event) {
  const value = handleEmailInput(event.target.value)
  event.target.value = value
  validateField('email', value)
}

function inputClass(field) {
  if (!touched[field]) return ''
  return errors[field] ? 'input-error' : 'input-valid'
}

function updateAvatar(file) {
  avatarName.value = file?.name || ''
  validateField('avatar', file)
}

function handleAvatar(event) {
  updateAvatar(event.target.files?.[0] || null)
}

function handleAvatarDrop(event) {
  const files = event.dataTransfer?.files
  if (!files?.length) return

  try {
    avatarInput.value.files = files
  } catch {
    // Browsers can block assigning a dropped FileList. The click picker still works.
  }

  updateAvatar(files[0])
}

function validateForm(form) {
  const data = new FormData(form)

  validateField('firstName', data.get('FirstName') || '')
  validateField('lastName', data.get('LastName') || '')
  validateField('username', data.get('UserName') || '')
  validateField('dob', data.get('dob') || '')
  validateField('email', data.get('Email') || '')
  validateField('about', data.get('About') || '')
  validateField('password', data.get('Password') || '')
  validateField('avatar', data.get('Avatar'))

  return Object.values(errors).every((error) => !error)
}

async function sendData(event) {
  event.preventDefault()
  const form = event.target

  if (!validateForm(form)) return

  try {
    const result = await registerUser(new FormData(form))

    if (!result.status) {
      addNotification(`Failed to register: ${result.message}`, 'error')
      return
    }

    addNotification('Account created successfully.', 'success')
    signingUp.value = false
  } catch (error) {
    addNotification(`Failed to register: ${error.message}`, 'error')
  }
}

async function loggUser(event) {
  event.preventDefault()
  loginErrors.identifier = ''
  loginErrors.password = ''
  loginErrors.form = ''

  const formData = new FormData(event.target)
  const identifier = formData.get('Identifier')
  const password = formData.get('Pass')

  if (!identifier) loginErrors.identifier = 'Email or username is required'
  if (!password) loginErrors.password = 'Password is required'
  if (loginErrors.identifier || loginErrors.password) return

  try {
    const result = await loggingSession({ Identifier: identifier, Pass: password })
    if (result.status) router.replace('/home')
  } catch (error) {
    loginErrors.form = error.message || 'Invalid email, username, or password'
  }
}
</script>

<template>
  <section class="auth-section" aria-label="Account access">
    <div class="auth-card">
      <nav class="auth-tabs" aria-label="Account access">
        <button class="auth-tab" :class="{ 'auth-tab--active': !signingUp }" type="button" @click="signingUp = false">
          Sign in
        </button>
        <button class="auth-tab" :class="{ 'auth-tab--active': signingUp }" type="button" @click="signingUp = true">
          Create account
        </button>
      </nav>

      <div v-if="!signingUp" class="auth-panel">
        <form class="auth-form" @submit.prevent="loggUser">
          <div class="input-group">
            <label for="login-email">EMAIL OR USERNAME *</label>
            <InputHolder id="login-email" type="text" name="Identifier" :min-length="3" :max-length="75" place-holder="noa@orbit.app" required autocomplete="username" />
            <span v-if="loginErrors.identifier" class="input-error-message">{{ loginErrors.identifier }}</span>
          </div>

          <div class="input-group">
            <label for="login-password">PASSWORD *</label>
            <InputHolder id="login-password" type="password" name="Pass" :min-length="8" :max-length="100" place-holder="••••••••••" required autocomplete="current-password" />
            <span v-if="loginErrors.password" class="input-error-message">{{ loginErrors.password }}</span>
          </div>

          <span v-if="loginErrors.form" class="input-error-message" role="alert">{{ loginErrors.form }}</span>

          <label class="remember-row">
            <input type="checkbox" checked />
            <span>Keep me signed in on this device <small>session cookie</small></span>
          </label>

          <button class="auth-submit" type="submit">Let's go →</button>
        </form>

        <p class="auth-footer">Don't have an account? <button type="button" @click="signingUp = true">Create account</button></p>
      </div>

      <div v-else class="auth-panel">
        <form class="auth-form" @submit.prevent="sendData">
          <div class="input-group">
            <label for="signup-email">EMAIL *</label>
            <InputHolder id="signup-email" type="email" name="Email" :min-length="5" :max-length="75" place-holder="noa@orbit.app" :class="inputClass('email')" autocomplete="email" required @input="handleEmail" />
            <span v-if="touched.email && errors.email" class="input-error-message">{{ errors.email }}</span>
          </div>

          <div class="input-group">
            <label for="signup-password">PASSWORD *</label>
            <InputHolder id="signup-password" type="password" name="Password" :min-length="8" :max-length="75" place-holder="••••••••••" :class="inputClass('password')" autocomplete="new-password" required @input="validateField('password', $event.target.value)" />
            <span v-if="touched.password && errors.password" class="input-error-message">{{ errors.password }}</span>
          </div>

          <div class="input-row">
            <div class="input-group">
              <label for="first-name">FIRST NAME *</label>
              <InputHolder id="first-name" type="text" name="FirstName" :min-length="2" :max-length="15" place-holder="Noa" :class="inputClass('firstName')" autocomplete="given-name" required @input="handleName('firstName', $event)" />
              <span v-if="touched.firstName && errors.firstName" class="input-error-message">{{ errors.firstName }}</span>
            </div>

            <div class="input-group">
              <label for="last-name">LAST NAME *</label>
              <InputHolder id="last-name" type="text" name="LastName" :min-length="2" :max-length="15" place-holder="Ferreira" :class="inputClass('lastName')" autocomplete="family-name" required @input="handleName('lastName', $event)" />
              <span v-if="touched.lastName && errors.lastName" class="input-error-message">{{ errors.lastName }}</span>
            </div>
          </div>

          <div class="input-row">
            <div class="input-group">
              <label for="dob">DATE OF BIRTH *</label>
              <InputHolder id="dob" type="date" name="dob" :class="inputClass('dob')" required @input="validateField('dob', $event.target.value)" />
              <span v-if="touched.dob && errors.dob" class="input-error-message">{{ errors.dob }}</span>
            </div>

            <div class="input-group">
              <label for="username">NICKNAME <span>— OPTIONAL</span></label>
              <InputHolder id="username" type="text" name="UserName" :min-length="3" :max-length="12" place-holder="@noa.png" :class="inputClass('username')" autocomplete="nickname" @input="handleUsername" />
              <span v-if="touched.username && errors.username" class="input-error-message">{{ errors.username }}</span>
            </div>
          </div>

          <div class="input-group">
            <label for="avatar">AVATAR <span>— OPTIONAL · JPG, PNG, GIF</span></label>
            <label class="upload-zone" for="avatar" @dragover.prevent @drop.prevent="handleAvatarDrop">
              <span class="upload-zone__icon" aria-hidden="true">☁</span>
              <span>{{ avatarName || 'Drop an image or click to browse' }}</span>
              <input id="avatar" ref="avatarInput" type="file" name="Avatar" accept="image/jpeg,image/png,image/gif" @change="handleAvatar" />
            </label>
            <span v-if="touched.avatar && errors.avatar" class="input-error-message">{{ errors.avatar }}</span>
          </div>

          <div class="input-group">
            <label for="about">ABOUT ME <span>— OPTIONAL</span></label>
            <textarea id="about" name="About" maxlength="1000" placeholder="Trail runner, pixel-art hobbyist, espresso before noon…" :class="inputClass('about')" @input="validateField('about', $event.target.value)"></textarea>
            <span v-if="touched.about && errors.about" class="input-error-message">{{ errors.about }}</span>
          </div>

          <label class="remember-row">
            <input type="checkbox" checked />
            <span>Keep me signed in on this device <small>session cookie</small></span>
          </label>

          <button class="auth-submit" type="submit">Create my orbit →</button>
        </form>

        <p class="auth-footer">Already orbiting? <button type="button" @click="signingUp = false">Sign in</button></p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.auth-section {
  width: 100%;
  max-width: 44rem;
  margin-inline: auto;
}

.auth-card {
  padding: 2.25rem 3rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  background: var(--color-surface);
  box-shadow: var(--shadow-raised);
}

.auth-tabs {
  display: flex;
  gap: 1.75rem;
  margin-bottom: 1.25rem;
  border-bottom: 1px solid var(--color-border);
}

.auth-tab {
  min-height: 2.75rem;
  padding: 0 0 .75rem;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: .875rem;
  font-weight: 600;
}

.auth-tab:hover,
.auth-tab--active {
  color: var(--color-text);
}

.auth-tab--active {
  border-color: var(--color-coral);
}

.auth-form,
.input-group {
  display: grid;
  gap: .5rem;
}

.auth-form { gap: 1rem; }
.input-group { gap: .25rem; }
.input-row { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: .75rem; }
.input-group { min-width: 0; }

.input-group > label {
  color: var(--color-text-muted);
  font-family: var(--font-meta);
  font-size: .6875rem;
  letter-spacing: .08em;
}

.input-group > label span { color: var(--color-text-faint); }

.input-group textarea {
  width: 100%;
  min-height: 4.5rem;
  padding: .75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  outline: 0;
  background: var(--color-input);
  color: var(--color-text);
  resize: vertical;
}

.input-group textarea:focus {
  border-color: var(--color-violet);
  box-shadow: var(--focus-ring);
}

.upload-zone {
  display: flex;
  align-items: center;
  gap: .75rem;
  min-height: 3.5rem;
  padding: .75rem;
  border: 1px dashed var(--color-violet);
  border-radius: var(--radius-small);
  background: rgb(124 92 255 / 5%);
  color: var(--color-text-soft);
  cursor: pointer;
}

.upload-zone:hover {
  background: rgb(124 92 255 / 11%);
}

.upload-zone__icon {
  color: var(--color-violet);
  font-size: 1.25rem;
  line-height: 1;
}

.upload-zone input {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
}

.remember-row {
  display: flex;
  align-items: flex-start;
  gap: .5rem;
  color: var(--color-text-muted);
  font-size: .75rem;
}

.remember-row input {
  width: 1rem;
  height: 1rem;
  margin: 0;
  accent-color: var(--color-violet);
}

.remember-row small {
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: .6875rem;
}

.auth-submit {
  width: 100%;
  min-height: 3.25rem;
  border: 0;
  border-radius: 999px;
  background: var(--gradient-action);
  color: var(--color-text);
  cursor: pointer;
  font-weight: 700;
}

.auth-submit:hover { filter: brightness(1.08); }
.auth-footer { margin: .5rem 0 0; color: var(--color-text-muted); font-size: .8125rem; text-align: center; }
.auth-footer button { border: 0; background: transparent; color: var(--color-violet-soft); cursor: pointer; font: inherit; }
.auth-footer button:hover { color: var(--color-text); }
.input-error-message { color: var(--color-coral); font-size: .75rem; }

@media (max-width: 30rem) {
  .input-row { grid-template-columns: 1fr; }
}
</style>
