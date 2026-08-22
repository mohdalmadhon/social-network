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
.auth-page {
  width: 100%;
  min-height: 100vh;
  background: var(--color-background);
}

.login-side {
  display: grid;
  width: 100%;
  min-width: 0;
  padding: var(--space-5) var(--space-3) var(--space-6);
  place-items: center;
}

@media (min-width: 48rem) {
  .login-side {
    padding: var(--space-6);
  }
}

@media (min-width: 64rem) {
  .auth-page {
    display: grid;
    grid-template-columns: minmax(28rem, 36%) minmax(0, 1fr);
  }

  .login-side {
    min-height: 100vh;
    padding: clamp(2rem, 4vw, 5rem);
  }
}
</style>
