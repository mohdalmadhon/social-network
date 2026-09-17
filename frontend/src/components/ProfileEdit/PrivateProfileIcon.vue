<script setup>
import IconGlyph from '@/components/layout/IconGlyph.vue'
</script>

<template>
<<<<<<< HEAD
    <div class="private-container">
        <p>This account is private</p>
    </div>
</template>

<style scoped>
.private-container {
    display: flex;
    margin-top: 100px;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: clamp(24px, 5vw, 40px) clamp(16px, 4vw, 24px);
    border: 2px solid var(--main-color);
    border-radius: 7px;
    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
    text-align: center;
}

.icon {
    font-size: clamp(24px, 4vw, 32px);
}

.private-container p {
    margin: 0;
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(11px, 1.6vw, 13px);
    font-weight: 600;
}
</style>
=======
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
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
