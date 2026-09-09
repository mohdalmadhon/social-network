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
    margin-bottom: clamp(18px, 3vw, 25px);
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--color-cyan);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
    text-transform: uppercase;
}

h3 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: clamp(20px, 3vw, 25px);
}

.links-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.link-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: var(--space-5);
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-left: 3px solid var(--color-amber);
    border-radius: var(--radius-small);
    background: var(--color-surface);
    color: var(--color-text);
    text-decoration: none;
    transition: 0.15s ease;
}

.link-item:nth-of-type(2) {
    border-left-color: var(--color-blue);
}

.link-item:nth-of-type(3) {
    border-left-color: var(--color-cyan);
}

.link-item:nth-of-type(4) {
    border-left-color: var(--color-magenta);
}

.link-item:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-raised);
    border-color: var(--color-violet);
}

.link-item:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.link-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 10px);
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.link-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-display);
    font-size: clamp(11px, 1.7vw, 13px);
    word-break: break-all;
}

.empty {
    margin: 0;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 10px;
}

@media (max-width: 500px) {
    .link-item {
        grid-template-columns: 1fr;
        gap: var(--space-2);
        border-left-width: 3px;
    }
}
</style>