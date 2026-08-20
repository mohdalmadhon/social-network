<script setup>
import { ref } from 'vue'

const activeTab = ref('create')

const showPassword = ref(false)

const showSignInPassword = ref(false)

const form = ref({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    username: '',
    avatar: null,
    about: ''
})

const logger = ref({
    identifier: '',
    password: ''
})

const errors = ref({
    signinIdentifier: '',
    signinPassword: '',
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    username: '',
    avatar: '',
    about: ''
})

const notification = ref({
    show: false,
    message: '',
    type: 'success'
})

let notificationTimer = null

function showNotification(message, type = 'success', duration = 3500) {
    if (notificationTimer) {
        clearTimeout(notificationTimer)
        notificationTimer = null
    }
    notification.value.show = false
    requestAnimationFrame(() => {
        notification.value.message = message
        notification.value.type = type
        notification.value.show = true
        notificationTimer = setTimeout(() => {
            notification.value.show = false
        }, duration)
    })
}

function switchTab(tab) {
    activeTab.value = tab
}

function errorClass(error) {
    return error ? 'error-show' : 'error-hide'
}

function clearError(field) {
    errors.value[field] = ''
}

async function validateEmail() {
    const value = form.value.email.trim()

    if (!value) {
        errors.value.email = 'Email is required'
        return false
    }

    if (value.length > 70) {
        errors.value.email = 'Email is too long. Max length is 70 characters'
        return false
    }

    if (value.length < 5) {
        errors.value.email = 'Email is too short'
        return false
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(value)) {
        errors.value.email = 'Email is not valid'
        return false
    }

    try {
        const resp = await fetch(
            `/api/register/checkEmail?email=${encodeURIComponent(value)}`
        )

        const text = await resp.text()

        if (!resp.ok) {
            errors.value.email = 'Unable to check email'
            return false
        }

        const result = JSON.parse(text)

        if (result.available === false) {
            errors.value.email = 'This email is already registered'
            return false
        }

        errors.value.email = ''
        return true
    } catch (err) {
        errors.value.email = 'Unable to check email'
        return false
    }
}

