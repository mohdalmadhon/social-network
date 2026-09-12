<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const route = useRoute()
const router = useRouter()

const code = computed(() => {
  return Number(route.query.code) || 500
})

const message = computed(() => {
  return route.query.message || 'Something went wrong.'
})

const title = computed(() => {
  const titles = {
    400: 'Bad Request',
    401: 'Unauthorized',
    403: 'Access Denied',
    404: 'Page Not Found',
    408: 'Request Timeout',
    429: 'Too Many Requests',
    500: 'Something Went Wrong',
    502: 'Bad Gateway',
    503: 'Service Unavailable',
  }

  return titles[code.value] || 'Something Went Wrong'
})

function goHome() {
  router.push('/home-feed')
}

function goBack() {
  router.back()
}
</script>

<template>
  <main class="error-page">
    <div class="error-background">
      <span class="orb orb--one"></span>
      <span class="orb orb--two"></span>
      <span class="orb orb--three"></span>
    </div>

    <section class="error-card">
      <div class="error-icon">
        <span>!</span>
      </div>

      <p class="error-code">
        {{ code }}
      </p>

      <h1>
        {{ title }}
      </h1>

      <p class="error-message">
        {{ message }}
      </p>

      <div class="error-actions">
        <button
          class="error-button error-button--primary"
          type="button"
          @click="goHome"
        >
          <IconGlyph name="home" :size="17" />
          <span>Go home</span>
        </button>

        <button
          class="error-button error-button--secondary"
          type="button"
          @click="goBack"
        >
          <IconGlyph name="arrowLeft" :size="17" />
          <span>Go back</span>
        </button>
      </div>

      <p class="error-footer">
        orbit
        <span>·</span>
        Something went off course.
      </p>
    </section>
  </main>
</template>

<style scoped>
@import '@/styles/global.css';
@import '@/styles/variables.css';

.error-page {
  position: relative;
  display: grid;
  min-height: 100vh;
  overflow: hidden;
  place-items: center;
  padding: var(--space-5);
  background:
    radial-gradient(
      circle at 50% 40%,
      rgb(124 92 255 / 10%),
      transparent 35%
    ),
    var(--color-background);
  isolation: isolate;
}

.error-background {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.orb {
  position: absolute;
  display: block;
  border: 1px solid rgb(124 92 255 / 12%);
  border-radius: 50%;
}

.orb--one {
  top: -12rem;
  right: -8rem;
  width: 32rem;
  height: 32rem;
}

.orb--two {
  bottom: -18rem;
  left: -12rem;
  width: 42rem;
  height: 42rem;
  border-color: rgb(99 230 190 / 8%);
}

.orb--three {
  top: 50%;
  left: 50%;
  width: 22rem;
  height: 22rem;
  transform: translate(-50%, -50%);
  border-color: rgb(255 107 125 / 6%);
}

.error-card {
  position: relative;
  z-index: 2;

  display: flex;
  flex-direction: column;
  align-items: center;

  width: min(100%, 32rem);

  padding: clamp(
    var(--space-6),
    6vw,
    var(--space-10)
  );

  background: rgb(20 23 42 / 88%);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-large);

  box-shadow:
    0 2rem 5rem rgb(0 0 0 / 25%),
    0 0 0 1px rgb(255 255 255 / 2%);

  backdrop-filter: blur(1.5rem);

  text-align: center;
}

.error-icon {
  display: grid;
  width: 4.5rem;
  height: 4.5rem;
  margin-bottom: var(--space-4);
  place-items: center;
  border-radius: 50%;
  background: rgb(124 92 255 / 10%);
  border: 1px solid rgb(124 92 255 / 25%);
}

.error-icon span {
  display: grid;
  width: 2.75rem;
  height: 2.75rem;
  place-items: center;
  border: 2px solid var(--color-violet);
  border-radius: 50%;
  color: var(--color-violet);
  font-family: var(--font-display);
  font-size: 1.3rem;
  font-weight: 700;
}

.error-code {
  margin: 0;
  color: var(--color-violet);
  font-family: var(--font-meta);
  font-size: clamp(4rem, 15vw, 7rem);
  font-weight: 700;
  line-height: 0.9;
  letter-spacing: -0.06em;
  opacity: 0.18;
}

.error-card h1 {
  margin: var(--space-3) 0 0;
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: clamp(1.5rem, 5vw, 2rem);
  font-weight: 700;
}

.error-message {
  max-width: 25rem;
  margin: var(--space-3) 0 0;
  color: var(--color-text-muted);
  font-size: 0.95rem;
  line-height: 1.7;
}

.error-actions {
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: var(--space-2);
  margin-top: var(--space-6);
}

.error-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: var(--touch-target);
  padding-inline: var(--space-5);
  border-radius: var(--radius-small);
  font: inherit;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
}

.error-button span {
  margin-left: var(--space-2);
}

.error-button--primary {
  border: 0;
  background: var(--gradient-action);
  color: white;
}

.error-button--secondary {
  border: 1px solid var(--color-border);
  background: var(--color-input);
  color: var(--color-text-muted);
}

.error-footer {
  display: flex;
  gap: var(--space-2);
  margin: var(--space-6) 0 0;
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.7rem;
}

@media (min-width: 30rem) {
  .error-actions {
    flex-direction: row;
  }

  .error-button {
    flex: 1;
  }
}
</style>
