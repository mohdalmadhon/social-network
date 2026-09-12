<script setup>
import { reactive, ref } from 'vue'
const signingUp = ref(false)

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
    validateAvatar
} from '@/helpers/validators/registration.js'

import InputHolder from './InputHolder.vue'
import { router } from '@/router/router.js'
import { addNotification } from '@/data/notifications.js'

const loginErrors = reactive({
    identifier: '',
    password: '',
    form: ''
})

const errors = reactive({
    firstName: '',
    lastName: '',
    username: '',
    dob: '',
    email: '',
    about: '',
    password: ''
})

const touched = reactive({
    firstName: false,
    lastName: false,
    username: false,
    dob: false,
    email: false,
    about: false,
    password: false
})

function handleAvatar(e) {
    const file = e.target.files[0]

    validateField('avatar', file)
}

function validateField(field, value) {
    touched[field] = true

    switch (field) {
        case 'avatar':
            errors[field] = validateAvatar(value)
            break
        case 'firstName':
        case 'lastName':
            errors[field] = validateName(value)
            break

        case 'username':
            errors[field] = validateUsername(value)
            break

        case 'email':
            errors[field] = validateEmail(value)
            break

        case 'dob':
            errors[field] = validateDOB(value)
            break

        case 'password':
            errors[field] = validatePassword(value)
            break

        case 'about':
            errors[field] = validateAbout(value)
            break
    }
}

function handleName(field, e) {
    const value = handleNameInput(e.target.value)

    e.target.value = value

    validateField(field, value)
}

function handleUsername(e) {
    const value = handleUsernameInput(e.target.value)

    e.target.value = value

    validateField('username', value)
}

function handleEmail(e) {
    const value = handleEmailInput(e.target.value)

    e.target.value = value

    validateField('email', value)
}

function inputClass(field) {
    if (!touched[field]) {
        return ''
    }

    return errors[field] ? 'input-error' : 'input-valid'
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
    return Object.values(errors).every(error => !error)
}

async function sendData(e) {
    e.preventDefault()

    const form = e.target

    if (!validateForm(form)) {
        return
    }

    const registerData = new FormData(form)

    try {
        const result = await registerUser(registerData)

        if (result.status) {
            addNotification('REGISTERED SUCCESSFULLY', 'success')
            signingUp.value = false
        } else {
            addNotification('Failed to register: ' + result.message, 'error')
        }
    } catch (error) {
        addNotification('Failed to register: ' + error.message, 'error')
        console.error(error)
    }
}

async function loggUser(e) {
    e.preventDefault()

    loginErrors.identifier = ''
    loginErrors.password = ''
    loginErrors.form = ''

    const form = e.target
    const formData = new FormData(form)

    const identifier = formData.get('Identifier')
    const password = formData.get('Pass')

    if (!identifier) {
        loginErrors.identifier = 'Identifier is required'
    }

    if (!password) {
        loginErrors.password = 'Password is required'
    }

    if (loginErrors.identifier || loginErrors.password) {
        return
    }

    try {
        const result = await loggingSession({
            Identifier: identifier,
            Pass: password
        })
        if(!result.status) return;
        router.replace("/home");
    } catch (error) {
        loginErrors.form = error.message || 'Invalid identifier or password'
    }
}
</script>

