<template>
    <section class="profile-section">
        <div class="section-heading">
            <p class="eyebrow">PROFILE</p>
            <h2>About</h2>
        </div>

        <div class="about-card">
            <div class="about-row">
                <span class="label">Username</span>
                <strong>@alexanderson</strong>
            </div>

            <div class="about-row">
                <span class="label">Email</span>
                <strong>alex@example.com</strong>
            </div>

            <div class="about-row">
                <span class="label">Date of birth</span>
                <strong>March 18, 2002</strong>
            </div>

            <div class="about-row">
                <span class="label">Joined</span>
                <strong>August 2026</strong>
            </div>

            <div class="about-description">
                <span class="label">About me</span>

                <p>
                    Building things, learning new stuff and enjoying
                    the little moments. Always looking for something
                    interesting to create.
                </p>
            </div>
        </div>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 6.25rem;
}

.section-heading {
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.about-card {
    padding: 0.625rem var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.about-row {
    display: grid;
    grid-template-columns: 11.25rem 1fr;
    gap: var(--space-5);
    padding: var(--space-4) 0;
    border-bottom: 1px solid var(--color-border);
}

.label {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
}

.about-row strong {
    color: var(--color-text);
    font-size: 0.75rem;
}

.about-description {
    padding: var(--space-4) 0;
}

.about-description p {
    margin: var(--space-3) 0 0;
    color: var(--color-text-soft);
    font-size: 0.8125rem;
    line-height: 1.7;
}

@media (max-width: 34.375rem) {
    .about-row {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>