function validatePassword() {
    const value = form.value.password

    if (!value) {
        errors.value.password = 'Password is required'
        return false
    }

    if (value.length < 8) {
        errors.value.password = 'Password must be at least 8 characters long'
        return false
    }

    if (value.length > 75) {
        errors.value.password = 'Password must be less than 75 characters long'
        return false
    }

    if (!/[0-9]/.test(value)) {
        errors.value.password = 'Password must contain a number'
        return false
    }

    if (!/[!@#$%*&^?]/.test(value)) {
        errors.value.password = 'Password must contain a special character: ! @ # $ % * & ^ ?'
        return false
    }

    if (!/[a-zA-Z]/.test(value)) {
        errors.value.password = 'Password must contain alphabetic characters'
        return false
    }

    if (/[^a-zA-Z0-9!@#$%*&^?]/.test(value)) {
        errors.value.password = 'Password contains forbidden characters'
        return false
    }

    errors.value.password = ''
    return true
}

function validateFirstName() {
    const value = form.value.firstName.trim()

    if (!value) {
        errors.value.firstName = 'First name is required'
        return false
    }

    if (value.length < 3) {
        errors.value.firstName = 'First name is too short'
        return false
    }

    if (value.length > 25) {
        errors.value.firstName = 'First name is too long. Max length is 25 characters'
        return false
    }

    if (!/^[A-Za-zÀ-ÖØ-öø-ÿ' -]+$/.test(value)) {
        errors.value.firstName = 'First name contains invalid characters'
        return false
    }

    errors.value.firstName = ''
    return true
}

function validateLastName() {
    const value = form.value.lastName.trim()

    if (!value) {
        errors.value.lastName = 'Last name is required'
        return false
    }

    if (value.length < 3) {
        errors.value.lastName = 'Last name is too short'
        return false
    }

    if (value.length > 25) {
        errors.value.lastName = 'Last name is too long. Max length is 25 characters'
        return false
    }

    if (!/^[A-Za-zÀ-ÖØ-öø-ÿ' -]+$/.test(value)) {
        errors.value.lastName = 'Last name contains invalid characters'
        return false
    }

    errors.value.lastName = ''
    return true
}

function validateDob() {
    const value = form.value.dob

    if (!value) {
        errors.value.dob = 'Date of birth is required'
        return false
    }

    const date = new Date(value + 'T00:00:00')
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    if (Number.isNaN(date.getTime())) {
        errors.value.dob = 'Date of birth is not valid'
        return false
    }

    if (date > today) {
        errors.value.dob = 'Date of birth cannot be in the future'
        return false
    }

    errors.value.dob = ''
    return true
}

async function validateUsername() {
    const value = form.value.username.trim()

    if (!value) {
        errors.value.username = 'Nickname is required'
        return false
    }

    if (value.length < 3) {
        errors.value.username = 'Nickname must be at least 3 characters'
        return false
    }

    if (value.length > 15) {
        errors.value.username = 'Nickname is too long. Max length is 15 characters'
        return false
    }

    if (!/^[A-Za-z0-9_]+$/.test(value)) {
        errors.value.username = 'Nickname can only contain letters, numbers and underscores'
        return false
    }

    try {
        const resp = await fetch(
            `/api/register/checkUsername?username=${encodeURIComponent(value)}`
        )

        if (!resp.ok) {
            errors.value.username = 'Unable to check nickname'
            return false
        }

        const result = await resp.json()

        if (result.usernameAvilable === true) {
            errors.value.username = 'This nickname is already taken'
            return false
        }
    } catch (err) {
        errors.value.username = 'Unable to check nickname'
        return false
    }

    errors.value.username = ''
    return true
}

function validateAvatar() {
    const file = form.value.avatar

    if (!file) {
        errors.value.avatar = ''
        return true
    }

    const allowedTypes = [
        'image/jpeg',
        'image/png',
        'image/gif'
    ]

    if (!allowedTypes.includes(file.type)) {
        errors.value.avatar = 'Avatar must be JPG, PNG or GIF'
        return false
    }

    const maxSize = 5 * 1024 * 1024
    if (file.size > maxSize) {
        errors.value.avatar = 'Avatar must be smaller than 5MB'
        return false
    }

    errors.value.avatar = ''
    return true
}

function validateAbout() {
    const value = form.value.about.trim()

    if (value.length > 100) {
        errors.value.about = 'About me must be 100 characters or less'
        return false
    }

    errors.value.about = ''
    return true
}

async function validateRegisterForm() {
    const emailValid = await validateEmail()
    const passwordValid = validatePassword()
    const firstNameValid = validateFirstName()
    const lastNameValid = validateLastName()
    const dobValid = validateDob()
    const usernameValid = await validateUsername()
    const avatarValid = validateAvatar()
    const aboutValid = validateAbout()

    return (
        emailValid &&
        passwordValid &&
        firstNameValid &&
        lastNameValid &&
        dobValid &&
        usernameValid &&
        avatarValid &&
        aboutValid
    )
}

async function submitRegister(e) {
    e.preventDefault()

    const valid = await validateRegisterForm()
    if (!valid) {
        showNotification('Please fix the highlighted fields', 'error')
        return
    }

    const data = new FormData()
    data.append('email', form.value.email.trim())
    data.append('password', form.value.password)
    data.append('firstName', form.value.firstName.trim())
    data.append('lastName', form.value.lastName.trim())
    data.append('dob', form.value.dob)
    data.append('username', form.value.username.trim())
    data.append('about', form.value.about.trim())

    if (form.value.avatar) {
        data.append('avatar', form.value.avatar)
    }

    try {
        const resp = await fetch('/api/register/submit', {
            method: 'POST',
            body: data
        })

        if (!resp.ok) {
            showNotification('Something went wrong. Please try again', 'error')
            return
        }

        const result = await resp.json()

        if (!result.status) {
            showNotification(result.message || 'Registration failed', 'error')
            return
        }

        showNotification('Account created successfully', 'success')
        resetRegisterForm()
        switchTab('signin')
    } catch (err) {
        showNotification('Network error. Please try again', 'error')
    }
}

function resetRegisterForm() {
    form.value.email = ''
    form.value.password = ''
    form.value.firstName = ''
    form.value.lastName = ''
    form.value.dob = ''
    form.value.username = ''
    form.value.avatar = null
    form.value.about = ''

    errors.value.email = ''
    errors.value.password = ''
    errors.value.firstName = ''
    errors.value.lastName = ''
    errors.value.dob = ''
    errors.value.username = ''
    errors.value.avatar = ''
    errors.value.about = ''
}

function handleAvatar(event) {
    const file = event.target.files[0] || null
    form.value.avatar = file
    validateAvatar()
}

function validateSigninIdentifier() {
    const value = logger.value.identifier.trim()

    if (!value) {
        errors.value.signinIdentifier = 'Email or nickname is required'
        return false
    }

    errors.value.signinIdentifier = ''
    return true
}

function validateSigninPassword() {
    const value = logger.value.password

    if (!value) {
        errors.value.signinPassword = 'Password is required'
        return false
    }

    errors.value.signinPassword = ''
    return true
}

function validateLoginForm() {
    const identifierValid = validateSigninIdentifier()
    const passwordValid = validateSigninPassword()

    return identifierValid && passwordValid
}

function resetLoginForm() {
    logger.value.identifier = ''
    logger.value.password = ''

    errors.value.signinIdentifier = ''
    errors.value.signinPassword = ''
}

async function submitLogin() {
    const valid = validateLoginForm()
    if (!valid) {
        showNotification('Please fix the highlighted fields', 'error')
        return
    }

    const data = {
        identifier: logger.value.identifier.trim(),
        password: logger.value.password
    }

    try {
        const resp = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(data)
        })

        if (!resp.ok) {
            errors.value.signinIdentifier = 'Identifier or password is not correct'
            errors.value.signinPassword = 'Identifier or password is not correct'
            showNotification('Something went wrong. Please try again', 'error')
            return
        }

        const result = await resp.json()

        if (result.status === false) {
            errors.value.signinIdentifier = 'Identifier or password is not correct'
            errors.value.signinPassword = 'Identifier or password is not correct'
            showNotification('Identifier or password is not correct', 'error')
            return
        }

        showNotification('Logged in', 'success')
        resetLoginForm()
    } catch (err) {
        showNotification('Something went wrong. Please try again', 'error')
    }
}
</script>

<template>
    <div class="login-container">
        <Transition name="notif-fade">
            <div v-if="notification.show" class="notification" :class="notification.type">
                <span class="notification-icon">
                    {{ notification.type === 'success' ? '✓' : '!' }}
                </span>
                <span class="notification-text">{{ notification.message }}</span>
            </div>
        </Transition>

        <div class="auth-tabs">
            <button type="button" class="tab-button" :class="{ active: activeTab === 'signin' }"
                @click="switchTab('signin')">
                Sign in
            </button>
            <button type="button" class="tab-button" :class="{ active: activeTab === 'create' }"
                @click="switchTab('create')">
                Create account
            </button>
        </div>
        <div class="form-wrapper">
            <Transition name="form-fade" mode="out-in">
                <form v-if="activeTab === 'signin'" key="signin" class="signin-form" @submit.prevent="submitLogin">
                    <h1 class="form-title">Sign in</h1>
                    <div class="login-row single">
                        <label for="signin-identifier">
                            EMAIL OR NICKNAME *
                        </label>
                        <input type="text" name="identifier" id="signin-identifier" v-model="logger.identifier" required
                            @blur="validateSigninIdentifier" @input="clearError('signinIdentifier')">
                        <p class="error-message" :class="errorClass(errors.signinIdentifier)">
                            {{ errors.signinIdentifier }}
                        </p>
                    </div>
                    <div class="login-row single">
                        <label for="signin-password">
                            PASSWORD *
                        </label>
                        <div class="password-wrapper">
                            <input :type="showSignInPassword ? 'text' : 'password'" name="password" id="signin-password"
                                v-model="logger.password" required @blur="validateSigninPassword"
                                @input="clearError('signinPassword')">
                            <button type="button" class="password-toggle"
                                :aria-label="showSignInPassword ? 'Hide password' : 'Show password'"
                                @click="showSignInPassword = !showSignInPassword">
                                ◌
                            </button>
                        </div>
                        <p class="error-message" :class="errorClass(errors.signinPassword)">
                            {{ errors.signinPassword }}
                        </p>
                    </div>
                    <label class="remember">
                        <input type="checkbox" checked>
                        <span class="checkbox"></span>
                        <span>
                            Keep me signed in on this device
                            <span class="session-text">— session cookie</span>
                        </span>
                    </label>
                    <button type="submit" class="submit-button">
                        Sign in →
                    </button>
                    <button type="button" class="bottom-signin">
                        Forgot your password?
                    </button>
                </form>
                <form v-else key="create" @submit.prevent="submitRegister">
                    <h1 class="form-title">Create account</h1>
                    <div class="login-row single">
                        <label for="email">EMAIL *</label>
                        <input type="email" name="email" id="email" v-model="form.email" maxlength="70" required
                            @blur="validateEmail" @input="clearError('email')">
                        <p class="error-message" :class="errorClass(errors.email)">
                            {{ errors.email }}
                        </p>
                    </div>
                    <div class="login-row single">
                        <label for="password">PASSWORD *</label>
                        <div class="password-wrapper">
                            <input :type="showPassword ? 'text' : 'password'" name="password" id="password"
                                v-model="form.password" maxlength="75" required @blur="validatePassword"
                                @input="clearError('password')">
                            <button type="button" class="password-toggle"
                                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                                @click="showPassword = !showPassword">
                                ◌
                            </button>
                        </div>
                        <p class="error-message" :class="errorClass(errors.password)">
                            {{ errors.password }}
                        </p>
                    </div>
                    <div class="login-row two-columns">
                        <div>
                            <label for="fName">FIRST NAME *</label>
                            <input type="text" name="fName" id="fName" v-model="form.firstName" maxlength="25" required
                                @blur="validateFirstName" @input="clearError('firstName')">
                            <p class="error-message" :class="errorClass(errors.firstName)">
                                {{ errors.firstName }}
                            </p>
                        </div>
                        <div>
                            <label for="lName">LAST NAME *</label>
                            <input type="text" name="lName" id="lName" v-model="form.lastName" maxlength="25" required
                                @blur="validateLastName" @input="clearError('lastName')">
                            <p class="error-message" :class="errorClass(errors.lastName)">
                                {{ errors.lastName }}
                            </p>
                        </div>
                    </div>
                    <div class="login-row two-columns">
                        <div>
                            <label for="dob">DATE OF BIRTH *</label>
                            <input type="date" name="dob" id="dob" v-model="form.dob" required @change="validateDob">
                            <p class="error-message" :class="errorClass(errors.dob)">
                                {{ errors.dob }}
                            </p>
                        </div>
                        <div>
                            <label for="username">
                                NICKNAME *
                            </label>
                            <input type="text" name="username" id="username" v-model="form.username" maxlength="15"
                                placeholder="@noa.png" required @blur="validateUsername"
                                @input="clearError('username')">
                            <p class="error-message" :class="errorClass(errors.username)">
                                {{ errors.username }}
                            </p>
                        </div>
                    </div>
                    <div class="login-row single">
                        <label for="avatar">
                            AVATAR — OPTIONAL · JPG, PNG, GIF
                        </label>
                        <label for="avatar" class="avatar-upload">
                            <input type="file" name="avatar" id="avatar" accept=".jpg,.jpeg,.png,.gif"
                                @change="handleAvatar">
                            <span class="upload-icon">☁</span>
                            <div>
                                <strong>
                                    Drop an image or click to browse
                                </strong>
                            </div>
                        </label>
                        <p class="error-message" :class="errorClass(errors.avatar)">
                            {{ errors.avatar }}
                        </p>
                    </div>
                    <div class="login-row single">
                        <label for="about">
                            ABOUT ME — OPTIONAL
                        </label>
                        <textarea name="about" id="about" v-model="form.about" maxlength="100"
                            placeholder="Trail runner, pixel-art hobbyist, espresso before noon..."
                            @blur="validateAbout" @input="clearError('about')"></textarea>
                        <p class="error-message" :class="errorClass(errors.about)">
                            {{ errors.about }}
                        </p>
                    </div>
                    <label class="remember">
                        <input type="checkbox" checked>
                        <span class="checkbox"></span>
                        <span>
                            Keep me signed in on this device
                            <span class="session-text">— session cookie</span>
                        </span>
                    </label>
                    <button type="submit" class="submit-button">
                        Create my orbit →
                    </button>
                    <button type="button" class="bottom-signin" @click="switchTab('signin')">
                        Already orbiting? Sign in
                    </button>
                </form>
            </Transition>
        </div>
    </div>
</template>

<style>
* {
    box-sizing: border-box;
}

.login-container {
    position: relative;
    width: min(580px, calc(100% - 40px));
    height: 650px;
    margin: 24px auto;
    padding: 34px 39px 10px;
    background: #171b2d;
    border: 1px solid #303754;
    border-radius: 23px;
    color: #f5f5ff;
    font-family: Arial, sans-serif;
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.25);
    z-index: 5;
}

.notification {
    position: absolute;
    top: -18px;
    left: 50%;
    transform: translate(-50%, 0);
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 220px;
    max-width: calc(100% - 40px);
    padding: 12px 18px;
    border-radius: 12px;
    background: #1f2440;
    border: 1px solid #303754;
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.35);
    font-size: 12px;
    font-weight: 600;
    color: #f5f5ff;
    z-index: 20;
}

