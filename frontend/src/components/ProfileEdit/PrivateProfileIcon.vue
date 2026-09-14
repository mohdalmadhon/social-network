<script setup>
import IconGlyph from '@/components/layout/IconGlyph.vue'
</script>

<template>
  <section class="private-profile orbit-surface" aria-labelledby="private-profile-heading">
    <span class="private-profile__icon"><IconGlyph name="lock" :size="24" /></span>
    <h2 id="private-profile-heading">This profile is private</h2>
    <p>Follow this member and wait for approval to view their posts and profile activity.</p>
  </section>
</template>

<style scoped>
.private-profile { display: grid; justify-items: center; gap: var(--space-3); padding: var(--space-7) var(--space-5); text-align: center; }
.private-profile__icon { display: grid; width: 3.5rem; height: 3.5rem; place-items: center; border-radius: 50%; background: var(--color-surface-amber); color: var(--color-amber-soft); }
.private-profile h2 { margin: 0; font-family: var(--font-display); font-size: 1.35rem; letter-spacing: 0; }
.private-profile p { max-width: 34rem; margin: 0; color: var(--color-text-muted); line-height: 1.6; }
</style>
