<template>
    <div class="private-container">
        <p>This account is private</p>
    </div>
</template>

<style scoped>
.private-container {
    display: flex;
    margin-top: var(--space-7);
    flex-direction: column;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-6) var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    text-align: center;
}

.icon {
    font-size: clamp(1.5rem, 4vw, 2rem);
}

.private-container p {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: 0.8125rem;
    font-weight: 600;
}
</style>