.notification.success {
    border-color: rgba(94, 224, 158, 0.4);
}

.notification.error {
    border-color: rgba(255, 98, 136, 0.4);
}

.notification-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    font-size: 12px;
    flex-shrink: 0;
}

.notification.success .notification-icon {
    background: rgba(94, 224, 158, 0.15);
    color: #5ee09e;
}

.notification.error .notification-icon {
    background: rgba(255, 98, 136, 0.15);
    color: #ff6288;
}

.notification-text {
    line-height: 1.4;
}

.notif-fade-enter-active,
.notif-fade-leave-active {
    transition:
        opacity 0.35s ease,
        transform 0.35s ease;
}

.notif-fade-enter-from,
.notif-fade-leave-to {
    opacity: 0;
    transform: translate(-50%, -14px);
}

.notif-fade-enter-to,
.notif-fade-leave-from {
    opacity: 1;
    transform: translate(-50%, 0);
}

.auth-tabs {
    display: flex;
    gap: 84px;
    margin-bottom: 25px;
}

.tab-button {
    position: relative;
    border: 0;
    background: none;
    color: #69749a;
    font-size: 14px;
    font-weight: 600;
    padding: 0 0 10px;
    cursor: pointer;
}

.tab-button.active {
    color: #ffffff;
}

.tab-button.active::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 3px;
    border-radius: 5px;
    background: linear-gradient(90deg, #805cff, #ff6288);
}

