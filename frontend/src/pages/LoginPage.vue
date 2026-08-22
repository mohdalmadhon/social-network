<script setup>
import LoginForm from '@/components/login/LoginForm.vue'
import Sider from '@/components/login/Sider.vue'
</script>

<template>
    <main class="auth-page">
        <Sider />
        <section class="login-side">
            <LoginForm />
        </section>
    </main>
</template>

<style>
* {
    box-sizing: border-box;
}

html,
body,
#app {
    width: 100%;
    min-height: 100%;
    margin: 0;
}

.auth-page {
    display: flex;
    min-height: 100vh;
    width: 100%;
    background: #0f1020;
}

.login-side {
    width: 70%;
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 30px;
}

@media (max-width: 900px) {
    .auth-page {
        flex-direction: column;
    }

    .login-side {
        width: 100%;
    }
}
</style>