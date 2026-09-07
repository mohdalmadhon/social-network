<template>
    <section class="profile-section">
        <div class="section-heading">
            <div>
                <p class="eyebrow">PEOPLE</p>
                <h2>Friends</h2>
            </div>

            <span class="count">342 friends</span>
        </div>

        <div class="people-grid">
            <article class="person-card">
                <div class="person-avatar">J</div>
                <h3>James Wilson</h3>
                <p>@jameswilson</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">S</div>
                <h3>Sarah Miller</h3>
                <p>@sarahmiller</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">M</div>
                <h3>Michael Brown</h3>
                <p>@michaelbrown</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">E</div>
                <h3>Emma Davis</h3>
                <p>@emmadavis</p>
                <button>View profile</button>
            </article>
        </div>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 6.25rem;
}

.section-heading {
    display: flex;
    justify-content: space-between;
    align-items: end;
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.count {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

.people-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--space-4);
}

.person-card {
    padding: var(--space-5) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    text-align: center;
    box-shadow: var(--shadow-raised);
}

.person-avatar {
    width: 4.6875rem;
    height: 4.6875rem;
    margin: 0 auto var(--space-3);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.875rem;
}

.person-card h3 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.8125rem;
}

.person-card p {
    margin: var(--space-1) 0 var(--space-4);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.5rem;
}

.person-card button {
    width: 100%;
    padding: var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-surface-raised);
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.6875rem;
    font-weight: 600;
    cursor: pointer;
    transition: filter 0.15s ease;
}

.person-card button:hover {
    filter: brightness(1.15);
}

@media (max-width: 53.125rem) {
    .people-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
</style>
