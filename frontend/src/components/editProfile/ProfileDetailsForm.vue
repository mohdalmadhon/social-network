<script setup>
import { getUserData } from '@/api/users'
import { userData } from '@/stores/userData'
import { onMounted, ref } from 'vue'
import FormField from './FormField.vue'
import FormTextarea from './FormTextarea.vue'

const formData = ref({
  firstName: '',
  lastName: '',
  username: '',
  email: '',
  about: '',
  password: ''
})

const errors = ref({
  firstName: '',
  lastName: '',
  username: '',
  email: '',
  password: ''
})

onMounted(async () => {
  await getUserData()

  formData.value.firstName = userData.value.firstName || ''
  formData.value.lastName = userData.value.lastName || ''
  formData.value.username = userData.value.username || ''
  formData.value.email = userData.value.email || ''
  formData.value.about = userData.value.about || ''
})

function validateName(field) {
  let value = formData.value[field]

  value = value.replace(/[^a-zA-Z]/g, '')

  if (value.length > 0) {
    value =
      value.charAt(0).toUpperCase() +
      value.slice(1).toLowerCase()
  }

  formData.value[field] = value

  if (value.length === 0) {
    errors.value[field] = 'This field is required.'
  } else if (value.length < 3) {
    errors.value[field] = 'Must be at least 3 characters.'
  } else if (value.length > 15) {
    errors.value[field] = 'Must be no more than 15 characters.'
  } else {
    errors.value[field] = ''
  }
}


function validateUsername() {
  let value = formData.value.username.toLowerCase()

  value = value.replace(/[^a-z0-9@#$%_-]/g, '')

  formData.value.username = value

  errors.value.username = ''
}

function validateEmail() {
  const value = formData.value.email

  if (value.length === 0) {
    errors.value.email = 'Email is required.'
    return
  }

  if (value.length > 50) {
    errors.value.email = 'Email must be no more than 50 characters.'
    return
  }

  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  if (!emailPattern.test(value)) {
    errors.value.email = 'Please enter a valid email address.'
    return
  }

  errors.value.email = ''
}

function validatePassword() {
  const value = formData.value.password

  if (value.length === 0) {
    errors.value.password = ''
    return
  }

  if (value.length < 8) {
    errors.value.password = 'Password must be at least 8 characters.'
  } else if (value.length > 75) {
    errors.value.password = 'Password must be no more than 75 characters.'
  } else {
    errors.value.password = ''
  }
}


const notification = ref({
  show: false,
  type: '',
  message: ''
})

let notificationTimer = null

function showNotification(type, message) {
  notification.value = {
    show: true,
    type,
    message
  }

  clearTimeout(notificationTimer)

  notificationTimer = setTimeout(() => {
    notification.value.show = false
  }, 4000)
}

function hideNotification() {
  notification.value.show = false
  clearTimeout(notificationTimer)
}

async function submitForm() {
  Object.keys(errors.value).forEach(field => {
    if (field === 'firstName' || field === 'lastName') {
      validateName(field)
    } else if (field === 'username') {
      validateUsername()
    } else if (field === 'email') {
      validateEmail()
    } else if (field === 'password') {
      validatePassword()
    }
  })

  const hasErrors = Object.values(errors.value).some(
    error => error !== ''
  )

  if (hasErrors) {
    showNotification(
      'error',
      'Please fix the errors before saving your changes.'
    )
    return
  }

  const currentData = {
    firstName: userData.value.firstName || '',
    lastName: userData.value.lastName || '',
    username: userData.value.username || '',
    email: userData.value.email || '',
    about: userData.value.about || ''
  }

  const formValues = {
    firstName: formData.value.firstName || '',
    lastName: formData.value.lastName || '',
    username: formData.value.username || '',
    email: formData.value.email || '',
    about: formData.value.about || ''
  }

  const isSame = Object.keys(currentData).every(
    key => currentData[key] === formValues[key]
  )

  if (isSame && !formData.value.password) {
    showNotification(
      'error',
      'No changes were made.'
    )
    return
  }

  try {
    const resp = await fetch('/api/user', {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        firstName: formData.value.firstName,
        lastName: formData.value.lastName,
        username: formData.value.username,
        email: formData.value.email,
        about: formData.value.about,
        password: formData.value.password
      })
    })

    if (!resp.ok) {
      throw new Error(`Request failed with status ${resp.status}`)
    }

    const result = await resp.json()

    if (!result.status) {
      showNotification(
        'error',
        result.message || 'Failed to update your profile.'
      )
      return
    }

    showNotification(
      'success',
      result.message || 'Your profile has been updated successfully.'
    )

    userData.value = {
      ...userData.value,
      firstName: formData.value.firstName,
      lastName: formData.value.lastName,
      username: formData.value.username,
      email: formData.value.email,
      about: formData.value.about
    }

    formData.value.password = ''
  } catch (err) {
    console.error(err)

    showNotification(
      'error',
      'Something went wrong while updating your profile.'
    )
  }
}
</script>

