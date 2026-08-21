<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'

defineProps({
  activePage: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <AuthenticatedLayout :active-page="activePage">
    <section class="placeholder orbit-surface">
      <p class="orbit-meta">Orbit page</p>
      <h1>{{ title }}</h1>
      <p>This page entry is ready for its owning feature. No unfinished behavior is being simulated.</p>
    </section>
  </AuthenticatedLayout>
</template>

<style scoped>
.placeholder {
  max-width: 42rem;
  padding: var(--space-5);
}

.placeholder h1 {
  margin: var(--space-2) 0;
  font-family: var(--font-display);
  font-size: clamp(2rem, 7vw, 2.75rem);
}

.placeholder p:last-child {
  margin-bottom: 0;
  color: var(--color-text-muted);
}
</style>