<template>
    <section class="auth-section">
        <div class="wrapper">
            <div class="card-switch">
                <div class="switch">
                    <nav class="auth-tabs" aria-label="Account access">
                      <button type="button" :aria-pressed="!signingUp" @click="signingUp = false">Log in</button>
                      <button type="button" :aria-pressed="signingUp" @click="signingUp = true">Sign up</button>
                    </nav>

                    <div class="flip-card__inner">
                        <div v-if="!signingUp" class="flip-card__front">
                            <div class="title">
                                Log in
                            </div>

                            <p class="form-subtitle">
                                Welcome back.
                            </p>

                            <form @submit.prevent="loggUser" class="flip-card__form">

                                <div class="input-group">
                                    <label for="login-email">
                                        Identifier
                                    </label>

                                    <InputHolder id="login-email" type="text" name="Identifier" :minLength="5"
                                        :maxLength="75" placeHolder="email@example.com" required />

                                    <span v-if="loginErrors.identifier" class="input-error-message">
                                        {{ loginErrors.identifier }}
                                    </span>
                                </div>

                                <div class="input-group">
                                    <label for="login-password">
                                        Password
                                    </label>

                                    <InputHolder id="login-password" type="password" name="Pass" :minLength="8"
                                        :maxLength="100" placeHolder="Password" required />

                                    <span v-if="loginErrors.password" class="input-error-message">
                                        {{ loginErrors.password }}
                                    </span>
                                </div>

                                <span v-if="loginErrors.form" class="input-error-message login-error">
                                    {{ loginErrors.form }}
                                </span>

                                <button class="flip-card__btn" type="submit">
                                    Let's go!
                                </button>

                            </form>

                            <p class="form-footer">
                                Don't have an account?
                                <button type="button" @click="signingUp = true">Sign up</button>
                            </p>
                        </div>

                        <div v-else class="flip-card__back">
                            <div class="title">
                                Sign up
                            </div>

                            <p class="form-subtitle">
                                Create your account.
                            </p>

                            <form class="flip-card__form" @submit.prevent="sendData">
                                <div class="input-row">
                                    <div class="input-group">
                                        <label for="first-name">
                                            First name
                                        </label>

                                        <InputHolder id="first-name" type="text" name="FirstName" :minLength="2"
                                            :maxLength="15" placeHolder="First name" :class="inputClass('firstName')"
                                            @input="handleName('firstName', $event)" required />

                                        <span v-if="touched.firstName && errors.firstName" class="input-error-message">
                                            {{ errors.firstName }}
                                        </span>
                                    </div>

                                    <div class="input-group">
                                        <label for="last-name">
                                            Last name
                                        </label>

                                        <InputHolder id="last-name" type="text" name="LastName" :minLength="2"
                                            :maxLength="15" placeHolder="Last name" :class="inputClass('lastName')"
                                            @input="handleName('lastName', $event)" required />

                                        <span v-if="touched.lastName && errors.lastName" class="input-error-message">
                                            {{ errors.lastName }}
                                        </span>
                                    </div>
                                </div>

                                <div class="input-row">
                                    <div class="input-group">
                                        <label for="username">
                                            Username
                                        </label>

                                        <InputHolder id="username" type="text" name="UserName" :minLength="3"
                                            :maxLength="12" placeHolder="Username" :class="inputClass('username')"
                                            @input="handleUsername" />

                                        <span v-if="touched.username && errors.username" class="input-error-message">
                                            {{ errors.username }}
                                        </span>
                                    </div>

                                    <div class="input-group">
                                        <label for="dob">
                                            Date of birth
                                        </label>

                                        <InputHolder id="dob" type="date" name="dob" :class="inputClass('dob')"
                                            @input="validateField('dob', $event.target.value)" required />

                                        <span v-if="touched.dob && errors.dob" class="input-error-message">
                                            {{ errors.dob }}
                                        </span>
                                    </div>
                                </div>

                                <div class="input-group">
                                    <label for="signup-email">
                                        Email
                                    </label>

                                    <InputHolder id="signup-email" type="email" name="Email" :minLength="5"
                                        :maxLength="75" placeHolder="email@example.com" :class="inputClass('email')"
                                        @input="handleEmail" required />

                                    <span v-if="touched.email && errors.email" class="input-error-message">
                                        {{ errors.email }}
                                    </span>
                                </div>

                                <div class="input-group">
                                    <label for="avatar">
                                        Avatar
                                    </label>

                                    <InputHolder id="avatar" type="file" name="Avatar" accept="image/*"
                                        :class="inputClass('avatar')" @change="handleAvatar" />

                                    <span v-if="touched.avatar && errors.avatar" class="input-error-message">
                                        {{ errors.avatar }}
                                    </span>
                                </div>

                                <div class="input-group">
                                    <label for="about">
                                        About
                                    </label>

                                    <textarea id="about" name="About" maxlength="1000"
                                        placeholder="Tell us a little about yourself..." :class="inputClass('about')"
                                        @input="validateField('about', $event.target.value)"></textarea>

                                    <span v-if="touched.about && errors.about" class="input-error-message">
                                        {{ errors.about }}
                                    </span>
                                </div>

                                <div class="input-group">
                                    <label for="signup-password">
                                        Password
                                    </label>

                                    <InputHolder id="signup-password" type="password" name="Password" :minLength="8"
                                        :maxLength="75" placeHolder="Password" :class="inputClass('password')"
                                        @input="validateField('password', $event.target.value)" required />

                                    <span v-if="touched.password && errors.password" class="input-error-message">
                                        {{ errors.password }}
                                    </span>
                                </div>

                                <button class="flip-card__btn" type="submit">
                                    Confirm!
                                </button>
                            </form>

                            <p class="form-footer">
                                Already have an account?
                                <button type="button" @click="signingUp = false">Log in</button>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<style scoped>
.auth-section { width: 100%; max-width: 36rem; margin-inline: auto; }
.wrapper, .card-switch, .switch, .flip-card__inner { width: 100%; }
.auth-tabs { display: flex; gap: .75rem; margin-bottom: 1.5rem; }
.auth-tabs button, .form-footer button { border: 0; background: transparent; color: var(--color-text-muted); cursor: pointer; min-height: 44px; }
.auth-tabs button { padding: .6rem 1.25rem; border-bottom: 2px solid transparent; }
.auth-tabs button[aria-pressed="true"] { color: var(--color-text); border-color: var(--color-violet); }
.flip-card__front, .flip-card__back { padding: clamp(1.25rem, 3vw, 2.5rem); border: 1px solid var(--color-border); border-radius: var(--radius-large); background: var(--color-surface); }
.title { font-family: var(--font-display); font-size: 1.8rem; color: var(--color-text); }
.form-subtitle { margin: .4rem 0 1.5rem; color: var(--color-text-muted); }
.flip-card__form, .input-group { display: grid; gap: .5rem; }
.flip-card__form { gap: 1rem; }
.input-row { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 1rem; }
.input-group { min-width: 0; }
.input-group label { font-size: .875rem; color: var(--color-text-soft); }
textarea, :deep(input) { width: 100%; min-width: 0; min-height: 44px; padding: .7rem; background: var(--color-input); color: var(--color-text); border: 1px solid var(--color-border); border-radius: var(--radius-small); font: inherit; }
textarea { min-height: 5rem; resize: vertical; }
.flip-card__btn { min-height: 44px; border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; font: inherit; cursor: pointer; }
.form-footer { margin: 1rem 0 0; color: var(--color-text-muted); font-size: .875rem; }
.form-footer button { color: var(--color-mint); }
.input-error-message { color: var(--color-coral); font-size: .8125rem; }
@media (max-width: 30rem) { .input-row { grid-template-columns: 1fr; } }
</style>
