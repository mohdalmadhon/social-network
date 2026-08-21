<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
</script>

<template>
  <AuthenticatedLayout active-page="home">
    <section class="page-introduction orbit-surface">
      <p class="orbit-meta">Home feed</p>
      <h1>Welcome to your orbit</h1>
      <p>The shared page structure is ready. The post composer and feed cards come next.</p>
    </section>
  </AuthenticatedLayout>
</template>

<style scoped>
.page-introduction {
  max-width: 48rem;
  padding: var(--space-5);
}

.page-introduction h1 {
  margin: var(--space-2) 0;
  font-family: var(--font-display);
  font-size: clamp(2rem, 7vw, 2.75rem);
  line-height: 1.15;
}

.page-introduction p:last-child {
  max-width: 55ch;
  margin: 0;
  color: var(--color-text-muted);
}
</style>
