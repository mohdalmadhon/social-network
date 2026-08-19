<script setup>
import { ref } from 'vue'

const activeTab = ref('create')
const showPassword = ref(false)
const showSignInPassword = ref(false)

const emailInput = document.getElementById('email');

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

function switchTab(tab) {
    activeTab.value = tab
}

function errorClass(error) {
    return error ? 'error-show' : 'error-hide'
}

async function validateEmail() {
    const value = email.value.trim()

    if (value.length < 3) {
        errors.value.email = 'Email is too short'
        return
    }
    if (value.length > 75) {
        errors.value.email = "Email is too long. Max Length is 75 characters"
        return
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if(!emailRegex.test(value)) {
        errors.value.email = "Email is not valid"
        return
    }
    
    try {
        const resp = await fetch(`/api/checkEmail?email=${encodeURIComponent(value)}`);
        if(!resp.ok) {
            console.error(err);
            // error notification
            return;
        }

        const result = await resp.json();
    } catch (err) {
        // error notification
    }

    errors.value.email = ''
}

</script>

<template>
    <div class="login-container">
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
                <form v-if="activeTab === 'signin'" key="signin" class="signin-form">
                    <h1 class="form-title">Sign in</h1>

                    <div class="login-row single">
                        <label for="signin-identifier">
                            EMAIL OR NICKNAME *
                        </label>

                        <input type="text" name="identifier" id="signin-identifier" required>

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
                                required>

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

                <form v-else key="create">
                    <h1 class="form-title">Create account</h1>

                    <div class="login-row single">
                        <label for="email">EMAIL *</label>

                        <input type="email" name="email" id="email" required @change="validateEmail">

                        <p class="error-message" :class="errorClass(errors.email)">
                            {{ errors.email }}
                        </p>
                    </div>

                    <div class="login-row single">
                        <label for="password">PASSWORD *</label>

                        <div class="password-wrapper">
                            <input :type="showPassword ? 'text' : 'password'" name="password" id="password" required>

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

                            <input type="text" name="fName" id="fName" required>

                            <p class="error-message" :class="errorClass(errors.firstName)">
                                {{ errors.firstName }}
                            </p>
                        </div>

                        <div>
                            <label for="lName">LAST NAME *</label>

                            <input type="text" name="lName" id="lName" required>

                            <p class="error-message" :class="errorClass(errors.lastName)">
                                {{ errors.lastName }}
                            </p>
                        </div>
                    </div>

                    <div class="login-row two-columns">
                        <div>
                            <label for="dob">DATE OF BIRTH *</label>

                            <input type="date" name="dob" id="dob" required>

                            <p class="error-message" :class="errorClass(errors.dob)">
                                {{ errors.dob }}
                            </p>
                        </div>

                        <div>
                            <label for="username">
                                NICKNAME — OPTIONAL
                            </label>

                            <input type="text" name="username" id="username" placeholder="@noa.png">

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
                            <input type="file" name="avatar" id="avatar" accept=".jpg,.jpeg,.png,.gif">

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

                        <textarea name="about" id="about" maxlength="100"
                            placeholder="Trail runner, pixel-art hobbyist, espresso before noon..."></textarea>

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