.form-wrapper {
    width: 100%;
}

.form-wrapper form {
    width: 100%;
}

.form-wrapper form.signin-form {
    min-height: 520px;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.form-fade-enter-active,
.form-fade-leave-active {
    transition:
        opacity 0.25s ease,
        transform 0.25s ease;
}

.form-fade-enter-from {
    opacity: 0;
    transform: translateY(8px);
}

.form-fade-leave-to {
    opacity: 0;
    transform: translateY(-8px);
}

.form-title {
    margin: 0 0 24px;
    color: #ffffff;
    font-size: 24px;
    font-weight: 700;
    text-align: center;
}

.login-row {
    margin-bottom: 14px;
}

.login-row.single {
    display: flex;
    flex-direction: column;
}

.login-row.two-columns {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 9px;
}

.login-row label {
    display: block;
    margin-bottom: 6px;
    color: #69749a;
    font-size: 10px;
    font-weight: 500;
    letter-spacing: 0.8px;
}

.login-row input,
.login-row textarea {
    width: 100%;
    border: 1px solid transparent;
    border-radius: 7px;
    outline: none;
    background: #222741;
    color: #f4f3ff;
    font-family: inherit;
    font-size: 12px;
    transition: 0.2s ease;
}

.login-row input {
    height: 38px;
    padding: 0 13px;
}

.login-row textarea {
    height: 61px;
    padding: 12px 13px;
    resize: none;
}

.login-row input:focus,
.login-row textarea:focus {
    border-color: #835cff;
    box-shadow: 0 0 0 1px rgba(131, 92, 255, 0.15);
}

.login-row input::placeholder,
.login-row textarea::placeholder {
    color: #59627f;
}

.password-wrapper {
    position: relative;
}

.password-wrapper input {
    padding-right: 42px;
}

.password-toggle {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
    width: 28px;
    height: 28px;
    border: 0;
    background: transparent;
    color: #69749a;
    cursor: pointer;
    font-size: 17px;
}

.password-toggle:hover {
    color: #ffffff;
}

.error-message {
    margin: 5px 0 0;
    color: #ff6288;
    font-size: 12px;
    line-height: 1.4;
    overflow: hidden;
    transition:
        opacity 0.2s ease,
        max-height 0.2s ease,
        transform 0.2s ease,
        margin 0.2s ease;
}

.error-hide {
    max-height: 0;
    margin-top: 0;
    opacity: 0;
    transform: translateY(-4px);
}

.error-show {
    max-height: 40px;
    opacity: 1;
    transform: translateY(0);
}

.avatar-upload {
    height: 48px;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 14px;
    border: 2px dashed #805cff;
    border-radius: 8px;
    background: #20243c;
    color: #d8d7e9;
    cursor: pointer;
    transition: 0.2s ease;
}

.avatar-upload:hover {
    background: #252a45;
    border-color: #a27cff;
}

.avatar-upload input {
    display: none;
}

.upload-icon {
    color: #8b62ff;
    font-size: 22px;
}

.avatar-upload strong {
    display: block;
    font-size: 11px;
    font-weight: 500;
    color: #e8e7f2;
}

#about {
    min-height: 61px;
}

