<script setup>
import AuthContainer from '@/components/auth/AuthContainer.vue';
import WelcomeSection from '@/components/auth/WelcomeSection.vue';
</script>

<template>
    <main class="auth-page">
        <WelcomeSection />
        <section class="login-side">
            <AuthContainer />
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
    display: grid;
    min-height: 100vh;
    width: 100%;
    background: var(--color-background);
}

.login-side {
    width: 100%;
    min-height: auto;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: clamp(1.5rem, 4vw, 3rem);
}

@media (min-width: 56.25rem) {
    .auth-page { grid-template-columns: minmax(20rem, .8fr) minmax(28rem, 1.2fr); }
    .login-side { min-height: 100vh; }
}
</style>
