<script setup>
import IconGlyph from '@/components/layout/IconGlyph.vue'
</script>

<template>
    <div class="private-container">
        <p><span class="icon"><IconGlyph name="lock" :size="16" /></span>This account is private</p>
    </div>
</template>

<style scoped>
@import '../../styles/global.css';
@import '../../styles/variables.css';

.private-container {
    display: flex;
    margin-top: 100px;
    flex-direction: column;
    align-items: center;
    gap: var(--space-3);
    padding: clamp(24px, 5vw, 40px) clamp(16px, 4vw, 24px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-coral);
    box-shadow: var(--shadow-raised);
    border-top: 3px solid transparent;
    border-image: var(--gradient-sunset) 1;
    text-align: center;
}

.icon {
    margin-right: var(--space-2);
    color: var(--color-amber);
    font-size: clamp(14px, 2vw, 16px);
}

.private-container p {
    margin: 0;
    display: flex;
    align-items: center;
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: clamp(11px, 1.6vw, 13px);
    font-weight: 600;
}
</style>
