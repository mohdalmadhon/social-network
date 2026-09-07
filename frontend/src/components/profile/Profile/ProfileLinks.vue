<script setup>
defineProps({
    about: {
        type: Object,
        required: true
    }
});
</script>

<template>
    <div class="content-section">
        <div class="content-heading">
            <p class="eyebrow">SOCIAL</p>
            <h3>Links</h3>
        </div>

        <div class="links-list">
            <a
                v-if="about.Website"
                :href="about.Website"
                target="_blank"
                rel="noopener noreferrer"
                class="link-item"
            >
                <span>Website</span>
                <p>{{ about.Website }}</p>
            </a>

            <a
                v-if="about.Linkedin"
                :href="about.Linkedin"
                target="_blank"
                rel="noopener noreferrer"
                class="link-item"
            >
                <span>LinkedIn</span>
                <p>{{ about.Linkedin }}</p>
            </a>

            <a
                v-if="about.Twitter"
                :href="about.Twitter"
                target="_blank"
                rel="noopener noreferrer"
                class="link-item"
            >
                <span>Twitter / X</span>
                <p>{{ about.Twitter }}</p>
            </a>

            <a
                v-if="about.Instgram"
                :href="about.Instgram"
                target="_blank"
                rel="noopener noreferrer"
                class="link-item"
            >
                <span>Instagram</span>
                <p>{{ about.Instgram }}</p>
            </a>

            <p
                v-if="!about.Website && !about.Linkedin && !about.Twitter && !about.Instgram"
                class="empty"
            >
                No social links provided.
            </p>
        </div>
    </div>
</template>

<style scoped>
.content-section {
    width: 100%;
}

.content-heading {
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
}

h3 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.5625rem;
    color: var(--color-text);
}

.links-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.link-item {
    display: grid;
    grid-template-columns: minmax(5.625rem, 8.125rem) 1fr;
    gap: var(--space-5);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-raised);
    color: var(--color-text-soft);
    text-decoration: none;
    transition: border-color 0.15s ease, transform 0.15s ease;
}

.link-item:hover {
    transform: translateY(-1px);
    border-color: var(--color-violet);
}

.link-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
}

.link-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    word-break: break-all;
}

.empty {
    margin: 0;
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

@media (max-width: 31.25rem) {
    .link-item {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>
