<script setup>
import { reactive } from 'vue'

import { registerUser } from '@/api/auth/auth.js'

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

        console.log(result)
    } catch (error) {
        console.error(error)
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

                            <form class="flip-card__form">
                                <div class="input-group">
                                    <label for="login-email">
                                        Email
                                    </label>

                                    <InputHolder id="login-email" type="email" name="Email" :minLength="5"
                                        :maxLength="75" placeHolder="email@example.com" required />
                                </div>

                                <div class="input-group">
                                    <label for="login-password">
                                        Password
                                    </label>

                                    <InputHolder id="login-password" type="password" name="Password" :minLength="8"
                                        :maxLength="100" placeHolder="Password" required />
                                </div>

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
.auth-section {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px;
    background: var(--page-background);
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
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 13px;
    font-weight: 600;
    text-decoration: underline;
}

.card-side::after {
    position: absolute;
    content: "Sign up";
    left: 75px;
    top: 0;
    width: 100px;
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 13px;
    font-weight: 600;
    text-decoration: none;
}

.toggle:checked~.card-side::before {
    text-decoration: none;
}

.toggle:checked~.card-side::after {
    text-decoration: underline;
}

.slider {
    position: absolute;
    top: 0;
    left: 50%;
    width: 50px;
    height: 22px;
    transform: translateX(-50%);
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--bg-color);
    box-shadow: 4px 4px var(--main-color);
    cursor: pointer;
    transition: 0.3s;
}

.slider::before {
    position: absolute;
    content: "";
    width: 20px;
    height: 20px;
    left: -2px;
    bottom: 2px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--bg-color);
    box-shadow: 0 3px 0 var(--main-color);
    transition: 0.3s;
}

.toggle:checked+.slider {
    background: var(--input-focus);
}

.toggle:checked+.slider::before {
    transform: translateX(30px);
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
    background: var(--bg-color);
    border: 2px solid var(--main-color);
    border-radius: 8px;
    box-shadow: 7px 7px var(--main-color);
    backface-visibility: hidden;
    -webkit-backface-visibility: hidden;
}

.flip-card__back {
    transform: rotateY(180deg);
    overflow-y: auto;
}

.title {
    margin: 0 0 2px;
    color: var(--main-color);
    font-family: "Liter", serif;
    font-size: 32px;
    font-weight: 900;
}

.form-subtitle {
    margin: 0 0 20px;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
}

.flip-card__form {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.input-row {
    width: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
}

.input-group {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 7px;
}

.input-group label {
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
    font-weight: 600;
}

.input-group :deep(.form-input) {
    width: 100%;
}

.input-group textarea {
    width: 100%;
    min-height: 75px;
    resize: vertical;
    padding: 10px 12px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    outline: none;
    background: var(--bg-color);
    box-shadow: 4px 4px var(--main-color);
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 12px;
    font-weight: 500;
}

.input-group textarea::placeholder {
    color: var(--font-color-sub);
    opacity: 0.7;
}

.input-group textarea:focus {
    border-color: var(--input-focus);
    box-shadow: 4px 4px var(--input-focus);
}

.input-group :deep(.input-error),
.input-group textarea.input-error {
    border-color: #e74c3c !important;
    box-shadow: 4px 4px #e74c3c !important;
}

.input-group :deep(.input-valid),
.input-group textarea.input-valid {
    border-color: #2ecc71 !important;
    box-shadow: 4px 4px #2ecc71 !important;
}

.input-error-message {
    color: #e74c3c;
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
    line-height: 1.3;
}

.flip-card__btn {
    align-self: center;
    width: 140px;
    height: 45px;
    margin-top: 8px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--input-focus);
    box-shadow: 4px 4px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: 0.15s;
}

.flip-card__btn:hover {
    transform: translate(-1px, -1px);
    box-shadow: 6px 6px var(--main-color);
}

.flip-card__btn:active {
    transform: translate(4px, 4px);
    box-shadow: 0 0 var(--main-color);
}

.form-footer {
    margin: 18px 0 0;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 10px;
}

.form-footer span {
    color: var(--input-focus);
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
    }

    .input-row {
        grid-template-columns: 1fr;
        gap: 15px;
    }
}
</style>