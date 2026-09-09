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
    margin-bottom: clamp(18px, 3vw, 25px);
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--color-mint);
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

.info-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.info-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: var(--space-5);
    padding-bottom: var(--space-3);
    border-bottom: 1px solid var(--color-border);
}

.info-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 10px);
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.info-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-display);
    font-size: clamp(12px, 1.8vw, 14px);
    line-height: 1.5;
    word-break: break-word;
}

@media (max-width: 500px) {
    .info-item {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>