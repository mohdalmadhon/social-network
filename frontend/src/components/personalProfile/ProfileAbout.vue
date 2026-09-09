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
    scroll-margin-top: 100px;
}

.section-heading {
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 9px;
    letter-spacing: 2px;
    text-transform: uppercase;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 29px;
}

.about-card {
    padding: var(--space-3) var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-violet);
    box-shadow: var(--shadow-raised);
}

.about-row {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: var(--space-5);
    padding: var(--space-5) 0;
    border-bottom: 1px solid var(--color-border);
}

.about-row:nth-child(1) strong {
    color: var(--color-cyan-soft);
}

.about-row:nth-child(2) strong {
    color: var(--color-blue-soft);
}

.about-row:nth-child(3) strong {
    color: var(--color-amber-soft);
}

.about-row:nth-child(4) strong {
    color: var(--color-mint-soft);
}

.label {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 9px;
    font-weight: 600;
}

.about-row strong {
    color: var(--color-text);
    font-size: 12px;
}

.about-description {
    padding: var(--space-5) 0;
}

.about-description p {
    margin: var(--space-3) 0 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 13px;
    line-height: 1.7;
}

@media (max-width: 550px) {
    .about-row {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>