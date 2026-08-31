<script setup>
import { reactive } from 'vue'

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
                <label class="switch">
                    <input type="checkbox" class="toggle">

                    <span class="slider"></span>

                    <span class="card-side"></span>

                    <div class="flip-card__inner">
                        <div class="flip-card__front">
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
                                <span>Sign up</span>
                            </p>
                        </div>

                        <div class="flip-card__back">
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
                                <span>Log in</span>
                            </p>
                        </div>
                    </div>
                </label>
            </div>
        </div>
    </section>
</template>

<style scoped>
.login-error {
    text-align: center;
    margin-top: -5px;
}

.auth-section {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px;
    background:
        radial-gradient(circle at 80% 30%, rgba(75, 63, 160, 0.14), transparent 40%),
        #100d2b;
}

.wrapper {
    width: 100%;
    max-width: 520px;
}

.card-switch {
    width: 100%;
}

.switch {
    position: relative;
    width: 100%;
    min-height: 700px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.toggle {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
}

.card-side {
    position: absolute;
    top: 0;
    width: 60px;
    height: 22px;
}

.card-side::before {
    position: absolute;
    content: "Log in";
    left: -90px;
    top: 0;
    width: 100px;
    color: #ffffff;
    font-family: Arial, sans-serif;
    font-size: 13px;
    font-weight: 600;
}

.card-side::after {
    position: absolute;
    content: "Sign up";
    left: 75px;
    top: 0;
    width: 100px;
    color: #69749a;
    font-family: Arial, sans-serif;
    font-size: 13px;
    font-weight: 600;
}

.toggle:checked~.card-side::before {
    color: #69749a;
}

.toggle:checked~.card-side::after {
    color: #ffffff;
}

.slider {
    position: absolute;
    top: 0;
    left: 50%;
    width: 50px;
    height: 24px;
    transform: translateX(-50%);
    border: 1px solid #303754;
    border-radius: 20px;
    background: #222741;
    cursor: pointer;
    transition: 0.3s;
}

.slider::before {
    position: absolute;
    content: "";
    width: 18px;
    height: 18px;
    left: 2px;
    top: 2px;
    border-radius: 50%;
    background: linear-gradient(135deg, #805cff, #ff6288);
    transition: 0.3s;
}

.toggle:checked+.slider {
    background: #222741;
    border-color: #835cff;
}

.toggle:checked+.slider::before {
    transform: translateX(26px);
}

.flip-card__inner {
    position: relative;
    width: 480px;
    height: 650px;
    margin-top: 45px;
    perspective: 1000px;
    transition: transform 0.8s;
    transform-style: preserve-3d;
}

.toggle:checked~.flip-card__inner {
    transform: rotateY(180deg);
}

.flip-card__front,
.flip-card__back {
    position: absolute;
    width: 100%;
    height: 100%;
    padding: 35px 40px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background: #171b2d;
    border: 1px solid #303754;
    border-radius: 23px;
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.25);
    color: #f5f5ff;
    font-family: Arial, sans-serif;
    backface-visibility: hidden;
    -webkit-backface-visibility: hidden;
}

.flip-card__back {
    transform: rotateY(180deg);
    overflow-y: auto;
}

.title {
    margin: 0 0 2px;
    color: #ffffff;
    font-family: Arial, sans-serif;
    font-size: 26px;
    font-weight: 700;
}

.form-subtitle {
    margin: 0 0 20px;
    color: #69749a;
    font-family: Arial, sans-serif;
    font-size: 12px;
}

.flip-card__form {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.input-row {
    width: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 9px;
}

.input-group {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.input-group label {
    color: #69749a;
    font-family: Arial, sans-serif;
    font-size: 10px;
    font-weight: 500;
    letter-spacing: 0.8px;
    text-transform: uppercase;
}

.input-group :deep(.form-input) {
    width: 100%;
}

.input-group textarea {
    width: 100%;
    min-height: 61px;
    resize: none;
    padding: 12px 13px;
    border: 1px solid transparent;
    border-radius: 7px;
    outline: none;
    background: #222741;
    color: #f4f3ff;
    font-family: inherit;
    font-size: 12px;
    transition: 0.2s ease;
}

.input-group textarea::placeholder {
    color: #59627f;
}

.input-group textarea:focus {
    border-color: #835cff;
    box-shadow: 0 0 0 1px rgba(131, 92, 255, 0.15);
}

.input-group :deep(.input-error),
.input-group textarea.input-error {
    border-color: #ff6288 !important;
    box-shadow: 0 0 0 1px rgba(255, 98, 136, 0.15) !important;
}

.input-group :deep(.input-valid),
.input-group textarea.input-valid {
    border-color: #35d68a !important;
    box-shadow: 0 0 0 1px rgba(53, 214, 138, 0.15) !important;
}

.input-error-message {
    color: #ff6288;
    font-family: Arial, sans-serif;
    font-size: 11px;
    font-weight: 400;
    line-height: 1.4;
}

.flip-card__btn {
    align-self: center;
    width: 100%;
    height: 44px;
    margin-top: 8px;
    border: 0;
    border-radius: 25px;
    background: linear-gradient(90deg, #805cff, #ff6288);
    color: white;
    font-family: inherit;
    font-size: 13px;
    font-weight: 700;
    cursor: pointer;
    transition:
        transform 0.15s ease,
        filter 0.15s ease;
}

.flip-card__btn:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

.flip-card__btn:active {
    transform: translateY(0);
}

.form-footer {
    margin: 18px 0 0;
    color: #69749a;
    font-family: Arial, sans-serif;
    font-size: 11px;
}

.form-footer span {
    color: #9b7cff;
    font-weight: 600;
}

@media (max-width: 900px) {
    .auth-section {
        min-height: 600px;
        padding: 60px 30px;
    }

    .flip-card__inner {
        width: 440px;
    }
}

@media (max-width: 550px) {
    .auth-section {
        padding: 50px 20px;
    }

    .flip-card__inner {
        width: min(400px, 90vw);
    }

    .flip-card__front,
    .flip-card__back {
        padding: 30px 25px;
        border-radius: 18px;
    }

    .input-row {
        grid-template-columns: 1fr;
        gap: 12px;
    }
}
</style>
