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
            <p class="eyebrow">PERSONAL</p>
            <h3>Information</h3>
        </div>

        <div class="info-list">
            <div class="info-item">
                <span>Work</span>
                <p>{{ about.Work || 'Not provided' }}</p>
            </div>

            <div class="info-item">
                <span>Education</span>
                <p>{{ about.Education || 'Not provided' }}</p>
            </div>

            <div class="info-item">
                <span>Hobbies</span>
                <p>{{ about.Hobbies || 'Not provided' }}</p>
            </div>

            <div class="info-item">
                <span>Interests</span>
                <p>{{ about.Intrests || 'Not provided' }}</p>
            </div>

            <div class="info-item">
                <span>Travel</span>
                <p>{{ about.Travel || 'Not provided' }}</p>
            </div>
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

.info-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
}

.info-item {
    display: grid;
    grid-template-columns: minmax(5.625rem, 8.125rem) 1fr;
    gap: var(--space-5);
    padding-bottom: var(--space-4);
    border-bottom: 1px solid var(--color-border);
}

.info-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
}

.info-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.875rem;
    line-height: 1.5;
    word-break: break-word;
}

@media (max-width: 31.25rem) {
    .info-item {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>