.remember {
    display: flex;
    align-items: center;
    gap: 7px;
    margin: 4px 0 14px;
    color: #9da4bd;
    font-size: 10px;
    cursor: pointer;
}

.remember input {
    display: none;
}

.checkbox {
    width: 15px;
    height: 15px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    background: #222741;
    border: 1px solid #59627f;
}

.remember input:checked+.checkbox {
    background: #805cff;
    border-color: #805cff;
}

.remember input:checked+.checkbox::after {
    content: "✓";
    color: white;
    font-size: 11px;
    font-weight: bold;
}

.session-text {
    color: #69749a;
}

.submit-button {
    width: 100%;
    height: 44px;
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

.submit-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

.submit-button:active {
    transform: translateY(0);
}

.bottom-signin {
    display: block;
    margin: 10px auto 0;
    border: 0;
    background: transparent;
    color: #69749a;
    font-family: inherit;
    font-size: 10px;
    cursor: pointer;
}

.bottom-signin:hover {
    color: #9b7cff;
}

@media (max-width: 560px) {
    .login-container {
        width: calc(100% - 24px);
        height: 650px;
        padding: 28px 22px 10px;
        border-radius: 18px;
    }

    .auth-tabs {
        gap: 45px;
    }

    .login-row.two-columns {
        display: block;
    }

    .login-row.two-columns>div {
        margin-bottom: 12px;
    }
}
</style>