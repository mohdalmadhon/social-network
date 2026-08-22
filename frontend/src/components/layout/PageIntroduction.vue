<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <section class="page-introduction orbit-surface">
    <p class="orbit-meta">{{ label }}</p>
    <h1>{{ title }}</h1>
    <p>{{ description }}</p>
  </section>
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