<template>
  <transition name="notification">
    <div v-if="notification.show" class="notification" :class="`notification--${notification.type}`" role="alert">
      <span class="notification__icon">
        {{ notification.type === 'success' ? '✓' : '!' }}
      </span>

      <span class="notification__message">
        {{ notification.message }}
      </span>

      <button type="button" class="notification__close" @click="hideNotification" aria-label="Close notification">
        ×
      </button>
    </div>
  </transition>

  <form class="details" @submit.prevent="submitForm">
    <div class="details__top">
      <div>
        <p class="details__eyebrow">USER INFORMATION</p>
        <h2 class="details__title">Personal details</h2>
      </div>
    </div>

    <section class="details__section">
      <h3 class="details__section-title">Name</h3>

      <div class="details__grid details__grid--two">
        <FormField id="first-name" v-model="formData.firstName" label="First name" placeholder="First name"
          autocomplete="given-name" maxlength="15" :error="errors.firstName" @input="validateName('firstName')" />

        <FormField id="last-name" v-model="formData.lastName" label="Last name" placeholder="Last name"
          autocomplete="family-name" maxlength="15" :error="errors.lastName" @input="validateName('lastName')" />
      </div>
    </section>

    <div class="details__divider"></div>

    <section class="details__section">
      <h3 class="details__section-title">Account details</h3>

      <div class="details__grid details__grid--two">
        <FormField id="username" v-model="formData.username" label="Username" placeholder="username" prefix="@"
          autocomplete="username" :error="errors.username" @input="validateUsername" />

        <FormField id="email" v-model="formData.email" label="Email" type="email" placeholder="you@example.com"
          autocomplete="email" maxlength="50" :error="errors.email" @input="validateEmail" />
      </div>
    </section>

    <div class="details__divider"></div>

    <section class="details__section">
      <h3 class="details__section-title">About you</h3>

      <FormTextarea id="about" v-model="formData.about" label="About me" placeholder="Introduce yourself..."
        :maxlength="160" hint="Keep it short and personal." />
    </section>

    <div class="details__divider"></div>

    <section class="details__section">
      <div>
        <h3 class="details__section-title">Password</h3>
        <p class="details__section-caption">
          Leave this empty if you don't want to change your password.
        </p>
      </div>

      <FormField id="password" v-model="formData.password" label="New password" type="password"
        placeholder="Enter a new password" autocomplete="new-password" maxlength="75" :error="errors.password"
        @input="validatePassword" />
    </section>

    <div class="details__actions">
      <button class="btn btn--secondary" type="button">
        Cancel
      </button>

      <button class="btn btn--primary" type="submit">
        <span>Save changes</span>
        <span class="btn__arrow">→</span>
      </button>
    </div>
  </form>
</template>

<style scoped>
@import '../../styles/variables.css';

.notification {
  position: fixed;
  top: var(--space-5);
  right: var(--space-5);
  z-index: 1000;

  display: flex;
  align-items: center;
  gap: var(--space-3);

  min-width: 18rem;
  max-width: 28rem;

  padding: var(--space-3) var(--space-4);

  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);

  color: var(--color-text);
  box-shadow: var(--shadow-raised);
}

.notification--success {
  border-color: rgb(34 197 94 / 40%);
}

.notification--error {
  border-color: rgb(239 68 68 / 40%);
}

.notification__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;

  flex: 0 0 auto;

  width: 1.5rem;
  height: 1.5rem;

  border-radius: 50%;

  font-size: 0.75rem;
  font-weight: 700;
}

.notification--success .notification__icon {
  background: rgb(34 197 94 / 12%);
  color: #22c55e;
}

.notification--error .notification__icon {
  background: rgb(239 68 68 / 12%);
  color: #ef4444;
}

.notification__message {
  flex: 1;
  font-size: 0.82rem;
  line-height: 1.4;
}

.notification__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;

  width: 1.5rem;
  height: 1.5rem;

  padding: 0;

  border: 0;
  background: transparent;

  color: var(--color-text-muted);

  font-size: 1.2rem;
  cursor: pointer;
}

.notification__close:hover {
  color: var(--color-text);
}

.notification-enter-active,
.notification-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.notification-enter-from,
.notification-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 39.99rem) {
  .notification {
    top: var(--space-3);
    right: var(--space-3);
    left: var(--space-3);

    min-width: 0;
    max-width: none;
  }
}

.details {
  position: relative;
  overflow: hidden;
  padding: var(--space-6);
  background:
    linear-gradient(145deg,
      rgb(124 92 255 / 5%),
      transparent 32%),
    var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-large);
  box-shadow: var(--shadow-raised);
}

.details::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: var(--gradient-action);
}

.details__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: var(--space-6);
}

.details__eyebrow {
  margin-bottom: var(--space-1);
  color: var(--color-violet);
  font-family: var(--font-meta);
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.details__title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: 1.4rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.details__section {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.details__section-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 700;
}

.details__section-caption {
  margin-top: 0.2rem;
  color: var(--color-text-faint);
  font-size: 0.8rem;
  line-height: 1.5;
}

.details__divider {
  height: 1px;
  margin: var(--space-6) 0;
  background:
    linear-gradient(90deg,
      transparent,
      var(--color-border),
      transparent);
}

.details__grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-4);
}

.details__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-6);
  padding-top: var(--space-5);
  border-top: 1px solid var(--color-border);
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  min-height: var(--touch-target);
  padding: 0 var(--space-5);
  border: 1px solid transparent;
  border-radius: var(--radius-small);
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  transition:
    transform 0.15s ease,
    filter 0.15s ease,
    border-color 0.15s ease,
    background 0.15s ease;
}

.btn:hover {
  transform: translateY(-1px);
}

.btn:active {
  transform: translateY(0);
}

.btn:focus-visible {
  outline: none;
  box-shadow: var(--focus-ring);
}

.btn--secondary {
  background: var(--color-input);
  border-color: var(--color-border);
  color: var(--color-text-soft);
}

.btn--secondary:hover {
  background: var(--color-surface-raised);
  border-color: var(--color-text-faint);
  color: var(--color-text);
}

.btn--primary {
  background: var(--gradient-action);
  color: #fff;
  box-shadow: var(--shadow-raised);
}

.btn--primary:hover {
  filter: brightness(1.08);
}

.btn__arrow {
  font-size: 1.1rem;
  transition: transform 0.15s ease;
}

.btn--primary:hover .btn__arrow {
  transform: translateX(3px);
}

@media (min-width: 40rem) {
  .details__grid--two {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 39.99rem) {
  .details {
    padding: var(--space-4);
  }

  .details__actions {
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
  }
}

@media (min-width: 64rem) {
  .details {
    padding: var(--space-7);
  }
}
</style>